package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"
const countDownStart = 3

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


func Countdown(writer io.Writer, sleeper Sleeper)  {
	for i:=countDownStart; i>0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(writer, i)
	}
	sleeper.Sleep()
	fmt.Fprint(writer, finalWord)
}

func main() {
	// sleeper := &DefaultSleeper{}
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
