package service

import (
	"auth-service/src/models"
	"context"
)

type AuthService interface {
	GenerateToken(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error)
	SaveUser(ctx context.Context, user models.UserCreateDto) (*models.UserDto, error)
	GetUserByToken(ctx context.Context, tokenString string) (*models.UserDto, error)
	LogInThroughSocialNetwork(ctx context.Context, user models.SocialNetworkUser) (*models.AuthResponse, error)
}
