schema "app" {}

table "places" {
  schema = schema.app
  # Using a spatial/geography type as an example of something that may not be mapped.
  column "loc" { type = geography }
}

