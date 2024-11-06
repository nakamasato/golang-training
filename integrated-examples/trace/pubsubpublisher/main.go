package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"cloud.google.com/go/pubsub"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

// https://cloud.google.com/stackdriver/docs/instrumentation/setup/go
func setupOpenTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown combines shutdown functions from multiple OpenTelemetry
	// components into a single function.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Configure Context Propagation to use the default W3C traceparent format
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	// Option1: OpenTelemetry Google Cloud Trace Exporter
	// https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/main/exporter/trace/README.md
    texporter, err := texporter.New()
    if err != nil {
        log.Fatalf("unable to set up tracing: %v", err)
    }

	// Option2: Configure Trace Export to send spans as OTLP
	// texporter, err := autoexport.NewSpanExporter(ctx)
	// if err != nil {
	// 	err = errors.Join(err, shutdown(ctx))
	// 	return
	// }

	tp := trace.NewTracerProvider(trace.WithBatcher(texporter))
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
	otel.SetTracerProvider(tp)

	// Configure Metric Export to send metrics as OTLP
	mreader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		err = errors.Join(err, shutdown(ctx))
		return
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(mreader),
	)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
	otel.SetMeterProvider(mp)

	return shutdown, nil
}

func main() {
	ctx := context.Background()

	shutdown, err := setupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up OpenTelemetry", slog.Any("error", err))
		os.Exit(1)
	}
	defer shutdown(ctx)

	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("PROJECT_ID must be set")
	}
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic("helloworld")

	msg := &pubsub.Message{Data: []byte("hello world")}
	if msg.Attributes == nil {
		msg.Attributes = make(map[string]string)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))

	res := topic.Publish(ctx, msg)
	id, err := res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("published %s\n", id)
}
