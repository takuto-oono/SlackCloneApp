package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/token"
)

func Authenticate(c *gin.Context) (uint32, error) {
	tokenString := token.GetTokenFromContext(c)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found jwt token"})
		return 0, fmt.Errorf("not found token from context")
	}
	return token.GetUserIdFromToken(tokenString)
}
