package services

import (
	// "fmt"
	"time"
	// "strings"
	"go-boilerplate/models"
	"go-boilerplate/repositories"
	// "errors"
	// "github.com/google/uuid"
)

type CarsService struct {
	userRepository repositories.UserRepository
	carsRepository repositories.CarsRepository
}

func NewInstanceOfCarsService(userRepository repositories.UserRepository, carsRepository repositories.CarsRepository) CarsService {
	return CarsService{userRepository: userRepository, carsRepository: carsRepository}
}

func (c *CarsService) GetAll(session models.Session, query models.ListCarQuery) ([]models.Car, error) {
	cars, err := c.carsRepository.List(session.Email, query)
	if err != nil {
		return []models.Car{}, err
	}
	return cars, nil
}

func (c *CarsService) GetByID(session models.Session, carID string) (models.Car, error) {
	car, err := c.carsRepository.Get(session.Email, carID)
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (c *CarsService) Create(session models.Session, body models.CreateCar) error {
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

func (c *CarsService) Update(session models.Session, carID string, body models.UpdateCar) error {
	// Update car
	err := c.carsRepository.Update(session.Email, carID, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *CarsService) Delete(session models.Session, carID string) error {
	// Delete car
	err := c.carsRepository.Delete(session.Email, carID)
	if err != nil {
		return err
	}
	return nil
}
