package main

import (
	"my-todo/internal/adapter/middleware"
	"my-todo/internal/adapter/router"
)

func main() {
	router := router.SetupRouter()
	router.Use(middleware.ErrorHandler())
	router.Run(":8080")
}
