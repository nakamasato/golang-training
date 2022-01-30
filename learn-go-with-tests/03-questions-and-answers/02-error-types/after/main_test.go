package main

import (
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

		got, isStatusErr := err.(BadStatusError)

		if !isStatusErr {
			t.Fatalf("was not a BadStatusError, got %T", err)
		}

		want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
