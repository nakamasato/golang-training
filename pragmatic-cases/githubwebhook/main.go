package main

import (
	"fmt"
	"os"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
)

const (
	path = "/webhooks"
)

func main() {
	hook, _ := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				fmt.Printf("ErrEventNotFound: %v", err)
			}
		}
		fmt.Println("received event: %v", payload)
		switch payload := payload.(type) {

		case github.CheckRunPayload:
			fmt.Printf("CheckRunPayload: %+v\n", payload)

		case github.StatusPayload:
			fmt.Printf("StatusPayload: %+v\n", payload)

		case github.PullRequestPayload:
			fmt.Printf("PullRequestPayload: %+v\n", payload)
		}
	})
	http.ListenAndServe(":8080", nil)
}
