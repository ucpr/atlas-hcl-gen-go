package main

import (
	"fmt"
	"strings"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

func toSchemaEvaluatorFunc(target string) (schemahcl.EvalFunc, error) {
	switch target {
	case "mysql":
		return mysql.EvalHCL, nil
	case "postgres", "postgresql":
		return postgres.EvalHCL, nil
	case "sqlite":
		return sqlite.EvalHCL, nil
	}
	return nil, fmt.Errorf("unsupported target database: %s", target)
}

// goTypeForColumn returns the Go type string for a given column,
// applying a smart-null policy: nullable columns are pointers (except slices).
func goTypeForColumn(c *schema.Column) string {
	base := baseGoType(c.Type)
	if c.Type != nil && c.Type.Null {
		// Do not pointer-ize slices or json.RawMessage (alias of []byte).
		if strings.HasPrefix(base, "[]") || base == "json.RawMessage" {
			return base
		}
		return "*" + base
	}
	return base
}

// baseGoType resolves a Go type for a column type without applying nullability wrapping.
func baseGoType(ct *schema.ColumnType) string {
	if ct == nil || ct.Type == nil {
		return "any"
	}
	switch t := ct.Type.(type) {
	case *schema.IntegerType:
		// Map by declared type name, and respect unsigned.
		// Fallback to platform-int if unknown.
		name := strings.ToLower(t.T)
		var v string
		switch name {
		case "tinyint":
			v = "int8"
		case "smallint":
			v = "int16"
		case "mediumint":
			v = "int32"
		case "bigint", "bigserial":
			v = "int64"
		case "serial", "integer", "int":
			v = "int"
		default:
			v = "int"
		}
		if t.Unsigned {
			if v == "int" {
				return "uint"
			}
			return "u" + v
		}
		return v
	case *schema.FloatType:
		// Prefer float64 by default; use float32 for low precision hints.
		if t.Precision > 0 && t.Precision <= 32 {
			return "float32"
		}
		return "float64"
	case *schema.DecimalType:
		// Lossless textual form by default.
		return "string"
	case *schema.StringType:
		return "string"
	case *schema.BoolType:
		return "bool"
	case *schema.TimeType:
		return "time.Time"
	case *schema.EnumType:
		return "string"
	case *schema.BinaryType:
		return "[]byte"
	case *schema.JSONType:
		// Default to json.RawMessage (alias of []byte) for flexibility.
		return "json.RawMessage"
	case *schema.UUIDType:
		return "string"
	case *schema.SpatialType:
		// Default to string representation.
		return "string"
	case *schema.UnsupportedType:
		return "any"
	default:
		return "any"
	}
}

// original code: https://github.com/ariga/atlas/blob/98bb7b9da852536523121754d19570c506ba69f7/sql/internal/specutil/spec.go#L165...L173
func hclBytesFunc(ev schemahcl.Evaluator) func(b []byte, v any, inp map[string]cty.Value) error {
	return func(b []byte, v any, inp map[string]cty.Value) error {
		parser := hclparse.NewParser()
		if _, diag := parser.ParseHCL(b, ""); diag.HasErrors() {
			return diag
		}
		return ev.Eval(parser, v, inp)
	}
}
