package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string    `json:"email",bson:"email"`
	Password string    `json:"password",bson:"password"`
	Name     string    `json:"name",bson:"name"`
	Created  time.Time `json:"created",bson:"created"`
}

type Session struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email   string             `json:"email",bson:"email"`
	Expiry  time.Time          `json:"expiry",bson:"expiry"`
	Created time.Time          `json:"created",bson:"created"`
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
