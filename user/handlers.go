package user

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go-boilerplate/logging"
	"strings"
)

type Handlers struct {
	logger      logging.Logger
	userServices Services
}

func NewInstanceOfUserHandlers(logger logging.Logger, userServices Services) *Handlers {
	return &Handlers{logger, userServices}
}

func (u *Handlers) SignIn(c *gin.Context) {
	var body SignInBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.userServices.SignIn(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed in", "token": token})
	return
}

func (u *Handlers) SignUp(c *gin.Context) {
	var body SignUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	token, err := u.userServices.SignUp(body)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Signed up", "token": token})
	return
}

func (u *Handlers) LogOut(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.ReplaceAll(authToken, "Bearer ", "")

	err := u.userServices.LogOut(authToken)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logged out"})
	return
}
