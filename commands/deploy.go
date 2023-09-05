package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	d "github.com/thecodekitchen/flo/templates/deployment"
	utils "github.com/thecodekitchen/flo/utils"
)

func GenerateDeploymentConfig(project string, registry string) {
	os.Chdir("app")
	os.WriteFile("kubernetes-manifest.yaml", d.KubernetesManifestBytes(project, registry), 0666)
	os.Chdir("..")
}

func BuildAndPushToRegistry(project string, registry string) {
	os.Chdir("app")
	build := exec.Command("docker", "build", "-t", registry+"/"+project, ".")
	build_err := build.Run()
	utils.PanicIf(build_err)

	push := exec.Command("docker", "push", registry+"/"+project)
	push_err := push.Run()
	utils.PanicIf(push_err)
	os.Chdir("..")
}

func DeployToLocalCluster(project string) {
	fmt.Println("Checking for existing cluster and replacing it if present ...")
	check_exists := exec.Command("kind", "get", "clusters")
	check_output, err := check_exists.CombinedOutput()
	utils.PanicIf(err)

	// Check if cluster already exists.
	// If it does, delete it and create a new one.
	cluster := strings.Replace(project, "_", "-", -1)
	if !strings.Contains(string(check_output), cluster) {
		create := exec.Command("kind", "create", "cluster", "-n", cluster)
		create_err := create.Run()
		utils.PanicIf(create_err)
	} else {
		delete := exec.Command("kind", "delete", "cluster", "-n", cluster)
		delete_err := delete.Run()
		utils.PanicIf(delete_err)
		create := exec.Command("kind", "create", "cluster", "-n", cluster)
		create_err := create.Run()
		utils.PanicIf(create_err)
	}

	// Install custom resource definitions for tikv cluster
	fmt.Println("Installing custom resource definitions ...")
	crd := exec.Command("kubectl", "create", "-f", "https://raw.githubusercontent.com/pingcap/tidb-operator/v1.4.5/manifests/crd.yaml")
	crd_err := crd.Run()
	utils.PanicIf(crd_err)

	// Add Pingcap Helm repo to locally available repos
	fmt.Println("Installing tidb cluster from Helm chart ...")
	add_pingcap_chart := exec.Command("helm", "repo", "add", "pingcap", "https://charts.pingcap.org")
	add_pingcap_chart_out, add_pingcap_chart_err := add_pingcap_chart.CombinedOutput()
	fmt.Printf(string(add_pingcap_chart_out))
	utils.PrintError(add_pingcap_chart_err)

	// Update Helm repos
	update := exec.Command("helm", "repo", "update")
	update_out, update_err := update.CombinedOutput()
	fmt.Printf(string(update_out))
	utils.PanicIf(update_err)

	// Install tidb cluster from Helm chart
	// helm install -n tidb-operator --create-namespace tidb-operator pingcap/tidb-operator --version v1.4.5
	install_tidb := exec.Command("helm", "install", "-n", "tidb-operator", "--create-namespace", "tidb-operator", "pingcap/tidb-operator", "--version", "v1.4.5")
	install_tidb_out, install_tidb_err := install_tidb.CombinedOutput()
	fmt.Printf(string(install_tidb_out))
	utils.PanicIf(install_tidb_err)

	fmt.Println("Waiting for tidb cluster to come online ...")
	// Wait for tidb pods to come online
	// kubectl wait --for=condition=Ready pods -l app.kubernetes.io/instance=tidb-operator --timeout=5m -n tidb-operator
	tidb_wait := exec.Command("kubectl", "wait", "-n", "tidb-operator", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/instance=tidb-operator", "--timeout=5m")
	tidb_wait_out, tidb_wait_err := tidb_wait.CombinedOutput()
	fmt.Printf(string(tidb_wait_out))
	utils.PanicIf(tidb_wait_err)

	// create tikv namespace
	fmt.Println("Deploying tikv pods ... ")
	tikv_ns := exec.Command("kubectl", "create", "ns", "tikv")
	tikv_ns_out, tikv_ns_err := tikv_ns.CombinedOutput()
	fmt.Printf(string(tikv_ns_out))
	utils.PanicIf(tikv_ns_err)
	// Deploy tikv pods
	// kubectl apply -n tikv -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.4.5/examples/basic/tidb-cluster.yaml
	tikv_deploy := exec.Command("kubectl", "apply", "-n", "tikv", "-f", "https://raw.githubusercontent.com/pingcap/tidb-operator/v1.4.5/examples/basic/tidb-cluster.yaml")
	tikv_deploy_out, tikv_deploy_err := tikv_deploy.CombinedOutput()
	fmt.Printf(string(tikv_deploy_out))
	utils.PanicIf(tikv_deploy_err)
	// Wait for tikv pods to come online
	// kubectl wait -n tikv --for=condition=Ready tidbcluster/basic --timeout=5m
	fmt.Println("Waiting for tikv pods to come online ... ")
	tikv_wait := exec.Command("kubectl", "wait", "-n", "tikv", "--for=condition=Ready", "tidbcluster/basic", "--timeout=5m")
	tikv_wait_out, tikv_wait_err := tikv_wait.CombinedOutput()
	fmt.Printf(string(tikv_wait_out))
	utils.PanicIf(tikv_wait_err)
	// Get tikv URL
	// export TIKV_URL=$(kubectl get -n tikv svc/basic-pd -o jsonpath='{.spec.clusterIP}:2379')
	get_tikv_url := exec.Command("kubectl", "get", "-n", "tikv", "svc/basic-pd", "-o", "jsonpath='{.spec.clusterIP}:2379'")
	get_tikv_url_out, get_tikv_url_err := get_tikv_url.CombinedOutput()
	utils.PanicIf(get_tikv_url_err)
	tikv_url := string(get_tikv_url_out)

	fmt.Println("Installing SurrealDB from Helm chart ...")
	// Add SurrealDB Helm repo
	// helm repo add surrealdb https://helm.surrealdb.com
	add_surrealdb_chart := exec.Command("helm", "repo", "add", "surrealdb", "https://helm.surrealdb.com")
	add_surrealdb_chart_err := add_surrealdb_chart.Run()
	utils.PrintError(add_surrealdb_chart_err)
	// Update Helm repos (repeating earlier update command)
	update_err = update.Run()
	utils.PanicIf(update_err)
	// install SurrealDB from Helm chart
	// helm install surrealdb-tikv surrealdb/surrealdb --set surrealdb.path=tikv://$TIKV_URL --set service.port=8000
	install_surrealdb := exec.Command("helm", "install", "surrealdb-tikv", "surrealdb/surrealdb", "--set", "surrealdb.path=tikv://"+tikv_url, "--set", "service.port=8000")
	install_surrealdb_err := install_surrealdb.Run()
	utils.PanicIf(install_surrealdb_err)
	// Deploy API defined in local kubernetes-manifest.yaml
	// kubectl apply -f kubernetes-manifest.yaml
	fmt.Println("Deploying API from local manifest ...")
	deploy_api := exec.Command("kubectl", "apply", "-f", "kubernetes-manifest.yaml")
	deploy_api_err := deploy_api.Run()
	utils.PanicIf(deploy_api_err)
	// Wait for API pod to come online
	// kubectl wait --for=condition=Ready pods -l app=api --timeout=5m
	fmt.Println("Waiting for API pod to come online ...")
	api_wait := exec.Command("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app=api", "--timeout=5m")
	api_wait_err := api_wait.Run()
	utils.PanicIf(api_wait_err)
	fmt.Println("Deployment successful! To expose the API on http://localhost:8080, run 'flo run -kind'.")
}
