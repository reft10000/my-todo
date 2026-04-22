package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique(),
		field.String("title").
			NotEmpty(),
		field.Enum("status").
			Values("pending", "in_progress", "done").
			Default("pending"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}
