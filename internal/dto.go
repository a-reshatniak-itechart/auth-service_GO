package internal

import "github.com/google/uuid"

type UserDto struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
}
type UserCreateDto struct {
	Password string `validate:"required"`
	Email    string `validate:"email"`
}
