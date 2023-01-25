package handler

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
	"backend/token"
)

type UserInput struct {
	Name     string `json:"name"`
	PassWord string `json:"password"`
}

func (uin *UserInput) validate() bool {
	if uin.Name == "" {
		return false
	}
	if uin.PassWord == "" {
		return false
	}
	return true
}

func SignUp(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !input.validate() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or password is blank"})
		return
	}

	user := models.NewUser(rand.Uint32(), input.Name, input.PassWord)
	err := user.Create()
	if err == nil {
		c.IndentedJSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": err})
	}
}

func Login(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !input.validate() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or password in blank"})
		return
	}
	user, err := models.GetUserByName(input.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if user.PassWord != input.PassWord {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong password"})
		return
	}

	// TODO generate Token
	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, token)
}

func GetCurrentUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	tokenString := token.GetTokenFromContext(c)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found jwt token"})
		return
	}
	userId, err := token.GetUserIdFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
