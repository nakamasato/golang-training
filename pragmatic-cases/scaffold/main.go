package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("main")
	exec.Command(
		"skaffold",
		"dev",
	)
}
