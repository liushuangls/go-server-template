package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/liushuangls/go-server-template/pkg/entutil"
)

var (
	OAuthPlatforms = []string{"google", "microsoft", "apple"}
)

type UserOAuth struct {
	ent.Schema
}

func (UserOAuth) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_oauth"},
	}
}

func (UserOAuth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entutil.TimeWithDelete{},
	}
}

func (UserOAuth) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").Positive(),
		field.Enum("platform").Values(OAuthPlatforms...),
		field.String("open_id").NotEmpty(),
		field.String("union_id").Default(""),
	}
}

func (UserOAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("oauth").Unique().Required().Field("user_id"),
	}
}

func (UserOAuth) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("open_id"),
	}
}
