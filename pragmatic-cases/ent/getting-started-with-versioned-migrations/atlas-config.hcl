variable "destructive" {
  type    = bool
  default = false
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
      drop_schema = !var.destructive
      drop_table  = !var.destructive
      drop_column = !var.destructive
      drop_index  = !var.destructive
      drop_foreign_key = !var.destructive
    }
  }
}
