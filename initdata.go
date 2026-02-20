package main

import _ "embed"

// defaultConfigYAML is the embedded template for initializing a config file.
//
//go:embed templates/atlas-hcl-gen-go.init.yaml
var defaultConfigYAML []byte
