# migrate (mysql)

1. Create `main.go`

    Create database connection pool

    ```go
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("failed to Open: %v\n", err)
	}
	defer db.Close()
    ```

    Create migrate driver & migrate instance

    ```go
	driver, err := mysql.WithInstance(db, &mysql.Config{
		DatabaseName: dbname,
	})
	if err != nil {
		log.Fatalf("failed to init driver: %v\n", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // sourceURL
		// "github://nakamasato/golang-training/pragmatic-cases/migrate/mysql/migrations#update-migrate",
		dbname, // databaseName
		driver, // databaseInstance
	)
    ```

    Execute migration

    ```go
    m.Up()
    m.Down()
    ```

1. Run mysql

    ```
    docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --rm mysql:8
    ```

1. Create `test_db` database

    ```
    docker exec -it $(docker ps | grep mysql | head -1 |awk '{print $1}') mysql -uroot -ppassword -e 'create database test_db;'
    ```

1. Run migrate

    ```
    MYSQL_DB_NAME=test_db MYSQL_HOST=localhost MYSQL_PASSWORD=password go run main.go
    ```

    ```
    Up
    Up finished
    show tables
    table: schema_migrations
    table: test_table
    -------------
    Down
    Down finished
    show tables
    table: schema_migrations
    ```
