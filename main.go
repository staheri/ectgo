package main

import (
	"db"
	"flag"
	"fmt"
	"instrument"
	"os"
	"schedtest"
	"strings"
	"util"
	"log"
	"path"
)

const WORD_CHUNK_LENGTH = 11

var CLOUTPATH = os.Getenv("GOPATH") + "/traces/clx"

var (
	flagCmd, flagOut, flagSrc, flagX, flagBase, flagApp, dbName, flagDB  string
	flagCons, flagAtrMode, flagN, flagTO, flagDepth, flagIter            int
	flagArgs                                                             []string
	flagMT, flagDebug                                                    bool
	validCategories    = []string{"CHNL", "GCMM", "GRTN", "MISC", "MUTX", "PROC", "SYSC", "WGCV", "SCHD", "BLCK"}
	validPrimeCmds     = []string{"word", "hac", "rr", "diff", "dineData", "cleanDB", "dev", "hb", "gtree", "cgraph", "resg","leakChecker"}
	validTestSchedCmds = []string{"test","execVis"}
	validSrc           = []string{"native", "x", "latest", "schedTest"}
)

func main() {
	fmt.Println("Initializing GOAT V.0.1 ...")
	parseFlags()

	if flagSrc == "schedTest" {

		fmt.Println("GOAT SchedTest mode ...")
		handleSchedTestCommands()

	} else {

		fmt.Println("GOAT Prime mode ...")
		// New App instance
		myapp := instrument.NewAppExec(flagApp, flagSrc, flagX, flagTO)
		// Obtain DB


		dbn, err := myapp.DBPointer()
		if err != nil {
			panic(err)
		}

		fmt.Println("Working DB: ",dbn)
		myapp.DBName = dbn
		handlePrimaryCommands(myapp.DBName)
		//fmt.Println(myapp.ToString())
	}
}

// Parse flags, execute app & store traces (if necessary), return app database handler
func parseFlags() {
	srcDescription := "native: execute the app and collect from scratch, latest: retrieve data from latest execution, x: retrieve data from specific execution (requires -x option)"
	// Parse flags
	flag.StringVar(&flagCmd, "cmd", "", "Commands: word, cl, rr, rg, diff")
	flag.StringVar(&flagBase, "baseX", "0", "Base execution for \"diff\" or \"schedTrace\" command")
	flag.StringVar(&flagOut, "outdir", "", "Output directory to write words and/or reports")
	flag.StringVar(&flagSrc, "src", "latest", srcDescription)
	flag.StringVar(&flagX, "x", "", "Execution version stored in database")
	flag.IntVar(&flagN, "n", 0, "Number of philosophers for dineData command")
	flag.IntVar(&flagCons, "cons", 1, "Number of consecutive elements for HAC & DIFF")
	flag.IntVar(&flagAtrMode, "atrmode", 0, "Modes for HAC & DIFF")
	flag.StringVar(&flagApp, "app", "", "Target application (*.go)")
	flag.StringVar(&flagDB, "db", "", "Specific SQL table name")
	flag.IntVar(&flagTO, "to", 0, "Timeout for deadlocks")
	flag.IntVar(&flagDepth, "depth", 0, "Max depth for rescheduling")
	flag.IntVar(&flagIter, "iter", 2, "Testing iteration")
	flag.BoolVar(&flagMT, "mt", false, "Measure Times")
	flag.BoolVar(&flagDebug, "debug", false, "Print debugging info")

	flag.Parse()

	// Check src validity
	if !util.Contains(validSrc, flagSrc) {
		util.PrintUsage()
		panic("Wrong source")
	}

	// Check prime cmd validity
	if flagSrc != "schedTest" && !util.Contains(validPrimeCmds, flagCmd) {
		util.PrintUsage()
		fmt.Printf("flagCMD: %s\n", flagCmd)
		panic("Wrong prime command")
	}

	// Check prime cmd validity
	if flagSrc == "schedTest" && !util.Contains(validTestSchedCmds, flagCmd) {
		util.PrintUsage()
		fmt.Printf("flagCMD: %s\n", flagCmd)
		panic("Wrong schedTest command")
	}

	if flagSrc == "schedTest" && flagCmd == "execVis" && flagDB == ""{
		util.PrintUsage()
		panic("DB name required")
	}
	// Check Outdir
	if flagOut == "" {
		flagOut = path.Dir(flagApp)
	}

	// Check app
	if flagApp == "" {
		util.PrintUsage()
		panic("App required")
	}

	// Check validity of categories
	for _, arg := range flagArgs {
		tl := strings.Split(arg, ",")
		for _, e := range tl {
			if !util.Contains(validCategories, e) {
				panic("Invalid category: " + e)
			}
		}
	}

	// diff command needs a base
	if flagCmd == "diff" && flagBase == "" {
		util.PrintUsage()
		panic("Undefined base for diff command!")
	}

	// dineData command needs N
	if flagCmd == "dineData" && flagN == 0 {
		util.PrintUsage()
		panic("Wrong N for dineData!")
	}

	// x command needs X value
	if flagSrc == "x" && flagX == "" {
		util.PrintUsage()
		panic("Needs X value!")
	}
	util.MeasureTime = flagMT
	util.Debug = flagDebug

	flagArgs = flag.Args()

	file, err := os.OpenFile(path.Dir(flagApp)+"/GOAT_log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
  if err != nil {
  	log.Fatal(err)
  }
  log.SetOutput(file)

}

// handle primary commands
func handlePrimaryCommands(dbName string) {
	switch flagCmd {
	case "word":
		for _, arg := range flagArgs {
			// For now, only one filter is allowed at a time
			if len(strings.Split(arg, ",")) != 1 {
				panic("Currently more than one filter is not allowed!")
			}
			db.WordData(dbName, flagOut, arg, WORD_CHUNK_LENGTH)
			//for _,e := range(tl){
			// TODO: Make db.WriteData compatible with combination of filters
			//}
		}

	case "hac":
		if len(flagArgs) > 0 {
			for _, arg := range flagArgs {
				tl := strings.Split(arg, ",")
				db.HAC(dbName, CLOUTPATH, flagOut, flagCons, flagAtrMode, tl...)
			}
		} else {
			var emptyList []string
			db.HAC(dbName, CLOUTPATH, flagOut, flagCons, flagAtrMode, emptyList...)
		}

	case "rr":
		for _, arg := range flagArgs {
			if len(strings.Split(arg, ",")) != 1 {
				panic("For rr, only one category is allowed")
			}
			switch arg {
			case "CHNL":
				db.ChannelReport(dbName, flagOut)
			case "MUTX":
				db.MutexReport(dbName)
				db.RWMutexReport(dbName)
			case "WGRP":
				db.WaitingGroupReport(dbName)
			default:
				panic("Wrong category for rr!")
			}
		}

	case "diff":
		baseDBName := db.Ops("x", util.AppName(flagBase), "13")
		for _, arg := range flagArgs {
			tl := strings.Split(arg, ",")
			db.JointHAC(dbName, baseDBName, CLOUTPATH, flagOut, flagCons, flagAtrMode, tl...)
		}
	case "dineData":
		db.DineData(dbName, flagOut+"/ch-chid", flagN, true, true)   // channel events only + channel ID
		db.DineData(dbName, flagOut+"/ch", flagN, true, false)       // channel events only
		db.DineData(dbName, flagOut+"/all-chid", flagN, false, true) // all events + channel ID (for channel events)
		db.DineData(dbName, flagOut+"/all", flagN, false, false)     // all events
	case "cleanDB":
		db.Ops("clean all", "", "0")
	case "hb":
		fmt.Println("HB DBNAME:", dbName)
		for _, arg := range flagArgs {
			tl := strings.Split(arg, ",")
			hbtable := db.HBTable(dbName, tl...)
			db.HBLog(dbName, hbtable, flagOut, true)
			//fmt.Println("****")
			db.HBLog(dbName, hbtable, flagOut, false)
		}
	case "gtree":
		db.Gtree(dbName, flagOut)
	case "cgraph":
		db.ChannelGraph(dbName, flagOut)
	case "resg":
		db.ResourceGraph(dbName, flagOut)
	case "leakChecker":
		db.Checker(dbName,true)
	}
}

// handle schedTest commands
func handleSchedTestCommands() {
	switch flagCmd {
	case "test":
		// for measuring overhead
		schedtest.NativeRun(flagApp)
		schedtest.SchedTest(flagApp, flagSrc, flagX, flagTO, flagDepth, flagIter)
	case "execVis":
		db.ExecVis(flagDB,flagOut)
	}
}
