package schema

import (
	"dndroller/internal/model"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Set holds the schema definition for the Set entity.
type Set struct {
	ent.Schema
}

// Fields of the Set.
func (Set) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("data", new(model.DiceSet)),
	}
}

// Edges of the Set.
func (Set) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("set").
			Unique().
			Required(),
	}
}
