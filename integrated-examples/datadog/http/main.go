package main

import (
	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start()
	defer tracer.Stop()
	mux := httptrace.NewServeMux(httptrace.WithServiceName("my-service"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// nolint
		w.Write([]byte("Hello World!\n"))
	})
	// nolint
	http.ListenAndServe(":8080", mux)
}
