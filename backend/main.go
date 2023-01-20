package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/handler"
	"backend/models"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	fmt.Println(models.DbConnection)

	// user handler (api test 2)
	r.GET("/users", handler.GetUsers)
	r.GET("/user/:id", handler.GetUser)
	r.POST("/user", handler.PostUser)
	r.PATCH("/user/:id", handler.UpdateUser)
	r.DELETE("/user/:id", handler.DeleteUser)

	// test api 1 handler
	r.GET("/test/:x", func(c *gin.Context) {
		x := c.Param("x")
		c.IndentedJSON(http.StatusOK, "Hello Golang"+x)
	})
	r.Run(":8080")
}
