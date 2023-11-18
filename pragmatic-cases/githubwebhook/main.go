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
	sugar := logger.Sugar()
	hook, err := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		sugar.Errorw("failed to create webhook", zap.Error(err))
		os.Exit(1)
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.CheckRunEvent, github.StatusEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				sugar.Errorw("failed to parse webhook", zap.Error(err))
			}
		}
		switch payload := payload.(type) {

		case github.CheckRunPayload:
			sugar.Info("CheckRunPayload", zap.String("Action", payload.Action), zap.String("Name", payload.CheckRun.Name), zap.String("Status", payload.CheckRun.Status))

		case github.StatusPayload:
			sugar.Info("StatusPayload", zap.String("State", payload.State), zap.String("sha", payload.Commit.Sha))

		case github.PullRequestPayload:
			var lables []string
			for _, label := range payload.PullRequest.Labels {
				lables = append(lables, label.Name)
			}
			sugar.Info("PullRequestPayload", zap.String("Action", payload.Action), zap.Int64("PR", payload.Number), zap.String("URL", payload.PullRequest.URL), zap.Strings("Labels", lables))
		default:
			sugar.Error("no action is defined for event", zap.Any("payload", payload))
		}
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		sugar.Errorw("failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
