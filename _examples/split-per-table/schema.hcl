schema "app" {}

table "users" {
  schema = schema.app
  column "id"         { type = int }
  column "name"       { type = text }
  column "created_at" { type = timestamp }
}

table "posts" {
  schema = schema.app
  column "id"         { type = int }
  column "user_id"    { type = int }
  column "title"      { type = text }
  column "created_at" { type = timestamp }
}

