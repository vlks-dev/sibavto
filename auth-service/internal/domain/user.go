package domain

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"  validate:"required"`
	Surname    string          `json:"surname"  validate:"required"`
	Patronymic string          `json:"patronymic"  validate:"required"`
	Email      string          `json:"email"  validate:"required,email"`
	Roles      map[string]bool `json:"roles"`
	Password   string          `json:"hash" validate:"required"`
}
