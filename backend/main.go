package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"backend/controllers"
	"backend/models"
	"backend/testdata"
)

func main() {
	gin.SetMode(gin.DebugMode)

	fmt.Println(models.DbConnection)
	switch os.Args[1:][0] {
	case "1":
		r := controllers.SetupRouter1()
		r.Run(":8080")
	case "2":
		r2 := controllers.SetupRouter2()
		r2.Run(":8000")
	case "testdata":
		testdata.TestDataMain()
	default:
		panic("don't run server")
	}
}
