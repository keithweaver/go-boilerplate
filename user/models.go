package user

import (
	"errors"
	"strings"
	"time"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
	Name     string    `json:"name" bson:"name"`
	Created  time.Time `json:"created" bson:"created"`
	VerifiedEmail bool `json:"verified" bson:"verified"`
	VerificationCode string `json:"-" bson:"validationCode"` // On sign up or requested, user is sent a temporary verification code
	VerificationExpiryTime time.Time `json:"-" bson:"validationExpiryTime" ` // The timeframe when the verification code is sent
	TrustedIPs []IP `json:"trustedIPs" bson:"trustedIPs"`
	InvalidIPs []IP `json:"invalidIPs" bson:"invalidIPs"`
	AccountLocked     bool      `json:"accountLocked" bson:"accountLocked"` // Stop new sign ins from happening
	KnownDevices      []Device  `json:"knownDevices" bson:"knownDevices"`
}

func (u *User) Greeting() string {
	if u.Name != "" {
		return u.Name
	}
	return "there"
}

type IP struct {
	Address string `json:"address" bson:"address"`
	LocationFound bool `json:"locationFound" bson:"locationFound"` // Boolean flag that indicates other location based attributes are set
	Latitude float64 `json:"latitude" bson:"latitude"` // Not required
	Longitude float64 `json:"longitude" bson:"longitude"` // Not required
	Country string `json:"country" bson:"country"`
	Region string `json:"region" bson:"region"`
	City string `json:"city" bson:"city"`
}

func (u *User) HasTrustedIP(ipAddress string) bool {
	for _, ip := range u.TrustedIPs {
		if ip.Address == ipAddress {
			return true
		}
	}
	return false
}

func (u *User) HasInvalidIPs(ipAddress string) bool {
	for _, ip := range u.InvalidIPs {
		if ip.Address == ipAddress {
			return true
		}
	}
	return false
}

type Device struct {
	Name            string `json:"name" bson:"name"`
	Mobile          bool   `json:"mobile" bson:"mobile"`
	Bot             bool   `json:"bot" bson:"bot"`
	Mozilla         string `json:"mozilla" bson:"mozilla"`
	Platform        string `json:"platform" bson:"platform"`
	OperatingSystem string `json:"operatingSystem" bson:"operatingSystem"`
	Engine          string `json:"engine" bson:"engine"`
	EngineVersion   string `json:"engineVersion" bson:"engineVersion"`
	Browser         string `json:"browser" bson:"browser"`
	BrowserVersion  string `json:"browserVersion" bson:"browserVersion"`
	ValidDevice     bool   `json:"validDevice" bson:"validDevice"` // Starts as true and user can change to false.
}

type Session struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email   string             `json:"email" bson:"email"`
	Expiry  time.Time          `json:"expiry" bson:"expiry"`
	Created time.Time          `json:"created" bson:"created"`
	Locked     bool               `json:"locked" bson:"locked"`
	UnlockCode string             `json:"unlockCode" bson:"unlockCode"` // Second layer of security, on suspicious signs in, emails code to confirm
	Device Device `json:"device" bson:"device,omitempty"`
}

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UnlockSessionBody struct {
	Code string `json:"code"`
}

func (b *UnlockSessionBody) Validate() error {
	if b.Code == "" {
		return errors.New("code is required")
	}
	return nil
}

type SendForgotPasswordBody struct {
	Email string `json:"email"`
}

func (b *SendForgotPasswordBody) Validate() error {
	if b.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (b *SendForgotPasswordBody) GetFormattedEmail() string {
	email := strings.Trim(b.Email, " ")
	email = strings.ToLower(email)
	return email
}

type ResetForgotPasswordBody struct {
	Email string `json:"email"`
	Code string `json:"code"`
	NewPassword string `json:"newPassword"`
}

func (b *ResetForgotPasswordBody) Validate() error {
	if b.Email == "" {
		return errors.New("email is required")
	}
	if b.Code == "" {
		return errors.New("code is required")
	}
	if b.NewPassword == "" {
		return errors.New("password is required")
	}
	return nil
}

func (b *ResetForgotPasswordBody) GetFormattedEmail() string {
	email := strings.Trim(b.Email, " ")
	email = strings.ToLower(email)
	return email
}

type ForgotPasswordCode struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Code string `json:"-" bson:"code"`
	Created time.Time `json:"created" bson:"created"`
	Expiry time.Time `json:"expiry" bson:"expiry"`
}