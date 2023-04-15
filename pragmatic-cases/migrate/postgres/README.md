# migrate (postgres)

1. `main.go`

    Create database connection pool:

    ```go
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres password=%s dbname=postgres sslmode=disable", host, passwd))
	if err != nil {
		log.Fatalf("failed to Open: %v\n", err)
	}
    ```

    Create migrate driver & migrate instance

    ```go
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to init driver: %v\n", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		// "file://migrations",
		"github://nakamasato/golang-training/pragmatic-cases/migrate/migrations#main",
		"postgres", driver)
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v\n", err)
	}
    ```

    Execute migration

    ```go
	m.Up()
	m.Down()
    ```

1. Run postgres
    ```
    docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_INITDB_ARGS="--encoding=UTF8 --no-locale" -e TZ=Asia/Tokyo -p 5432:5432 -d postgres
    ```
1. Run migrate

    ```
    POSTGRES_HOST=localhost POSTGRES_PASSWORD=postgres go run main.go
    ```

    ```
    Up
    Up finished
    select roles
    my_user
    -------------
    Down
    Down finished
    select roles
    ```
