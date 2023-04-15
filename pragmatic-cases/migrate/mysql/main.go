package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	. "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("MYSQL_HOST")
	passwd := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DB_NAME")
	config := Config{
		User:   "root",
		Passwd: passwd,
		Addr:   fmt.Sprintf("%s:%d", host, 3306),
		Net:    "tcp",
		DBName: dbname,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("failed to Open: %v\n", err)
	}
	defer db.Close()
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
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v\n", err)
	}
	fmt.Println("Up")
	err = m.Up()
	if err != nil {
		log.Fatalf("Up failed: %v\n", err)
	}
	fmt.Println("Up finished")

	fmt.Println("show tables")
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("failed to execute select role: %v\n", err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("table: %s\n", name)
	}
	defer rows.Close()

	fmt.Println("-------------")
	fmt.Println("Down")
	err = m.Down()
	if err != nil {
		log.Fatalf("Down failed: %v\n", err)
	}
	fmt.Println("Down finished")
	fmt.Println("show tables")
	rows, err = db.Query("SHOW TABLES;")
	if err != nil {
		log.Fatalf("failed to execute select role: %v\n", err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("table: %s\n", name)
	}
	defer rows.Close()

}
