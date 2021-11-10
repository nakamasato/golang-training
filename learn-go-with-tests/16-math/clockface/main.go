package main
import (
    "os"
    "time"

	"tmp/learn-go-with-tests/16-math"
)

func main() {
    t := time.Now()
    clockface.SVGWriter(os.Stdout, t)
}
