package cars

import (
	"context"

	// "database/sql"
	// "github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db             *mongo.Database
	collectionName string
}

func NewInstanceOfCarsRepository(db *mongo.Database) Repository {
	return Repository{db: db, collectionName: "cars"}
}

func (c *Repository) List(email string, query ListCarQuery) ([]Car, error) {
	filters := query.Filter(email)
	var cars []Car

	options := options.Find()

	// Add paging
	options.SetLimit(int64(query.Limit))
	options.SetSkip(int64((query.Page * query.Limit) - query.Limit))

	// Add timestamp
	options.SetSort(bson.M{"created": -1})

	cursor, err := c.db.Collection(c.collectionName).Find(context.Background(), filters, options)
	if err != nil {
		return []Car{}, err
	}

	for cursor.Next(context.Background()) {
		car := Car{}
		err := cursor.Decode(&car)
		if err != nil {
			//handle err
		} else {
			cars = append(cars, car)
		}
	}
	return cars, nil
}

func (c *Repository) Get(email string, carID string) (Car, error) {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return Car{}, err
	}

	filter := bson.M{"_id": docID, "email": email}

	var result Car
	err = c.db.Collection(c.collectionName).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return Car{}, err
	}
	return result, nil
}

func (c *Repository) Save(car Car) error {
	_, err := c.db.Collection(c.collectionName).InsertOne(context.TODO(), car)
	if err != nil {
		return err
	}
	// insertResult.InsertedID.(string)
	return nil
}

func (c *Repository) Update(email string, carID string, body UpdateCar) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	update := body.Update()
	if update == nil {
		return nil
	}
	filter := bson.M{"_id": docID, "email": email}

	_, err = c.db.Collection(c.collectionName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *Repository) Delete(email string, carID string) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": docID, "email": email}
	_, err = c.db.Collection(c.collectionName).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
