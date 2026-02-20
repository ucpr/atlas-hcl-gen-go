schema "app" {}

table "profiles" {
  schema = schema.app
  column "nickname" {
    type = text
    null = true
  }
  column "age" {
    type = int
    null = true
  }
  column "verified" {
    type = boolean
    null = true
  }
  column "signed_at" {
    type = timestamp
    null = true
  }
}
