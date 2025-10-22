package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Добавляем в токен всю необходимую информацию
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	// Подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
