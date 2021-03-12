package services

import (
	"errors"
	"fmt"
	"go-boilerplate/logging"
	"go-boilerplate/models"
	"go-boilerplate/repositories"
	"strings"
	"time"
	// "github.com/google/uuid"
)

type UserService struct {
	logger logging.Logger
	userRepository repositories.UserRepository
}

type UserServiceContract interface {
	SignUp(body models.SignUpBody) (string, error)
	IsValidPassword(password string) bool
	GetEncryptedPassword(password string) string
	SignIn(body models.SignInBody) (string, error)
	LogOut(authToken string) error
}

func NewInstanceOfUserService(logger logging.Logger, userRepository repositories.UserRepository) UserService {
	return UserService{logger, userRepository}
}

func (u *UserService) SignUp(body models.SignUpBody) (string, error) {
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
	newUser := models.User{
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

func (u *UserService) IsValidPassword(password string) bool {
	// TODO - Insert your password rules. Ex. Must have a digit, special
	// characters, is super long. I choose not to include password and security
	// since it's a challenge and is up to the developer to understand how it
	// works. You need to add to getEncryptedPassword.
	fmt.Println("---- IMPORTANT ----")
	fmt.Println("Your sign up validate is failing on purpose. Please open isValidPassword in services/user_service.go and read the message.")
	// After reading this message above, remove the "return false" and comment out
	// the next line at your own risk:
	// return true
	return false
}

func (u *UserService) GetEncryptedPassword(password string) string {
	// TODO - See comment in the isValidPassword
	return password
}

func (u *UserService) SignIn(body models.SignInBody) (string, error) {
	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")
	return u.signIn(emailTrimmed, body.Password)
}

func (u *UserService) signIn(email string, password string) (string, error) {
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
	newSession := models.Session{
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

func (u *UserService) LogOut(authToken string) error {
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
