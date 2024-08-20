package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/liushuangls/go-server-template/pkg/entutil"
)

type AppConfig struct {
	ent.Schema
}

func (AppConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entutil.TimeWithDelete{},
	}
}

func (AppConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.Text("value").NotEmpty(),
		field.Enum("value_type").Values("string", "int", "object").Default("string"),
		field.Enum("type").Values("client", "server").Default("client"),
		field.String("app_name").Default(""),
		field.String("app_version_gte").Default(""),
		field.String("app_version_lte").Default(""),
	}
}

func (AppConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (AppConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key", "app_name").Unique(),
	}
}
