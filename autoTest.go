package main

import (
	"flag"
	"fmt"
	"time"
	//"io/ioutil"
	"xmltest/btalkTest/TestEngine"
	"os"
	"path/filepath"
)

const (
	DEFAULT_EXTENSION = "xml"
	RESET_ITERATION = 1000
	SLEEP_BETWEEN_TEST = 100 // milisecond
	TIME_OUT_FOR_READ = 3
)

var timeOut int

// visit a file path
// if the file is a .xml (test file), call TestEngine to run the test
func visit(currentPath string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", currentPath)
	if !f.IsDir() {
		var fileName = f.Name()
		var strLen = len(fileName)
		if (strLen >= len(DEFAULT_EXTENSION) && 
			fileName[strLen-len(DEFAULT_EXTENSION):strLen] == DEFAULT_EXTENSION) {
			// run file with extension "xml"
			
			TestEngine.ExecuteTest(currentPath, timeOut)
		}
	}
	return err
}

func main() {
	var inputDir string
	flag.StringVar(&inputDir, "inputDir", "TestCase", "The inputDir to be parsed")
	flag.IntVar(&timeOut, "timeOut", TIME_OUT_FOR_READ, "The time (s) to wait when trying to read from server")
	flag.Parse()
	fmt.Println("Input directory: ", inputDir)
		
	iteration := 0
	for {
		if (iteration == RESET_ITERATION) {
			iteration = 0
		}
		iteration++
		fmt.Println("ITERATION NUMBER:", iteration)

		root := inputDir
		err := filepath.Walk(root, visit)
		if (err != nil) {
			fmt.Println(err)
		}

		time.Sleep(SLEEP_BETWEEN_TEST * time.Millisecond)
	}
}