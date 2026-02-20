schema "app" {}

table "users" {
  schema = schema.app
  column "id"         { type = int }
  column "name"       { type = text }
  column "created_at" { type = timestamp }
}

