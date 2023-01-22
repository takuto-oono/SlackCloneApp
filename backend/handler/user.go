package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/models"
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

	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	user := models.NewUser(uuid, input.Name, input.PassWord)
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

	c.IndentedJSON(http.StatusOK, user)
	
}
