# [Atlas](https://atlasgo.io/)

![](diagram.drawio.svg)

Atlas is a very interesting tool in that the status is obtained by the SQLs which are completely different for databases. e.g. [PostgreSQL](https://github.com/ariga/atlas/blob/2e0886e03c5862c54247f41f906f60d64f9c7eaf/sql/postgres/inspect.go#L728),  [MySQL](https://github.com/ariga/atlas/blob/2e0886e03c5862c54247f41f906f60d64f9c7eaf/sql/mysql/inspect.go#L631), etc.

(ref: [Programmatic inspection of databases in Go using Atlas](https://ariga.io/programmatic-inspection-of-databases-in-go-using-atlas/))

## [CLI](https://atlasgo.io/getting-started/)

1. Install

    ```
    brew install ariga/tap/atlas
    ```

1. Run mysql with docker

    ```
    docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --rm mysql:8
    ```

1. Create database `test_db`

    ```
    docker exec -it $(docker ps | grep mysql | head -1 |awk '{print $1}') mysql -uroot -ppassword -e 'create database test_db'
    ````

1. Prepare hcl

    ```hcl
    schema "test_db" {}

    table "users" {
      schema = schema.test_db
      column "id" {
        type = int
      }
      column "greeting" {
        type = text
      }
    }
    ```

1. Apply

    ```
    atlas schema apply -u "mysql://root:password@localhost:3306/test_db" --to file://schema.hcl
    ```
    ```
    -- Planned Changes:
    -- Create "users" table
    CREATE TABLE `test_db`.`users` (`id` int NOT NULL, `greeting` text NOT NULL);
    ✔ Apply
    ```

1. Inspect

    ```
    atlas schema inspect -u "mysql://root:password@localhost:3306/"
    ```
    ```
    table "users" {
      schema = schema.test_db
      column "id" {
        null = false
        type = int
      }
      column "greeting" {
        null = false
        type = text
      }
    }
    schema "test_db" {
      charset = "utf8mb4"
      collate = "utf8mb4_0900_ai_ci"
    }
    ```

## [Go](https://atlasgo.io/integrations/go-api)

1. Run mysql with docker
1. Create `test_db`
1. `InspectSchema`, `SchemaDiff`, and `ApplyChanges`

    ```
    InspectSchema: read schema from MySQL
    shema: test_db
    ----- table users ----
    col 0 name: id, type: int
    col 1 name: greeting, type: text
    EvalHCLBytes: read from schema.hcl
    SchemaDiff: compare schema objects
    no changes
    ```
