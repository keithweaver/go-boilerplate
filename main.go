package main

import (
	"context"
	"fmt"
	"go-boilerplate/auth"
	"go-boilerplate/logging"
	"go-boilerplate/middleware"

	//"strings"
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




func main() {
	fmt.Println("Starting...")
	logger := logging.NewLogger()

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
	userService := services.NewInstanceOfUserService(logger, userRepository)
	carsService := services.NewInstanceOfCarsService(logger, userRepository, carsRepository)

	// Handlers
	userHandler := handlers.NewInstanceOfUserHandler(logger, userService)
	carsHandler := handlers.NewInstanceOfCarsHandler(logger, carsService)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	healthAPI := router.Group("/")
	{
		healthAPI.GET("", handlers.HealthCheck)
		healthAPI.GET("health", handlers.HealthCheck)
	}

	userAPI := router.Group("/user")
	{
		userAPI.POST("/signin", userHandler.SignIn)
		userAPI.POST("/signup", userHandler.SignUp)
		userAPI.POST("/logout", auth.ValidateAuth(userRepository), userHandler.LogOut)
	}

	carsAPI := router.Group("/cars")
	{
		carsAPI.GET("/", auth.ValidateAuth(userRepository), carsHandler.GetAll)
		carsAPI.GET("/:id", auth.ValidateAuth(userRepository), carsHandler.GetByID)
		carsAPI.POST("/", auth.ValidateAuth(userRepository), carsHandler.Create)
		carsAPI.PUT("/:id", auth.ValidateAuth(userRepository), carsHandler.Update)
		carsAPI.DELETE("/:id", auth.ValidateAuth(userRepository), carsHandler.Delete)
	}

	router.Run(":8080")
}
