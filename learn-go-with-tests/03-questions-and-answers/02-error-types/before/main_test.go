package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

		svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusTeapot)
		}))
		defer svr.Close()

		_, err := DumbGetter(svr.URL)

		if err == nil {
			t.Fatal("expected an error")
		}

		want := fmt.Sprintf("did not get 200 from %s, got %d", svr.URL, http.StatusTeapot)
		got := err.Error()

		if got != want {
			t.Errorf(`got "%v", want "%v"`, got, want)
		}
	})
}
