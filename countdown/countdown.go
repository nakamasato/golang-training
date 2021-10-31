package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Sleeper interface {
    Sleep()
}

type DefaultSleeper struct {}

func (d *DefaultSleeper) Sleep() {
    time.Sleep(1 * time.Second)
}

type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}


func main() {
	// sleeper := &DefaultSleeper{}
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
// In main we will send to os.Stdout so our users see the countdown printed to the terminal.
// In test we will send to bytes.Buffer so our tests can capture what data is being generated.

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := 3; i>0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(out, i)
	}
	sleeper.Sleep()
	fmt.Fprint(out, "Go!")
}
