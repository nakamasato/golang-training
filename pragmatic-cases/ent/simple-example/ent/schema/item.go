package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
// id text NOT NULL PRIMARY KEY,
// name VARCHAR(50) NOT NULL,
// status SMALLINT NOT NULL,
// created_at TIMESTAMP NOT NULL
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StructTag(`json:"oid,omitempty"`),
		field.String("name"),
		field.Int("status"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("categories", Category.Type),
	}
}
