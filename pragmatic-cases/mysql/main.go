package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	db, err := Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow("select User from mysql.user")
	if row.Err() != nil {
		log.Fatal(row.Err())
	}

	if err := CreateMySQLUser(db, "test_user"); err != nil {
		log.Fatal(err)
	}

	res, err := CheckMySQLHasUser(db, "test_user")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("test_user exists: %t\n", res)
}

func Connect(ctx context.Context) (*sql.DB, error) {
	data := make(chan *sql.DB, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/")
				if err != nil {
					time.Sleep(10 * time.Second)
					continue
				}
				err = db.Ping()
				if err != nil {
					db.Close()
					continue
				}
				data <- db
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case db := <-data:
		return db, nil
	}
}

func CheckMySQLHasUser(db *sql.DB, mysqluser string) (bool, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM mysql.user where User = '" + mysqluser + "'")
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func CreateMySQLUser(db *sql.DB, mysqluser string) error {
	_, err := db.Exec("CREATE USER IF NOT EXISTS '" + mysqluser + "'")
	return err
}
