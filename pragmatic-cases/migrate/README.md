# [migrate](https://github.com/golang-migrate/migrate)

## Prerequisite

Run postgres on Docker

```
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_INITDB_ARGS="--encoding=UTF8 --no-locale" -e TZ=Asia/Tokyo -p 5432:5432 -d postgres
```

## Getting started (postgres)

1. Write `main.go`

	```go
	package main

	import (
		"database/sql"
		"fmt"
		"log"

		"github.com/golang-migrate/migrate/v4"
		"github.com/golang-migrate/migrate/v4/database/postgres"
		_ "github.com/golang-migrate/migrate/v4/source/file"
		_ "github.com/lib/pq"
	)

	func main() {
		host := os.Getenv("POSTGRES_HOST")
		passwd := os.Getenv("POSTGRES_PASSWORD")
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres password=%s dbname=postgres sslmode=disable", host, passwd))
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatalf("failed to init driver: %v\n", err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			log.Fatalf("failed to initialize migrate: %v\n", err)
		}
		fmt.Println("Up")
		m.Up()
		fmt.Println("Up finished")

		fmt.Println("select roles")
		rows, err := db.Query("SELECT rolname FROM pg_roles WHERE rolname = $1", "my_user")
		if err != nil {
			log.Fatalf("failed to execute select role: %v\n", err)
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Fatal(err)
			}
			fmt.Println(name)
		}
		defer rows.Close()

		fmt.Println("-------------")
		fmt.Println("Down")
		m.Down()
		fmt.Println("Down finished")
		fmt.Println("select roles")
		rows, err = db.Query("SELECT rolname FROM pg_roles WHERE rolname = $1", "my_user")
		if err != nil {
			log.Fatalf("failed to execute select role: %v\n", err)
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Fatal(err)
			}
			fmt.Println(name)
		}
		defer rows.Close()

	}
	```

1. Run
	```
	POSTGRES_HOST=localhost POSTGRES_PASSWORD=postgres go run main.go
	Up
	Up finished
	select roles
	my_user
	-------------
	Down
	Down finished
	select roles
	```

## CLI

1. Install

    ```
    brew install golang-migrate
    ```
1. Migrate with `migrate up` command
    ```
    migrate -source file://migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
    ```
1. Check
    ```
    docker exec -it postgres psql -U postgres
    psql (14.4 (Debian 14.4-1.pgdg110+1))
    Type "help" for help.

    postgres=#
    ```

    ```
    postgres=# SELECT rolname FROM pg_roles;
              rolname
    ---------------------------
     pg_database_owner
     pg_read_all_data
     pg_write_all_data
     pg_monitor
     pg_read_all_settings
     pg_read_all_stats
     pg_stat_scan_tables
     pg_read_server_files
     pg_write_server_files
     pg_execute_server_program
     pg_signal_backend
     postgres
     my_user
     (13 rows)
    ```

    ```
    postgres=# \dt
                   List of relations
     Schema |       Name        | Type  |  Owner
    --------+-------------------+-------+----------
     public | schema_migrations | table | postgres
    (1 row)

    postgres=# select * from schema_migrations;
     version | dirty
    ---------+-------
           1 | f
    (1 row)

    postgres=#
    ```
1. Migrate `down`

    ```
    migrate -source file://migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down
    Are you sure you want to apply all down migrations? [y/N]
    y
    Applying all down migrations
    1/d create_user (84.107458ms)
    ```

1. Check

    ```
    docker exec -it postgres psql -U postgres
    psql (14.4 (Debian 14.4-1.pgdg110+1))
    Type "help" for help.

    postgres=# SELECT rolname FROM pg_roles;
              rolname
    ---------------------------
     pg_database_owner
     pg_read_all_data
     pg_write_all_data
     pg_monitor
     pg_read_all_settings
     pg_read_all_stats
     pg_stat_scan_tables
     pg_read_server_files
     pg_write_server_files
     pg_execute_server_program
     pg_signal_backend
     postgres
    (12 rows)

    postgres=# select * from schema_migrations;
     version | dirty
    ---------+-------
    (0 rows)
    ```
