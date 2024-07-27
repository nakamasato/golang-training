package main

import (
	"net/http"
	"strings"
)

// CustomAdapter is a type alias for a function that wraps an http.Handler.
type CustomAdapter func(http.Handler) http.Handler

// ModifyPathForServiceA is a CustomAdapter that modifies the URL path.
func ModifyPathForServiceA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the URL path starts with /servicea/query/
		if strings.HasPrefix(r.URL.Path, "/servicea/query/") {
			// Modify the URL path
			r.URL.Path = strings.Replace(r.URL.Path, "/servicea/query", "/query", 1)
		}
		// Call the next handler with the modified request
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Your final handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, world! Path" + r.URL.Path)); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Wrap the final handler with the ModifyPath adapter
	http.Handle("/", ModifyPathForServiceA(finalHandler))

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
