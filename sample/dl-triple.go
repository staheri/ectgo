package main

import (
    "fmt"
    "sync"
)

var number = make(chan int)
var mutex = &sync.Mutex{}

func worker(wg *sync.WaitGroup, id int) {
    defer wg.Done()

    mutex.Lock()
    number <- id + <-number
    mutex.Unlock()
}

func main() {
    var wg sync.WaitGroup
    number <- 0
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(&wg, i)
    }

    wg.Wait()
    fmt.Println(<-number) // expected output: 0+1+2+3+4 = 10
}

