package controller

import (
	"auth-service/internal"
	"auth-service/internal/custom_validator"
	"context"
)

type Auth struct {
	authService internal.AuthService
	validator   custom_validator.CustomValidator
}

func NewAuthController(authService internal.AuthService) internal.AuthController {
	return Auth{authService: authService, validator: custom_validator.NewValidator()}
}

func (ac Auth) GenerateToken(ctx context.Context, req *internal.AuthRequest) (*internal.AuthResponse, error) {
	return ac.authService.GenerateToken(ctx, req)
}

func (ac Auth) SaveUser(ctx context.Context, req internal.UserCreateDto) (*internal.UserDto, error) {
	err := ac.validator.Validate(&req)
	if err != nil {
		return nil, err
	}
	return ac.authService.SaveUser(ctx, req)
}

func (ac Auth) GetUserByToken(ctx context.Context, tokenString string) (*internal.UserDto, error) {
	return ac.authService.GetUserByToken(ctx, tokenString)
}
