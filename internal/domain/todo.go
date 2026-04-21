package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Todo struct {
	ID        uuid.UUID
	Title     string
	Status    Status
	CreatedAt time.Time
}

func NewTodo(title string) (*Todo, error) {
	if !isValidTitle(title) {
		return nil, errors.New("title is required")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate uuid")
	}

	return &Todo{
		ID:        id,
		Title:     title,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}, nil
}

func isValidTitle(title string) bool {
	return strings.TrimSpace(title) != ""
}
