package dao

import (
	"auth-service/internal"
	"auth-service/internal/custom_error"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return UserDao{db}
}

func (ud UserDao) FindByEmail(ctx context.Context, email string) (internal.User, error) {
	user := internal.User{}
	err := ud.db.WithContext(ctx).First(&user, "email=?", email).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, &custom_error.AppError{
			Err:           fmt.Errorf("user not found"),
			Message:       "user not found",
			HttpErrorCode: 404,
		}
	}
	return user, err
}

func (ud UserDao) SaveUser(ctx context.Context, user internal.User) (internal.User, error) {
	err := ud.db.WithContext(ctx).Save(&user).Error
	return user, err
}
