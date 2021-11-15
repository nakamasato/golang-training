package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name) // Fprintf takes writer while Printf uses stdout as writer
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
    Greet(w, "world")
}

func main() {
	// Greet(os.Stdout, "Elodie")
	log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))) // you can open localhost:5000
}
