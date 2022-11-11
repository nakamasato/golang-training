package schema

import "entgo.io/ent"

// ItemCategory holds the schema definition for the ItemCategory entity.
type ItemCategory struct {
	ent.Schema
}

// Fields of the ItemCategory.
func (ItemCategory) Fields() []ent.Field {
	return nil
}

// Edges of the ItemCategory.
func (ItemCategory) Edges() []ent.Edge {
	return nil
}
