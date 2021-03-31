package sendgrid

import "os"


func GetSendgridAPIKey() string {
	return os.Getenv("SENDGRID_API_KEY")
}

func GetSenderEmail() string {
	return os.Getenv("SENDER_EMAIL")
}

func GetSenderName() string {
	return os.Getenv("SENDER_NAME")
}