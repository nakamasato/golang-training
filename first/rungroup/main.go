package main

import (
	"net"
	"net/http"
	"os"
	"fmt"

	"github.com/oklog/run"
)

func main() {
	var g run.Group
	ln, _ := net.Listen("tcp", ":8080")
	g.Add(func() error {
		return http.Serve(ln, nil)
	}, func(error) {
		ln.Close()
	})
	// go run main.go したらnetcatで確認できる (errorのケースは確認できてない)
	// nc -vz localhost 8080
	// Connection to localhost port 8080 [tcp/http-alt] succeeded!
	if err := g.Run(); err != nil {
		fmt.Printf("failed to Run")
		os.Exit(1)
	}
}
