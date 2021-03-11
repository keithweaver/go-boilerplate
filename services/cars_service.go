package services

import (
	"context"
	"go-boilerplate/logging"
	// "fmt"
	"time"
	// "strings"
	"go-boilerplate/models"
	"go-boilerplate/repositories"
	// "errors"
	// "github.com/google/uuid"
)

type CarsService struct {
	logger logging.Logger
	userRepository repositories.UserRepository
	carsRepository repositories.CarsRepository
}

func NewInstanceOfCarsService(logger logging.Logger, userRepository repositories.UserRepository, carsRepository repositories.CarsRepository) CarsService {
	return CarsService{logger, userRepository, carsRepository}
}

func (c *CarsService) GetAll(ctx context.Context, session models.Session, query models.ListCarQuery) ([]models.Car, error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "GetAll")

	cars, err := c.carsRepository.List(session.Email, query)
	if err != nil {
		return []models.Car{}, err
	}
	return cars, nil
}

func (c *CarsService) GetByID(ctx context.Context, session models.Session, carID string) (models.Car, error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "GetByID")

	car, err := c.carsRepository.Get(session.Email, carID)
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (c *CarsService) Create(ctx context.Context, session models.Session, body models.CreateCar) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Create")
	// Create new car object
	car := models.Car{
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

func (c *CarsService) Update(ctx context.Context, session models.Session, carID string, body models.UpdateCar) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Update")
	// Update car
	err := c.carsRepository.Update(session.Email, carID, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *CarsService) Delete(ctx context.Context, session models.Session, carID string) error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "Delete")

	// Delete car
	err := c.carsRepository.Delete(session.Email, carID)
	if err != nil {
		return err
	}
	return nil
}
