table "users" {
  schema = schema.test_db
  column "id" {
    null = false
    type = int
  }
  column "greeting" {
    null = false
    type = text
  }
}
schema "test_db" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
