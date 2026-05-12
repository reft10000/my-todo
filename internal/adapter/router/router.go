package router

import (
	"net/http"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/adapter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(todoHandler *handler.TodoHandler) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/todos", todoHandler.Create)
	router.GET("/todos", todoHandler.List)
	return router
}
