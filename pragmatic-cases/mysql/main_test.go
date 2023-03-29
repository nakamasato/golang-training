package main_test

import (
	"testing"
	. "tmp/pragmatic-cases/mysql"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMain(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT .* FROM mysql.user where User = 'root'").WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	res, err := CheckMySQLHasUser(db, "root")
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatalf("want true but got %t\n", res)
	}
}
