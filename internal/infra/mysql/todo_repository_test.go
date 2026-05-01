package mysql_test

import (
	"testing"

	"my-todo/internal/domain"
	"my-todo/internal/infra/mysql"

	"github.com/stretchr/testify/assert"
)

func TestTodoRepository_Create(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)

	todo, err := domain.NewTodo("買い物")
	assert.NoError(t, err)

	err = repo.Create(todo)
	assert.NoError(t, err)
}
