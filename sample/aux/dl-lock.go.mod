package main

import (
	"time"
	"runtime/trace"
	"os"
	"sync"
)

func main() {
	trace.Start(os.Stderr)
	defer func() {
		time.Sleep(50 * time.Millisecond)
		trace.Stop()
	}()
	ch1 := make(chan int)
	var m sync.Mutex

	go func() {
		time.Sleep(time.Millisecond)
		m.Lock()

		ch1 <- 1

		m.Unlock()

	}()

	go func() {
		m.Lock()
		m.Unlock()

		<-ch1

	}()

	time.Sleep(time.Second)
}
