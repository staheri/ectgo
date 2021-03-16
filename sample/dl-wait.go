package main

/*import (
    "fmt"
)

func main() {
    c1 := make(chan int)
    c1 <- 1
    go func() {fmt.Println(<- c1)}()
    //c1 <- 2
    //go func() {fmt.Println(<- c1)}()
    //c1 <- 3
    //go func() {fmt.Println(<- c1)}()
    //go func() {fmt.Println(<- c1)}()
}*/

import (
  "time"
  "sync"
  "fmt"
)

func main() {
    var group sync.WaitGroup
    var t = []int{1, 2, 3, 4}
    group.Add(len(t))
    for _,p := range t{
      go func(p int){
        fmt.Println(p)
        defer group.Done()
      }(p)
      group.Wait()
    }
    time.Sleep(time.Second)
}
