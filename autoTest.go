package main

import (
	"flag"
	"fmt"
	"time"
	//"io/ioutil"
	"os"
	"path/filepath"
	"xmltest/btalkTest/TestEngine"
)

const (
	DEFAULT_EXTENSION  = "xml"
	RESET_ITERATION    = 100000
	SLEEP_BETWEEN_TEST = 100 // milisecond
	TIME_OUT_FOR_READ  = 10
)

// DEBUG FLAGS
var DEBUG_COMPARE_POINTER bool = false
var DEBUG_COMPARE_SLICE bool = false
var DEBUG_COMPARE_STRUCT bool = false

var DEBUG_CLEANING_UP bool = false
var DEBUG_PARSING_MESSAGE bool = false

var DEBUG_SENDING_MESSAGE bool = false
var DEBUG_READING_MESSAGE bool = false
var DEBUG_IGNORING_MESSAGE bool = false

var DEBUG_CALLING_FUNC bool = false
var BYPASS_CONNECTION_SERVER bool = false

var timeOut int

var flagMap map[string]bool

// visit a file path
// if the file is a .xml (test file), call TestEngine to run the test
func visit(currentPath string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", currentPath)
	if !f.IsDir() {
		var fileName = f.Name()
		var strLen = len(fileName)
		if strLen >= len(DEFAULT_EXTENSION) &&
			fileName[strLen-len(DEFAULT_EXTENSION):strLen] == DEFAULT_EXTENSION {
			// run file with extension "xml"

			TestEngine.ExecuteTest(
				flagMap,
				currentPath, timeOut)
		}
	}
	return err
}

func main() {
	var inputDir string
	var numIteration int
	var sleepTime int
	flag.StringVar(&inputDir, "inputDir", "TestCase", "The inputDir to be parsed")
	flag.IntVar(&numIteration, "numItr", -1, "The number of iteration to run. Negative means infinite")

	flag.BoolVar(&DEBUG_COMPARE_POINTER, "cmp_ptr", false, "Debug compare pointer")
	flag.BoolVar(&DEBUG_COMPARE_SLICE, "cmp_slc", false, "Debug compare slice")
	flag.BoolVar(&DEBUG_COMPARE_STRUCT, "cmp_str", false, "Debug compare struct")

	flag.BoolVar(&DEBUG_CLEANING_UP, "clean", false, "Debug clean up")
	flag.BoolVar(&DEBUG_PARSING_MESSAGE, "parse", false, "Debug parse messages")

	flag.BoolVar(&DEBUG_READING_MESSAGE, "read", false, "Debug read messages")
	flag.BoolVar(&DEBUG_SENDING_MESSAGE, "send", false, "Debug send messages")
	flag.BoolVar(&DEBUG_IGNORING_MESSAGE, "ignore", false, "Debug ignore messages")

	flag.BoolVar(&DEBUG_CALLING_FUNC, "func", false, "Debug calling function")
	flag.BoolVar(&BYPASS_CONNECTION_SERVER, "bypass", false, "Bypass connection server")

	flag.IntVar(&timeOut, "timeOut", TIME_OUT_FOR_READ, "The time (s) to wait when trying to read from server")
	flag.IntVar(&sleepTime, "sleepTime", SLEEP_BETWEEN_TEST, "The time (ms) to wait between subsequent iterations")

	flag.Parse()
	flagMap = make(map[string]bool)
	flagMap["cmpPtr"] = DEBUG_COMPARE_POINTER
	flagMap["cmpSlc"] = DEBUG_COMPARE_SLICE
	flagMap["cmpStr"] = DEBUG_COMPARE_STRUCT

	flagMap["clean"] = DEBUG_CLEANING_UP
	flagMap["parse"] = DEBUG_PARSING_MESSAGE

	flagMap["send"] = DEBUG_SENDING_MESSAGE
	flagMap["read"] = DEBUG_READING_MESSAGE
	flagMap["ignore"] = DEBUG_IGNORING_MESSAGE

	flagMap["func"] = DEBUG_CALLING_FUNC
	flagMap["bypass"] = BYPASS_CONNECTION_SERVER

	fmt.Println("List of debug flag:", flagMap)

	fmt.Println("Input directory: ", inputDir)

	iteration := 0
	for {
		if iteration == numIteration {
			break
		}

		iteration++
		fmt.Println("ITERATION NUMBER:", iteration)

		root := inputDir
		err := filepath.Walk(root, visit)
		if err != nil {
			fmt.Println(err)
		}

		sleepDuration, _ := time.ParseDuration(fmt.Sprintf("%dms", sleepTime))
		time.Sleep(sleepDuration)
	}
}
