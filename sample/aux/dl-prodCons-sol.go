
package main

import(
  "time"
)
func main() {
  ch1, ch2 := make(chan int), make(chan int)
  go Producer(ch1)
  go Producer(ch2)
  Consumer(ch1, ch2)
  time.Sleep(time.Second)
}

func Producer(ch chan int){
  for i := 0; i<5; i++ {
    ch <- i
  }
  defer close(ch)
}

func Consumer(ch1, ch2 chan int){
  for i := range ch1{
    print (i)
  }
  for i := range ch2{
    print (i)
  }
}
