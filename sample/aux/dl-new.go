package main

import "fmt"
import "time"

func test_a(test_channel chan int) {
  test_channel <- 1
  return
}

func test() {
  test_channel := make(chan int)
  for i := 0; i < 10; i++ {
    go test_a(test_channel)
  }
  time.Sleep(time.Millisecond)
  for {
    fmt.Println(<-test_channel)
  }
}
func main() {
  test()
}
