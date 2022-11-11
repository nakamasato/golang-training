# Simple Example


1. Run postgres
    ```
    docker compose up
    ```
1. Create database
    ```
    docker exec -it postgres psql -U postgres -c 'create database ent_simple_example'
    ```
1. Create Item and ItemCategory
    ```
    go run entgo.io/ent/cmd/ent init Item
    go run entgo.io/ent/cmd/ent init ItemCategory
    ```

    ```
    tree
    .
    ├── README.md
    └── ent
        ├── generate.go
        └── schema
            ├── item.go
            └── itemcategory.go

    2 directories, 4 files
    ```

1. Define Fields in [ent/schema/item.go](ent/schema/item.go)

1. Generate go.
    ```
    go generate ./ent
    ```
1. Check DB
    ```
    docker exec -it postgres psql -U postgres -d ent_simple_example -c 'select * from items;'
     id | name | status | created_at
    ----+------+--------+------------
    (0 rows)
    ```
1. Create (Upsert) entity (postgres)

    ```go
    id, err := client.Item.Debug().
        Create().
        SetID("A0001") // UserIDがUnique指定されている場合
        SetName("a8m").
        SetStatus(1).
        OnConflict(
            sql.ConflictColumns(item.FieldID),
        ).
        UpdateNewValues().
        ID(ctx)
    ```

1. Query entity (postgres)

    ```go
    func QueryItem(ctx context.Context, client *ent.Client) (*ent.Item, error) {
        i, err := client.Item.
            Query().
            Where(item.Name("a8m")).
            Only(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed querying item: %w", err)
        }
        log.Println("item returned: ", i)
        return i, nil
    }
    ```
