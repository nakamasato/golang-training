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

1. Create a `ent.client` in `start/start.go`.

    ```go
    package main

    import (
        "context"
        "fmt"
        "log"

        "tmp/pragmatic-cases/ent/ent"

        _ "github.com/lib/pq"
    )

    func main() {
        client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres") // hardcoding
        if err != nil {
            log.Fatalf("failed opening connection to postgres: %v", err)
        }
        defer client.Close()
        // Run the auto migration tool.
        if err := client.Schema.Create(context.Background()); err != nil {
            log.Fatalf("failed creating schema resources: %v", err)
        }
    }

    func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
        u, err := client.User.
            Create().
            SetAge(30).
            SetName("a8m").
            Save(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed creating user: %w", err)
        }
        log.Println("user was created: ", u)
        return u, nil
    }
    ```

### Query entity

Add the following code to `start/start.go`

```go
func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
```

### Add first Edge (Relation)

Add another entity `Car`

```
go run entgo.io/ent/cmd/ent init Car Group
```



### Run
1. Run postgres with docker.

```
docker run --name postgres \
           -e POSTGRES_PASSWORD=password \
           -e POSTGRES_INITDB_ARGS="--encoding=UTF8 --no-locale" \
           -e TZ=Asia/Tokyo \
           -v postgresdb:/var/lib/postgresql/data \
           -p 5432:5432 \
           -d postgres
```
