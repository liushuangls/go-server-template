package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/liushuangls/go-server-template/pkg/entutil"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entutil.TimeWithDelete{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("nickname").Default(""),
		field.Enum("register_type").Values("oauth2", "email"),
		field.String("register_ip").NotEmpty(),
		field.String("register_region").NotEmpty(),
		field.String("email").Default(""),
		field.Bool("email_verified").Default(false),
		field.String("password").Default(""),
		field.String("avatar").Default(""),
		field.String("profile").Default(""),
		field.Time("bind_at").Nillable().Optional().SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("oauth", UserOAuth.Type),
	}
}

func (u User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
	}
}
