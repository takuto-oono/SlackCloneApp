package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"backend/config"
)

func GenerateToken(userId string) (string, error) {
	token_lifespan := 2

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config.SecretKey))

}
