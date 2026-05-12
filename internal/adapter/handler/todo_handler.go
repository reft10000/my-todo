package handler

import (
	"context"
	"net/http"
	"strconv"

	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TodoUsecaseInterface interface {
	Create(title string) (*domain.Todo, error)
	List(ctx context.Context, input usecase.ListInput) (*domain.ListResult[domain.Todo], error)
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

type listTodoResponse struct {
	Todos   []createTodoResponse `json:"todos"`
	Total   int                  `json:"total"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	HasNext bool                 `json:"has_next"`
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

func (h *TodoHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.Error(domain.ErrBadRequest)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		c.Error(domain.ErrBadRequest)
		return
	}

	result, err := h.usecase.List(c.Request.Context(), usecase.ListInput{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		c.Error(err)
		return
	}

	todos := make([]createTodoResponse, len(result.Todos))
	for i, t := range result.Todos {
		todos[i] = createTodoResponse{
			ID:     t.ID.String(),
			Title:  t.Title,
			Status: string(t.Status),
		}
	}

	c.JSON(http.StatusOK, listTodoResponse{
		Todos:   todos,
		Total:   result.Total,
		Page:    result.Page,
		Limit:   result.Limit,
		HasNext: result.HasNext,
	})
}
