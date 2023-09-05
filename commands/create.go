package commands

import (
	"os"
	"os/exec"

	d "github.com/thecodekitchen/flo/templates/deployment"
	f "github.com/thecodekitchen/flo/templates/flutter"
	p "github.com/thecodekitchen/flo/templates/flutter/pages"
	g "github.com/thecodekitchen/flo/templates/gofiles"
	utils "github.com/thecodekitchen/flo/utils"
)

func GenerateBackendApi(project string, lang string) {
	os.WriteFile("./models.json", d.TestModelsBytes(), 0666)
	if lang == "go" {
		generate_go_api(project)
	}
}

func GenerateFlutterApp(project string) {
	err := exec.Command("flutter", "create", project+"_fe").Run()
	utils.PanicIf(err)
	os.Chdir(project + "_fe")
	os.WriteFile("./.env", d.EnvFileBytes(), 0666)
	os.WriteFile("./lib/backend.dart", f.BackendFileBytes(), 0666)
	os.WriteFile("./lib/models.dart", f.ModelsFileBytes(), 0666)
	os.WriteFile("./lib/main.dart", f.MainFileBytes(), 0666)
	os.WriteFile("./lib/pages/admin_login_page.dart", p.AdminLoginPageBytes(), 0666)
	os.WriteFile("./lib/pages/admin_page.dart", p.AdminPageBytes(), 0666)
	os.WriteFile("./lib/pages/login_page.dart", p.LoginPageBytes(), 0666)
	os.WriteFile("./lib/pages/error_page.dart", p.ErrorPageBytes(), 0666)
	os.WriteFile("./lib/pages/home_page.dart", p.HomePageBytes(), 0666)
	os.WriteFile("./android/app/build.gradle", f.BuildGradleBytes(), 0666)
	os.WriteFile("./android/app/src/main/AndroidManifest.xml", f.AndroidManifestBytes(project), 0666)
	os.WriteFile("./ios/Runner/Info.plist", f.IosRunnerBytes(project), 0666)
	os.WriteFile("pubspec.yaml", f.PubspecFileBytes(project), 0666)
	pub_get := exec.Command("flutter", "pub", "get")
	err = pub_get.Run()
	utils.PanicIf(err)
}

func generate_go_api(project string) {
	err := os.Mkdir("app", 0750)
	utils.ExistsError(err)
	os.Chdir("./app")
	err = os.Mkdir("./models", 0750)
	utils.ExistsError(err)

	os.WriteFile("./models/models.go", g.ModelsFileBytes(), 0666)
	os.WriteFile("../app/main.go", g.MainFileBytes(project), 0666)
	os.WriteFile("../app/go.mod", g.ModFileBytes(project), 0666)
	os.WriteFile("Dockerfile", d.DockerfileBytes(), 0666)
	os.WriteFile(".env", d.EnvFileBytes(), 0666)

	cmd := exec.Command("go", "get")
	err = cmd.Run()
	utils.PanicIf(err)
	os.Chdir("..")
}
