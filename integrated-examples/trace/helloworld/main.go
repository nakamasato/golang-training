// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

type server struct {
	cli   *cloudtasks.Client
	queue *queue
}

type queue struct {
	name string
	url  string
}

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	// OpenTelemetry setup
	shutdown, err := setupOpenTelemetry(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer shutdown(ctx)

	// Cloud Tasks client
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	log.Print("starting server...")
	q := NewQueue()
	srv := &server{c, q}

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	router.Handle("/cloudtask", otelhttp.NewHandler(http.HandlerFunc(s.createCloudTaskHandler), "createCloudTaskHandler"))
	router.Handle("/helloworld", otelhttp.NewHandler(http.HandlerFunc(helloHandler), "helloHandler"))
	router.Handle("/", http.HandlerFunc(helloHandler))
	router.ServeHTTP(w, r)
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message struct {
		Data       []byte            `json:"data,omitempty"`
		Attributes map[string]string `json:"attributes"`
		ID         string            `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// Create a task in queue.name that sends an HTTP request to queue.url
// This is called by pubsub push subscription
func (s *server) createCloudTaskHandler(w http.ResponseWriter, r *http.Request) {
	// ctx, span := tracer.Start(r.Context(), "createCloudTask")
	// defer span.End()

	var m PubSubMessage
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("io.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// byte slice unmarshalling handles base64 decoding.
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Pub/Sub メッセージの attributes から trace 情報を抽出してコンテキストに設定
	carrier := propagation.MapCarrier(m.Message.Attributes)
	for k, v := range carrier {
		// traceparent is stored in googclient_traceparent
		fmt.Printf("carrier key: %s, value: %s\n", k, v)
	}
	for k, v := range r.Header {
		// traceparent and tracestate is stored in the header from Cloud Tasks
		fmt.Printf("header key: %s, value: %s\n", k, v)
	}

	// r.Header.Set("traceparent", carrier.Get("googclient_traceparent")) // Update traceparent with the value from Pub/Sub message

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	if s.queue == nil {
		fmt.Println("queue is empty")
		fmt.Fprint(w, "skipped creating Cloud Tasks task")
		return
	}

	fmt.Println(r.Header.Get("traceparent"))
	fmt.Println(r.Header.Get("tracestate"))
	fmt.Println(s.queue.url)
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: s.queue.name,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					Url:        s.queue.url,
					HttpMethod: cloudtaskspb.HttpMethod_GET,
					Headers: map[string]string{
						// "traceparent": r.Header.Get("traceparent"), // これにするとCloud Run /cloudtask から先がつながるが、 Pubsub Publishがつながらない
						// "tracestate":  r.Header.Get("tracestate"), // 上とセット
						// "traceparent": carrier.Get("googclient_traceparent"), // これはCloud PubSubの標準Tracingを使う場合で、 pubsub publishはつながるが Cloud Run /cloudtaskが繋がらない <- Cloud RunがRoot Spanになってしまう
						"traceparent": carrier.Get("traceparent"), // これはCloud PubsubPublish側でCustom SpanをRoot Spanにするタイプで、pubsub publishはつながるが Cloud Run /cloudtaskが繋がらない <- Cloud RunがRoot Spanになってしまう
					},
					Body:                []byte{},
					AuthorizationHeader: nil,
				},
			},
		},
		ResponseView: 0,
	}
	resp, err := s.cli.CreateTask(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	fmt.Fprint(w, "created Cloud Tasks task")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func NewQueue() *queue {
	projectID := os.Getenv("PROJECT_ID")
	locationID := os.Getenv("LOCATION_ID")
	queueID := os.Getenv("QUEUE_ID")
	targetURL := os.Getenv("CLOUD_TASK_TARGET_URL")
	if projectID == "" || locationID == "" || queueID == "" || targetURL == "" {
		return nil
	}
	return &queue{
		fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID),
		targetURL,
	}
}

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
	texporter, err := texporter.New(
		// Disable spans created by the exporter.
		texporter.WithTraceClientOptions(
			[]option.ClientOption{option.WithTelemetryDisabled()},
		),
	)
	if err != nil {
		log.Fatalf("unable to set up tracing: %v", err)
	}

	// Option2: Configure Trace Export to send spans as OTLP
	// texporter, err := autoexport.NewSpanExporter(ctx)
	// if err != nil {
	// 	err = errors.Join(err, shutdown(ctx))
	// 	return
	// }

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(texporter))
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
