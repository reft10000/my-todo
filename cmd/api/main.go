package main

import (
	"my-todo/internal/adapter/router"
)

func main() {
	router := router.SetupRouter()
	router.Run(":8080")
}
