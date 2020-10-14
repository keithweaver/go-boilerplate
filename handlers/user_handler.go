package handlers

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go-boilerplate/models"
	"go-boilerplate/services"
	"strings"
)

type UserHandler struct {
	UserService services.UserService
}

func NewInstanceOfUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
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

	token, err := u.UserService.SignIn(body)
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

	token, err := u.UserService.SignUp(body)
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

	err := u.UserService.LogOut(authToken)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logged out"})
	return
}
