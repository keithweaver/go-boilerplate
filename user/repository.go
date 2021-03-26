package user

import (
	"context"
	// "common"
	// "fmt"
	"time"

	// "database/sql"
	// "github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db                 *mongo.Database
	usersCollection    string
	sessionsCollection string
}

func NewInstanceOfUserRepository(db *mongo.Database) Repository {
	return Repository{db: db, usersCollection: "users", sessionsCollection: "sessions"}
}

func (u *Repository) GetUserByEmail(email string) (bool, User, error) {
	var user User
	filter := bson.M{"email": email}
	count, err := u.db.Collection(u.usersCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, User{}, err
	}
	if count != 1 {
		return false, User{}, nil
	}
	err = u.db.Collection(u.usersCollection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false, User{}, err
	}

	return true, user, nil
}

func (u *Repository) DoesUserExist(email string) (bool, error) {
	exists, _, err := u.GetUserByEmail(email)
	return exists, err
}

func (u *Repository) SaveUser(user User) error {
	_, err := u.db.Collection(u.usersCollection).InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) SaveSession(session Session) (string, error) {
	insertResult, err := u.db.Collection(u.sessionsCollection).InsertOne(context.TODO(), session)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (u *Repository) GetSessionById(token string) (bool, Session, error) {
	docID, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		return false, Session{}, err
	}

	var session Session
	filter := bson.M{
		"_id": docID,
		"expiry": bson.M{
			"$gte": time.Now(),
		},
	}

	count, err := u.db.Collection(u.sessionsCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, Session{}, err
	}
	if count != 1 {
		return false, Session{}, nil
	}
	err = u.db.Collection(u.sessionsCollection).FindOne(context.TODO(), filter).Decode(&session)
	if err != nil {
		return false, Session{}, err
	}
	return true, session, nil
}

func (u *Repository) MarkSessionAsExpired(authToken string) error {
	filter := bson.M{"_id": authToken}
	update := bson.M{"$set": bson.M{"expiry": time.Now()}}
	_, err := u.db.Collection(u.sessionsCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) UpdateAccountLocked(email string, accountLock bool) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"accountLocked": accountLock}}
	_, err := u.db.Collection(u.usersCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) UpdateOrAddTrustedIPToUser(email string, newIP IP) error {
	// Check if the IP address already exists
	filter := bson.M{"email": email, "trustedIPs.address": newIP.Address}
	count, err := u.db.Collection(u.usersCollection).CountDocuments(context.TODO(), filter)
	if err != nil {
		return err
	}

	if count > 0 {
		// Update the current IP information (Docs: https://docs.mongodb.com/manual/reference/operator/update/positional/#update-values-in-an-array)
		update := bson.M{"$set": bson.M{ "trustedIPs.$": newIP}}
		_, err = u.db.Collection(u.usersCollection).UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	}

	// Push new one (Docs: https://docs.mongodb.com/manual/reference/operator/update/push/#examples)
	update := bson.M{"$push": bson.M{"trustedIPs": newIP}}
	_, err = u.db.Collection(u.usersCollection).UpdateOne(context.Background(), bson.M{"email": email}, update)
	if err != nil {
		return err
	}

	return nil
}

func (u *Repository) UnlockSession(authToken string) error {
	filter := bson.M{"_id": authToken}
	update := bson.M{"$set": bson.M{"locked": false}}
	_, err := u.db.Collection(u.sessionsCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) UpdatePassword(email string, newPassword string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": newPassword}}
	_, err := u.db.Collection(u.usersCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}