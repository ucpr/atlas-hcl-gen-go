package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds generator options loaded from YAML.
type Config struct {
	Dialect string `yaml:"dialect"`
	Package string `yaml:"package"`
	Tag     string `yaml:"tag"`

	// Future options (not fully wired yet but reserved for config shape).
	Null                string `yaml:"null"`    // smart|pointer|sqlnull
	Decimal             string `yaml:"decimal"` // string|big.Rat
	UUID                string `yaml:"uuid"`    // string|bytes16
	JSON                string `yaml:"json"`    // raw|bytes|string
	StrictTypes         bool   `yaml:"strict_types"`
	MysqlTinyint1AsBool bool   `yaml:"mysql_tinyint1_as_bool"`
	Enum                string `yaml:"enum"` // string|named
}

func loadConfig(path string) (Config, error) {
	var c Config
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("failed to parse config: %w", err)
	}
	return c, nil
}
