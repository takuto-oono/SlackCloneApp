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
	if controllerUtils.IsExistUserSameUsernameAndPassword(u.Name, u.PassWord) {
		c.JSON(http.StatusConflict, gin.H{"message": "already exist same username and password"})
		return
	}

	// dbに登録
	if err := u.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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

	// usernameとpasswordからIDを特定
	u, err := models.GetUserByNameAndPassword(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// jwtTokenを作成
	token, err := token.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
