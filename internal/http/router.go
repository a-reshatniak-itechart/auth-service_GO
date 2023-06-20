package http

import (
	_ "auth-service/docs" //is needed for swagger
	"auth-service/internal"
	"auth-service/internal/application_context"
	"auth-service/internal/custom_error"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

// @Summary      generates token
// @Description  get token by creds
// @Tags         token
// @Accept       json
// @Produce      json
// @Param AuthRequest body models.AuthRequest true "auth request body"
// @Success      200  {object}  models.AuthResponse
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /token [post]
func getToken(context *gin.Context, controller internal.AuthController) {
	authRequest := &internal.AuthRequest{}
	err := context.BindJSON(authRequest)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	obj, saveErr := controller.GenerateToken(context, authRequest)
	resolveResponse(obj, saveErr, context)
}

// @Summary      saves new user
// @Description  saves new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param UserCreateDto body models.UserCreateDto true "UserCreateDto"
// @Success      200  {object}  models.UserDto
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /user [post]
func saveUser(context *gin.Context, controller internal.AuthController) {
	user := &internal.UserCreateDto{}
	err := context.BindJSON(user)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	obj, saveErr := controller.SaveUser(context, *user)
	resolveResponse(obj, saveErr, context)
}

// @Summary      Returns user dto
// @Description  get existing user by valid token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param token query string true "valid jwt"
// @Success      200  {object}  models.UserDto
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /user [get]
func getUserByToken(context *gin.Context, controller internal.AuthController) {
	token := context.Query("token")
	if &token == nil {
		context.AbortWithError(http.StatusBadRequest, fmt.Errorf("jwt is null"))
		return
	}
	obj, err := controller.GetUserByToken(context, token)
	resolveResponse(obj, err, context)
}

// @Summary      Google auth
// @Description  Redirects user to google auth page
// @Tags         auth
// @Accept       json
// @Success      302
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /auth/google [get]
func redirectToGoogleAuthPage(context *gin.Context) {
	req := context.Request
	v := req.URL.Query()
	v.Add("provider", "google")
	req.URL.RawQuery = v.Encode()

	gothic.BeginAuthHandler(context.Writer, context.Request)
}

// @Summary      Google auth
// @Description  Returns google auth url
// @Tags         auth
// @Accept       json
// @Success      200  {string}  string
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /auth/google/url [get]
func getGoogleAuthPage(context *gin.Context) {
	req := context.Request
	v := req.URL.Query()
	v.Add("provider", "google")
	req.URL.RawQuery = v.Encode()
	url, err := gothic.GetAuthURL(context.Writer, context.Request)
	var appErr *custom_error.AppError = nil
	if err != nil {
		appErr = &custom_error.AppError{Message: "google auth failed", Err: err, HttpErrorCode: 400}
	}

	resolveResponse(url, appErr, context)
}

// @Summary      Completes auth through google
// @Description  authenticate exiting or create a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param code query string true
// @Param scope query string true
// @Param authuser query string true
// @Param prompt query string true
// @Success      200  {object}  models.AuthResponse
// @Failure      400  {object}  custom_error.AppError
// @Failure      404  {object}  custom_error.AppError
// @Failure      500  {object}  custom_error.AppError
// @Router       /auth/google/callback [get]
func getGoogleAuthCallback(context *gin.Context, authService internal.AuthService) {
	user, err := gothic.CompleteUserAuth(context.Writer, context.Request)
	if err != nil {
		resolveResponse(
			nil,
			&custom_error.AppError{Message: "Login through google is failed", HttpErrorCode: 400, Err: err},
			context,
		)
		return
	}

	authResponse, err := authService.LogInThroughSocialNetwork(
		context,
		internal.SocialNetworkUser{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	)

	resolveResponse(authResponse, err, context)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func InitRoutes(r *gin.Engine) {
	logger := application_context.ResolveLogger()
	logger.Info("Initiating routes")
	controller := application_context.ResolveAuthController()
	authService := application_context.ResolveAuthService()

	r.POST("/token", func(context *gin.Context) {
		getToken(context, controller)
	})

	r.POST("/user", func(context *gin.Context) {
		saveUser(context, controller)
	})

	r.GET("/user", func(context *gin.Context) {
		getUserByToken(context, controller)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/auth/google", redirectToGoogleAuthPage)

	r.GET("/auth/google/url", getGoogleAuthPage)

	r.GET("/auth/google/callback", func(context *gin.Context) {
		getGoogleAuthCallback(context, authService)
	})

	r.GET("/metrics", prometheusHandler())
	logger.Info("Routes initiated")
}

func resolveResponse(obj any, err error, context *gin.Context) {
	if err == nil {
		context.JSON(http.StatusOK, obj)
		return
	}

	appErr, ok := err.(*custom_error.AppError)
	if !ok {
		appErr = &custom_error.AppError{
			Message:       "internal error",
			HttpErrorCode: http.StatusInternalServerError,
			Err:           err,
		}
	}

	context.JSON(appErr.HttpErrorCode, appErr.Message)
}
