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

func GetTokenFromHeader(c *gin.Context) string {
	return c.Request.Header.Get("Authorization")
}

func GetTokenFromQueryParams(c *gin.Context) string {
	return c.Request.URL.Query().Get("token")
}

func GetTokenFromContext(c *gin.Context) string {
	token := GetTokenFromHeader(c)
	if token == "" {
		token = GetTokenFromQueryParams(c)
	}
	if len(strings.Split(token, "\"")) == 3 {
		return strings.Split(token, "\"")[1]
	}
	return token
}

func GetUserIdFromToken(tokenString string) (uint32, error) {
	jwtToken, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("err")
		}
		return []byte(config.Config.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uint(uid)), nil
	}
	return 0, fmt.Errorf("error in GetUserIdFromToken func")
}
