## blog-site

More realistic blog schema demonstrating multiple tables, enums, JSON, UUIDs, timestamps, and nullable columns.

What it shows
- Named enum type + consts for post status (`enum: named`)
- `json.RawMessage` mapping for JSONB metadata
- Smart null handling (nullable timestamps -> `*time.Time`, binary `[]byte` kept non-pointer)
- UUID as `string`

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```


