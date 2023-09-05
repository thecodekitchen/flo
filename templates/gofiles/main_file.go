package gofiles

func MainFileBytes(project string) []byte {
	return []byte(`package main

	import (
		"net/http"
		"os"
		"strings"
		m "` + project + `/models"
	
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
			var json m.Book
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
			var json m.User
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
	}`)
}
