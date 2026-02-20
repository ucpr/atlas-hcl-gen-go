schema "blog" {}

table "users" {
  schema = schema.blog
  column "id"           { type = uuid }
  column "email"        { type = text }
  column "display_name" { type = text }
  column "bio" {
    type = text
    null = true
  }
  column "is_admin" {
    type = boolean
    default = false
  }
  column "created_at"   { type = timestamptz }
  column "updated_at"   { type = timestamptz }
}

table "posts" {
  schema = schema.blog
  column "id"           { type = serial }
  column "title"        { type = text }
  column "slug"         { type = text }
  column "status"       { type = text }
  column "author_id"    { type = uuid }
  column "published_at" {
    type = timestamptz
    null = true
  }
  column "metadata" {
    type = jsonb
    null = true
  }
  column "rating" {
    type = numeric(3,1)
    null = true
  }
  column "cover_image" {
    type = bytea
    null = true
  }
  column "created_at"   { type = timestamptz }
  column "updated_at"   { type = timestamptz }
}

table "comments" {
  schema = schema.blog
  column "id"         { type = serial }
  column "post_id"    { type = int }
  column "author_id" {
    type = uuid
    null = true
  }
  column "body"       { type = text }
  column "created_at" { type = timestamptz }
}

table "tags" {
  schema = schema.blog
  column "id"   { type = serial }
  column "name" { type = text }
}

table "post_tags" {
  schema = schema.blog
  column "post_id" { type = int }
  column "tag_id"  { type = int }
}
