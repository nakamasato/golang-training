package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	service     = "httptrace-demo"
	environment = "development"
)

func NewJaegerTracerProvider(service, environment, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
		)),
	)
	return tp, nil
}

func main() {
	ctx := context.Background()
	tp, err := NewJaegerTracerProvider(service, environment, "http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp) // register tp as the global trace provider

	clientTrace := otelhttptrace.NewClientTrace(ctx)
	ctx = httptrace.WithClientTrace(ctx, clientTrace)
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://journal.lampetty.net/", nil)

	_, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
