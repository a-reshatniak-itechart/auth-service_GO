package util

import (
	"auth-service/internal"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

// todo required refactoring
const hmacSampleSecret = "secret"

func GenerateJwt(user internal.User) (string, string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(10 * time.Minute).Unix(),
		"authorized": true,
		"user":       user.Email,
	})

	refreshToken := "refreshToken"
	stringToken, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		//todo handle error
		log.Fatal(err)
	}

	log.Print("GenerateJwt. Token:" + stringToken)

	return stringToken, refreshToken
}

func VerifyJwt(tokenString string) (string, error) {
	log.Print("VerifyJwt. Token:" + tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	email := fmt.Sprintf("%v", claims["user"])
	fmt.Print("User email: " + email)
	return email, nil
}
