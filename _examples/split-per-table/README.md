## split-per-table

Generate one `.go` file per table.

Run:

```bash
mkdir -p out
atlas-hcl-gen-go -i schema.hcl -o out --config atlas-hcl-gen-go.yaml
# writes: out/users.go, out/posts.go
```

