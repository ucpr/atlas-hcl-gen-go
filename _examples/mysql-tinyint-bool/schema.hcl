schema "app" {}

table "flags" {
  schema = schema.app
  column "enabled" { type = tinyint(1) }
}

