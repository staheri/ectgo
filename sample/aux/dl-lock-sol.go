package main

import (
  "time"
  "sync"
  "fmt"
)

func main() {
  ch1 := make(chan int)
  var m sync.Mutex

  // goroutine 1
  go func() {
    m.Lock()
    select{
      case ch1 <- 1:
      default:
    }
    m.Unlock()
  }()

  // goroutine 2
  go func() {
    m.Lock() // block here
    m.Unlock()
    fmt.Println(<- ch1)
  }()
  time.Sleep(time.Second)
}
