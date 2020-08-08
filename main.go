package main

import (
	"context"
	"fmt"
	"strings"
	// "time"

	"github.com/gin-gonic/gin"
	"go-boilerplate/db"
	"go-boilerplate/handlers"
	"go-boilerplate/repositories"
	"go-boilerplate/services"
	// "github.com/gin-contrib/cors"
	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func ValidateAuth(userRepository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" {
			c.AbortWithStatus(403)
			return
		}

		authToken = strings.ReplaceAll(authToken, "Bearer ", "")

		found, _, err := userRepository.GetSessionById(authToken)
		if err != nil {
			fmt.Printf("err :: %+v\n", err)
			c.AbortWithStatus(403)
			return
		}
		if !found {
			fmt.Println("Not found")
			c.AbortWithStatus(403)
			return
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("Starting...")

	dbName := "demoDB" // move to env
	client, err := db.CreateDatabaseConnection(dbName)
	if err != nil {
		fmt.Println("Failed to connect to DB")
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(dbName)

	// Repositories
	userRepository := repositories.NewInstanceOfUserRepository(db)
	carsRepository := repositories.NewInstanceOfCarsRepository(db)

	// Services
	userService := services.NewInstanceOfUserService(userRepository)
	carsService := services.NewInstanceOfCarsService(userRepository, carsRepository)

	// Handlers
	userHandler := handlers.NewInstanceOfUserHandler(userService)
	carsHandler := handlers.NewInstanceOfCarsHandler(carsService)

	router := gin.Default()
	router.Use(CORSMiddleware())

	healthAPI := router.Group("/")
	{
		healthAPI.GET("", handlers.HealthCheck)
		healthAPI.GET("health", handlers.HealthCheck)
	}

	userAPI := router.Group("/user")
	{
		userAPI.POST("/signin", userHandler.SignIn)
		userAPI.POST("/signup", userHandler.SignUp)
		userAPI.POST("/logout", ValidateAuth(userRepository), userHandler.LogOut)
	}

	carsAPI := router.Group("/cars")
	{
		carsAPI.GET("/", ValidateAuth(userRepository), carsHandler.GetAll)
		carsAPI.GET("/:id", ValidateAuth(userRepository), carsHandler.GetByID)
		carsAPI.POST("/", ValidateAuth(userRepository), carsHandler.Create)
		carsAPI.PUT("/:id", ValidateAuth(userRepository), carsHandler.Update)
		carsAPI.DELETE("/:id", ValidateAuth(userRepository), carsHandler.Delete)
	}

	router.Run(":8080")
}
