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

## 1. Create Item

![](item.drawio.svg)

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
                SetID("item_id_1").
                SetName("Item 1").
                SetStatus(1).
                OnConflict(
                    sql.ConflictColumns(item.FieldID),
                ).
                Update(func(u *ent.ItemUpsert) {
                    u.SetName("Item 1")
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
                Where(item.Name("Item 1")).
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

## 2. Create Edge (`Item` -> `Categories`) [O2M]

![](item-categories.drawio.svg)

1. Create `Category`
    ```
    go run entgo.io/ent/cmd/ent init Category
    ```
1. Create Schema [ent/schema/category.go](ent/schema/category.go)

    ```go
    // Fields of the Category.
    func (Category) Fields() []ent.Field {
        return []ent.Field{
            field.String("id").
                StructTag(`json:"oid,omitempty"`),
            field.String("name"),
        }
    }
    ```
1. Add `categories` Edge to Item schema. [ent/schema/item.go](ent/schema/item.go)

    ```go
    // Edges of the Item.
    func (Item) Edges() []ent.Edge {
        return []ent.Edge{
            edge.To("category", Category.Type),
        }
    }
    ```
1. Go generate
    ```
    go generate ./ent
    ```
1. Add `UpsertCategory` and `QueryCategory` func to [start/start.go](start/start.go)

    ```go
    func UpsertCategory(ctx context.Context, client *ent.Client) (string, error) {
        id, err := client.Debug().Category.
            Create().
            SetID("category_id_1"). // UserIDがUnique指定されている場合
            SetName("Category 1").
            OnConflict(
                sql.ConflictColumns(item.FieldID),
            ).
            Update(func(u *ent.CategoryUpsert) {
                u.SetName("Category 1")
            }).
            ID(ctx)
        return id, err
    }

    func QueryCategory(ctx context.Context, client *ent.Client) (*ent.Category, error) {
        c, err := client.Debug().Category.
            Query().
            Where(category.Name("Category 1")).
            Only(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed querying item: %w", err)
        }
        log.Println("category returned: ", c)
        return c, nil
    }
    ```
1. Run
    ```
    go run start/start.go
    2022/11/11 10:40:23 driver.Query: query=INSERT INTO "items" ("name", "status", "created_at", "id") VALUES ($1, $2, $3, $4) ON CONFLICT ("id") DO UPDATE SET "name" = $5, "status" = $6 RETURNING "id" args=[a8m 1 2022-11-11 10:40:23.994703 +0900 JST m=+0.105566626 A0001 a8m 1]
    id: A0001
    2022/11/11 10:40:23 driver.Query: query=SELECT DISTINCT "items"."id", "items"."name", "items"."status", "items"."created_at" FROM "items" WHERE "items"."name" = $1 LIMIT 2 args=[a8m]
    2022/11/11 10:40:24 item returned:  Item(id=A0001, name=a8m, status=1, created_at=Fri Nov 11 01:31:29 2022)
    2022/11/11 10:40:24 driver.Query: query=INSERT INTO "categories" ("name", "id") VALUES ($1, $2) ON CONFLICT ("id") DO UPDATE SET "name" = $3 RETURNING "id" args=[Category 1 category_id_1 Category 1]
    id: category_id_1
    2022/11/11 10:40:24 driver.Query: query=SELECT DISTINCT "categories"."id", "categories"."name" FROM "categories" WHERE "categories"."name" = $1 LIMIT 2 args=[Category 1]
    2022/11/11 10:40:24 category returned:  Category(id=category_id_1, name=Category 1)
    ```
1. Check DB
    ```
    docker exec -it postgres psql -U postgres -d ent_simple_example -c 'select * from categories;'
          id       |    name    | item_categories
    ---------------+------------+-----------------
     category_id_1 | Category 1 |
    (1 row)
    ```
1. Add `Edge` to the existing `Item`.

    ```go
    client.Debug().Item.UpdateOneID("A0001").AddCategory(category).Save(ctx)
    ```

1. Run `start/start.go`

    ```
    go run start/start.go
    2022/11/11 10:49:50 driver.Query: query=INSERT INTO "items" ("name", "status", "created_at", "id") VALUES ($1, $2, $3, $4) ON CONFLICT ("id") DO UPDATE SET "name" = $5, "status" = $6 RETURNING "id" args=[a8m 1 2022-11-11 10:49:50.084698 +0900 JST m=+0.112350835 A0001 a8m 1]
    id: A0001
    2022/11/11 10:49:50 driver.Query: query=SELECT DISTINCT "items"."id", "items"."name", "items"."status", "items"."created_at" FROM "items" WHERE "items"."name" = $1 LIMIT 2 args=[a8m]
    2022/11/11 10:49:50 item returned:  Item(id=A0001, name=a8m, status=1, created_at=Fri Nov 11 01:31:29 2022)
    2022/11/11 10:49:50 driver.Query: query=INSERT INTO "categories" ("name", "id") VALUES ($1, $2) ON CONFLICT ("id") DO UPDATE SET "name" = $3 RETURNING "id" args=[Category 1 category_id_1 Category 1]
    id: category_id_1
    2022/11/11 10:49:50 driver.Query: query=SELECT DISTINCT "categories"."id", "categories"."name" FROM "categories" WHERE "categories"."name" = $1 LIMIT 2 args=[Category 1]
    2022/11/11 10:49:50 category returned:  Category(id=category_id_1, name=Category 1)
    2022/11/11 10:49:50 driver.Tx(59aef2b3-be70-4b55-89f8-8eae7b01c2e5): started
    2022/11/11 10:49:50 Tx(59aef2b3-be70-4b55-89f8-8eae7b01c2e5).Exec: query=UPDATE "categories" SET "item_categories" = $1 WHERE "id" = $2 AND "item_categories" IS NULL args=[A0001 category_id_1]
    2022/11/11 10:49:50 Tx(59aef2b3-be70-4b55-89f8-8eae7b01c2e5).Query: query=SELECT "id", "name", "status", "created_at" FROM "items" WHERE "id" = $1 args=[A0001]
    2022/11/11 10:49:50 Tx(59aef2b3-be70-4b55-89f8-8eae7b01c2e5): committed
    ```

1. Check schema

    ```
     go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema
    Category:
            +-------+--------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
            | Field |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |       StructTag       | Validators |
            +-------+--------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
            | id    | string | false  | false    | false    | false   | false         | false     | json:"oid,omitempty"  |          0 |
            | name  | string | false  | false    | false    | false   | false         | false     | json:"name,omitempty" |          0 |
            +-------+--------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+

    Item:
            +------------+-----------+--------+----------+----------+---------+---------------+-----------+-----------------------------    +------------+
            |   Field    |   Type    | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |          StructTag          |     Validators |
            +------------+-----------+--------+----------+----------+---------+---------------+-----------+-----------------------------    +------------+
            | id         | string    | false  | false    | false    | false   | false         | false     | json:"oid,    omitempty"        |          0 |
            | name       | string    | false  | false    | false    | false   | false         | false     | json:"name,    omitempty"       |          0 |
            | status     | int       | false  | false    | false    | false   | false         | false     | json:"status,    omitempty"     |          0 |
            | created_at | time.Time | false  | false    | false    | true    | false         | false     | json:"created_at,    omitempty" |          0 |
            +------------+-----------+--------+----------+----------+---------+---------------+-----------+-----------------------------    +------------+
            +------------+----------+---------+---------+----------+--------+----------+
            |    Edge    |   Type   | Inverse | BackRef | Relation | Unique | Optional |
            +------------+----------+---------+---------+----------+--------+----------+
            | categories | Category | false   |         | O2M      | false  | true     |
            +------------+----------+---------+---------+----------+--------+----------+
    ```
1. Check `categories` table. `item_categories` field is added.

    ```sql
    docker exec -it postgres psql -U postgres -d ent_simple_example -c 'select * from categories;'
          id       |    name    | item_categories
    ---------------+------------+-----------------
     category_id_1 | Category 1 | A0001
    (1 row)
    ```

    ```sql
    docker exec -it postgres psql -U postgres -d ent_simple_example -c '\d categories'
                          Table "public.categories"
         Column      |       Type        | Collation | Nullable | Default
    -----------------+-------------------+-----------+----------+---------
     id              | character varying |           | not null |
     name            | character varying |           | not null |
     item_categories | character varying |           |          |
    Indexes:
        "categories_pkey" PRIMARY KEY, btree (id)
    Foreign-key constraints:
        "categories_items_categories" FOREIGN KEY (item_categories) REFERENCES items(id) ON DELETE SET NULL
    ```

## 3. Add Inverse Edge (`Category` -> `Item`)

![](item-categories-backref.drawio.svg)

1. Add `BackRef`

    ```go
    // Edges of the Category.
    func (Category) Edges() []ent.Edge {
        return []ent.Edge{
            // Create an inverse-edge called "items" of type `Items`
            // and reference it to the "categories" edge (in Item schema)
            // explicitly using the `Ref` method.
            edge.From("items", Item.Type).
                Ref("categories").
                Unique(),
        }
    }
    ```
1. Go generate
    ```
    go generate ./ent
    ```
1. Add `QueryCategoryForItem`

    ```go
    func QueryCategoryForItem(ctx context.Context, category *ent.Category) error {
        items, err := category.QueryItems().All(ctx)
        if err != nil {
            return fmt.Errorf("failed querying user categories: %w", err)
        }
        // Query the inverse edge.
        for _, i := range items {
            category, err := i.QueryCategory().Only(ctx)
            if err != nil {
                return fmt.Errorf("failed querying item %q category: %w", i.Name, err)
            }
            log.Printf("item %q category: %q\n", i.Name, category.Name)
        }
        return nil
    }
    ```

1. Run
    ```
    go run start/start.go
    ```

    ```
    2022/11/11 14:20:52 driver.Query: query=INSERT INTO "items" ("name", "status", "created_at", "id") VALUES ($1, $2, $3, $4) ON CONFLICT ("id") DO UPDATE SET "name" = $5, "status" = $6 RETURNING "id" args=[Item 1 1 2022-11-11 14:20:52.592921 +0900 JST m=+0.118637209 item_id_1 Item 1 1]
    id: item_id_1
    2022/11/11 14:20:52 driver.Query: query=SELECT DISTINCT "items"."id", "items"."name", "items"."status", "items"."created_at" FROM "items" WHERE "items"."name" = $1 LIMIT 2 args=[Item 1]
    2022/11/11 14:20:52 item returned:  Item(id=item_id_1, name=Item 1, status=1, created_at=Fri Nov 11 05:16:09 2022)
    2022/11/11 14:20:52 driver.Query: query=INSERT INTO "categories" ("name", "id") VALUES ($1, $2) ON CONFLICT ("id") DO UPDATE SET "name" = $3 RETURNING "id" args=[Category 1 category_id_1 Category 1]
    id: category_id_1
    2022/11/11 14:20:52 driver.Query: query=SELECT DISTINCT "categories"."id", "categories"."name" FROM "categories" WHERE "categories"."name" = $1 LIMIT 2 args=[Category 1]
    2022/11/11 14:20:52 category returned:  Category(id=category_id_1, name=Category 1)
    2022/11/11 14:20:52 driver.Tx(12204e0b-c154-432b-a6ac-5a5098e01cc9): started
    2022/11/11 14:20:52 Tx(12204e0b-c154-432b-a6ac-5a5098e01cc9).Exec: query=UPDATE "categories" SET "item_category" = $1 WHERE "id" = $2 AND "item_category" IS NULL args=[A0001 category_id_1]
    2022/11/11 14:20:52 Tx(12204e0b-c154-432b-a6ac-5a5098e01cc9): rollbacked
    2022/11/11 14:20:52 driver.Query: query=SELECT DISTINCT "items"."id", "items"."name", "items"."status", "items"."created_at" FROM "items" JOIN (SELECT "item_category" FROM "categories" WHERE "id" = $1) AS "t1" ON "items"."id" = "t1"."item_category" args=[category_id_1]
    2022/11/11 14:20:52 driver.Query: query=SELECT DISTINCT "categories"."id", "categories"."name" FROM "categories" WHERE "item_category" = $1 LIMIT 2 args=[A0001]
    2022/11/11 14:20:52 item "a8m" category: "Category 1"
    ```
