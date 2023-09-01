package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type Account struct {
	UserID    int
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
	LastLogin time.Time
}

func (a *Account) String() string {
	return fmt.Sprintf("uid:%d, username: %s", a.UserID, a.Username)
}

type server struct {
	db *sql.DB
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(handler))
	router.Handle("/get", http.HandlerFunc(s.getHandler))
	router.ServeHTTP(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func connectWithConnectorIAMAuthN() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser                 = mustGetenv("DB_IAM_USER")              // e.g. 'service-account-name@project-id.iam'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}

func connectWithConnector() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep passwords and other secrets safe.
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}

func connect() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}
	var dsn string
	if os.Getenv("CLOUD_SQL_WITH_IAM_AUTH") == "true" {
		// https://cloud.google.com/sql/docs/postgres/samples/cloud-sql-postgres-databasesql-auto-iam-authn
		return connectWithConnectorIAMAuthN()
	} else if os.Getenv("CLOUD_SQL_WITH_BUILT_IN_USER") == "true" {
		// https://cloud.google.com/sql/docs/postgres/connect-run#go
		return connectWithConnector()
	} else { // local postgres or with SQL Auth Proxy
		dsn = fmt.Sprintf("host=%s user=%s password=%s database=%s",
			mustGetenv("DB_HOST"),
			mustGetenv("DB_USER"),
			mustGetenv("DB_PASS"),
			mustGetenv("DB_NAME"),
		)

		pgxCfg, err := pgx.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed %v", err)
		}
		db, err := sql.Open("pgx", stdlib.RegisterConnConfig(pgxCfg))
		if err != nil {
			return nil, fmt.Errorf("failed to open database %v", err)
		}
		return db, nil
	}
}

func main() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := &server{db}

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}

func (s *server) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting!\n")
	if accounts, err := getRows(s.db); err != nil {
		fmt.Printf("getRows failed %v\n", err)
	} else {
		fmt.Fprintf(w, "Got accounts: %s\n", accounts)
	}
}

func getRows(db *sql.DB) ([]*Account, error) {
	fmt.Println("getRows")
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		u := &Account{}
		if err := rows.Scan(&u.UserID, &u.Username, &u.Password, &u.Email, &u.CreatedOn, &u.LastLogin); err != nil {
			return nil, err
		}
		accounts = append(accounts, u)
	}
	return accounts, nil
}
