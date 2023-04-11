package main

import (
	_ "auth-service/docs" //is needed for swagger
	"auth-service/src/application_context"
	"auth-service/src/controller"
	"auth-service/src/custom_error"
	"auth-service/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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
func getToken(context *gin.Context, controller controller.AuthController) {
	authRequest := &models.AuthRequest{}
	err := context.BindJSON(authRequest)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	obj, saveErr := controller.GenerateToken(authRequest)
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
func saveUser(context *gin.Context, controller controller.AuthController) {
	user := &models.UserCreateDto{}
	err := context.BindJSON(user)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	obj, saveErr := controller.SaveUser(*user)
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
func getUserByToken(context *gin.Context, controller controller.AuthController) {
	token := context.Query("token")
	if &token == nil {
		context.AbortWithError(http.StatusBadRequest, fmt.Errorf("jwt is null"))
		return
	}
	obj, err := controller.GetUserByToken(token)
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
		appErr = &custom_error.AppError{Message: "google auth failed", Error: err, HttpErrorCode: 400}
	}

	resolveResponse(url, appErr, context)
}
func InitRoutes(r *gin.Engine) {
	controller := application_context.ResolveAuthController()
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
		user, err := gothic.CompleteUserAuth(context.Writer, context.Request)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Google user:")
		fmt.Println(user)
	})

}

func resolveResponse(obj any, err *custom_error.AppError, context *gin.Context) {
	if err != nil {
		context.JSON(err.HttpErrorCode, err.Message)
	} else {
		context.JSON(http.StatusOK, obj)
	}
}
