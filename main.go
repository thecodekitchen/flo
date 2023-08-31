package main

import (
	"flag"

	gen "github.com/thecodekitchen/flo/generators"
)

func main() {
	flag.Parse()

	switch command := flag.Arg(0); command {
	case "create":
		{

			if len(flag.Args()) < 2 {
				panic("Need to specify a project name")
			}
			project := flag.Arg(1)
			// Extra argument after project name will be
			// interpreted as a backend language choice
			// default is Go.
			if len(flag.Args()) > 2 {
				lang := flag.Arg(2)
				gen.GenerateBackendApi(project, lang)
			} else {
				gen.GenerateBackendApi(project, "go")
			}
			gen.GenerateFlutterApp(project)
		}
	default:
		{
			panic("invalid command")
		}
	}
}
