package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

const maxConcurrency = 2
const batchSize = 5

func main() {

	var candidates [97]string
	for i := 0; i < 97; i++ {
		candidates[i] = fmt.Sprintf("email%d", i)
	}

	ctx := context.Background()

	target := make(chan string, batchSize) // channel
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(maxConcurrency)

		defer close(target)
		for _, email := range candidates {
			e := email
			g.Go(func() error {
				// check if sent
				if !isSent(ctx, e) {
					target <- e
				}
				return nil
			})
		}
		err := g.Wait()
		if err != nil {
			return err
		}
		return nil
	})

	targets := make([]string, 0, batchSize)
	eg.Go(func() error {
		for t := range target {
			targets = append(targets, t)
			if len(targets) == batchSize {
				err := send(ctx, targets)
				if err != nil {
					return err
				}
				targets = make([]string, 0, batchSize)
			}
		}
		err := send(ctx, targets)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}

func isSent(ctx context.Context, email string) bool {
	sleepDuration := time.Duration(rand.Intn(50)) * time.Millisecond
	time.Sleep(sleepDuration)
	fmt.Printf("is sent: %s (latency: %v)\n", email, sleepDuration)
	return false
}

func send(ctx context.Context, emails []string) error {
	sleepDuration := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(sleepDuration)
	fmt.Printf("send email: %s (latency: %v)\n", emails, sleepDuration)
	return nil
}
