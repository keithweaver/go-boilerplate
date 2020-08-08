package repositories

import (
	"context"
	// "errors"
	// "fmt"
	// "time"
	// "log"
	"go-boilerplate/models"
	// "database/sql"
	// "github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CarsRepository struct {
	db *mongo.Database
}

func NewInstanceOfCarsRepository(db *mongo.Database) CarsRepository {
	return CarsRepository{db: db}
}

func (c *CarsRepository) List(email string, query models.ListCarQuery) ([]models.Car, error) {
	filters := query.Filter(email)
	var cars []models.Car

	options := options.Find()

	// Add paging
	options.SetLimit(int64(query.Limit))
	options.SetSkip(int64((query.Page * query.Limit) - query.Limit))

	// Add timestamp
	options.SetSort(bson.M{"created": -1})

	cursor, err := c.db.Collection("cars").Find(context.Background(), filters, options)
	if err != nil {
		return []models.Car{}, err
	}

	for cursor.Next(context.Background()) {
		car := models.Car{}
		err := cursor.Decode(&car)
		if err != nil {
			//handle err
		} else {
			cars = append(cars, car)
		}
	}
	return cars, nil
}

func (c *CarsRepository) Get(email string, carID string) (models.Car, error) {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return models.Car{}, err
	}

	filter := bson.M{"_id": docID, "email": email}

	var result models.Car
	err = c.db.Collection("cars").FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.Car{}, err
	}
	return result, nil
}

func (c *CarsRepository) Save(car models.Car) error {
	_, err := c.db.Collection("cars").InsertOne(context.TODO(), car)
	if err != nil {
		return err
	}
	// insertResult.InsertedID.(string)
	return nil
}

func (c *CarsRepository) Update(email string, carID string, body models.UpdateCar) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	update := body.Update()
	if update == nil {
		return nil
	}
	filter := bson.M{"_id": docID, "email": email}

	_, err = c.db.Collection("cars").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *CarsRepository) Delete(email string, carID string) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": docID, "email": email}
	_, err = c.db.Collection("cars").DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
