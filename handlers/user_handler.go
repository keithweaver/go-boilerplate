package handlers

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go-boilerplate/logging"
	"go-boilerplate/models"
	"go-boilerplate/services"
	"strings"
)

type UserHandler struct {
	logger logging.Logger
	userService services.UserService
}

func NewInstanceOfUserHandler(logger logging.Logger, userService services.UserService) *UserHandler {
	return &UserHandler{logger, userService}
}

func (u *UserHandler) SignIn(c *gin.Context) {
	var body models.SignInBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.userService.SignIn(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed in", "token": token})
	return
}

func (u *UserHandler) SignUp(c *gin.Context) {
	var body models.SignUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.userService.SignUp(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed up", "token": token})
	return
}

func (u *UserHandler) LogOut(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.ReplaceAll(authToken, "Bearer ", "")

	err := u.userService.LogOut(authToken)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logged out"})
	return
}
