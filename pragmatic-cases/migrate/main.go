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
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("failed to Open: %v\n", err)
	}
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
