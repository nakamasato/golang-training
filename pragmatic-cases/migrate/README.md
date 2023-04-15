# [migrate](https://github.com/golang-migrate/migrate)

## Concept

1. Source:
	```go
	m, err := migrate.NewWithDatabaseInstance(
		"github://nakamasato/golang-training/pragmatic-cases/migrate/postgres/migrations#main", // sourceUrl
		"postgres", //databaseName
		driver, // databaseInstance
	)
	```

	1. file: `"file://migrations"`
	1. github: `"github://<owner>/<repo>/path#ref"`
	1. ...
1. Database
	1. postgres
	1. mysql
	1. ...
## Getting Started (postgres)

![](postgres/README.md)

## Getting Started (mysql)

![](mysql/README.md)
