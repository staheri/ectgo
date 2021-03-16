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
    case res:= <- ch1:
        fmt.Println(res)
      default:
    }
    m.Unlock()
  }()

  // goroutine 2
  go func() {
    m.Lock() // block here
    m.Unlock()
    ch1 <- 1

  }()
  time.Sleep(time.Second)
}
