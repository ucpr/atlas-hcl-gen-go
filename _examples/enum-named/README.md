## enum-named

Enum column mapped to a named Go type with constants.

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

Expected (snippet):

```go
type UsersStatus string

const (
    UsersStatusActive   UsersStatus = "active"
    UsersStatusInactive UsersStatus = "inactive"
)

type Users struct {
    Status UsersStatus `db:"status"`
}
```

