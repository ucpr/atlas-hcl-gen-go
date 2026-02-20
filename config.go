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
    // Use NullPolicy as canonical key to avoid YAML 'null' literal ambiguity.
    NullPolicy string `yaml:"null_policy"` // smart|pointer|sqlnull
    // Deprecated: 'null' key in YAML. Parsed manually for compatibility.
    Null                string `yaml:"-"`
    Decimal             string `yaml:"decimal"` // string|big.Rat
    UUID                string `yaml:"uuid"`    // string|bytes16
    JSON                string `yaml:"json"`    // raw|bytes|string
    StrictTypes         bool   `yaml:"strict_types"`
    MySQLTinyint1AsBool bool   `yaml:"mysql_tinyint1_as_bool"`
    Enum                string `yaml:"enum"` // string|named
    // SplitPerTable: when true, generate one Go file per table
    // under the output directory. The CLI -o path is treated as a
    // directory (or its parent directory if a file is provided).
    SplitPerTable bool `yaml:"split_per_table"`
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
	// Manual parse to support legacy 'null' key and avoid YAML ambiguity.
	var raw map[string]any
	if err := yaml.Unmarshal(b, &raw); err == nil {
		if v, ok := raw["null"]; ok {
			if s, ok := v.(string); ok && s != "" {
				c.Null = s
			}
		}
	}
	return c, nil
}
