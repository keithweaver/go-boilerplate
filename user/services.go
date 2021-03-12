package user

import (
	"errors"
	"fmt"
	"go-boilerplate/logging"
	"strings"
	"time"
	// "github.com/google/uuid"
)

type Services struct {
	logger         logging.Logger
	userRepository Repository
}

type ServiceContract interface {
	SignUp(body SignUpBody) (string, error)
	IsValidPassword(password string) bool
	GetEncryptedPassword(password string) string
	SignIn(body SignInBody) (string, error)
	LogOut(authToken string) error
}

func NewInstanceOfUserServices(logger logging.Logger, userRepository Repository) Services {
	return Services{logger, userRepository}
}

func (u *Services) SignUp(body SignUpBody) (string, error) {
	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")

	// Verify password meets sign up requirements
	if !u.IsValidPassword(body.Password) {
		return "", errors.New("error: Your password does not meet requirements.")
	}

	// Check for user
	userExists, err := u.userRepository.DoesUserExist(emailTrimmed)
	if err != nil {
		return "", err
	}
	if userExists {
		// Just try signing them in
		return u.signIn(emailTrimmed, body.Password)
	}

	// Sign up user
	newUser := User{
		Email:    emailTrimmed,
		Password: u.GetEncryptedPassword(body.Password),
		Name:     body.Name,
		Created:  time.Now(),
	}
	err = u.userRepository.SaveUser(newUser)
	if err != nil {
		return "", err
	}

	// Sign in user
	return u.signIn(emailTrimmed, body.Password)
}

func (u *Services) IsValidPassword(password string) bool {
	// TODO - Insert your password rules. Ex. Must have a digit, special
	// characters, is super long. I choose not to include password and security
	// since it's a challenge and is up to the developer to understand how it
	// works. You need to add to getEncryptedPassword.
	fmt.Println("---- IMPORTANT ----")
	fmt.Println("Your sign up validate is failing on purpose. Please open isValidPassword in services/services.go and read the message.")
	// After reading this message above, remove the "return false" and comment out
	// the next line at your own risk:
	// return true
	return false
}

func (u *Services) GetEncryptedPassword(password string) string {
	// TODO - See comment in the isValidPassword
	return password
}

func (u *Services) SignIn(body SignInBody) (string, error) {
	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")
	return u.signIn(emailTrimmed, body.Password)
}

func (u *Services) signIn(email string, password string) (string, error) {
	// Encrypt password
	encryptedPassword := u.GetEncryptedPassword(password)

	// Grab user
	found, user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New("error: Unauthorized")
	}

	if user.Password != encryptedPassword {
		return "", errors.New("error: Unauthorized")
	}

	// Create session
	now := time.Now()
	expiryDate := now.AddDate(0, 0, 1)
	newSession := Session{
		Email:   email,
		Created: now,
		Expiry:  expiryDate,
	}

	token, err := u.userRepository.SaveSession(newSession)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *Services) LogOut(authToken string) error {
	// Grab session
	found, _, err := u.userRepository.GetSessionById(authToken)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("error: Unauthorized")
	}

	// Mark as expired
	err = u.userRepository.MarkSessionAsExpired(authToken)
	if err != nil {
		return err
	}

	return nil
}
