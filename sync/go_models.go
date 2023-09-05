package sync

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	utils "github.com/thecodekitchen/flo/utils"
)

func GenerateGoModels(models_json []byte) {
	os.Chdir("./app")
	// Consult  comment on this step in flutter_models.go for
	// a detailed explanation of this models map format.
	var models_map map[string]map[string]string
	err := json.Unmarshal(models_json, &models_map)
	utils.PanicIf(err)

	// start building the models file string
	models_file_string := `package models

type BaseModel struct {
	ID  string ` + "`" + `json:"id,omitempty"` + "`" + `
}`
	for model_sig, model_map := range models_map {
		extension := ""
		model_name := model_sig
		if strings.Contains(model_sig, "(") && strings.Contains(model_sig, ")") {
			ext_regex := regexp.MustCompile(`\(([^)]+)\)`)
			matches := ext_regex.FindStringSubmatch(model_sig)
			extension = matches[1]
			base_regex := regexp.MustCompile(`([^(]+)`)
			base_matches := base_regex.FindStringSubmatch(model_sig)
			model_name = base_matches[1]
		}
		models_file_string += `
type ` + model_name + ` struct {`
		// If there is an "extension", Go models will simply embed the parent struct
		if extension != "" {
			models_file_string += `
	` + extension
		} else {
			// If there is no extension, embed the BaseModel
			models_file_string += `
	BaseModel`
		}
		for field_sig, type_sig := range model_map {
			field_name, optional, list := analyze_field_sig(field_sig)
			// make necessary conversions for Go types
			if type_sig == "float" {
				type_sig = "float32"
			}
			go_field_name := strings.ToUpper(field_name[:1]) + field_name[1:]
			if list {
				type_sig = "[]" + type_sig
			}
			if optional {
				models_file_string += `
	` + go_field_name + ` ` + type_sig + ` ` + "`" + `json:"` + field_name + `"` + "`"
			} else {
				models_file_string += `
	` + go_field_name + ` ` + type_sig + ` ` + "`" + `json:"` + field_name + `", binding: "required"` + "`"
			}
		}
		models_file_string += `
}`
	}
	os.WriteFile("models/models.go", []byte(models_file_string), 0666)
}
