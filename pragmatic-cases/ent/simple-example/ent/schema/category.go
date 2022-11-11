package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
	return []ent.Edge{
		// Create an inverse-edge called "items" of type `Items`
		// and reference it to the "categories" edge (in Item schema)
		// explicitly using the `Ref` method.
		edge.From("items", Item.Type).
			Ref("categories"),
		//   Unique(), Not add Unique to make it M2M. Otherwise, O2M
	}
}
