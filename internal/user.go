package internal

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Password     string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"uniqueIndex;not null"`
	FirstName    string    `gorm:"type:varchar(255)"`
	LastName     string    `gorm:"type:varchar(255)"`
	RefreshToken string    `gorm:"type:varchar(255)"`
}

type UserDatabase interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	SaveUser(ctx context.Context, user User) (User, error)
}
