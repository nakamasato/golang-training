package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StructTag(`json:"oid,omitempty"`),
		field.String("name"),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return nil
}
