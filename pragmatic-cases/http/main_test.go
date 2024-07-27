package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Send a request with /servicea/query/123 and check if the path is modified
func TestModifyPathForServiceA(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/servicea/query/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rec := httptest.NewRecorder()

	// Create a new handler by wrapping the final handler with the ModifyPath adapter
	handler := ModifyPathForServiceA(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world! Path" + r.URL.Path))
	}))

	// Serve the HTTP request to the handler
	handler.ServeHTTP(rec, req)

	// Check the status code
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, world! Path/query/123"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

// Send a request with /serviceb/query/123 and confirm the path is not changed
func TestModifyPathForServiceB(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/serviceb/query/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rec := httptest.NewRecorder()

	// Create a new handler by wrapping the final handler with the ModifyPath adapter
	handler := ModifyPathForServiceA(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world! Path" + r.URL.Path))
	}))

	// Serve the HTTP request to the handler
	handler.ServeHTTP(rec, req)

	// Check the status code
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, world! Path/serviceb/query/123"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
