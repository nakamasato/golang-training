package main

import (
	"fmt"
	"sync"
)

type NaiveCouner struct {
	num int32
}

func (c *NaiveCouner) Inc(wg *sync.WaitGroup) {
	defer wg.Done()
	c.num += 1
}

type Counter struct {
	lock *sync.Cond
	num  int32
}

func (c *Counter) Inc(wg *sync.WaitGroup) {
	defer wg.Done()
	c.lock.L.Lock()
	defer c.lock.L.Unlock()
	c.num += 1
	c.lock.Signal() // Broadcast or Signal is necessary to to awaken Wait().
	// Broadcast wakes all goroutines waiting on c.
	// Signal wakes one goroutine waiting on c, if there is any.
}

func (c *Counter) Print10() {
	c.lock.L.Lock()
	for !c.condition() {
		fmt.Println("Print10: Waiting")
	    c.lock.Wait()
		// Wait atomically unlocks c.L and suspends execution of the calling goroutine.
		// After later resuming execution, Wait locks c.L before returning.
		// Unlike in other systems, Wait cannot return unless awoken by Broadcast or Signal.
		fmt.Println("Print10: Finished waiting")
	}
	fmt.Printf("Print10: num reached %d\n", c.num)
	c.lock.L.Unlock()
}

func (c *Counter) condition() bool {
	return c.num > 0 && c.num % 10 == 0
}

func main() {
	var wg sync.WaitGroup
	cnt := Counter{lock: sync.NewCond(&sync.Mutex{})}
	naiveCnt := NaiveCouner{}
	maxNum := 1000
	fmt.Printf("Increment %d times\n", maxNum)
	expectedVal := 0

	for i := 1; i <= maxNum; i++ {
		expectedVal += 1
		wg.Add(2)
		go cnt.Inc(&wg)
		go naiveCnt.Inc(&wg)
	}
	cnt.Print10()
	wg.Wait()
	fmt.Printf("Counter.num: %d, NaiveCounter.num: %d, expected: %d \n", cnt.num, naiveCnt.num, expectedVal)
}
