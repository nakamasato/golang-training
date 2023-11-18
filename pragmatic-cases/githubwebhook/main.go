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
	hook, err := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.CheckRunEvent, github.StatusEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				fmt.Printf("ErrEventNotFound: %v\n", err)
			}
		}
		switch payload := payload.(type) {

		case github.CheckRunPayload:
			fmt.Printf("CheckRunPayload(Action: %s, Name: %s, Status: %s)\n", payload.Action, payload.CheckRun.Name, payload.CheckRun.Status)

		case github.StatusPayload:
			fmt.Printf("StatusPayload(State: %s, sha: %s)\n", payload.State, payload.Commit.Sha)

		case github.PullRequestPayload:
			var lables []string
			for _, label := range payload.PullRequest.Labels {
				lables = append(lables, label.Name)
			}

			fmt.Printf("PullRequestPayload(Action: %s, PR: %d, URL: %s, Labels: %q)\n",
				payload.Action, payload.Number, payload.PullRequest.URL, lables)
		default:
			fmt.Printf("no action is defined for event: %v\n", payload)
		}
	})
	http.ListenAndServe(":8080", nil)
}
