package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/models"
)

type userInput struct {
	Name     string `json:"name"`
	PassWord string `json:"password"`
}

func GetUsers(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	users, err := models.GetAllUsers()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Param("id")
	user, err := models.GetUserById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func PostUser(c *gin.Context) {
	// test curl cmd
	// curl -X POST -H "Content-Type: application/json" -d '{"name": "abc", "password":"pass"}' http://localhost:8000/user
	c.Header("Access-Control-Allow-Origin", "*")

	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error message": err})
		return
	}
	
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	user := models.NewUser(uuid, input.Name, input.PassWord)
	err := user.Create()
	if err == nil {
		c.IndentedJSON(http.StatusOK, input)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": err})
	}
}

func UpdateUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error message": "not found id"})
		return
	}

	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error message": "not found user information"})
		return
	}
	user := models.NewUser(id, input.Name, input.PassWord)

	if err := user.UpdateUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error message": "not found id"})
		return
	}
	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error message": "not found user information"})
		return
	}
	user := models.NewUser(id, input.Name, input.PassWord)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error message": "not found user information"})
		return
	}
	if err := user.DeleteUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
