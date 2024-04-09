package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Animal holds the schema definition for the Animal entity.
type Animal struct {
	ent.Schema
}

// Fields of the Animal.
func (Animal) Fields() []ent.Field {
	return []ent.Field{
		field.String("species"),
		field.Int("age"),
		field.String("name"),
	}
}

// Edges of the Animal.
func (Animal) Edges() []ent.Edge {
	return nil
}
