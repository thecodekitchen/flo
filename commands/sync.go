package commands

import (
	sync "github.com/thecodekitchen/flo/sync"
)

func GenerateFloModelFiles(models_json []byte, project string, lang string) {
	sync.GenerateFlutterModels(models_json, project)
	switch lang {
	case "go":
		{
			sync.GenerateGoModels(models_json)
		}
	default:
		{
			panic("invalid backend language")
		}
	}
}
