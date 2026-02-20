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

// goTypeForColumn returns the Go type string for a given column
// using the provided config and dialect. It applies nullability policy
// and may error if strict_types is enabled and the type is unsupported.
func goTypeForColumn(c *schema.Column, conf Config, dialect string, tableName string) (string, error) {
	base, supported := baseGoType(c.Type, conf, strings.ToLower(dialect))
	// If enum and config requires named enums, override base with type name.
	if c.Type != nil && c.Type.Type != nil {
		if _, ok := c.Type.Type.(*schema.EnumType); ok && strings.ToLower(conf.Enum) == "named" {
			base = enumTypeName(tableName, c.Name)
			supported = true
		}
	}
	if conf.StrictTypes && !supported {
		// Build a DB type description if possible.
		dbt := "unknown"
		if c.Type != nil && c.Type.Type != nil {
			switch t := c.Type.Type.(type) {
			case *schema.IntegerType:
				dbt = strings.ToLower(t.T)
			case *schema.FloatType:
				dbt = strings.ToLower(t.T)
			case *schema.DecimalType:
				dbt = strings.ToLower(t.T)
			case *schema.StringType:
				dbt = strings.ToLower(t.T)
			case *schema.BoolType:
				dbt = strings.ToLower(t.T)
			case *schema.TimeType:
				dbt = strings.ToLower(t.T)
			case *schema.EnumType:
				dbt = strings.ToLower(t.T)
			case *schema.BinaryType:
				dbt = strings.ToLower(t.T)
			case *schema.JSONType:
				dbt = strings.ToLower(t.T)
			case *schema.UUIDType:
				dbt = strings.ToLower(t.T)
			case *schema.SpatialType:
				dbt = strings.ToLower(t.T)
			case *schema.UnsupportedType:
				dbt = strings.ToLower(t.T)
			}
		}
		return "", fmt.Errorf("unsupported or ambiguous type: %s", dbt)
	}

	// Apply nullability policy.
	isNullable := c.Type != nil && c.Type.Null
	tp := applyNullPolicy(base, isNullable, conf)
	return tp, nil
}

// baseGoType resolves a Go type for a column type without applying nullability wrapping.
// It returns (typeName, supported).
func baseGoType(ct *schema.ColumnType, conf Config, dialect string) (string, bool) {
	if ct == nil || ct.Type == nil {
		return "any", false
	}
	switch t := ct.Type.(type) {
	case *schema.IntegerType:
		name := strings.ToLower(t.T)
		// MySQL specific: tinyint(1) => bool (approximate by T name)
		if dialect == "mysql" && conf.MysqlTinyint1AsBool && name == "tinyint" {
			return "bool", true
		}
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
				return "uint", true
			}
			return "u" + v, true
		}
		return v, true
	case *schema.FloatType:
		if t.Precision > 0 && t.Precision <= 32 {
			return "float32", true
		}
		return "float64", true
	case *schema.DecimalType:
		switch strings.ToLower(conf.Decimal) {
		case "big.rat":
			return "big.Rat", true
		default:
			return "string", true
		}
	case *schema.StringType:
		return "string", true
	case *schema.BoolType:
		return "bool", true
	case *schema.TimeType:
		return "time.Time", true
	case *schema.EnumType:
		return "string", true
	case *schema.BinaryType:
		return "[]byte", true
	case *schema.JSONType:
		switch strings.ToLower(conf.JSON) {
		case "bytes":
			return "[]byte", true
		case "string":
			return "string", true
		default:
			// raw
			return "json.RawMessage", true
		}
	case *schema.UUIDType:
		switch strings.ToLower(conf.UUID) {
		case "bytes16":
			return "[16]byte", true
		default:
			return "string", true
		}
	case *schema.SpatialType:
		return "string", true
	case *schema.UnsupportedType:
		return "any", false
	default:
		return "any", false
	}
}

// applyNullPolicy wraps base type according to nullability and configuration.
func applyNullPolicy(base string, isNullable bool, conf Config) string {
	if !isNullable {
		return base
	}
	// Slices are kept as-is (e.g., []byte), and json.RawMessage too.
	if strings.HasPrefix(base, "[]") || base == "json.RawMessage" {
		return base
	}
	switch strings.ToLower(conf.Null) {
	case "pointer":
		return "*" + base
	case "sqlnull":
		// Try mapping basic types to sql.Null*
		switch base {
		case "string":
			return "sql.NullString"
		case "bool":
			return "sql.NullBool"
		case "time.Time":
			return "sql.NullTime"
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			return "sql.NullInt64"
		case "float32", "float64":
			return "sql.NullFloat64"
		default:
			// For non-basic types, fall back to pointerization.
			return "*" + base
		}
	default: // smart
		return "*" + base
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
