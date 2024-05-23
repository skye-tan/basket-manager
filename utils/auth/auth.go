package auth

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	custom_error "github.com/skye-tan/basket-manager/utils"
)

var secret_key = make([]byte, 16)

func GenerateSecretKey() {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i := range secret_key {
		secret_key[i] = letters[rand.Intn(len(letters))]
	}
}

func CreateToken(user_id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user_id,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})

	token_string, err := token.SignedString(secret_key)
	if err != nil {
		return "", err
	}

	return token_string, nil
}

func VerifyToken(token_string string) (uint, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New(custom_error.INVALID_TOKEN)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New(custom_error.INVALID_TOKEN)
	}

	return uint(claims["user_id"].(float64)), nil
}
