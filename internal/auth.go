package internal

import "context"

type AuthRequest struct {
	Password string
	Email    string
}

type AuthResponse struct {
	Jwt     string
	Refresh string
}

type SocialNetworkUser struct {
	Email     string
	FirstName string
	LastName  string
}

type AuthService interface {
	GenerateToken(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	SaveUser(ctx context.Context, user UserCreateDto) (*UserDto, error)
	GetUserByToken(ctx context.Context, tokenString string) (*UserDto, error)
	LogInThroughSocialNetwork(ctx context.Context, user SocialNetworkUser) (*AuthResponse, error)
}
