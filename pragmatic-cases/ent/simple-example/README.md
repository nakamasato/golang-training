# Simple Example

## Prepare

1. Run postgres
    ```
    docker compose up
    ```
1. Create database
    ```
    docker exec -it postgres psql -U postgres -c 'drop database if exists ent_simple_example'
    docker exec -it postgres psql -U postgres -c 'create database ent_simple_example'
    ```

## Create Item

1. Create Item
    ```
    go run entgo.io/ent/cmd/ent init Item
    ```

    ```
    tree
    .
    ├── README.md
    └── ent
        ├── generate.go
        └── schema
            └── item.go

    2 directories, 3 files
    ```

1. Define Fields in [ent/schema/item.go](ent/schema/item.go)
    ```go
    // Fields of the Item.
    // id text NOT NULL PRIMARY KEY,
    // name VARCHAR(50) NOT NULL,
    // status SMALLINT NOT NULL,
    // created_at TIMESTAMP NOT NULL
    func (Item) Fields() []ent.Field {
    	return []ent.Field{
    		field.String("id").
    			StructTag(`json:"oid,omitempty"`),
    		field.String("name"),
    		field.Int("status"),
    		field.Time("created_at").
    			Default(time.Now),
    	}
    }
    ```
1. Update [ent/generate.go](ent/generate.go) to use Upsert
    ```go
    //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./schema
    ```
1. Generate go.
    ```
    go generate ./ent
    ```
1. Create a script `start/start.go`
    1. Create Schema
        ```go
        client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=ent_simple_example password=postgres sslmode=disable") // hardcoding
        if err != nil {
            log.Fatalf("failed opening connection to postgres: %v", err)
        }
        defer client.Close()
        // Run the auto migration tool.
        if err := client.Schema.Create(context.Background()); err != nil {
            log.Fatalf("failed creating schema resources: %v", err)
        }
        ```
    1. Create (Upsert) entity (postgres)

        ```go
        func UpsertItem(ctx context.Context, client *ent.Client) (string, error) {
            id, err := client.Debug().Item.
                Create().
                SetID("A0001").
                SetName("a8m").
                SetStatus(1).
                OnConflict(
                    sql.ConflictColumns(item.FieldID),
                ).
                Update(func(u *ent.ItemUpsert) {
                    u.SetName("a8m")
                    u.SetStatus(1)
                }).
                ID(ctx)
            return id, err
        }
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
1. Run [start/start.go](start/start.go) (rerunnable)
    ```
    go run start/start.go
    2022/11/11 10:20:34 driver.Query: query=INSERT INTO "items" ("name", "status", "created_at", "id") VALUES ($1, $2, $3, $4) ON CONFLICT ("id") DO UPDATE SET "name" = $5, "status" = $6 RETURNING "id" args=[a8m 1 2022-11-11 10:20:34.8543 +0900 JST m=+0.083714918 A0001 a8m 1]
    id: A0001
    2022/11/11 10:20:34 driver.Query: query=SELECT DISTINCT "items"."id", "items"."name", "items"."status", "items"."created_at" FROM "items" WHERE "items"."name" = $1 LIMIT 2 args=[a8m]
    2022/11/11 10:20:34 item returned:  Item(id=A0001, name=a8m, status=1, created_at=Fri Nov 11 01:20:06 2022)
    ```
1. Check db

    ```
    docker exec -it postgres psql -U postgres -d ent_simple_example -c 'select * from items;'
      id   | name | status |          created_at
    -------+------+--------+-------------------------------
     A0001 | a8m  |      1 | 2022-11-11 01:20:06.742316+00
    (1 row)
    ```

## Create Edge

1. Create `Category`
    ```
    go run entgo.io/ent/cmd/ent init Category
    ```
