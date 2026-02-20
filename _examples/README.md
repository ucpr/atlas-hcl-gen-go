## Examples

Each directory contains:
- `schema.hcl`: Atlas HCL schema
- `atlas-hcl-gen-go.yaml`: generator config (dialect, options)
- `README.md`: how to run and expected output snippet

Run any example by:

```bash
cd _examples/<name>
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

