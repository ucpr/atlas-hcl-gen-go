package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"ariga.io/atlas/sql/schema"
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
	var hclPath, outPath, target, tag, pkg string
	var showVersion bool
	flag.StringVar(&hclPath, "i", "", "input file path")
	flag.StringVar(&outPath, "o", "", "output file path")
	flag.StringVar(&target, "t", "mysql", "target database")
	flag.StringVar(&tag, "tag", "db", "tag name")
	flag.StringVar(&pkg, "package", "main", "package name")
	flag.BoolVar(&showVersion, "version", false, "print version information and exit")
	flag.Parse()

	if showVersion {
		v := BuildVersion
		if v == "" {
			v = "unknown"
		}
		r := BuildRevision
		if r == "" {
			r = "unknown"
		}
		ts := BuildTimestamp
		if ts == "" {
			ts = "unknown"
		}
		fmt.Printf("atlas-hcl-gen-go: \n\tversion: %s\n\trevision: %s\n\tbuilt: %s\n", v, r, ts)
		return nil
	}

	b, err := os.ReadFile(hclPath)
	if err != nil {
		return err
	}

	// parse hcl schema
	ev, err := toSchemaEvaluatorFunc(strings.ToLower(target))
	if err != nil {
		return err
	}
	var s schema.Schema
	if err := hclBytesFunc(ev)(b, &s, nil); err != nil {
		return err
	}

	// generate go code
	pb, err := generate(s, input{
		hclPath: hclPath,
		outPath: outPath,
		pkg:     pkg,
		tag:     tag,
	})
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// save to file
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
