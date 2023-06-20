package controller

import (
	"auth-service/src/models"
	"context"
)

type AuthController interface {
	GenerateToken(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error)
	SaveUser(ctx context.Context, req models.UserCreateDto) (*models.UserDto, error)
	GetUserByToken(ctx context.Context, tokenString string) (*models.UserDto, error)
}
