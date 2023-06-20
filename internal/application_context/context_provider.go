package application_context

import (
	"auth-service/internal"
	"auth-service/internal/util"
	"github.com/golobby/container/v3"
)

func ResolveAuthController() internal.AuthController {
	var authController internal.AuthController
	containerErr := container.Resolve(&authController)
	if containerErr != nil {
		panic("AuthController impl is not fount")
	}

	return authController
}

func ResolveAuthService() internal.AuthService {
	var authService internal.AuthService
	containerErr := container.Resolve(&authService)
	if containerErr != nil {
		panic("AuthService impl is not fount")
	}

	return authService
}

func ResolveLogger() util.CustomLogger {
	var logger util.CustomLogger
	containerErr := container.Resolve(&logger)
	if containerErr != nil {
		panic("Custom logger impl is not fount")
	}

	return logger
}
