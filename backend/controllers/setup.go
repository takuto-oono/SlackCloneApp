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
	workspace.PATCH("/rename/:workspace_id", RenameWorkspaceName)
	workspace.DELETE("/delete_user", DeleteUserFromWorkSpace)
	workspace.GET("/get_by_user", GetWorkspacesByUserId)

	channel := api.Group("/channel")
	channel.POST("/create", CreateChannel)
	channel.POST("/add_user", AddUserInChannel)
	channel.DELETE("/delete_user/:workspace_id", DeleteUserFromChannel)
	channel.DELETE("/delete", DeleteChannel)
	channel.GET("/get_user_belongs_in_workspace/:workspace_id", GetChannelsByUser)

	message := api.Group("/message")
	message.POST("/send", SendMessage)
	message.GET("/get_from_channel/:channel_id", GetAllMessagesFromChannel)

	return r
}
