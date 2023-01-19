package main

import (
	"fmt"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"

	"backend/handler"
	"backend/models"
)

func main() {
	r := gin.Default()
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
	r.Run(":8000")
}
