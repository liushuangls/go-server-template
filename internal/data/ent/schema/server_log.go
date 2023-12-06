package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/liushuangls/go-server-template/pkg/entutil"
)

type ServerLog struct {
	ent.Schema
}

func (ServerLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entutil.CreateTime{},
	}
}

func (ServerLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").Default(0),
		field.String("ip").Default(""),
		field.String("method").Default(""),
		field.String("path").Default(""),
		field.Text("query").Default(""),
		field.Text("body").Default(""),
		field.Enum("level").Values("ERROR", "WARN", "INFO", "DEBUG", "PANIC").Default("ERROR"),
		field.Enum("from").Values("api", "log").Default("api"),
		field.Text("err_msg"),
		field.Text("resp_err_msg").Default(""),
		field.Int("code").Optional().Default(0),
		field.JSON("extra", map[string]any{}),
	}
}

func (ServerLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (ServerLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("level"),
		index.Fields("path"),
		index.Fields("code"),
	}
}
