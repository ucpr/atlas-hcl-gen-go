package main

import (
    "path/filepath"
    "testing"

    "ariga.io/atlas/sql/schema"
    "github.com/stretchr/testify/assert"
)

func Test_generatePerTable_Simple(t *testing.T) {
    t.Parallel()

    sc := schema.Schema{
        Tables: []*schema.Table{
            schema.NewTable("users").AddColumns(
                schema.NewStringColumn("id", "string"),
                schema.NewIntColumn("age", "int"),
            ),
            schema.NewTable("posts").AddColumns(
                schema.NewIntColumn("id", "int"),
                schema.NewStringColumn("title", "string"),
            ),
        },
    }

    in := input{
        hclPath: "input.hcl",
        outPath: "out",
        pkg:     "model",
        tag:     "db",
        dialect: "postgres",
        conf:    Config{Null: "smart", Decimal: "string", JSON: "raw", UUID: "string"},
    }

    outDir := "outdir"
    m, err := generatePerTable(sc, in, outDir)
    assert.NoError(t, err)

    // Expect 2 files
    assert.Len(t, m, 2)
    _, okUsers := m[filepath.Join(outDir, "users.go")]
    _, okPosts := m[filepath.Join(outDir, "posts.go")]
    assert.True(t, okUsers)
    assert.True(t, okPosts)

    // Sanity: content should contain respective struct names
    assert.Contains(t, string(m[filepath.Join(outDir, "users.go")]), "type Users struct")
    assert.Contains(t, string(m[filepath.Join(outDir, "posts.go")]), "type Posts struct")
}

