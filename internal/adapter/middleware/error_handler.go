package middleware

import (
	"errors"
	"my-todo/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if errors.As(err, &domain.ErrInvalidTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
	}
}
