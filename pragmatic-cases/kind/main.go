package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kindName := "test"
	err := deleteKindCluster(ctx, kindName)
	if err != nil {
		return
	}

	cmd := exec.CommandContext(
		ctx,
		"kind",
		"create",
		"cluster",
		"--name",
		kindName,
		"--kubeconfig",
		"./.kubeconfig",
	)
	data := make(chan struct{})
	go func() {
		fmt.Println("start creating kind cluster")
		if err := cmd.Run(); err != nil {
			fmt.Println("failed to create kind cluster")
		}
		fmt.Println("successfully created kind cluster")
		data <- struct{}{}
	}()

	ticker := time.NewTicker(5 * time.Second)
	out:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("kind create cluster failed")
			err = deleteKindCluster(ctx, kindName)
			if err != nil {
				return
			}
			break out
		case <-data:
			fmt.Println("kind create cluster succeeded")
			break out
		case t := <-ticker.C:
			cmd := exec.CommandContext(
				ctx,
				"kind",
				"get",
				"clusters",
			)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				fmt.Println("faied to get kind clusters")
			}
			if strings.Contains(out.String(), kindName) {
				fmt.Printf("kind cluster '%s' exists %s\n", kindName, t)
			} else {
				fmt.Printf("kind cluster '%s' not exists %s\n", kindName, t)
			}
		}
	}
}

func deleteKindCluster(ctx context.Context, kindName string) error {
	cmd := exec.CommandContext(
		ctx,
		"kind",
		"delete",
		"cluster",
		"--name",
		kindName,
	)
	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to delete kind cluster")
		return err
	}
	fmt.Println("successfully deleted kind cluster")
	return nil
}
