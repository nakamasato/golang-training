# [ent](https://entgo.io/)


## [Getting started](https://entgo.io/docs/getting-started/)

### Install

1. Install

    ```
    go install entgo.io/ent/cmd/ent@latest
    go get entgo.io/ent/cmd/ent
    ```

1. Change dir

    ```
    cd pragmatic-cases/ent
    ```

### Create schema

1. Generate the schema for User

    ```
    go run entgo.io/ent/cmd/ent init User
    ```

    ```
    .
    ├── README.md
    └── ent
        ├── generate.go
        └── schema
            └── user.go

    2 directories, 3 files
    ```

1. Add two fields `age` and `name` by changing `Field` function.

    ```go
    func (User) Fields() []ent.Field {
        return []ent.Field{
            field.Int("age").
                Positive(),
            field.String("name").
                Default("unknown"),
        }
    }
    ```

1. Generate go.

    ```
    go generate ./ent
    ```

    ```
    tree ent
    ent
    ├── client.go
    ├── config.go
    ├── context.go
    ├── ent.go
    ├── enttest
    │   └── enttest.go
    ├── generate.go
    ├── hook
    │   └── hook.go
    ├── migrate
    │   ├── migrate.go
    │   └── schema.go
    ├── mutation.go
    ├── predicate
    │   └── predicate.go
    ├── runtime
    │   └── runtime.go
    ├── runtime.go
    ├── schema
    │   └── user.go
    ├── tx.go
    ├── user
    │   ├── user.go
    │   └── where.go
    ├── user.go
    ├── user_create.go
    ├── user_delete.go
    ├── user_query.go
    └── user_update.go

    7 directories, 22 files
    ```

### Create entity (postgres)
