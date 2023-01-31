package main

import (
	"backend/handler"
)

func main() {
	r := handler.SetupRouter()
	r.Run(":8080")
}
