package main

import (
	"fmt"

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

func toGoTypeString(ct *schema.ColumnType) string {
	// TODO: support more types
	// https://atlasgo.io/atlas-schema/hcl-types
	switch ct.Type.(type) {
	case *schema.IntegerType:
		return "int"
	case *schema.FloatType:
		return "float"
	case *schema.StringType:
		return "string"
	case *schema.BoolType:
		return "bool"
	case *schema.TimeType:
		return "time.Time"
	case *schema.EnumType:
		return "string"
	case *schema.BinaryType, *schema.JSONType:
		return "[]byte"
	}
	return "any"
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
