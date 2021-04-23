package main
import (
  "sync"
  "os"
  "strconv"
)

func main() {
  n,err := strconv.Atoi(os.Args[1])
  if err!= nil{
    panic(err)
  }

  container := &Container{stop:make(chan struct{})}
  go Monitor(container)
  for i :=0; i<n;i++{
    go StatusChange(container)
  }
}

type Container struct{
  sync.Mutex
  stop  chan struct{}
}

func Monitor(cnt *Container){
  for{
    select{
    case <- cnt.stop:
      return
    default:
      cnt.Lock()
      cnt.Unlock()
    }
  }
}

func StatusChange(cnt *Container){
  cnt.Lock()
  defer cnt.Unlock()
  cnt.stop <- struct{}{}
}
