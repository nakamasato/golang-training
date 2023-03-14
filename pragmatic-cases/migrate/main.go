package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	passwd := os.Getenv("POSTGRES_PASSWORD")
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres password=%s dbname=postgres sslmode=disable", host, passwd))
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
	err = m.Up()
	if err != nil {
		log.Fatalf("Up failed: %v\n", err)
	}
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
	err = m.Down()
	if err != nil {
		log.Fatalf("Down failed: %v\n", err)
	}
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
