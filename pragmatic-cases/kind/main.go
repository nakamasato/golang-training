package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kindName := "test"
	defer deleteKindCluster(ctx, kindName)
	err := deleteKindCluster(ctx, kindName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createKindCluster(ctx, kindName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// select {
	// case <-ctx.Done():
	// 	fmt.Println("ctx canceled")
	// 	err := deleteKindCluster(ctx, kindName)
	// 	if err != nil {
	// 		fmt.Println("failed to delete cluster")
	// 	}
	// }
}

func createKindCluster(ctx context.Context, kindName string) error {
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
	// cmd.Path = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("start creating kind cluster")
	return cmd.Run()
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
