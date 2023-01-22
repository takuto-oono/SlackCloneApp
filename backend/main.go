package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/handler"
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

	user.POST("/signUp", handler.SignUp)
	return r
}

func main() {
	r := SetupRouter()
	r.Run(":8080")
}
