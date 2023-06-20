package dao

import (
	"auth-service/src/models"
	"context"
)

type UserDatabase interface {
	FindByEmail(ctx context.Context, email string) (models.User, error)
	SaveUser(ctx context.Context, user models.User) (models.User, error)
}
