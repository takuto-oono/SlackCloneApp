package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"backend/token"
)

func Authenticate(c *gin.Context) (uint32, error) {
	tokenString := token.GetTokenFromContext(c)
	if tokenString == "" {
		return 0, fmt.Errorf("token not found from context")
	}
	return token.GetUserIdFromToken(tokenString)
}
