## null-sqlnull

Use `sql.Null*` for nullable basic types (`null: sqlnull`).

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

Expected (snippet):

```go
type Profiles struct {
    Nickname      sql.NullString `db:"nickname"`
    Age          sql.NullInt64   `db:"age"`
    Verified     sql.NullBool    `db:"verified"`
    SignedAt     sql.NullTime    `db:"signed_at"`
}
```

