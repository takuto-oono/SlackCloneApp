package controllers

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

func SignUp(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	// bodyの情報を取得
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 必要な項目が取得できているか確認
	if u.Name == "" || u.PassWord == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username pr password is blank"})
		return
	}

	// IDを確定
	u.ID = rand.Uint32()

	// 既に同じuserNameとpasswordの組み合わせのユーザーが存在しないかを確認
	b, err := u.IsExistUserSameUsernameAndPassword()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if b {
		c.JSON(http.StatusBadRequest, gin.H{"message": "already exist same username and password"})
		return
	}

	// dbに登録
	if err := u.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, u)
}

func Login(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// var input UserInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }
	// if !input.validate() {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "username or password in blank"})
	// 	return
	// }
	// user, err := models.GetUserByName(input.Name)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }
	// if user.PassWord != input.PassWord {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "wrong password"})
	// 	return
	// }

	// token, err := token.GenerateToken(user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }

	// bodyの情報を取得
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 必要な情報が取得できているか確認
	if u.Name == "" || u.PassWord == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or password is blank"})
		return
	}

	// usernameとpasswordからIDを特定
	if err := u.GetUserByNameAndPassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := token.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token, "user_id": u.ID, "username": u.Name})
}

func GetCurrentUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
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
