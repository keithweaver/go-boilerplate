package cars

import (
	"context"
	"go-boilerplate/logging"
	"go-boilerplate/user"

	// "fmt"
	"time"
	// "common"
	// "github.com/google/uuid"
)

type Services struct {
	logger         logging.Logger
	userRepository user.Repository
	carsRepository Repository
}

func NewInstanceOfCarsServices(logger logging.Logger, userRepository user.Repository, carsRepository Repository) Services {
	return Services{logger, userRepository, carsRepository}
}

func (c *Services) GetAll(ctx context.Context, session user.Session, query ListCarQuery) ([]Car, error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "GetAll")

	cars, err := c.carsRepository.List(session.Email, query)
	if err != nil {
		return []Car{}, err
	}
	return cars, nil
}

func (c *Services) GetByID(ctx context.Context, session user.Session, carID string) (Car, error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "GetByID")

	car, err := c.carsRepository.Get(session.Email, carID)
	if err != nil {
		return Car{}, err
	}
	return car, nil
}

func (c *Services) Create(ctx context.Context, session user.Session, body CreateCar) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Create")
	// Create new car object
	car := Car{
		Make:    body.Make,
		Model:   body.Model,
		Year:    body.Year,
		Status:  "",
		Created: time.Now(),
		Email:   session.Email,
	}
	err := c.carsRepository.Save(car)
	if err != nil {
		return err
	}
	return nil
}

func (c *Services) Update(ctx context.Context, session user.Session, carID string, body UpdateCar) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Update")
	// Update car
	err := c.carsRepository.Update(session.Email, carID, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Services) Delete(ctx context.Context, session user.Session, carID string) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Delete")

	// Delete car
	err := c.carsRepository.Delete(session.Email, carID)
	if err != nil {
		return err
	}
	return nil
}
