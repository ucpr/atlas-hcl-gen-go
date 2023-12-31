## atlas-hcl-gen-go

This is a command to generate Go struct from [Atlas HCL Schema](https://atlasgo.io/atlas-schema/hcl).

## Usage

Install the command.

```sh
go install github.com/ucpr/atlas-hcl-gen-go
```

Generate Go struct from Atlas HCL Schema.

```sh
atlas-hcl-gen-go -f schema.hcl -o output.go
```

## Example

The input data uses the schema written in HCL below.

```hcl:schema.hcl
schema "market" {}

table "users" {
  schema = schema.market
  column "name" {
    type = text
  }
  column "updated_at" {
    type = int
  }
  column "created_at" {
    type = int
  }
}
```

Execute the command.

```sh
atlas-hcl-gen-go -i schema.hcl -o output.go
```

The following Go struct will be generated.

```go:output.go
// Code generated by github.com/ucpr/atlas-hcl-gen-go. DO NOT EDIT.
// atlas-hcl-gen-go: 922707f-dirty
// source: testdata/schema.hcl

package main

type Users struct {
	Name      int `db:"name"`
	UpdatedAt int `db:"updated_at"`
	CreatedAt int `db:"created_at"`
}
```


## TODO

- [ ] Support some types
- [x] Support some RDBMS schemas (MySQL, PostgreSQL, SQLite, ...)

## Contributing

Contributions of any kind welcome!

## License

[MIT License](LICENSE)
