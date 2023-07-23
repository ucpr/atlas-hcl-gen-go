package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/tools/imports"
)

var (
	// inject by ldflags
	BuildVersion   = ""
	BuildRevision  = ""
	BuildTimestamp = ""
)

func main() {
	log.SetPrefix("atlas-hcl-gen-go: ")
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var hclPath, outPath string
	flag.StringVar(&hclPath, "f", "", "input file path")
	flag.StringVar(&outPath, "o", "", "output file path")
	flag.Parse()

	b, err := os.ReadFile(hclPath)
	if err != nil {
		return err
	}

	var s schema.Schema
	ev := mysql.EvalHCL
	if err := hclBytesFunc(ev)(b, &s, nil); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "// Code generated by github.com/ucpr/atlas-hcl-gen-go. DO NOT EDIT.\n")
	fmt.Fprintf(buf, "// atlas-hcl-gen-go: %s\n", BuildVersion)
	fmt.Fprintf(buf, "// source: %s\n\n", hclPath)
	fmt.Fprintf(buf, "package main\n\n")

	for i := range s.Tables {
		table := s.Tables[i]
		fmt.Fprintf(buf, "type %s struct {\n", toCamelCase(table.Name))
		for j := range table.Columns {
			column := table.Columns[j]
			tp := toGoTypeString(column.Type)
			fmt.Fprintf(buf, "\t%s\t%s\t`json:\"%s\"`\n", toCamelCase(column.Name), tp, column.Name)
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	pb, err := imports.Process(outPath, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("failed to format: %w", err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(pb); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func toGoTypeString(ct *schema.ColumnType) string {
	// TODO: support more types
	switch ct.Type.(type) {
	case *schema.IntegerType:
		return "int"
	case *schema.FloatType:
		return "float"
	case *schema.StringType:
		return "string"
	}
	return "any"
}

func toCamelCase(s string) string {
	var result strings.Builder
	upperNext := true

	for _, r := range s {
		if !unicode.IsLetter(r) {
			upperNext = true
			continue
		}
		if upperNext {
			result.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			result.WriteRune(r)
		}
	}

	// handle the first character of lower camel case
	if len(s) > 0 && unicode.IsLower(rune(s[0])) {
		resultString := result.String()
		return string(unicode.ToUpper(rune(resultString[0]))) + resultString[1:]
	}

	return result.String()
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
