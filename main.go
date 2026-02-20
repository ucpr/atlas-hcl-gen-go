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

type cliArgs struct {
	hclPath    string
	outPath    string
	target     string
	tag        string
	pkg        string
	configPath string
}

func run() error {
	var args cliArgs
	var showVersion bool
	var initConfig bool

	flag.StringVar(&args.hclPath, "i", "", "input file path")
	flag.StringVar(&args.outPath, "o", "", "output file path")
	flag.StringVar(&args.target, "t", "mysql", "target database")
	flag.StringVar(&args.tag, "tag", "db", "tag name")
	flag.StringVar(&args.pkg, "package", "main", "package name")
	flag.StringVar(&args.configPath, "config", "", "config file path (YAML)")
	flag.BoolVar(&showVersion, "version", false, "print version information and exit")
	flag.BoolVar(&initConfig, "init", false, "create an example config file and exit")
	flag.Parse()

	// Which flags were explicitly provided (for precedence decisions).
	setFlags := map[string]bool{}
	flag.CommandLine.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })

	if showVersion {
		return runVersion()
	}
	if initConfig {
		return runInit(args.configPath)
	}
	return runGenerate(args, setFlags)
}

// runVersion prints version metadata.
func runVersion() error {
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

// runInit writes an example config file to the given path (or default name).
func runInit(configPath string) error {
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

// runGenerate evaluates the HCL schema and generates Go code using merged config.
func runGenerate(args cliArgs, setFlags map[string]bool) error {
	// Load config if provided or if default file exists.
	var conf Config
	configPath := args.configPath
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
	target := args.target
	tag := args.tag
	pkg := args.pkg
	if !setFlags["t"] && conf.Dialect != "" {
		target = conf.Dialect
	}
	if !setFlags["package"] && conf.Package != "" {
		pkg = conf.Package
	}
	if !setFlags["tag"] && conf.Tag != "" {
		tag = conf.Tag
	}

	// Read input HCL file.
	b, err := os.ReadFile(args.hclPath)
	if err != nil {
		return err
	}

	// Parse HCL schema using dialect evaluator.
	ev, err := toSchemaEvaluatorFunc(strings.ToLower(target))
	if err != nil {
		return err
	}
	var s schema.Schema
	if err := hclBytesFunc(ev)(b, &s, nil); err != nil {
		return err
	}

	// Generate Go code.
	pb, err := generate(s, input{
		hclPath: args.hclPath,
		outPath: args.outPath,
		pkg:     pkg,
		tag:     tag,
		dialect: strings.ToLower(target),
		conf:    conf,
	})
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// Save to file (ensure directory exists).
	if dir := filepath.Dir(args.outPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create output dir: %w", err)
		}
	}
	f, err := os.Create(args.outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(pb); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
