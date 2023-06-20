package controller

import (
	"auth-service/src/custom_validator"
	"auth-service/src/models"
	"auth-service/src/service"
	"context"
)

type AuthControllerImpl struct {
	authService service.AuthService
	validator   custom_validator.CustomValidator
}

func NewAuthController(authService service.AuthService) AuthController {
	return AuthControllerImpl{authService: authService, validator: custom_validator.NewValidator()}
}

func (ac AuthControllerImpl) GenerateToken(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	return ac.authService.GenerateToken(ctx, req)
}

func (ac AuthControllerImpl) SaveUser(ctx context.Context, req models.UserCreateDto) (*models.UserDto, error) {
	err := ac.validator.Validate(&req)
	if err != nil {
		return nil, err
	}
	return ac.authService.SaveUser(ctx, req)
}

func (ac AuthControllerImpl) GetUserByToken(ctx context.Context, tokenString string) (*models.UserDto, error) {
	return ac.authService.GetUserByToken(ctx, tokenString)
}
