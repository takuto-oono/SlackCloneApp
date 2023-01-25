package token

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"backend/config"
)

func GenerateToken(userId uint32) (string, error) {
	token_lifespan := 2

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config.SecretKey))
}

func GetTokenFromContext(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	fmt.Println(token)
	if len(strings.Split(token, " ")) == 2 {
		fmt.Println(strings.Split(token, " ")[1])
		return strings.Split(token, " ")[1]
	}
	return ""
}

func GetUserIdFromToken(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(1)
			return nil, fmt.Errorf("err")
		}
		return []byte(config.Config.SecretKey), nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uint(uid)), nil
	}
	return 0, nil
}
