package auth

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	custom_error "github.com/skye-tan/basket-manager/utils"
)

type CustomClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

var secret_key = make([]byte, 16)

func GenerateSecretKey() {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i := range secret_key {
		secret_key[i] = letters[rand.Intn(len(letters))]
	}
}

func CreateToken(user_id uint) (string, error) {
	claims := &CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_string, err := token.SignedString(secret_key)

	return token_string, err
}

func VerifyToken(token_string string) (uint, error) {
	token, err := jwt.ParseWithClaims(token_string, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})
	if err != nil {
		return 0, errors.New(custom_error.INVALID_TOKEN)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New(custom_error.INVALID_TOKEN)
}
