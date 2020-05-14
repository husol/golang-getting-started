package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang-getting-started/handlers"
	"golang-getting-started/models"
	"golang-getting-started/repositories"
	"log"
	"os"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	config := &aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("DYNAMO_DB_REGION")),
		Endpoint:    aws.String(os.Getenv("DYNAMO_DB_ENDPOINT")),
	}
	repositories.InitDb(config)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Signup route
	e.POST("/users/sign_up", handlers.SignUp())
	// Login route
	e.POST("/users/sign_in", handlers.SignIn())

	// Configure middleware with the custom claims type
	configJWT := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SIGNING_KEY")),
	}

	// Restricted group
	rPost := e.Group("/posts")
	rPost.Use(middleware.JWTWithConfig(configJWT))

	//Post route
	rPost.GET("", handlers.GetPosts())
	rPost.GET("/:id", handlers.GetPostById())
	rPost.POST("", handlers.CreatePost())
	rPost.PUT("/:id", handlers.UpdatePost())
	rPost.DELETE("/:id", handlers.RemovePost())
	rPost.GET("/me", handlers.GetMyPosts())

	e.Logger.Fatal(e.Start(":9090"))
}
