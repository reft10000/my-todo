package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// MaxTitleLength はDBの title カラム (varchar(255)) に合わせた文字数上限
const MaxTitleLength = 255

type Todo struct {
	ID        uuid.UUID
	Title     string
	Status    Status
	CreatedAt time.Time
}

func NewTodo(title string) (*Todo, error) {
	if !isValidTitle(title) {
		return nil, ErrInvalidTitle
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

func (t *Todo) ChangeTitle(title string) error {
	if !isValidTitle(title) {
		return ErrInvalidTitle
	}
	t.Title = title
	return nil
}

func (t *Todo) ChangeStatus(status Status) error {
	if !status.isValid() {
		return ErrInvalidStatus
	}
	t.Status = status
	return nil
}

func (s Status) isValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func isValidTitle(title string) bool {
	return strings.TrimSpace(title) != "" && utf8.RuneCountInString(title) <= MaxTitleLength
}
