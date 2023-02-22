package controllers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
	"backend/token"
)

func SignUp(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	// bodyの情報を取得
	// if err := c.ShouldBindJSON(&u); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }

	// // 必要な項目が取得できているか確認
	// if u.Name == "" || u.PassWord == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "username pr password is blank"})
	// 	return
	// }

	ui, err := controllerUtils.InputSignUpAndLogin(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// IDを確定
	u := models.NewUser(rand.Uint32(), ui.Name, ui.Password)

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
	// bodyの情報を取得

	input, err := controllerUtils.InputSignUpAndLogin(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	u := models.NewUser(0, input.Name, input.Password)

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
