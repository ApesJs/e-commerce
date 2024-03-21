package fakers

import (
	"e-commerce/app/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		ID:            uuid.New().String(),
		FirstName:     faker.FirstName(),
		LastName:      faker.LastName(),
		Email:         faker.Email(),
		Password:      faker.Password(),
		RememberToken: faker.Jwt(),
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     gorm.DeletedAt{},
	}
}
