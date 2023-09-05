package sync

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	utils "github.com/thecodekitchen/flo/utils"
)

func GenerateFlutterModels(models_json []byte, project string) {
	os.Chdir(project + "_fe")
	// Unmarshal JSON as a map of maps.
	// NOTE: This supports only two nested layers.
	// The first layer is the names of the models.
	// The second layer is the fields in each model.
	// Nest by reference to other models for further layers.
	var models_map map[string]map[string]string
	err := json.Unmarshal(models_json, &models_map)
	utils.PanicIf(err)
	model_file_string := `class BaseModel {
	String? type = 'BaseModel';
	BaseModel({this.type = 'BaseModel'});
	Map toJson() => {
		'type': 'BaseModel'
	};
}`
	for model_sig, model_map := range models_map {
		extension := ""
		model_name := model_sig
		super_attributes := map[string]string{}
		if strings.Contains(model_sig, "(") && strings.Contains(model_sig, ")") {
			ext_regex := regexp.MustCompile(`\(([^)]+)\)`)
			ext_matches := ext_regex.FindStringSubmatch(model_sig)
			extension = ext_matches[1]
			base_regex := regexp.MustCompile(`([^(]+)`)
			base_matches := base_regex.FindStringSubmatch(model_sig)
			model_name = base_matches[1]
		}
		if extension == "" {
			model_file_string += `
class ` + model_name + ` extends BaseModel {`
		} else {
			super_attributes = find_super_attributes(models_map, extension)
			model_file_string += `
class ` + model_name + ` extends ` + extension + ` {`
		}
		// declare model attributes
		for field_sig, type_sig := range model_map {
			field_type, _ := convert_type_signature(type_sig)
			field_name, optional, list := analyze_field_sig(field_sig)

			if optional {
				if list {
					model_file_string += `
	List<` + field_type + `>? ` + field_name + `;`
				} else {
					model_file_string += `
	` + field_type + `? ` + field_name + `;`
				}
			} else if list {
				model_file_string += `
	List<` + field_type + `> ` + field_type + `;`
			} else {
				model_file_string += `
	` + field_type + ` ` + field_name + `;`
			}
		}

		// declare main constructor
		model_file_string += `
	` + model_name + ` ({
		super.type = '` + model_name + `',`
		for field_sig := range super_attributes {
			field_name, optional, _ := analyze_field_sig(field_sig)
			if optional {
				model_file_string += `
		super.` + field_name + `,`
			} else {
				model_file_string += `
		required super.` + field_name + `,`
			}
		}
		for field_sig := range model_map {
			field_name, optional, _ := analyze_field_sig(field_sig)
			if optional {
				model_file_string += `
		this.` + field_name + `,`
			} else {
				model_file_string += `
		required this.` + field_name + `,`
			}
		}
		model_file_string += `
	)};`
		// declare toJson method
		model_file_string += `
	Map toJson() {`
		// unpack nested model lists and convert items to JSON
		for field_sig, type_sig := range super_attributes {
			field_name, optional, list := analyze_field_sig(field_sig)
			field_type, is_model := convert_type_signature(type_sig)
			if is_model {
				if optional {
					if list {
						model_file_string += `
		List<` + field_type + `> ` + field_name + `List = [];
		for (` + field_type + ` item in ` + field_name + `??[]) {`
						model_file_string += `
			` + field_name + `List.add(item.toJson());
		}`
					}
				} else if list {
					model_file_string += `
		List<` + field_type + `> ` + field_name + `List = [];
		for (` + field_type + ` item in ` + field_name + `) {
			` + field_name + `List.add(item.toJson());
		}`
				}
			}
		}
		for field_sig, type_sig := range model_map {
			field_name, optional, list := analyze_field_sig(field_sig)
			field_type, is_model := convert_type_signature(type_sig)
			if is_model {
				if optional {
					if list {
						model_file_string += `
		List<` + field_type + `> ` + field_name + `List = [];
		for (` + field_type + ` item in ` + field_name + `??[]) {`
						model_file_string += `
			` + field_name + `List.add(item.toJson());
		}`
					}
				} else if list {
					model_file_string += `
		List<` + field_type + `> ` + field_name + `List = [];
		for (` + field_type + ` item in ` + field_name + `) {
			` + field_name + `List.add(item.toJson());
		}`
				}
			}
		}
		// return Map
		model_file_string += `
		return {
			'type': '` + model_name + `',`
		for field_sig := range super_attributes {
			field_name, _, _ := analyze_field_sig(field_sig)
			model_file_string += `
			'` + field_name + `': ` + field_name + `,`
		}
		for field_sig := range model_map {
			field_name, _, _ := analyze_field_sig(field_sig)
			model_file_string += `
			'` + field_name + `': ` + field_name + `,`
		}
		model_file_string += `
		};
	}`

		// Declare static fromJson method
		model_file_string += `
	static ` + model_name + ` fromJson(Map json) {`
		// unpack Lists of nested models and construct
		// Dart typed objects for them
		for field_sig, type_sig := range super_attributes {
			field_name, _, list := analyze_field_sig(field_sig)
			field_type, is_model := convert_type_signature(type_sig)

			if list && is_model {
				model_file_string += `
		List<` + field_type + `> ` + field_name + `List = []
		for (Map<dynamic, dynamic> item in json['` + field_name + `']) {
			importsList.add(` + field_type + `.fromJson(item));
		}`
			}
		}
		for field_sig, type_sig := range model_map {
			field_name, _, list := analyze_field_sig(field_sig)
			field_type, is_model := convert_type_signature(type_sig)

			if list && is_model {
				model_file_string += `
		List<` + field_type + `> ` + field_name + `List = []
		for (Map<dynamic, dynamic> item in json['` + field_name + `']) {
			importsList.add(` + field_type + `.fromJson(item));
		}`
			}
		}

		// Construct return statement for fromJson method
		model_file_string += `
		
		return ` + model_name + `(
			type: '` + model_name + `,`
		for field_sig := range super_attributes {
			field_name, _, list := analyze_field_sig(field_sig)
			if list {
				model_file_string += `
			` + field_name + `: ` + field_name + `List,`
			} else {
				model_file_string += `
			` + field_name + `: json['` + field_name + `'],`
			}
		}
		for field_sig := range model_map {
			field_name, _, list := analyze_field_sig(field_sig)
			if list {
				model_file_string += `
			` + field_name + `: ` + field_name + `List,`
			} else {
				model_file_string += `
			` + field_name + `: json['` + field_name + `'],`
			}
		}
		model_file_string += `
		);
	}
}`
	}
	os.WriteFile("./lib/models.dart", []byte(model_file_string), 0666)
	os.Chdir("..")
}

func convert_type_signature(type_sig string) (string, bool) {
	// return value of true indicates an extension
	switch type_sig {
	case "string":
		{
			return "String", false
		}
	case "float", "float32", "float64":
		{
			return "double", false
		}
	case "int":
		{
			return type_sig, false
		}
	case "bool":
		{
			return type_sig, false
		}
	default:
		{
			return type_sig, true
		}
	}
}

func analyze_field_sig(field_sig string) (string, bool, bool) {
	field_name := field_sig
	optional := false
	list := false
	if strings.HasPrefix(field_sig, "?") {
		optional = true
		field_name = strings.TrimPrefix(field_name, "?")
	}
	if strings.HasSuffix(field_sig, "[]") {
		list = true
		field_name = strings.TrimSuffix(field_name, "[]")
	}
	return field_name, optional, list
}

func find_super_attributes(models_map map[string]map[string]string, extension string) map[string]string {
	if models_map[extension] == nil {
		// parent might itself be an extension. Check for this.
		next_extension := ""
		model := map[string]string{}
		for k, v := range models_map {
			if strings.Contains(k, extension) {
				model = v
				ext_regex := regexp.MustCompile(`\(([^)]+)\)`)
				matches := ext_regex.FindStringSubmatch(k)
				if len(matches) >= 2 {
					next_extension = matches[1]
				}
				parent_model := find_super_attributes(models_map, next_extension)
				for super_k, super_v := range parent_model {
					model[super_k] = super_v
				}
			}
		}
		if next_extension != "" && len(model) > 0 {
			return model
		} else {
			panic("invalid extension: " + extension)
		}
	} else {
		return models_map[extension]
	}
}
