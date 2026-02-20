package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	var configPath string
	var showVersion bool
	var initConfig bool
	flag.StringVar(&hclPath, "i", "", "input file path")
	flag.StringVar(&outPath, "o", "", "output file path")
	flag.StringVar(&target, "t", "mysql", "target database")
	flag.StringVar(&tag, "tag", "db", "tag name")
	flag.StringVar(&pkg, "package", "main", "package name")
	flag.StringVar(&configPath, "config", "", "config file path (YAML)")
	flag.BoolVar(&showVersion, "version", false, "print version information and exit")
	flag.BoolVar(&initConfig, "init", false, "create an example config file and exit")
	flag.Parse()

	// Track which flags were explicitly set to apply CLI > config precedence.
	setFlags := map[string]bool{}
	flag.CommandLine.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })

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

	if initConfig {
		// Choose output path: explicit --config path or default file in CWD.
		out := configPath
		if out == "" {
			out = "atlas-hcl-gen-go.yaml"
		}
		if _, err := os.Stat(out); err == nil {
			return fmt.Errorf("config already exists: %s", out)
		}
		if dir := filepath.Dir(out); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create config dir: %w", err)
			}
		}
		if err := os.WriteFile(out, defaultConfigYAML, 0o644); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
		fmt.Printf("wrote config: %s\n", out)
		return nil
	}

	// Load config if provided or if default file exists.
	var conf Config
	if configPath == "" {
		if _, err := os.Stat("atlas-hcl-gen-go.yaml"); err == nil {
			configPath = "atlas-hcl-gen-go.yaml"
		}
	}
	if configPath != "" {
		c, err := loadConfig(configPath)
		if err != nil {
			return err
		}
		conf = c
	}

	// Merge precedence: CLI > config > defaults.
	// Apply config values only if the corresponding CLI flag was not explicitly set.
	if !setFlags["t"] && conf.Dialect != "" {
		target = conf.Dialect
	}
	if !setFlags["package"] && conf.Package != "" {
		pkg = conf.Package
	}
	if !setFlags["tag"] && conf.Tag != "" {
		tag = conf.Tag
	}

	b, err := os.ReadFile(hclPath)
	if err != nil {
		return err
	}

	// parse hcl schema
	// Support both legacy 'target' and new 'dialect' naming (via config).
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
	// Ensure output directory exists if a path is given.
	if dir := filepath.Dir(outPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create output dir: %w", err)
		}
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
