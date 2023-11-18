package main

import (
	"os"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"
)

const (
	path = "/webhooks"
)

func main() {
	logger := zap.NewExample()
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()
	hook, err := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		logger.Error("failed to create webhook", zap.Error(err))
		os.Exit(1)
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.CheckRunEvent, github.StatusEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				logger.Error("failed to parse webhook", zap.Error(err))
			}
		}

		switch payload := payload.(type) {

		case github.CheckRunPayload:
			logger.Info("CheckRunPayload", zap.String("Action", payload.Action), zap.String("Repo", payload.Repository.Name), zap.String("sha", payload.CheckRun.HeadSHA), zap.String("RequestURI", r.RequestURI))

		case github.StatusPayload:
			logger.Info("StatusPayload", zap.String("Repo", payload.Repository.Name), zap.String("State", payload.State), zap.String("sha", payload.Commit.Sha), zap.String("RequestURI", r.RequestURI))

		case github.PullRequestPayload:
			var labels []string
			for _, label := range payload.PullRequest.Labels {
				labels = append(labels, label.Name)
			}
			logger.Info("PullRequestPayload", zap.String("Action", payload.Action), zap.String("Repo", payload.Repository.Name), zap.Int64("PR", payload.Number), zap.Strings("Labels", labels), zap.String("RequestURI", r.RequestURI))
		default:
			logger.Error("no action is defined for event", zap.Any("payload", payload))
		}
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
