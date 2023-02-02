package handler

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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Headers", "Content-Type"},
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
	workspace.POST("/add_user", AddUserWorkspace)

	return r
}
