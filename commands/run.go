package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thecodekitchen/flo/utils"
)

func RunBackend(project string, kind bool) {
	if kind {
		run := exec.Command("kubectl", "port-forward", "svc/api-service", "8080")
		run_out, run_err := run.CombinedOutput()
		fmt.Printf(string(run_out))
		utils.PanicIf(run_err)
	} else {
		os.Chdir("app")
		run := exec.Command("go", "run", "main.go", "&")
		run_out, run_err := run.CombinedOutput()
		fmt.Printf(string(run_out))
		utils.PanicIf(run_err)
	}
}

func RunFrontend(project string) {
	os.Chdir(project + "_fe")
	run := exec.Command("flutter", "run", "-d", "chrome")
	run_out, err := run.CombinedOutput()
	fmt.Printf(string(run_out))
	utils.PanicIf(err)
}
