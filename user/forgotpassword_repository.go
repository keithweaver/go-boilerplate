package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ForgotPasswordRepository struct {
	db                 *mongo.Database
	forgotPasswordCollection    string
}

func NewInstanceOfForgotPasswordRepository(db *mongo.Database) ForgotPasswordRepository {
	return ForgotPasswordRepository{db: db, forgotPasswordCollection: "forgotPasswordCodes"}
}

func (r *ForgotPasswordRepository) Save(forgotPassword ForgotPasswordCode) error {
	_, err := r.db.Collection(r.forgotPasswordCollection).InsertOne(context.TODO(), forgotPassword)
	if err != nil {
		return err
	}
	return nil
}

func (r *ForgotPasswordRepository) Exists(email string, code string) (bool, error) {
	filter := bson.M{"email": email, "code": code}
	count, err := r.db.Collection(r.forgotPasswordCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ForgotPasswordRepository) MarkCodeAsComplete(email string, code string) error {
	filter := bson.M{"email": email, "code": code}
	update := bson.M{"$set": bson.M{"expiry": time.Now()}}
	_, err := r.db.Collection(r.forgotPasswordCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}