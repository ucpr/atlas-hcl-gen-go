package main

import (
	"fmt"
	"os"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	b, err := os.ReadFile("testdata/schema.hcl")
	if err != nil {
		panic(err)
	}

	var s schema.Schema
	ev := mysql.EvalHCL
	if err := HCLBytesFunc(ev)(b, &s, nil); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", s)

	for i := range s.Tables {
		table := s.Tables[i]
		fmt.Printf("%+v\n", table)
		for j := range table.Columns {
			fmt.Printf("  %+v\n", table.Columns[j])
		}
	}
}

// Reference: https://github.com/ariga/atlas/blob/98bb7b9da852536523121754d19570c506ba69f7/sql/internal/specutil/spec.go#L165
func HCLBytesFunc(ev schemahcl.Evaluator) func(b []byte, v any, inp map[string]cty.Value) error {
	return func(b []byte, v any, inp map[string]cty.Value) error {
		parser := hclparse.NewParser()
		if _, diag := parser.ParseHCL(b, ""); diag.HasErrors() {
			return diag
		}
		return ev.Eval(parser, v, inp)
	}
}
