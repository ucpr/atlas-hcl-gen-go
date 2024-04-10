package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"ariga.io/atlas/sql/schema"
	"golang.org/x/tools/imports"
)

type input struct {
	hclPath string
	outPath string
	pkg     string
	tag     string
}

func generate(s schema.Schema, in input) ([]byte, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "// Code generated by github.com/ucpr/atlas-hcl-gen-go. DO NOT EDIT.\n")
	fmt.Fprintf(buf, "// atlas-hcl-gen-go: %s\n", BuildVersion)
	fmt.Fprintf(buf, "// source: %s\n\n", in.hclPath)
	fmt.Fprintf(buf, "package %s\n\n", in.pkg)

	for i := range s.Tables {
		table := s.Tables[i]
		fmt.Fprintf(buf, "type %s struct {\n", toCamelCase(table.Name))
		for j := range table.Columns {
			column := table.Columns[j]
			tp := toGoTypeString(column.Type)
			fmt.Fprintf(buf, "\t%s\t%s\t`%s:\"%s\"`\n", toCamelCase(column.Name), tp, in.tag, column.Name)
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	pb, err := imports.Process(in.outPath, buf.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to format: %w", err)
	}

	return pb, nil
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
