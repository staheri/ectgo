package main

import (
  "time"
  "sync"
  "fmt"
  "runtime"
)

func main() {
  // Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()

  ch1 := make(chan int)
  var m sync.Mutex

  // goroutine 1
  go func() {
    m.Lock()
    ch1 <- 1 // block here
    m.Unlock()
  }()

  // goroutine 2
  go func() {
    m.Lock() // block here
    runtime.Gosched()
    m.Unlock()
    fmt.Println(<- ch1)
  }()

  time.Sleep(time.Second)
  // Capture ending number of goroutines.
  endingGs := runtime.NumGoroutine()

  // Report the results.
  fmt.Println("========================================")
  fmt.Println("Number of goroutines before:", startingGs)
  fmt.Println("Number of goroutines after :", endingGs)
  fmt.Println("Number of goroutines leaked:", endingGs-startingGs)
}
