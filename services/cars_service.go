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

func (c *CarsService) GetAll(authToken string, query models.ListCarQuery) ([]models.Car, error) {
	// Get session
	_, session, err := c.userRepository.GetSessionById(authToken)
	if err != nil {
		return []models.Car{}, err
	}

	cars, err := c.carsRepository.List(session.Email, query)
	if err != nil {
		return []models.Car{}, err
	}
	return cars, nil
}

func (c *CarsService) GetByID(authToken string, carID string) (models.Car, error) {
	// Get session
	_, session, err := c.userRepository.GetSessionById(authToken)
	if err != nil {
		return models.Car{}, err
	}

	car, err := c.carsRepository.Get(session.Email, carID)
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (c *CarsService) Create(authToken string, body models.CreateCar) error {
	// Get session
	_, session, err := c.userRepository.GetSessionById(authToken)
	if err != nil {
		return err
	}

	// Create new car object
	car := models.Car{
		Make:    body.Make,
		Model:   body.Model,
		Year:    body.Year,
		Status:  "",
		Created: time.Now(),
		Email:   session.Email,
	}
	err = c.carsRepository.Save(car)
	if err != nil {
		return err
	}
	return nil
}

func (c *CarsService) Update(authToken string, carID string, body models.UpdateCar) error {
	// Get session
	_, session, err := c.userRepository.GetSessionById(authToken)
	if err != nil {
		return err
	}

	// Update car
	err = c.carsRepository.Update(session.Email, carID, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *CarsService) Delete(authToken string, carID string) error {
	// Get session
	_, session, err := c.userRepository.GetSessionById(authToken)
	if err != nil {
		return err
	}

	// Delete car
	err = c.carsRepository.Delete(session.Email, carID)
	if err != nil {
		return err
	}
	return nil
}
