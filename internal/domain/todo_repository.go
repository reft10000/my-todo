package domain

import (
	"context"

	"github.com/google/uuid"
)

type ListParams struct {
	Pagination PaginationParams
	Filter     *FilterParams
	Sort       *SortParams
}

type PaginationParams struct {
	Page  int
	Limit int
}

type FilterParams struct {
	Status *string
	UserID *string
}

type SortParams struct {
	Field string
	Desc  bool
}

type ListResult[t Todo] struct {
	Todos   []t
	Total   int
	Page    int
	Limit   int
	HasNext bool
}

type TodoRepository interface {
	Create(todo *Todo) error
	List(ctx context.Context, params ListParams) (ListResult[Todo], error)
	FindByID(ctx context.Context, id uuid.UUID) (*Todo, error)
}
