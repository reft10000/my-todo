package handler

import (
	"net/http"

	"my-todo/internal/domain"

	"github.com/gin-gonic/gin"
)

type TodoUsecaseInterface interface {
	Create(title string) (*domain.Todo, error)
}

type TodoHandler struct {
	usecase TodoUsecaseInterface
}

func NewTodoHandler(uc TodoUsecaseInterface) *TodoHandler {
	return &TodoHandler{usecase: uc}
}

type createTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

type createTodoResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}
	todo, err := h.usecase.Create(req.Title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createTodoResponse{
		ID:     todo.ID.String(),
		Title:  todo.Title,
		Status: string(todo.Status),
	})
}
