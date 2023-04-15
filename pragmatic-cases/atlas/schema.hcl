schema "test_db" {}

table "users" {
  schema = schema.test_db
  column "id" {
    type = int
  }
  column "greeting" {
    type = text
  }
}
