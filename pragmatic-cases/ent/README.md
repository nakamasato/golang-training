# Ent

## Contents

1. [Getting Started](getting-started)
1. [Simple Example](simple-example)

## Tips

1. Postgres [GIN index](https://www.postgresql.org/docs/current/gin.html) ([ref](https://github.com/ent/ent/blob/a85a22931191a4964a9a14a9083f025c1c74d37b/doc/md/schema-indexes.md))

    ```go
    index.Fields("c5").
          Annotations(
              entsql.IndexTypes(map[string]string{
                  dialect.MySQL:    "FULLTEXT",
                  dialect.Postgres: "GIN",
              }),
          ),
    ```
1. Upsert add `--feature sql/upsert` to ent/generate.go (ref: [#4272](https://github.com/ent/ent/issues/4272), [doc](https://entgo.io/docs/feature-flags#upsert))
1. Postgres `JSONB`

    ```go
    field.JSON("metadata", map[string]interface{}{}).
        SchemaType(map[string]string{
            "postgres": "jsonb", // Specify PostgreSQL type
        }).Default(map[string]interface{}{}),
    ```

## Ref

1. [Mixin](https://entgo.io/docs/schema-mixin) : A Mixin allows you to create reusable pieces of ent.Schema code that can be injected into other schemas using composition.
