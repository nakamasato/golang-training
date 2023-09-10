# [ent versioned migrations](https://entgo.io/docs/versioned-migrations#moving-from-auto-migration-to-versioned-migrations)


## Overview

1. Define Ent schema in `ent/schema/xxx.go`
1. Generate with `go generate ./ent` -> generate migration files
1. Apply migration files with `atlas migrate apply` command

    ```
    atlas migrate apply \
        --dir "file://ent/migrate/migrations" \
        --url "postgres://postgres:pass@localhost:5432/database?search_path=public&sslmode=disable"
    ```

## Prerequisite

1. Code preparation[Getting Started](../getting-started/README.md)
1. atlas cli: `curl -sSf https://atlasgo.sh | sh`

## Steps

### 1. Create a migration generation script

Add `--feature sql/versioned-migration` to [ent/generate.go](ent/generate.go)

```diff
- //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
+ //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./schema
```

### 2. Run `go generate ./ent`

```
GOWORK=off go generate ./ent
```

The following methods are added to [ent/migrate/migrate.go](ent/migrate/migrate.go)

```go
// Diff compares the state read from a database connection or migration directory with
// the state defined by the Ent schema. Changes will be written to new migration files.
func Diff(ctx context.Context, url string, opts ...schema.MigrateOption) error {
	return NamedDiff(ctx, url, "changes", opts...)
}

// NamedDiff compares the state read from a database connection or migration directory with
// the state defined by the Ent schema. Changes will be written to new named migration files.
func NamedDiff(ctx context.Context, url, name string, opts ...schema.MigrateOption) error {
	return schema.Diff(ctx, url, name, Tables, opts...)
}

// Diff creates a migration file containing the statements to resolve the diff
// between the Ent schema and the connected database.
func (s *Schema) Diff(ctx context.Context, opts ...schema.MigrateOption) error {
	migrate, err := schema.NewMigrate(s.drv, opts...)
	if err != nil {
		return fmt.Errorf("ent/migrate: %w", err)
	}
	return migrate.Diff(ctx, Tables...)
}

// NamedDiff creates a named migration file containing the statements to resolve the diff
// between the Ent schema and the connected database.
func (s *Schema) NamedDiff(ctx context.Context, name string, opts ...schema.MigrateOption) error {
	migrate, err := schema.NewMigrate(s.drv, opts...)
	if err != nil {
		return fmt.Errorf("ent/migrate: %w", err)
	}
	return migrate.NamedDiff(ctx, name, Tables...)
}
```

### 3. Generate Versioned Migration Files

1. Start postgres container

    ```
    docker run --name migration --rm -p 5432:5432 -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=test -d postgres
    ```

1. Add `ent/migrate.main.go`

    ```go
    //go:build ignore

    package main

    import (
    	"context"
    	"log"
    	"os"

    	"tmp/pragmatic-cases/ent/getting-started-with-versioned-migrations/ent/migrate"

    	atlas "ariga.io/atlas/sql/migrate"
    	"entgo.io/ent/dialect"
    	"entgo.io/ent/dialect/sql/schema"
    	_ "github.com/lib/pq" // postgres driver
    )

    func main() {
    	ctx := context.Background()
    	// Create a local migration directory able to understand Atlas migration file format for replay.
    	dir, err := atlas.NewLocalDir("ent/migrate/migrations")
    	if err != nil {
    		log.Fatalf("failed creating atlas migration directory: %v", err)
    	}
    	// Migrate diff options.
    	opts := []schema.MigrateOption{
    		schema.WithDir(dir),                         // provide migration directory
    		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
    		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
    		schema.WithFormatter(atlas.DefaultFormatter),
    	}
    	if len(os.Args) != 2 {
    		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
    	}
    	// Generate migrations using Atlas support for Postgres (note the Ent dialect option passed above).
    	err = migrate.NamedDiff(ctx, os.Getenv("DSN"), os.Args[1], opts...) // "postgres://postgres:pass@localhost:5432/test?sslmode=disable"
    	if err != nil {
    		log.Fatalf("failed generating migration file: %v", err)
    	}
    }
    ```

    There's another option to generate the migration files:

    ```
    atlas migrate diff migration_name \
        --dir "file://ent/migrate/migrations" \
        --to "ent://ent/schema" \
        --dev-url "docker://postgres/15/test?search_path=public"
    ```

1. Create a dir `ent/migrate/migrations`

    ```
    mkdir -p ent/migrate/migrations
    ```

1. Trigger migration generation by executing `go run -mod=mod ent/migrate/main.go <name>`

    ```
    DSN="postgres://postgres:pass@localhost:5432/test?sslmode=disable" GOWORK=off go run -mod=mod ent/migrate/main.go init_db
    ```

    `<name>`: You can use any name. I used `init_db` as this is the initial generation of migration files.

1. Two files are generated

    ```
    ls ent/migrate/migrations
    20230910044615_init_db.sql atlas.sum
    ```

    `cat ent/migrate/migrations/20230910044615_init_db.sql`:

    ```sql
    -- Create "users" table
    CREATE TABLE "users" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "age" bigint NOT NULL, "name" character varying NOT NULL DEFAULT 'unknown', PRIMARY KEY ("id"));
    -- Create "cars" table
    CREATE TABLE "cars" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "model" character varying NOT NULL, "registered_at" timestamptz NOT NULL, "user_cars" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "cars_users_cars" FOREIGN KEY ("user_cars") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
    -- Create "groups" table
    CREATE TABLE "groups" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, PRIMARY KEY ("id"));
    -- Create "group_users" table
    CREATE TABLE "group_users" ("group_id" bigint NOT NULL, "user_id" bigint NOT NULL, PRIMARY KEY ("group_id", "user_id"), CONSTRAINT "group_users_group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "group_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
    ```

    `cat ent/migrate/migrations/atlas.sum`:

    ```
    h1:cLF1KcKbNRhepzh2GaVjPFEMBCfrnrVGX79hIJiFz50=
    20230910044615_init_db.sql h1:0saPbV9Qj/VfaaEiHhsOKKXT5LUNUE1cxDUPtezLyV0=
    ```

1. (Optional) Run `atlas migrate diff` to generate migration files: (no changes)

    ```
    atlas migrate diff migration_name \
    --dir "file://ent/migrate/migrations" \
    --to "ent://ent/schema" \
    --dev-url "docker://postgres/15/test?search_path=public"
    The migration directory is synced with the desired state, no changes to be made
    ```

## 4. Verify and lint migrations

```
atlas migrate lint \
  --dev-url="docker://postgres/15/test?search_path=public" \
  --dir="file://ent/migrate/migrations" \
  --latest=1
```

## 5. Apply migration files

```
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:pass@localhost:5432/database?search_path=public&sslmode=disable"
```

```
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
Migrating to version 20230910044615 (1 migrations in total):

  -- migrating version 20230910044615
    -> CREATE TABLE "users" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "age" bigint NOT NULL, "name" character varying NOT NULL DEFAULT 'unknown', PRIMARY KEY ("id"));
    -> CREATE TABLE "cars" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "model" character varying NOT NULL, "registered_at" timestamptz NOT NULL, "user_cars" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "cars_users_cars" FOREIGN KEY ("user_cars") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
    -> CREATE TABLE "groups" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, PRIMARY KEY ("id"));
    -> CREATE TABLE "group_users" ("group_id" bigint NOT NULL, "user_id" bigint NOT NULL, PRIMARY KEY ("group_id", "user_id"), CONSTRAINT "group_users_group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "group_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
  -- ok (20.090291ms)

  -------------------------
  -- 33.010958ms
  -- 1 migrations
  -- 4 sql statements
```

If you rerun, no change will be made.

```
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
No migration files to execute
```

## 5. Run app

```
DSN="postgres://postgres:pass@localhost:5432/test?sslmode=disable" go run start/start.go
```

## Tips: commands

1. `atlas migrate validate --dir file://ent/migrate/migrations`:
1. `atlas migrate hash --dir file://ent/migrate/migrations`

