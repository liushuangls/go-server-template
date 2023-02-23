// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "update_time", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "delete_time", Type: field.TypeTime, Nullable: true, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "email", Type: field.TypeString},
		{Name: "password", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_create_time",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[1]},
			},
			{
				Name:    "user_delete_time",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[3]},
			},
			{
				Name:    "user_email",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		UsersTable,
	}
)

func init() {
}
