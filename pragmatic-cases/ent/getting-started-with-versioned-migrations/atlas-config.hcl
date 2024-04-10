variable "skip_drop" {
  type    = bool
  default = true
  description = "Skip drop operations (schema, table, column, index, foreign key)"
}

env "local" {
  migration {
    dir = "file://ent/migrate/migrations"
    revisions_schema = "public"
  }
  src = "ent://ent/schema"
  url = "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  diff {
    skip {
      drop_schema = var.skip_drop
      drop_table  = var.skip_drop
      drop_column = var.skip_drop
      drop_index  = var.skip_drop
      drop_foreign_key = var.skip_drop
    }
  }
}
