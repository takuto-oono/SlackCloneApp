package main

import (
	"backend/controllers"
)

func main() {
	r := controllers.SetupRouter()
	r.Run(":8080")
}
