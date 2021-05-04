package main

import (
  "github.com/staheri/goat/evaluate"
	"flag"
	"fmt"
	"log"
	"os"
	//"path/filepath"
	_"bufio"
)



var (
	flagPath            string
	flagArgs            []string
	flagVerbose         bool
)

func main(){
	fmt.Println("Initializing ECTGO V.0.1 ...")

	// set log
	file, err := os.OpenFile("ECTGO_log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	check(err)
  log.SetOutput(file)

	parseFlags()
  evaluate.EvaluateOverhead(flagPath,100,[]int{1,2,4,16,64,256,512,1024,2048})
}



func parseFlags() {
	flag.StringVar(&flagPath, "path", "", "Target folder (*.go)")
	flag.BoolVar(&flagVerbose, "verb", false, "Print verbose info")

	flag.Parse()

	flagArgs = flag.Args()
}

func check(err error){
	if err != nil{
		panic(err)
	}
}
