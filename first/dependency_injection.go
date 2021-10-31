package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// io.Writer which we know both os.Stdout and bytes.Buffer implement.
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
    Greet(w, "world")
}

func main() {
    // Greet(os.Stdout, "Elodie")
	log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler)))
}

// Printf -> Fprintf(writer, format, a...) -> p := newPrinter() w.Write(p.buf)

// It returns the number of bytes written and any write error encountered.
// func Printf(format string, a ...interface{}) (n int, err error) {
//     return Fprintf(os.Stdout, format, a...)
// }
