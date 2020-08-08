package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go-boilerplate/models"
	"go-boilerplate/services"
	"strconv"
	"strings"
)

type CarsHandler struct {
	CarsService services.CarsService
}

func NewInstanceOfCarsHandler(carsService services.CarsService) *CarsHandler {
	return &CarsHandler{CarsService: carsService}
}

func (u *CarsHandler) GetAll(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	if authToken != "" {
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "25")
	make := c.DefaultQuery("make", "")
	model := c.DefaultQuery("model", "")
	year := c.DefaultQuery("year", "0")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	query := models.ListCarQuery{
		Page:  pageInt,
		Limit: limitInt,
		Make:  make,
		Model: model,
		Year:  yearInt,
	}

	v := validator.New()
	if err := v.Struct(query); err != nil {
		fmt.Print("Validation failed.")
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	cars, err := u.CarsService.GetAll(authToken, query)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if cars == nil {
		cars = []models.Car{}
	}
	c.JSON(200, gin.H{"message": "Cars retrieved", "cars": cars})
	return
}

func (u *CarsHandler) GetByID(c *gin.Context) {
	carsID := c.Param("id")

	authToken := c.Request.Header.Get("Authorization")
	if authToken != "" {
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	}

	car, err := u.CarsService.GetByID(authToken, carsID)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Car retrieved", "car": car})
	return
}

func (u *CarsHandler) Create(c *gin.Context) {
	var body models.CreateCar
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := body.Valid(); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	authToken := c.Request.Header.Get("Authorization")
	if authToken != "" {
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err := u.CarsService.Create(authToken, body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Created car"})
	return
}

func (u *CarsHandler) Update(c *gin.Context) {
	carsID := c.Param("id")

	var body models.UpdateCar
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := body.Valid(); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	authToken := c.Request.Header.Get("Authorization")
	if authToken != "" {
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err := u.CarsService.Update(authToken, carsID, body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Updated car"})
	return
}

func (u *CarsHandler) Delete(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	if authToken != "" {
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	}

	carsID := c.Param("id")

	err := u.CarsService.Delete(authToken, carsID)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted car"})
	return
}
