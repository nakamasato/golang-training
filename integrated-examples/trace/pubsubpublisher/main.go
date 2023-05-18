package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("PROJECT_ID must be set")
	}
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic("helloworld")
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
	})
	id, err := res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("published %s\n", id)
}
