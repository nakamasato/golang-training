package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kindName := "test"
	isLazy := true

	if !isLazy && checkKindCluster(ctx, kindName) {
		err := deleteKindCluster(ctx, kindName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if !(isLazy && checkKindCluster(ctx, kindName)) {
		err := createKindCluster(ctx, kindName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("do something")

	if !isLazy && checkKindCluster(ctx, kindName) {
		err := deleteKindCluster(ctx, kindName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
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

func checkKindCluster(ctx context.Context, kindName string) bool {
	out, err := exec.CommandContext(
		ctx,
		"kind",
		"get",
		"clusters",
	).Output()

	if err != nil {
		return false
	}
	clusters := strings.Split(string(out), "\n")
	return stringInSlice(kindName, clusters)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
