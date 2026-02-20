## strict-types

Enable `strict_types: true` to fail fast on unsupported/ambiguous database types.

Run:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
# Expect a non-zero exit with an error like:
#   unsupported or ambiguous type: geography
```

