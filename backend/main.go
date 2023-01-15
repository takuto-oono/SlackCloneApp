package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/:x", func(c *gin.Context) {
		x, err := strconv.Atoi(c.Param("x"))
		if err != nil {
			c.JSON(400, gin.H{
				"message": "error convert int",
			})
		} else {
			c.JSON(200, gin.H{
				"x": x + 1,
			})
		}
	})

	router.Run(":8000")
}
