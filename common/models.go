package common

import "github.com/gin-gonic/gin"

type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func ReturnErrorResponse(c *gin.Context, err *Error) {
	if err == nil {
		c.JSON(500, gin.H{
			"message": "Error: Internal server error",
		})
		return
	}
	message := err.Message
	if err.StatusCode == 403 && message == ""{
		message = "Error: Unauthorized"
	} else if err.StatusCode == 400 && message == ""{
		message = "Error: Invalid payload"
	} else if err.StatusCode == 500 && message == "" {
		message = "Error: Internal server error"
	}
	c.JSON(err.StatusCode, gin.H{
		"message": message,
	})
}