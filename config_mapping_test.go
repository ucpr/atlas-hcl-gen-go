package main

import (
	"testing"

	"ariga.io/atlas/sql/schema"
	"github.com/stretchr/testify/assert"
)

func Test_NullPoliciesAndConfigMappings(t *testing.T) {
	t.Parallel()

	mk := func(ct schema.Type, null bool) *schema.Column {
		return &schema.Column{Type: &schema.ColumnType{Type: ct, Null: null}}
	}

	// sqlnull: nullable string -> sql.NullString; non-null remains string
	{
		c := Config{Null: "sqlnull"}
		colN := mk(&schema.StringType{T: "text"}, true)
		colNN := mk(&schema.StringType{T: "text"}, false)
		gotN, err := goTypeForColumn(colN, c, "postgres", "t")
		assert.NoError(t, err)
		gotNN, err := goTypeForColumn(colNN, c, "postgres", "t")
		assert.NoError(t, err)
		assert.Equal(t, "sql.NullString", gotN)
		assert.Equal(t, "string", gotNN)
	}

	// pointer: always pointer for nullable (except slices)
	{
		c := Config{Null: "pointer"}
		col := mk(&schema.IntegerType{T: "int"}, true)
		got, err := goTypeForColumn(col, c, "postgres", "t")
		assert.NoError(t, err)
		assert.Equal(t, "*int", got)
		// bytes slice remains []byte
		bs := mk(&schema.BinaryType{T: "bytea"}, true)
		gotB, err := goTypeForColumn(bs, c, "postgres", "t")
		assert.NoError(t, err)
		assert.Equal(t, "[]byte", gotB)
	}

	// decimal: big.Rat
	{
		c := Config{Decimal: "big.Rat"}
		b, ok := baseGoType(&schema.ColumnType{Type: &schema.DecimalType{T: "numeric"}}, c, "postgres")
		assert.True(t, ok)
		assert.Equal(t, "big.Rat", b)
		// nullable with smart -> *big.Rat
		col := mk(&schema.DecimalType{T: "numeric"}, true)
		tp, err := goTypeForColumn(col, c, "postgres", "t")
		assert.NoError(t, err)
		assert.Equal(t, "*big.Rat", tp)
	}

	// json: bytes and string
	{
		cBytes := Config{JSON: "bytes"}
		b, ok := baseGoType(&schema.ColumnType{Type: &schema.JSONType{T: "json"}}, cBytes, "postgres")
		assert.True(t, ok)
		assert.Equal(t, "[]byte", b)
		cStr := Config{JSON: "string"}
		b2, ok := baseGoType(&schema.ColumnType{Type: &schema.JSONType{T: "jsonb"}}, cStr, "postgres")
		assert.True(t, ok)
		assert.Equal(t, "string", b2)
	}

	// uuid: bytes16
	{
		c := Config{UUID: "bytes16"}
		b, ok := baseGoType(&schema.ColumnType{Type: &schema.UUIDType{T: "uuid"}}, c, "postgres")
		assert.True(t, ok)
		assert.Equal(t, "[16]byte", b)
		// nullable with smart -> *[16]byte
		col := mk(&schema.UUIDType{T: "uuid"}, true)
		tp, err := goTypeForColumn(col, c, "postgres", "t")
		assert.NoError(t, err)
		assert.Equal(t, "*[16]byte", tp)
	}
}

func Test_StrictTypesAndMySQLTinyintAsBool(t *testing.T) {
	t.Parallel()

	mk := func(ct schema.Type, null bool) *schema.Column {
		return &schema.Column{Type: &schema.ColumnType{Type: ct, Null: null}}
	}

	// strict_types: true errors on unsupported
	{
		c := Config{StrictTypes: true}
		col := mk(&schema.UnsupportedType{T: "weird"}, false)
		_, err := goTypeForColumn(col, c, "postgres", "t")
		assert.Error(t, err)
	}

	// mysql tinyint(1) as bool (approximate by T name)
	{
		c := Config{MySQLTinyint1AsBool: true}
		col := mk(&schema.IntegerType{T: "tinyint"}, false)
		tp, err := goTypeForColumn(col, c, "mysql", "t")
		assert.NoError(t, err)
		assert.Equal(t, "bool", tp)
	}
}
