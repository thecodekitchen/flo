package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/thecodekitchen/flo/commands"
	"github.com/thecodekitchen/flo/utils"
)

func main() {
	registry := ""
	kind := false
	lang := ""
	flag.StringVar(&lang, "lang", "go", "The language for your backend")
	flag.BoolVar(&kind, "kind", false, "set this flag in order to run the backend from a local kubernetes cluster via port-forwarding.")
	flag.StringVar(&registry, "registry", "",
		`The name of your container registry that Kubernetes can pull from.
Required in order to run 'flo deploy' in your project.
To configure this later, run 'flo register <registry_name>'.
All deployment configuration files will be regenerated,
so any changes you made to them prior to registering will be lost.`)

	flag.Parse()

	switch command := flag.Arg(0); command {
	case "create":
		{

			if len(flag.Args()) < 2 {
				panic("Need to specify a project name")
			}
			project := flag.Arg(1)
			fmt.Printf(`Welcome to Flo!
Creating ` + project + ` back end ...`)

			// Extra argument after project name will be
			// interpreted as a backend language choice
			// default is Go.
			os.Mkdir(project, 0750)
			os.Chdir(project)

			cmd.GenerateBackendApi(project, lang)
			if registry == "" {
				fmt.Println("No artifact registry specified. Will not generate Kubernetes manifest.")
				fmt.Println("To generate deployment configuration with an artifact registry later, run 'flo register <new_registry_name>' from the project root.")
			} else {
				fmt.Println("Creating deployment configuration for artifact registry ", registry)
				cmd.GenerateDeploymentConfig(project, registry)
			}

			fmt.Printf(`Creating ` + project + ` front end ...`)
			cmd.GenerateFlutterApp(project)
			fmt.Printf(`Your application is built!

To test it, run

cd ` + project + `
flo run

To set up Supabase authentication, first set up a new Supabase project.
Then, go to the Auth/Url Configuration tab in the Supabase console and 
copy your project's url and anon key into the .env file. While you're there,
you should add a custom scheme to your allowed redirect urls so that authentication
redirects can refer back to desktop and mobile front end applications.`)

		}
	case "deploy":
		{
			project_dir, err := os.Getwd()
			utils.PanicIf(err)
			project := filepath.Base(project_dir)
			cmd.BuildAndPushToRegistry(project, registry)
			cmd.DeployToLocalCluster(project)
		}
	case "run":
		{
			server := "front"
			if len(flag.Args()) > 1 {
				if flag.Arg(1) != "front" && flag.Arg(1) != "back" {
					panic("Run 'flo run front' to run frontend or 'flo run back' to run backend.")
				}
				server = flag.Arg(1)
			}
			project_dir, err := os.Getwd()
			utils.PanicIf(err)

			project := filepath.Base(project_dir)

			if server == "back" {
				cmd.RunBackend(project, kind)
			}
			if server == "front" {
				cmd.RunFrontend(project)
			}
		}
	case "sync":
		{
			project_dir, err := os.Getwd()
			utils.PanicIf(err)
			project := filepath.Base(project_dir)
			filename := "models.json"
			if !strings.HasSuffix(filename, ".json") {
				panic("invalid filename for model specification. Must be a .json file.")
			}
			models_file, err := os.ReadFile("./" + filename)
			utils.PanicIf(err)
			cmd.GenerateFloModelFiles(models_file, project, lang)
		}
	default:
		{
			panic("invalid command")
		}
	}
}
