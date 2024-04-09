# [ent versioned migrations](https://entgo.io/docs/versioned-migrations#moving-from-auto-migration-to-versioned-migrations)


## Overview

1. Define Ent schema in `ent/schema/xxx.go`
1. Generate with `go generate ./ent` -> generate migration files
1. Apply migration files with `atlas migrate apply` command (`atlas` cli: `brew install ariga/tap/atlas`)

    ```
    atlas migrate apply \
        --dir "file://ent/migrate/migrations" \
        --url "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
    ```

## Prerequisite

1. Code preparation[Getting Started](../getting-started/README.md)
1. atlas cli: `curl -sSf https://atlasgo.sh | sh`

## Steps

Under `getting-started-with-versioned-migrations` directory


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

1. Add `ent/migrate/main.go`

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

1. Create a dir `ent/migrate/migrations`

    ```
    mkdir -p ent/migrate/migrations
    ```

1. Trigger migration generation by executing `go run -mod=mod ent/migrate/main.go <name>`

    ```
    DSN="postgres://postgres:pass@localhost:5432/test?sslmode=disable" GOWORK=off go run -mod=mod ent/migrate/main.go init_db
    ```

    `<name>`: You can use any name. I used `init_db` as this is the initial generation of migration files.

    There's another option to generate the migration files:

    ```
    atlas migrate diff init_db \
        --dir "file://ent/migrate/migrations" \
        --to "ent://ent/schema" \
        --dev-url "docker://postgres/15/test?search_path=public"
    ```


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

### 4. Verify and lint migrations

```
atlas migrate lint \
  --dev-url="docker://postgres/15/test?search_path=public" \
  --dir="file://ent/migrate/migrations" \
  --latest=1
```

### 5. Apply migration files

Apply the migration files:

```
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
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

### 6. Run app


## From auto migration to versioned migrations

### 0. Prepare

```
rm -rf ent/migrate/migrations/
mkdir ent/migrate/migrations/
```

### 1. Create empty DB

```
docker run --name migration --rm -p 5532:5432 -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=test -d postgres
```

### 2. Enable auto migration

```go
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
```

### 3. Run application

```
DSN="postgres://postgres:pass@localhost:5532/test?sslmode=disable" go run start/start.go
```

```
2023/12/11 16:08:56 user was created:  User(id=1, age=30, name=a8m)
2023/12/11 16:08:56 driver.Query: query=SELECT DISTINCT "users"."id", "users"."age", "users"."name" FROM "users" WHERE "users"."name" = $1 LIMIT 2 args=[a8m]
2023/12/11 16:08:56 user returned:  User(id=1, age=30, name=a8m)
2023/12/11 16:08:56 car was created:  Car(id=1, model=Tesla, registered_at=Mon Dec 11 16:08:56 2023)
2023/12/11 16:08:56 car was created:  Car(id=2, model=Ford, registered_at=Mon Dec 11 16:08:56 2023)
2023/12/11 16:08:56 user was created:  User(id=2, age=30, name=a8m)
2023/12/11 16:08:56 returned cars: [Car(id=1, model=Tesla, registered_at=Mon Dec 11 07:08:56 2023) Car(id=2, model=Ford, registered_at=Mon Dec 11 07:08:56 2023)]
2023/12/11 16:08:56 Car(id=2, model=Ford, registered_at=Mon Dec 11 07:08:56 2023)
2023/12/11 16:08:56 car "Tesla" owner: "a8m"
2023/12/11 16:08:56 car "Ford" owner: "a8m"
2023/12/11 16:08:56 The graph was created successfully
2023/12/11 16:08:56 cars returned: [Car(id=3, model=Tesla, registered_at=Mon Dec 11 07:08:56 2023) Car(id=4, model=Mazda, registered_at=Mon Dec 11 07:08:56 2023)]
2023/12/11 16:08:56 cars returned: [Car(id=3, model=Tesla, registered_at=Mon Dec 11 07:08:56 2023) Car(id=5, model=Ford, registered_at=Mon Dec 11 07:08:56 2023)]
2023/12/11 16:08:56 groups returned: [Group(id=2, name=GitHub) Group(id=1, name=GitLab)]
```

### 4. Check db

```
psql --host localhost --port 5532 test -U postgres
Password for user postgres:
psql (15.3, server 16.1 (Debian 16.1-1.pgdg120+1))
WARNING: psql major version 15, server major version 16.
         Some psql features might not work.
Type "help" for help.

test=# \dt
            List of relations
 Schema |    Name     | Type  |  Owner
--------+-------------+-------+----------
 public | cars        | table | postgres
 public | group_users | table | postgres
 public | groups      | table | postgres
 public | users       | table | postgres
(4 rows)

test=#
```

### 5. Use versioned migrations

```go
	// if err := client.Schema.Create(ctx); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }
```

### 6. Generate migration files

```
atlas migrate diff init \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

```
ls ent/migrate/migrations
20231211071237_init.sql atlas.sum
```

### 7. Check

```
atlas schema apply \
  --url "postgresql://postgres:pass@localhost:5532/test?search_path=public&sslmode=disable" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

```
Schema is synced, no changes to be made
```

### 8. Change ent schema

Add `email` field to `User` schema:

```diff
+ 		field.String("email"),
```

### 9. Generate migration files

```
atlas migrate diff add_email_to_user \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

`ent/migrate/migrations/20231211071759_add_email_to_user.sql` is generated:

```
cat ent/migrate/migrations/20231211071759_add_email_to_user.sql
-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "email" character varying NOT NULL;
```

### 10. Apply migration files

```
atlas schema apply \
  --url "postgresql://postgres:pass@localhost:5532/test?search_path=public&sslmode=disable" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

Error:
```
atlas schema apply \
  --url "postgresql://postgres:pass@localhost:5532/test?search_path=public&sslmode=disable" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
-- Planned Changes:
-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "email" character varying NOT NULL;
✔ Apply
Error: modify "users" table: pq: column "email" of relation "users" contains null values
```

### 11. Fix the schema and generate the files

```go
		field.String("email").
			Default("unknown"),
```

```
atlas migrate diff add_email_to_user \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

```
diff ent/migrate/migrations/20231211071759_add_email_to_user.sql ent/migrate/migrations/20231211072624_add_email_to_user.sql
2c2
< ALTER TABLE "users" ADD COLUMN "email" character varying NOT NULL;
---
> ALTER TABLE "users" ALTER COLUMN "email" SET DEFAULT 'unknown';
```

This change is not correct.

### 12. Remove the two migration files manually

```
rm ent/migrate/migrations/20231211071759_add_email_to_user.sql
rm ent/migrate/migrations/20231211072624_add_email_to_user.sql
```

### 13. Try generating again

```
atlas migrate diff add_email_to_user \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

```
You have a checksum error in your migration directory.
This happens if you manually create or edit a migration file.
Please check your migration files and run

'atlas migrate hash'

to re-hash the contents and resolve the error

Error: checksum mismatch
```

### 14. Fix the hash

```
atlas migrate hash --dir file://ent/migrate/migrations
```

### 15. Regenerate

```
atlas migrate diff add_email_to_user \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

Now it's successfully generated:

```
ls ent/migrate/migrations
20231211071237_init.sql              20231211073102_add_email_to_user.sql atlas.sum
```

```
cat ent/migrate/migrations/20231211073102_add_email_to_user.sql
-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "email" character varying NOT NULL DEFAULT 'unknown';
```

### 16. Apply

```
atlas schema apply \
  --url "postgresql://postgres:pass@localhost:5532/test?search_path=public&sslmode=disable" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

```
ALTER TABLE "users" ADD COLUMN "email" character varying NOT NULL DEFAULT 'unknown';
-- Drop "atlas_schema_revisions" table
DROP TABLE "atlas_schema_revisions";
✔ Apply
```

## Configuration file `atlas.hcl`

> [!NOTE]
> The default is `atlas.hcl` but if I use the name, all the commands above references the configuration file, so in this case I use `atlas-config.hcl`.

### 1. Basic configuration


The command we used above can be simplified with the configuration file.

```
atlas migrate diff migration_name \
    --dir "file://ent/migrate/migrations" \
    --to "ent://ent/schema" \
    --dev-url "docker://postgres/15/test?search_path=public"
```

We can configure a local environment in `atlas-config.hcl`:

```hcl
env "local" {
  migration {
    dir = "file://ent/migrate/migrations"
    revisions_schema = "public"
  }
  src = "ent://ent/schema"
  url = "postgres://postgres:pass@localhost:5432/test?search_path=public&sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"
}
```

We can run the same command with the configuration file:

```
atlas migrate diff migration_name --config 'file://atlas-config.hcl' --env local
```

### 2. Destructive skip

https://atlasgo.io/versioned/diff#diff-policy

Let's consider a case to drop a field. (e.g. Remove `field.Time("registered_at"),` from `Car` schema `ent/schema/car.go`)

generate ent

```
GOWORK=off go generate ./ent
```

Now generate the migration files:

```
atlas migrate diff remove_registered_at_from_car --config 'file://atlas-config.hcl' --env local
```

`ent/migrate/migrations/20240408074945_remove_registered_at_from_car.sql` is generated.

Lint fails with the following error:

```
atlas migrate lint --config 'file://atlas-config.hcl' --env local --latest 1
Analyzing changes from version 20230910044615 to 20240408074945 (1 migration in total):

  -- analyzing version 20240408074945
    -- destructive changes detected:
      -- L2: Dropping non-virtual column "registered_at"
         https://atlasgo.io/lint/analyzers#DS103
    -- suggested fix:
      -> Add a pre-migration check to ensure column "registered_at" is NULL before
         dropping it
  -- ok (254.541µs)

  -------------------------
  -- 146.386084ms
  -- 1 version with errors
  -- 1 schema change
  -- 1 diagnostic
```

This is safe to prevent unintended drop. But if you are sure, you can skip this check with the configuration file and `--var "destructive=true"`:

```hcl
variable "destructive" {
  type    = bool
  default = false
}

env "local" {
  ...
  diff {
    skip {
      drop_schema = !var.destructive
      drop_table  = !var.destructive
    }
  }
}
```

```
atlas migrate lint --config 'file://atlas-config.hcl' --env local --latest 1 --var "destructive=true"
```

## Example

### Add new Schema

1. Create schema
    ```
    GOWORK=off go run -mod=mod entgo.io/ent/cmd/ent new Animal
    ```
1. Add fields/edges (manual) `ent/schema/animal.go`

    ```go
    func (Animal) Fields() []ent.Field {
      return []ent.Field{
        field.String("species"),
        field.Int("age"),
        field.String("name"),
      }
    }
    ```
1. Generate ent

    ```
    GOWORK=off go generate ./ent
    ```

1. Generate migartiion files

    ```
    atlas migrate diff add_animal --config 'file://atlas-config.hcl' --env local
    ```

1. Apply migration files

    ```
    atlas migrate apply --config 'file://atlas-config.hcl' --env local
    ```

1. Remove schema

    ```
    rm ent/schema/animal.go
    ```

1. Generate ent

    ```
    GOWORK=off go generate ./ent
    ```


## FAQ

1. Why is `atlas_schema_revisions` dropped?

## Tips: commands

1. `atlas migrate validate --dir file://ent/migrate/migrations`:
1. `atlas migrate hash --dir file://ent/migrate/migrations`
