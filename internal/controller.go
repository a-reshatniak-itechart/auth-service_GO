package internal

import (
	"context"
)

type AuthController interface {
	GenerateToken(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	SaveUser(ctx context.Context, req UserCreateDto) (*UserDto, error)
	GetUserByToken(ctx context.Context, tokenString string) (*UserDto, error)
}
