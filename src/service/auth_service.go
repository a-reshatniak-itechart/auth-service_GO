package service

import (
	"auth-service/src/custom_error"
	"auth-service/src/dao"
	"auth-service/src/models"
	"auth-service/src/util"
	"context"
	"fmt"
	"github.com/devfeel/mapper"
)

type AuthServiceImpl struct {
	db     dao.UserDatabase
	mapper mapper.IMapper
	logger util.CustomLogger
}

func NewAuthService(db dao.UserDatabase, logger util.CustomLogger) AuthService {
	return AuthServiceImpl{db, mapper.NewMapper(), logger}
}

func (as AuthServiceImpl) GenerateToken(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
	as.logger.Info("Generate token request started")
	user, appErr := as.getUserByEmail(ctx, req.Email)

	if appErr != nil {
		return nil, appErr
	}
	if util.CheckPasswordHash(req.Password, user.Password) != true {
		return nil, &custom_error.AppError{
			Err:           fmt.Errorf("wrong password"),
			Message:       "wrong password",
			HttpErrorCode: 401,
		}
	}

	token, refreshToken := util.GenerateJwt(*user)
	user.RefreshToken = refreshToken
	_, err := as.db.SaveUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	as.logger.Info("Generate token request ended")
	return &models.AuthResponse{Jwt: token, Refresh: refreshToken}, nil
}

func (as AuthServiceImpl) SaveUser(ctx context.Context, userCreateDto models.UserCreateDto) (*models.UserDto, error) {
	as.logger.Info("Save user request started")
	hash, err := util.HashPassword(userCreateDto.Password)
	if err != nil {
		return nil, &custom_error.AppError{
			Err:           err,
			Message:       "hashing password error",
			HttpErrorCode: 400,
		}
	}
	userCreateDto.Password = hash
	userToSave := models.User{}

	_ = as.mapper.Mapper(&userCreateDto, &userToSave)

	savedUser, err := as.db.SaveUser(ctx, userToSave)
	if err != nil {
		return nil, &custom_error.AppError{
			Err:           err,
			Message:       "can't save user",
			HttpErrorCode: 500,
		}
	}

	dto := &models.UserDto{}
	_ = as.mapper.Mapper(&savedUser, dto)
	as.logger.Info("Save user request ended")
	return dto, nil
}

func (as AuthServiceImpl) getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	as.logger.Info("Get user by email started. User email: " + email)
	user, err := as.db.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	as.logger.Info("Get user by email ended. User email: " + email)
	return &user, nil
}

func (as AuthServiceImpl) GetUserByToken(ctx context.Context, tokenString string) (*models.UserDto, error) {
	as.logger.Info("Get user by token started")
	email, err := util.VerifyJwt(tokenString)
	if err != nil {
		return nil, &custom_error.AppError{
			Err:           err,
			Message:       "Jwt is invalid",
			HttpErrorCode: 403,
		}
	}

	user, err := as.db.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	dto := &models.UserDto{}
	_ = as.mapper.Mapper(&user, dto)
	as.logger.Info("Get user by token ended")
	return dto, nil
}

func (as AuthServiceImpl) LogInThroughSocialNetwork(ctx context.Context, user models.SocialNetworkUser) (*models.AuthResponse, error) {
	as.logger.Info("LogIn trough social network started. User: " + fmt.Sprintf("%v", user))
	dbUser, _ := as.db.FindByEmail(ctx, user.Email)
	if &dbUser == nil {
		dbUser = models.User{Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}
	}

	token, refreshToken := util.GenerateJwt(dbUser)
	dbUser.RefreshToken = refreshToken
	_, err := as.db.SaveUser(ctx, dbUser)
	if err != nil {
		return nil, err
	}
	as.logger.Info("LogIn trough social network ended")
	return &models.AuthResponse{Jwt: token, Refresh: refreshToken}, nil
}
