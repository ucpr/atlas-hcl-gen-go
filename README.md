## ☄️ atlas-hcl-gen-go

Generate Go structs from [Atlas HCL Schema](https://atlasgo.io/atlas-schema/hcl), with configurable type mapping, nullability policies, and optional enum type generation.

### Features
- Config-first: drive behavior from YAML (dialect, package, tags, policies)
- Rich type mapping for MySQL / PostgreSQL / SQLite
- Null handling strategies: smart | pointer | sqlnull
- Decimal, JSON, UUID mapping options (stdlib-first)
- Optional named enum types + const values (`enum: named`)
- MySQL `TINYINT(1) → bool` option
- Strict type checking (`strict_types`)

## Install

```bash
go install github.com/ucpr/atlas-hcl-gen-go@latest
```

## Quick Start

1) Initialize config (documents fields and accepted values):

```bash
atlas-hcl-gen-go -init    # writes ./atlas-hcl-gen-go.yaml
```

2) Create a minimal HCL schema (e.g., `schema.hcl`):

```hcl
schema "app" {}

table "users" {
  schema = schema.app
  column "id"        { type = int }
  column "name"      { type = text }
  column "created_at"{ type = timestamp }
}
```

3) Generate code:

```bash
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

Example output (snippet):

```go
type Users struct {
    Id        int       `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
}
```

## Init Config

Initialize a config file with documented fields and accepted values.

```bash
# Writes ./atlas-hcl-gen-go.yaml (does not overwrite if exists)
atlas-hcl-gen-go -init

# Or specify a custom path
atlas-hcl-gen-go -init --config configs/atlas-hcl-gen-go.yaml
```

Use a config when generating.

```bash
atlas-hcl-gen-go -i schema.hcl -o output.go --config atlas-hcl-gen-go.yaml
```

## Config Options

- dialect: target dialect (`mysql`|`postgres`|`sqlite`)
- package: Go package name for generated file
- tag: struct tag key (e.g., `db`)
- null_policy: null strategy (`smart`|`pointer`|`sqlnull`)
- decimal: DECIMAL/NUMERIC mapping (`string`|`big.Rat`)
- uuid: UUID mapping (`string`|`bytes16`)
- json: JSON/JSONB mapping (`raw`|`bytes`|`string`)
- strict_types: error on unsupported/ambiguous types (bool)
- mysql_tinyint1_as_bool: map MySQL `TINYINT` to `bool` (bool)
- enum: ENUM mapping (`string`|`named`)

### ENUM: named type + consts

- When `enum: named`, the generator emits a named Go type per enum column and const values.
- Type name: `<Table><Column>` (CamelCase). Example: `UsersStatus`.
- Const names: `<Type><ValueCamel>` with string literal assignments.
- Struct field type becomes the named enum type; null policy applies (e.g., `*UsersStatus`).

## Examples

See `_examples/` for runnable samples:
- basic-postgres: minimal Postgres schema and defaults
- enum-named: enum column with `enum: named`, emits type + consts
- mysql-tinyint-bool: map `tinyint(1)` to `bool`
- null-sqlnull: use `sql.Null*` for nullable basics
- decimal-bigrat: map DECIMAL/NUMERIC to `big.Rat`
- strict-types: fail on unsupported or ambiguous types
 - blog-site: realistic blog schema (users, posts, comments, tags) with enums, JSON, UUID, and nullable timestamps

Run an example:

```bash
cd _examples/basic-postgres
atlas-hcl-gen-go -i schema.hcl -o model.go --config atlas-hcl-gen-go.yaml
```

## Version

```bash
atlas-hcl-gen-go -version
```

## Contributing

Issues and PRs welcome.

## License

[MIT License](LICENSE)
