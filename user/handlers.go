package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"go-boilerplate/common"

	// "fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"go-boilerplate/logging"
	"strings"
)

type Handlers struct {
	logger       logging.Logger
	userServices Services
}

func NewInstanceOfUserHandlers(logger logging.Logger, userServices Services) *Handlers {
	return &Handlers{logger, userServices}
}

func (u *Handlers) SignIn(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "SignIn")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

	// Capture IP
	clientIP := c.ClientIP()

	ctx = context.WithValue(ctx, logging.CtxClientIP, clientIP)

	// Capture User Agent header
	var userAgent *user_agent.UserAgent
	if c.Request.Header["User-Agent"] != nil && len(c.Request.Header["User-Agent"]) > 0 {
		userAgent = user_agent.New(c.Request.Header["User-Agent"][0])
	}


	var body SignInBody
	if err := c.ShouldBindJSON(&body); err != nil {
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400, Message: err.Error()})
		return
	}

	token, sessionLocked, err := u.userServices.SignIn(ctx, userAgent, clientIP, body)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Signed in", "token": token, "sessionLocked": sessionLocked})
	return
}

func (u *Handlers) SignUp(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "SignUp")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

	// Capture IP
	clientIP := c.ClientIP()

	ctx = context.WithValue(ctx, logging.CtxClientIP, clientIP)

	// Capture User Agent header
	var userAgent *user_agent.UserAgent
	if c.Request.Header["User-Agent"] != nil && len(c.Request.Header["User-Agent"]) > 0 {
		userAgent = user_agent.New(c.Request.Header["User-Agent"][0])
	}

	var body SignUpBody
	if err := c.ShouldBindJSON(&body); err != nil {
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400, Message: err.Error()})
		return
	}

	token, sessionLocked, err := u.userServices.SignUp(ctx, userAgent, clientIP, body)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Signed up", "token": token, "sessionLocked": sessionLocked})
	return
}

func (u *Handlers) LogOut(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "LogOut")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())
	ctx = context.WithValue(ctx, logging.CtxClientIP, c.ClientIP())

	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.ReplaceAll(authToken, "Bearer ", "")

	err := u.userServices.LogOut(ctx, authToken)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{"message": "Logged out"})
	return
}

func (u *Handlers) UnlockSession(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "UnlockSession")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

	// Capture IP
	clientIP := c.ClientIP()
	ctx = context.WithValue(ctx, logging.CtxClientIP, clientIP)

	// Get auth token
	if c.Request.Header["Authorization"] == nil || len(c.Request.Header["Authorization"]) == 0 {
		u.logger.Warning(ctx, "no auth token provided", errors.New("unauthorized"))
		common.ReturnErrorResponse(c, &common.Error{
			StatusCode: 403,
		})
		return
	}
	authToken := c.Request.Header["Authorization"][0]
	authToken = strings.ReplaceAll(authToken, "Bearer ", "")
	authToken = strings.Trim(authToken, " ")

	if authToken == "" {
		u.logger.Warning(ctx, "no auth token provided 2", errors.New("unauthorized"))
		common.ReturnErrorResponse(c, &common.Error{
			StatusCode: 403,
		})
		return
	}

	// Get code from request body
	var body UnlockSessionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Warning(ctx, "invalid request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400})
		return
	}

	if err := body.Validate(); err != nil {
		u.logger.Warning(ctx, "validation failed on request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400, Message: err.Error()})
		return
	}

	// Call the service
	err := u.userServices.UnlockSession(ctx, clientIP, authToken, body)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Session unlocked"})
	return
}

func (u *Handlers) SendForgotPassword(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "SendForgotPassword")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

	// Capture IP
	clientIP := c.ClientIP()
	ctx = context.WithValue(ctx, logging.CtxClientIP, clientIP)

	// Get code from request body
	var body SendForgotPasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Warning(ctx, "invalid request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400})
		return
	}

	if err := body.Validate(); err != nil {
		u.logger.Warning(ctx, "validation failed on request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400, Message: err.Error()})
		return
	}

	// Call the service
	err := u.userServices.SendForgotPassword(ctx, clientIP, body)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Email sent"})
	return
}

func (u *Handlers) ForgotPassword(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logging.CtxDomain, "user")
	ctx = context.WithValue(ctx, logging.CtxHandlerMethod, "ForgotPassword")
	ctx = context.WithValue(ctx, logging.CtxRequestID, uuid.New().String())

	// Capture IP
	clientIP := c.ClientIP()
	ctx = context.WithValue(ctx, logging.CtxClientIP, clientIP)

	// Get code from request body
	var body ResetForgotPasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		u.logger.Warning(ctx, "invalid request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400})
		return
	}

	if err := body.Validate(); err != nil {
		u.logger.Warning(ctx, "validation failed on request body", err)
		common.ReturnErrorResponse(c, &common.Error{StatusCode: 400, Message: err.Error()})
		return
	}

	// Call the service
	err := u.userServices.ForgotPassword(ctx, clientIP, body)
	if err != nil {
		common.ReturnErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "Password has been reset"})
	return

}