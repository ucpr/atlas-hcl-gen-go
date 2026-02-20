## basic-postgres

Minimal Postgres schema with default mappings.

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

Expected (snippet):

```go
type Users struct {
    Id        int       `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
}
```

