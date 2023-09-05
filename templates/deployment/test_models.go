package deployment

func TestModelsBytes() []byte {
	return []byte(
		`{
	"User": {
		"name": "string",
		"age": "int"
	},
	"Book": {
		"isbn": "string",
		"seller": "string"
	}
}`)
}
