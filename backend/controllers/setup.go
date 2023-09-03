package controllers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/ws"
)

func settingRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Headers", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	return r
}

func SetupRouter1() *gin.Engine {
	r := settingRouter()
	api := r.Group("/api")

	user := api.Group("/user")
	user.POST("/signUp", SignUp)
	user.POST("/login", Login)
	user.GET("/currentUser", GetCurrentUser)
	user.GET("/all", GetAllUsers)

	workspace := api.Group("/workspace")
	workspace.GET("/get_by_user", GetWorkspacesByUserId)
	workspace.GET("/get_users/:workspace_id", GetUsersInWorkspace)
	workspace.POST("/create", CreateWorkspace)
	workspace.POST("/add_user", AddUserInWorkspace)
	workspace.PATCH("/rename/:workspace_id", RenameWorkspaceName)
	workspace.DELETE("/delete_user", DeleteUserFromWorkSpace)

	channel := api.Group("/channel")
	channel.GET("/get_by_user_and_workspace/:workspace_id", GetChannelsByUser)
	channel.GET("/:workspace_id", GetChannelsByWorkspace)
	channel.GET("/all_user/:channel_id", GetAllUsersInChannel)
	channel.POST("/create", CreateChannel)
	channel.POST("/add_user", AddUserInChannel)
	channel.DELETE("/delete_user/:workspace_id", DeleteUserFromChannel)
	channel.DELETE("/delete", DeleteChannel)

	message := api.Group("/message")
	message.GET("/get_from_channel/:channel_id", GetAllMessagesFromChannel)
	message.POST("/send", SendMessage)
	message.POST("/read_by_user/:message_id", ReadMessageByUser)
	message.PATCH("/edit/:message_id", EditMessage)

	dm := api.Group("/dm")
	dm.GET("/:dm_line_id", GetDMsInLine)
	dm.GET("/dm_lines/:workspace_id", GetDMLines)
	dm.POST("/send", SendDM)
	dm.PATCH("/:dm_id", EditDM)
	dm.DELETE("/:dm_id", DeleteDM)

	thread := api.Group("/thread")
	thread.GET("/by_user/:workspace_id", GetThreadsByUser)
	thread.POST("/post", PostThread)

	mention := api.Group("/mention")
	mention.GET("/by_user/:workspace_id", GetMessagesMentionedByUser)

	return r
}

func SetupRouter2() *gin.Engine {
	r := settingRouter()
	h := ws.NewHub()
	go h.Run()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "server2 OK"})
	})

	socket := r.Group("/websocket")
	socket.GET("/", func(ctx *gin.Context) {
		WsController(ctx, h)
	})
	return r
}
