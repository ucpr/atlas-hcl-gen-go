## decimal-bigrat

Map DECIMAL/NUMERIC to `big.Rat` for arithmetic-friendly, exact rational representation.

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

Expected (snippet):

```go
type Orders struct {
    Amount *big.Rat `db:"amount"` // nullable DECIMAL uses pointer under smart null
}
```

