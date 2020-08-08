package handlers

import (
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Healthy"})
}
