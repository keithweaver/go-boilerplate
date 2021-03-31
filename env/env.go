package env

import (
	"log"
	"os"
)

// Environment variables:
// - SENDGRID_API_KEY
// - SENDER_EMAIL
// - SENDER_NAME
// - DB_NAME
// - FRONTEND_DOMAIN


// VerifyRequiredEnvVarsSet checks that the minimum set of environment variables
// are set. If you remove this method from the main.go, there may be unexpected
// issues.
func VerifyRequiredEnvVarsSet() {
	if os.Getenv("DB_NAME") == "" {
		// TODO - Add note to documentation
		log.Fatal("Error: DB_NAME is required in environment variables")
	}
}
