package main

import (
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()

	r.GET("/test/:x", func(c *gin.Context) {
		x := c.Param("x")
		c.IndentedJSON(http.StatusOK, "Hello Golang" + x)
	})
	r.Run(":8000")
}
