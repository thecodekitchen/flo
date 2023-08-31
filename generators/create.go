package generators

import (
	"os"
	"os/exec"

	utils "github.com/thecodekitchen/flo/utils"
)

func GenerateBackendApi(project string, lang string) {
	if lang == "go" {
		generate_go_api(project)
	}
}

func GenerateFlutterApp(project string) {
	err := exec.Command("flutter", "create", project+"_fe").Run()
	utils.PanicIf(err)
	os.Chdir(project + "_fe")

}

func generate_go_api(project string) {
	err := os.Mkdir("../app", 0750)
	utils.ExistsError(err)
	os.Chdir("../app")
	err = os.Mkdir("./models", 0750)
	utils.ExistsError(err)
	os.WriteFile("./models/models.go", []byte(
		`package models

type User struct {
	ID       string `+"`"+`json:"id,omitempty"`+"`"+`
	User     string `+"`"+`json:"user" binding:"required"`+"`"+`
	Password string `+"`"+`json:"password" binding:"required"`+"`"+`
	Books    []Book `+"`"+`json:"projects" binding:"required"`+"`"+`
}
type Book struct {
	ID     string `+"`"+`json:"id,omitempty"`+"`"+`
	Type   string `+"`"+`json:"type"`+"`"+`
	Seller string `+"`"+`json:"seller" binding:"required"`+"`"+`
	ISBN   string `+"`"+`json:"isbn" binding:"required"`+"`"+`
}`), 0666)
	os.WriteFile("../app/main.go",
		[]byte(`package main

import (
	"net/http"
	"os"
	"strings"
	"`+project+`/models"

	"github.com/gin-gonic/gin"
	supa "github.com/nedpals/supabase-go"
	"github.com/surrealdb/surrealdb.go"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		supabaseUrl := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_KEY")

		supabase := supa.CreateClient(supabaseUrl, supabaseKey)
		user, err := supabase.Auth.User(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithError(1, err)
		}
		if user != nil {
			// check other headers, path parameters, or query parameters
			// to confirm scope access for whichever route group you're
			// attaching the middleware to.
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	// r.Use(AuthMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/save_book", func(c *gin.Context) {
		var json Book
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		surreal_url := os.Getenv("SURREALDB_URL")
		db, err := surrealdb.New(surreal_url)
		if err != nil {
			panic(err)
		}

		if _, err = db.Signin(gin.H{
			"user": "root",
			"pass": "surrealdb",
		}); err != nil {
			panic(err)
		}
		if _, err = db.Use("test", "test"); err != nil {
			panic(err)
		}
		data, err := db.Create("book", json)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"surreal_response": data})
	})
	r.POST("/create_user", func(c *gin.Context) {
		// parse JSON body as User struct or return BadRequest error
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		surreal_url := os.Getenv("SURREALDB_URL")
		db, err := surrealdb.New(surreal_url)
		if err != nil {
			panic(err)
		}

		if _, err = db.Signin(gin.H{
			"user": "root",
			"pass": "surrealdb",
		}); err != nil {
			panic(err)
		}
		if _, err = db.Use("test", "test"); err != nil {
			panic(err)
		}

		data, err := db.Create("user", json)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"surreal_response": data})
	})
	r.Run(":8080")
}`), 0666)
	os.WriteFile("../app/go.mod", []byte(
		`module `+project+`

go 1.18

require github.com/surrealdb/surrealdb.go v0.2.1

require (
	github.com/bytedance/sonic v1.10.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.9.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nedpals/postgrest-go v0.1.3 // indirect
	github.com/nedpals/supabase-go v0.3.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.4.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)`), 0666)

	cmd := exec.Command("go", "get")
	err = cmd.Run()
	utils.PanicIf(err)
}
