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

### Add first Edge (Relation) User -> Cars

1. Add two more entities `Car` and `Group`

    ```
    go run entgo.io/ent/cmd/ent init Car Group
    ```

    You might need to install printer `go get entgo.io/ent/cmd/internal/printer@v0.11.1`

1. Add fields to `Car` and `Group`.

    ```go
    func (Car) Fields() []ent.Field {
        return []ent.Field{
            field.String("model"),
            field.Time("registered_at"),
        }
    }
    ```

    ```go
    func (Group) Fields() []ent.Field {
        return []ent.Field{
            field.String("name").
                // Regexp validation for group name.
                Match(regexp.MustCompile("[a-zA-Z_]+$")),
        }
    }
    ```

1. Add the "cars" edge to the User schema.

    ```go
    // Edges of the User.
    func (User) Edges() []ent.Edge {
        return []ent.Edge{
            edge.To("cars", Car.Type),
        }
    }
    ```

1. Add `CreateCars` and `QueryCars` to `start/start.go`.

    ```go
    func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
        // Create a new car with model "Tesla".
        tesla, err := client.Car.
            Create().
            SetModel("Tesla").
            SetRegisteredAt(time.Now()).
            Save(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed creating car: %w", err)
        }
        log.Println("car was created: ", tesla)

        // Create a new car with model "Ford".
        ford, err := client.Car.
            Create().
            SetModel("Ford").
            SetRegisteredAt(time.Now()).
            Save(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed creating car: %w", err)
        }
        log.Println("car was created: ", ford)

        // Create a new user, and add it the 2 cars.
        a8m, err := client.User.
            Create().
            SetAge(30).
            SetName("a8m").
            AddCars(tesla, ford).
            Save(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed creating user: %w", err)
        }
        log.Println("user was created: ", a8m)
        return a8m, nil
    }

    func QueryCars(ctx context.Context, a8m *ent.User) error {
        cars, err := a8m.QueryCars().All(ctx)
        if err != nil {
            return fmt.Errorf("failed querying user cars: %w", err)
        }
        log.Println("returned cars:", cars)

        // What about filtering specific cars.
        ford, err := a8m.QueryCars().
            Where(car.Model("Ford")).
            Only(ctx)
        if err != nil {
            return fmt.Errorf("failed querying user cars: %w", err)
        }
        log.Println(ford)
        return nil
    }
    ```

### Add Inverse Edge (BackRef) Car -> User

1. Add Edges to Car `ent/schema/car.go`

    ```go
    // Edges of the Car.
    func (Car) Edges() []ent.Edge {
        return []ent.Edge{
            // Create an inverse-edge called "owner" of type `User`
            // and reference it to the "cars" edge (in User schema)
            // explicitly using the `Ref` method.
            edge.From("owner", User.Type).
                Ref("cars").
                // setting the edge to unique, ensure
                // that a car can have only one owner.
                Unique(),
        }
    }
    ```

1. Add `QueryCarUsers` to `start/start.go`.

    ```go
    func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
        cars, err := a8m.QueryCars().All(ctx)
        if err != nil {
            return fmt.Errorf("failed querying user cars: %w", err)
        }
        // Query the inverse edge.
        for _, c := range cars {
            owner, err := c.QueryOwner().Only(ctx)
            if err != nil {
                return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
            }
            log.Printf("car %q owner: %q\n", c.Model, owner.Name)
        }
        return nil
    }
    ```

### Create another Edge Group -> User

1. Add Edge to Group.
    ```go
    // Edges of the Group.
    func (Group) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("users", User.Type),
    }
    }
    ```

1. Update User.Edges.

    ```go
    // Edges of the User.
    func (User) Edges() []ent.Edge {
        return []ent.Edge{
            edge.To("cars", Car.Type),
            // Create an inverse-edge called "groups" of type `Group`
            // and reference it to the "users" edge (in Group schema)
            // explicitly using the `Ref` method.
            edge.From("groups", Group.Type).
                Ref("users"),
        }
    }
    ```

1. Generate ent assets

    ```
    go generate ./ent
    ```

1. Add `CreateGraph` in `start/start.go`

    ```go
    func CreateGraph(ctx context.Context, client *ent.Client) error {
        // First, create the users.
        a8m, err := client.User.
            Create().
            SetAge(30).
            SetName("Ariel").
            Save(ctx)
        if err != nil {
            return err
        }
        neta, err := client.User.
            Create().
            SetAge(28).
            SetName("Neta").
            Save(ctx)
        if err != nil {
            return err
        }
        // Then, create the cars, and attach them to the users created above.
        err = client.Car.
            Create().
            SetModel("Tesla").
            SetRegisteredAt(time.Now()).
            // Attach this car to Ariel.
            SetOwner(a8m).
            Exec(ctx)
        if err != nil {
            return err
        }
        err = client.Car.
            Create().
            SetModel("Mazda").
            SetRegisteredAt(time.Now()).
            // Attach this car to Ariel.
            SetOwner(a8m).
            Exec(ctx)
        if err != nil {
            return err
        }
        err = client.Car.
            Create().
            SetModel("Ford").
            SetRegisteredAt(time.Now()).
            // Attach this graph to Neta.
            SetOwner(neta).
            Exec(ctx)
        if err != nil {
            return err
        }
        // Create the groups, and add their users in the creation.
        err = client.Group.
            Create().
            SetName("GitLab").
            AddUsers(neta, a8m).
            Exec(ctx)
        if err != nil {
            return err
        }
        err = client.Group.
            Create().
            SetName("GitHub").
            AddUsers(a8m).
            Exec(ctx)
        if err != nil {
            return err
        }
        log.Println("The graph was created successfully")
        return nil
    }
    ```

1. Add `QueryGitHub` to `start/start.go`

    ```go
    func QueryGithub(ctx context.Context, client *ent.Client) error {
        cars, err := client.Group.
            Query().
            Where(group.Name("GitHub")). // (Group(Name=GitHub),)
            QueryUsers().                // (User(Name=Ariel, Age=30),)
            QueryCars().                 // (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
            All(ctx)
        if err != nil {
            return fmt.Errorf("failed getting cars: %w", err)
        }
        log.Println("cars returned:", cars)
        // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
        return nil
    }
    ```

1. Add `QueryArielCars` to `start/start.go`

    ```go
    func QueryArielCars(ctx context.Context, client *ent.Client) error {
        // Get "Ariel" from previous steps.
        a8m := client.User.
            Query().
            Where(
                user.HasCars(),
                user.Name("Ariel"),
            ).
            OnlyX(ctx)
        cars, err := a8m. // Get the groups, that a8m is connected to:
                    QueryGroups(). // (Group(Name=GitHub), Group(Name=GitLab),)
                    QueryUsers().  // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
                    QueryCars().   //
                    Where(         //
                car.Not( //  Get Neta and Ariel cars, but filter out
                    car.Model("Mazda"), //  those who named "Mazda"
                ), //
            ). //
            All(ctx)
        if err != nil {
            return fmt.Errorf("failed getting cars: %w", err)
        }
        log.Println("cars returned:", cars)
        // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Ford, RegisteredAt=<Time>),)
        return nil
    }
    ```

1. Add `QueryGroupWithUsers` to `start/start.go`

    ```go
    func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
        groups, err := client.Group.
            Query().
            Where(group.HasUsers()).
            All(ctx)
        if err != nil {
            return fmt.Errorf("failed getting groups: %w", err)
        }
        log.Println("groups returned:", groups)
        // Output: (Group(Name=GitHub), Group(Name=GitLab),)
        return nil
    }
    ```

1. Update `main()` in `start/start.go` to create records.

    1. Create schema: `client.Schema.Create(ctx)`
    1. Create entities and query them.

    ```go
    func main() {
        client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=password sslmode=disable") // hardcoding
        if err != nil {
            log.Fatalf("failed opening connection to postgres: %v", err)
        }
        defer client.Close()
        ctx := context.Background()
        // Run the auto migration tool.
        if err := client.Schema.Create(ctx); err != nil {
            log.Fatalf("failed creating schema resources: %v", err)
        }

        if _, err = CreateUser(ctx, client); err != nil {
            log.Fatal(err)
        }
        if _, err = QueryUser(ctx, client); err != nil {
            log.Fatal(err)
        }
        a8m, err := CreateCars(ctx, client)
        if err != nil {
            log.Fatal(err)
        }
        if err := QueryCars(ctx, a8m); err != nil {
            log.Fatal(err)
        }
        if err := QueryCarUsers(ctx, a8m); err != nil {
            log.Fatal(err)
        }
        if err := CreateGraph(ctx, client); err != nil {
            log.Fatal(err)
        }
        if err := QueryGithub(ctx, client); err != nil {
            log.Fatal(err)
        }
        if err := QueryArielCars(ctx, client); err != nil {
            log.Fatal(err)
        }
        if err := QueryGroupWithUsers(ctx, client); err != nil {
            log.Fatal(err)
        }
    }
    ```

### Run

1. Run postgres with docker.

    ```
    docker-compose up -d
    ```

1. Run `start/start.go`.

    ```
    go run start/start.go
    2022/07/15 14:40:26 user was created:  User(id=1, age=30, name=a8m)
    2022/07/15 14:40:26 user returned:  User(id=1, age=30, name=a8m)
    2022/07/15 14:40:26 car was created:  Car(id=1, model=Tesla, registered_at=Fri Jul 15 14:40:26 2022)
    2022/07/15 14:40:26 car was created:  Car(id=2, model=Ford, registered_at=Fri Jul 15 14:40:26 2022)
    2022/07/15 14:40:26 user was created:  User(id=2, age=30, name=a8m)
    2022/07/15 14:40:26 returned cars: [Car(id=1, model=Tesla, registered_at=Fri Jul 15 14:40:26 2022) Car(id=2, model=Ford, registered_at=Fri Jul 15 14:40:26 2022)]
    2022/07/15 14:40:26 Car(id=2, model=Ford, registered_at=Fri Jul 15 14:40:26 2022)
    2022/07/15 14:40:26 car "Tesla" owner: "a8m"
    2022/07/15 14:40:26 car "Ford" owner: "a8m"
    2022/07/15 14:40:26 The graph was created successfully
    2022/07/15 14:40:26 cars returned: [Car(id=3, model=Tesla, registered_at=Fri Jul 15 14:40:26 2022) Car(id=4, model=Mazda, registered_at=Fri Jul 15 14:40:26 2022)]
    2022/07/15 14:40:26 cars returned: [Car(id=3, model=Tesla, registered_at=Fri Jul 15 14:40:26 2022) Car(id=5, model=Ford, registered_at=Fri Jul 15 14:40:26 2022)]
    2022/07/15 14:40:26 groups returned: [Group(id=2, name=GitHub) Group(id=1, name=GitLab)]
    ```

1. Check the schema.

    ```
    docker exec -it postgres psql -U postgres -c '\dt'
                List of relations
     Schema |    Name     | Type  |  Owner
    --------+-------------+-------+----------
     public | cars        | table | postgres
     public | group_users | table | postgres
     public | groups      | table | postgres
     public | users       | table | postgres
    (4 rows)
    ```

1. Check records.

    ```
    docker exec -it postgres psql -U postgres -c 'select * from users'
     id | age | name
    ----+-----+-------
      1 |  30 | a8m
      2 |  30 | a8m
      3 |  30 | Ariel
      4 |  28 | Neta
    (4 rows)

    docker exec -it postgres psql -U postgres -c 'select * from cars'
     id | model |         registered_at         | user_cars
    ----+-------+-------------------------------+-----------
      1 | Tesla | 2022-07-15 14:40:26.003466+09 |         2
      2 | Ford  | 2022-07-15 14:40:26.015701+09 |         2
      3 | Tesla | 2022-07-15 14:40:26.074102+09 |         3
      4 | Mazda | 2022-07-15 14:40:26.078375+09 |         3
      5 | Ford  | 2022-07-15 14:40:26.081656+09 |         4
    (5 rows)

    docker exec -it postgres psql -U postgres -c 'select * from groups'
     id |  name
    ----+--------
      1 | GitLab
      2 | GitHub
    (2 rows)

    docker exec -it postgres psql -U postgres -c 'select * from group_users'
     group_id | user_id
    ----------+---------
            1 |       4
            1 |       3
            2 |       3
    (3 rows)
    ```
