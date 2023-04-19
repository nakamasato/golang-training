package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	. "github.com/go-sql-driver/mysql"
)

const hcl_str = `
table "users" {
	schema = schema.test_db
	column "id" {
		type = bigint
		unsigned = true
	}
	column "shouldnt" {
		type = bigint
	}
}
schema "test_db" {
}
`

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "set for dryrun")
}

func main() {
	flag.Parse()
	fmt.Println(dryRun)
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

	fmt.Println("InspectSchema: read schema from MySQL")

	// Inspect the schema of the connected database
	curSch, err := driver.InspectSchema(ctx, "", nil)
	if err != nil {
		log.Fatalf("failed inspecting schema: %s", err)
	}
	fmt.Printf("shema: %s\n", curSch.Name)
	for _, tbl := range curSch.Tables {
		fmt.Printf("----- table %s ----\n", tbl.Name)
		for i, col := range tbl.Columns {
			fmt.Printf("col %d name: %s, type: %s\n", i, col.Name, col.Type.Raw)
		}
	}

	// TODO: https://pkg.go.dev/ariga.io/atlas@v0.10.0/schemahcl#State.EvalFiles
	// https://github.com/ariga/atlas/blob/3e658c6bb46607404434135eb3c190fcfc58919b/internal/integration/hclsqlspec/hclsqlspec_test.go
	fmt.Println("EvalHCLBytes: read from schema.hcl")
	var s schema.Schema

	bytes, err := ioutil.ReadFile("schema.hcl")
	if err != nil {
		log.Fatal(err)
	}
	err = mysql.EvalHCLBytes(bytes, &s, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Schema diff
	fmt.Println("SchemaDiff: compare schema objects")
	changes, err := driver.SchemaDiff(curSch, &s)
	if err != nil {
		log.Fatalf("failed to scheme diff: %s", err)
	}
	if len(changes) == 0 {
		fmt.Println("no changes")
	} else if dryRun { // only plan
		fmt.Println("PlayChanges")
		// https://github.com/ariga/atlas/blob/6d8605ca50556d8e6d6b3884f04b07894529f87d/sql/mysql/migrate.go#L31
		plan, err := driver.PlanChanges(ctx, "test", changes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(plan)
	} else { // apply
		fmt.Println("ApplyChanges")
		err = driver.ApplyChanges(ctx, changes)
		if err != nil {
			log.Fatal(err)
		}
	}
}
