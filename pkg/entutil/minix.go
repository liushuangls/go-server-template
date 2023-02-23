package entutil

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type CreateTime struct{ ent.Schema }

// Fields of the time mixin.
func (CreateTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
	}
}

func (CreateTime) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("create_time"),
	}
}

type Time struct{ ent.Schema }

// Fields of the time mixin.
func (Time) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
	}
}

func (Time) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("create_time"),
	}
}

type TimeWithDelete struct{ ent.Schema }

// Fields of the time mixin.
func (TimeWithDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
		field.Time("delete_time").
			Nillable().
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
	}
}

func (TimeWithDelete) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("create_time"),
		index.Fields("delete_time"),
	}
}
