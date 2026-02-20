package main

import (
	"testing"

	"ariga.io/atlas/sql/schema"
	"github.com/stretchr/testify/assert"
)

func Test_EnumNamedGeneration(t *testing.T) {
	t.Parallel()

	sc := schema.Schema{
		Tables: []*schema.Table{
			schema.NewTable("users").AddColumns(
				&schema.Column{Name: "status", Type: &schema.ColumnType{Type: &schema.EnumType{T: "enum", Values: []string{"active", "inactive", "1st"}}}},
				&schema.Column{Name: "nickname", Type: &schema.ColumnType{Type: &schema.StringType{T: "text"}}},
			),
		},
	}

	out, err := generate(sc, input{
		hclPath: "schema.hcl",
		outPath: "out.go",
		pkg:     "model",
		tag:     "db",
		dialect: "postgres",
		conf:    Config{Enum: "named", Null: "smart"},
	})
	assert.NoError(t, err)
	code := string(out)
	// Expect type + consts
	assert.Contains(t, code, "type UsersStatus string")
	assert.Contains(t, code, "const (")
	assert.Contains(t, code, "UsersStatusActive   UsersStatus = \"active\"")
	assert.Contains(t, code, "UsersStatusInactive UsersStatus = \"inactive\"")
	// value starting with digit becomes prefixed ident
	assert.Contains(t, code, "UsersStatusN1st")
	assert.Contains(t, code, "= \"1st\"")
	// Struct field uses named type
	assert.Contains(t, code, "UsersStatus `db:\"status\"`")
}
