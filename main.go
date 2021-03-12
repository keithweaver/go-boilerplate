package main

import (
	"context"
	"fmt"
	"go-boilerplate/auth"
	"go-boilerplate/cars"
	"go-boilerplate/health"
	"go-boilerplate/logging"
	"go-boilerplate/middleware"
	"go-boilerplate/user"

	//"strings"
	// "time"

	"github.com/gin-gonic/gin"
	"go-boilerplate/db"
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
	userRepository := user.NewInstanceOfUserRepository(db)
	carsRepository := cars.NewInstanceOfCarsRepository(db)

	// Services
	userServices := user.NewInstanceOfUserServices(logger, userRepository)
	carsServices := cars.NewInstanceOfCarsServices(logger, userRepository, carsRepository)

	// Handlers
	userHandlers := user.NewInstanceOfUserHandlers(logger, userServices)
	carsHandlers := cars.NewInstanceOfCarsHandlers(logger, carsServices)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	healthAPI := router.Group("/")
	{
		healthAPI.GET("", health.Check)
		healthAPI.GET("health", health.Check)
	}

	userAPI := router.Group("/user")
	{
		userAPI.POST("/signin", userHandlers.SignIn)
		userAPI.POST("/signup", userHandlers.SignUp)
		userAPI.POST("/logout", auth.ValidateAuth(userRepository), userHandlers.LogOut)
	}

	carsAPI := router.Group("/cars")
	{
		carsAPI.GET("/", auth.ValidateAuth(userRepository), carsHandlers.GetAll)
		carsAPI.GET("/:id", auth.ValidateAuth(userRepository), carsHandlers.GetByID)
		carsAPI.POST("/", auth.ValidateAuth(userRepository), carsHandlers.Create)
		carsAPI.PUT("/:id", auth.ValidateAuth(userRepository), carsHandlers.Update)
		carsAPI.DELETE("/:id", auth.ValidateAuth(userRepository), carsHandlers.Delete)
	}

	router.Run(":8080")
}
