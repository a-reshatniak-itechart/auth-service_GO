package application_context

import (
	"auth-service/src/config"
	"auth-service/src/controller"
	"auth-service/src/dao"
	"auth-service/src/service"
	"auth-service/src/util"
	"github.com/golobby/container/v3"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"go.uber.org/zap"
)

func LoadContext(config *config.Config) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	customLogger := util.ZapCustomLogger{Logger: *logger}
	customLogger.Info("Loading context")

	db := dao.Configure(config)
	userDao := dao.NewUserDao(db)
	authService := service.NewAuthService(userDao, customLogger)
	authController := controller.NewAuthController(authService)

	goth.UseProviders(
		google.New(config.GoogleClientId, config.GoogleClientSecret, config.GoogleAuthCallbackUrl, "email", "profile"),
	)

	container.Singleton(func() controller.AuthController {
		return authController
	})

	container.Singleton(func() service.AuthService {
		return authService
	})

	container.Singleton(func() util.CustomLogger {
		return customLogger
	})

	customLogger.Info("Context loaded")
}
