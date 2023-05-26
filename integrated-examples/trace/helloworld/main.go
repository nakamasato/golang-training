// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type server struct {
	cli   *cloudtasks.Client
	queue *queue
}

type queue struct {
	taskName string
	url      string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	router.Handle("/cloudtask", http.HandlerFunc(s.createCloudTask))
	router.Handle("/helloworld", http.HandlerFunc(handler))
	router.Handle("/", http.HandlerFunc(handler))
	router.ServeHTTP(w, r)
}

func (s *server) createCloudTask(w http.ResponseWriter, r *http.Request) {
	if s.queue == nil {
		fmt.Println("queue is empty")
		fmt.Fprint(w, "skipped creating Cloud Tasks task")
		return
	}
	fmt.Println(s.queue.url)
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: s.queue.taskName,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					Url:        s.queue.url,
					HttpMethod: 0,
					Headers:    map[string]string{
						// traceparent header
						// tracestate header
					},
					Body:                []byte{},
					AuthorizationHeader: nil,
				},
			},
		},
		ResponseView: 0,
	}
	resp, err := s.cli.CreateTask(r.Context(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	fmt.Fprint(w, "created Cloud Tasks task")
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

func main() {
	ctx := context.Background()

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

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}
