package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-boilerplate/repositories"
	"strings"
)

func ValidateAuth(userRepository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" {
			c.AbortWithStatus(403)
			return
		}

		authToken = strings.ReplaceAll(authToken, "Bearer ", "")

		found, session, err := userRepository.GetSessionById(authToken)
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
		c.Set("session", session)

		c.Next()
	}
}
