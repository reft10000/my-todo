package mysql_test

import (
	"context"
	"testing"

	"my-todo/internal/domain"
	"my-todo/internal/infra/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoRepository_Create(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)

	todo, err := domain.NewTodo("買い物")
	assert.NoError(t, err)

	err = repo.Create(todo)
	assert.NoError(t, err)
}

func TestTodoRepository_List(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)
	ctx := context.Background()

	titles := []string{"買い物", "掃除", "洗濯"}
	for _, title := range titles {
		todo, err := domain.NewTodo(title)
		require.NoError(t, err)
		err = repo.Create(todo)
		require.NoError(t, err)
	}

	result, err := repo.List(ctx, domain.ListParams{
		Pagination: domain.PaginationParams{Page: 1, Limit: 2},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Todos))
	assert.Equal(t, 3, result.Total)
	assert.True(t, result.HasNext)
}
