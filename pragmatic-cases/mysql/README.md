# mysql

Check if `root` user exists by golang.

## Prerequisite

- Docker

## Check

```
make all
```

This command does the following:
1. Create MySQL container.
1. Run `main.go` to check if `root` user exists.
1. Delete MySQL container.

## Step by step

1. Run mysql with Docker.

    ```
    container_id=$(docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --rm mysql:5.7)
    ```

1. Run `main.go`.

    ```
    go run main.go
    root user exists: true
    ```

1. Clean up.

    ```
    docker rm -f $container_id
    ```

1. Run test
    ```
    go test .
    ```

## Ref
1. https://github.com/DATA-DOG/go-sqlmock
