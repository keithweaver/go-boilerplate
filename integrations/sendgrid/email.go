package sendgrid

import (
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email using the Sendgrid SDK. If this returns an error, you should return a 500 on
// the service & handler level.
func SendEmail(fullName string, email string, subject string, plainTextContent string, htmlContent string) error {
	apiKey := GetSendgridAPIKey()
	if apiKey == "" {
		return errors.New("error: no Sendgrid API key found in environment variables")
	}
	senderEmail := GetSenderEmail()
	if senderEmail == "" {
		return errors.New("error: no sender email found in environment variables")
	}
	senderName := GetSenderName()
	if senderName == "" {
		return errors.New("error: no sender name found in environment variables")
	}

	from := mail.NewEmail(senderName, senderEmail)
	to := mail.NewEmail(fullName, email)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)
	return err
}
