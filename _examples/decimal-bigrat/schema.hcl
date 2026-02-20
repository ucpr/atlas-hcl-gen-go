schema "app" {}

table "orders" {
  schema = schema.app
  column "amount" {
    type = decimal(20, 6)
    null = true
  }
}
