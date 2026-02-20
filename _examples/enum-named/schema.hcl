schema "app" {}

table "users" {
  schema = schema.app
  column "status" { type = enum("active", "inactive") }
}

