package main

import (
  //"runtime"
  "time"
  "sync"
  //"fmt"
)

func main() {
  ch1 := make(chan int)
  var m sync.Mutex

  // goroutine 1
  go func() {
    time.Sleep(time.Millisecond)
    m.Lock()
    //runtime.Gosched()
    ch1 <- 1 // block here
    //runtime.Gosched()
    m.Unlock()
    //runtime.Gosched()
  }()

  // goroutine 2
  go func() {
    m.Lock() // block here
    m.Unlock()
    //runtime.Gosched()
    <- ch1
    //runtime.Gosched()
  }()

  //time.Sleep(2*time.Second)
  time.Sleep(time.Second)
}
