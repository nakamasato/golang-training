package main

import (
	"os"
	"time"

	"tmp/learn-go-with-tests/01-go-fundamentals/16-math"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
