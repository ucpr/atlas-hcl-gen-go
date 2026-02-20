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
  
  primary_key {
    columns = [column.id]
  }

  unique "users_email_key" {
    columns = [column.email]
  }

  index "users_created_at_idx" {
    columns = [column.created_at]
  }
}

table "posts" {
  schema = schema.blog
  column "id"           { type = int }
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
  
  primary_key {
    columns = [column.id]
  }

  unique "posts_slug_key" {
    columns = [column.slug]
  }

  index "posts_author_id_idx" {
    columns = [column.author_id]
  }

  index "posts_published_at_idx" {
    columns = [column.published_at]
  }

  foreign_key "posts_author_id_fkey" {
    columns     = [column.author_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

table "comments" {
  schema = schema.blog
  column "id"         { type = int }
  column "post_id"    { type = int }
  column "author_id" {
    type = uuid
    null = true
  }
  column "body"       { type = text }
  column "created_at" { type = timestamptz }
  
  primary_key {
    columns = [column.id]
  }

  index "comments_post_id_idx" {
    columns = [column.post_id]
  }

  foreign_key "comments_post_id_fkey" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }

  foreign_key "comments_author_id_fkey" {
    columns     = [column.author_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
}

table "tags" {
  schema = schema.blog
  column "id"   { type = int }
  column "name" { type = text }
  
  primary_key {
    columns = [column.id]
  }

  unique "tags_name_key" {
    columns = [column.name]
  }
}

table "post_tags" {
  schema = schema.blog
  column "post_id" { type = int }
  column "tag_id"  { type = int }
  
  index "post_tags_post_id_idx" {
    columns = [column.post_id]
  }

  index "post_tags_tag_id_idx" {
    columns = [column.tag_id]
  }

  unique "post_tags_post_id_tag_id_key" {
    columns = [column.post_id, column.tag_id]
  }

  foreign_key "post_tags_post_id_fkey" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }

  foreign_key "post_tags_tag_id_fkey" {
    columns     = [column.tag_id]
    ref_columns = [table.tags.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
