package main

import (
	"testing"

	"ariga.io/atlas/sql/schema"
	"github.com/stretchr/testify/assert"
)

func Test_baseGoType_and_goTypeForColumn(t *testing.T) {
	t.Parallel()

	mk := func(name string, ct schema.Type, null bool) struct {
		name string
		col  *schema.Column
	} {
		return struct {
			name string
			col  *schema.Column
		}{name: name, col: &schema.Column{Type: &schema.ColumnType{Type: ct, Null: null}}}
	}

	patterns := []struct {
		name string
		col  *schema.Column
		base string
		full string
	}{
		// integers
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("tinyint", &schema.IntegerType{T: "tinyint"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"tinyint", c, "int8", "int8"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("smallint", &schema.IntegerType{T: "smallint"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"smallint", c, "int16", "int16"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("mediumint", &schema.IntegerType{T: "mediumint"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"mediumint", c, "int32", "int32"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("int", &schema.IntegerType{T: "int"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"int", c, "int", "int"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("bigint", &schema.IntegerType{T: "bigint"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"bigint", c, "int64", "int64"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("uint", &schema.IntegerType{T: "int", Unsigned: true}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"uint", c, "uint", "uint"}
		}(),
		// floats
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("float32", &schema.FloatType{T: "float", Precision: 24}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"float32", c, "float32", "float32"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("float64", &schema.FloatType{T: "double", Precision: 53}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"float64", c, "float64", "float64"}
		}(),
		// decimal
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("decimal", &schema.DecimalType{T: "decimal", Precision: 10, Scale: 2}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"decimal", c, "string", "string"}
		}(),
		// string/binary/bool
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("varchar", &schema.StringType{T: "varchar", Size: 255}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"varchar", c, "string", "string"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("binary", &schema.BinaryType{T: "bytea"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"binary", c, "[]byte", "[]byte"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("bool", &schema.BoolType{T: "boolean"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"bool", c, "bool", "bool"}
		}(),
		// time
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("timestamp", &schema.TimeType{T: "timestamp"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"timestamp", c, "time.Time", "time.Time"}
		}(),
		// enum
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("enum", &schema.EnumType{T: "enum", Values: []string{"a", "b"}}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"enum", c, "string", "string"}
		}(),
		// json
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("json", &schema.JSONType{T: "json"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"json", c, "json.RawMessage", "json.RawMessage"}
		}(),
		// uuid
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("uuid", &schema.UUIDType{T: "uuid"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"uuid", c, "string", "string"}
		}(),
		// spatial (fallback)
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("spatial", &schema.SpatialType{T: "point"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"spatial", c, "string", "string"}
		}(),
		// unsupported
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("unsupported", &schema.UnsupportedType{T: "whatever"}, false).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"unsupported", c, "any", "any"}
		}(),

		// nullability smart: value types pointerized, slices not
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("nullable-int64", &schema.IntegerType{T: "bigint"}, true).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"nullable-int64", c, "int64", "*int64"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("nullable-bytes", &schema.BinaryType{T: "blob"}, true).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"nullable-bytes", c, "[]byte", "[]byte"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("nullable-json", &schema.JSONType{T: "json"}, true).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"nullable-json", c, "json.RawMessage", "json.RawMessage"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("nullable-time", &schema.TimeType{T: "timestamp"}, true).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"nullable-time", c, "time.Time", "*time.Time"}
		}(),
		func() struct {
			name       string
			col        *schema.Column
			base, full string
		} {
			c := mk("nullable-string", &schema.StringType{T: "text"}, true).col
			return struct {
				name       string
				col        *schema.Column
				base, full string
			}{"nullable-string", c, "string", "*string"}
		}(),
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotBase := baseGoType(tt.col.Type)
			gotFull := goTypeForColumn(tt.col)
			assert.Equal(t, tt.base, gotBase)
			assert.Equal(t, tt.full, gotFull)
		})
	}
}
