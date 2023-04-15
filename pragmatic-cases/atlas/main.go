package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas/sql/mysql"
	. "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()

	// Prepare mysql.Config
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
	// Open a "connection" to mysql.
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("failed opening db: %s", err)
	}
	// Open an atlas driver.
	driver, err := mysql.Open(db)
	if err != nil {
		log.Fatalf("failed opening atlas driver: %s", err)
	}

	// Inspect the schema of the connected database
	curSch, err := driver.InspectSchema(ctx, "", nil)
	if err != nil {
		log.Fatalf("failed inspecting schema: %s", err)
	}
	fmt.Println(curSch.Name)
	for _, tbl := range curSch.Tables {
		fmt.Printf("----- table %s ----\n", tbl.Name)
		for i, col := range tbl.Columns {
			fmt.Printf("col %d: %s\n", i, col.Name)
		}
	}

	// TODO: https://pkg.go.dev/ariga.io/atlas@v0.10.0/schemahcl#State.EvalFiles
	// state := schemahcl.State{}
	// schema := &schema.Schema{}
	// err = state.EvalFiles([]string{"schema.hcl"}, schema, map[string]cty.Value{})

	// if err != nil {
	// 	log.Fatal(err)
	// }
	schema := curSch

	// Schema diff
	change, err := driver.SchemaDiff(curSch, schema)
	if err != nil {
		log.Fatalf("failed to scheme diff: %s", err)
	}
	fmt.Println(change)
}
