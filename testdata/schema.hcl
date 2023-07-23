schema "market" {}

table "users" {
  schema = schema.market
  column "name" {
    type = int
  }
  column "updatedAt" {
    type = int
  }
  column "createdAt" {
    type = int
  }
}

table "tokens" {
  schema = schema.market
  column "value" {
    type = int
  }
  column "updatedAt" {
    type = int
  }
  column "createdAt" {
    type = int
  }
}
