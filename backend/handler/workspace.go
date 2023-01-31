package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
	"backend/token"
)

type WorkspaceInput struct {
	Name string `json:"name"`
}

func CreateWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	tokenString := token.GetTokenFromContext(c)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found jwt token"})
		return
	}
	primaryOwnerId, err := token.GetUserIdFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var input WorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspace name"})
		return
	}
	w := models.NewWorkspace(0, input.Name, primaryOwnerId)
	if err := w.CreateWorkspace(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, w)
}
