package controllers

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/models"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Headers", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	fmt.Println(models.DbConnection)
	api := r.Group("/api")

	user := api.Group("/user")
	user.POST("/signUp", SignUp)
	user.POST("/login", Login)
	user.GET("/currentUser", GetCurrentUser)

	workspace := api.Group("/workspace")
	workspace.POST("/create", CreateWorkspace)
	workspace.POST("/add_user", AddUserInWorkspace)
	workspace.POST("/rename", RenameWorkspaceName)
	workspace.POST("/delete_user", DeleteUserFromWorkSpace)

	channel := api.Group("/channel")
	channel.POST("/create", CreateChannel)
	channel.POST("/add_user/:workspace_id", AddUserInChannel)
	channel.POST("/delete_user/:workspace_id", DeleteUserFromChannel)
	channel.POST("/delete", DeleteChannel)
	return r
}
