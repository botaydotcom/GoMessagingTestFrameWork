package main

import (
	"flag"
	"fmt"
	"time"
	//"io/ioutil"
	"errors"
	"os"
	"path/filepath"
	"xmltest/btalkTest/StressTestEngine"
	"math"
	"math/rand"
	"sync"
	"log"
	"regexp"
)

const (
	DEFAULT_EXTENSION  = "xml"
	RESET_ITERATION    = 1000
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

var DEBUG_STEP bool = false

var DEBUG_FAILURE bool = true

var timeOut int

var flagMap map[string]bool

var inputTestSuites []StressTestEngine.TestSuite
var arrivalRateInMillis float64

var totalConcurrentConnection = 0

type ClientResult struct {
	IsSuccessful bool  // success or not
	ExecuteTime time.Duration
	Reason error
}

/*
	Parameters:
	- event arrival rate (for Poisson distribution)
	- dir to read test cases
	- reading time out (not supported)

	- metric to print out (not supported)
		= total concurrent connections
		= average response time
		= success rate
	- interval between report result (not supported)
*/

var resultChannel chan ClientResult
var controlChannel chan int
var finishChannel chan int
var mutex *sync.Mutex
var rate int
var testDuration time.Duration

var errorAggregation map[string] int

func main() {
	var inputDir string
	var duration int
	mutex = new (sync.Mutex)
	// read the parameters:
	flag.IntVar(&rate, "rate", 5000, "The rate (number) of requests per minute")
	flag.StringVar(&inputDir, "inputDir", "StressTestCase", "The input directory to be parsed")
	flag.IntVar(&duration, "duration", 300, "The duration (s) of the test")
	flag.Parse()
	
	arrivalRateInMillis = float64(rate) / (60.0 * 1000.0) // convert rate to rate in millisecond
	fmt.Println(arrivalRateInMillis)

	//flag.IntVar(&timeOut, "timeOut", TIME_OUT_FOR_READ, "The time (s) to wait when trying to read from server")

	inputTestSuites = make([]StressTestEngine.TestSuite, 0)
	// read all input test case
	readInputDir(inputDir)

	
	testDuration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))


	controlChannel = make(chan int, 1)
	resultChannel = make(chan ClientResult, 1000)
	finishChannel = make(chan int, 0)
	// spawn processing routing
	go processingRoutine()
	// wait duration

	
	errorAggregation = make(map[string]int)
	statisicRoutine()
}

func readInputDir(inputDir string) {
	err := filepath.Walk(inputDir, visit)
	if err != nil {
		fmt.Println(err)
	}
}


// visit a file path
// if the file is a .xml (test file), 
// read the whole file into a TestSuite struct, and store it
func visit(currentPath string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", currentPath)
	if !f.IsDir() {
		var fileName = f.Name()
		var strLen = len(fileName)
		if strLen >= len(DEFAULT_EXTENSION) &&
			fileName[strLen-len(DEFAULT_EXTENSION):strLen] == DEFAULT_EXTENSION {
			// run file with extension "xml"

			ts, err := StressTestEngine.ReadXmlInput(currentPath)
			if err != nil {
				errorMsg := fmt.Sprintf("Cannot read input file because %v", err)
				log.Fatal(errorMsg)
			} 

			// add the test to the test suites
			inputTestSuites = append(inputTestSuites, *ts)

		}
	}
	return err
}

func getNextArrivalValue(rate float64) (float64) {
	return -math.Log(1.0 - rand.Float64()) / rate
}

func getNextArriveTimeInMillisecond(rate float64) (time.Duration) {
	duration := fmt.Sprintf("%.2fms", getNextArrivalValue(rate))
	result, _ := time.ParseDuration(duration)
	return result
}

func getRandomTestSuite() (StressTestEngine.TestSuite) {
	return inputTestSuites[rand.Intn(len(inputTestSuites))]
}

func processingRoutine() {
	totalTest := 0
	startTime := time.Now()
	for {
		elapsedTime := time.Since(startTime)
		if (elapsedTime <= testDuration) {
			nextTime := getNextArriveTimeInMillisecond(arrivalRateInMillis)
			time.Sleep(nextTime)
			testSuite := getRandomTestSuite()
			totalTest++
			go startAClient(&testSuite, resultChannel)//put channel here)
		} else {
			// stop
			break
		}
	}
	controlChannel <- totalTest
	<- finishChannel
}

func atomicChangeValue(valueToChange *int, delta int) {
	mutex.Lock()
	*valueToChange += delta
	mutex.Unlock()
}

func changeNumConnection(delta int) {
	atomicChangeValue(&totalConcurrentConnection, delta)
}

func startAClient(testSuite *StressTestEngine.TestSuite, resultChannel chan ClientResult) {
	if DEBUG_STEP  {
		fmt.Println("START A NEW CLIENT RUNNING:", testSuite.TestSuiteName)
	}
	addr := StressTestEngine.ResolveAddress(testSuite)
	startTestTime := time.Now()
	changeNumConnection(+1)
	testSuiteResult := StressTestEngine.ExecuteTestSuite(addr, testSuite)
	changeNumConnection(-1)
	duration := time.Since(startTestTime)
	result := ClientResult {
		IsSuccessful: testSuiteResult.IsCorrect,
		ExecuteTime: duration,
		Reason: testSuiteResult.Reason,
	}
	if DEBUG_STEP  {
		fmt.Println("FINISH A CLIENT RUNNING:", testSuite.TestSuiteName)
	}
	resultChannel <- result
}

type AggregateResult struct {
	NumCorrect int
	TotalFinished int
	TotalExecuteTime time.Duration
	AverageExecuteTime time.Duration
	TotalSuccessExecuteTime time.Duration
	AverageSuccessExecuteTime time.Duration
}

func statisicRoutine() {
	if DEBUG_STEP  {
		fmt.Println("START COLLECTING STATISTICS:")
	}
	finishSpawningClient := false
	finishExecuteClient := false
	 
	numClientFinished := 0

	var totalNum int
	defaultDuration, _ := time.ParseDuration("0s")
	aggregateResult := AggregateResult {
		NumCorrect: 0, 
		TotalFinished: 0,
		TotalExecuteTime: defaultDuration,
		AverageExecuteTime: defaultDuration,
		TotalSuccessExecuteTime: defaultDuration,
		AverageSuccessExecuteTime: defaultDuration,
	}

	for ((!finishSpawningClient) || (!finishExecuteClient)) {
		stop := false
		for !stop {
			select {
				case totalNum = <- controlChannel:
					fmt.Println("CLIENT WILL NOT BE SPAWNED ANY MORE!!! XXXXXXXXXXXXX")
					finishSpawningClient = true
				case clientResult := <- resultChannel:
					//fmt.Println("RECEIVING RESULT SIGNAL")
					numClientFinished++
					mergeResult(clientResult, &aggregateResult)
					if (finishSpawningClient) {
						if numClientFinished == totalNum {
							finishExecuteClient = true
						}
					}
				default: 
					//fmt.Println("NO SIGNAL")
					stop = true
			}
		}
		sleepDuration, _ := time.ParseDuration("1s")
		time.Sleep(sleepDuration)
		reportStatistic(totalConcurrentConnection, aggregateResult)
	}

	reportFinalStatisic(totalConcurrentConnection, aggregateResult, errorAggregation)
	finishChannel <- 1
}

func mergeResult(clientResult ClientResult, 
	aggregateResult *AggregateResult){
	aggregateResult.TotalFinished++
	aggregateResult.TotalExecuteTime = aggregateResult.TotalExecuteTime + clientResult.ExecuteTime
	totalNanosecond := aggregateResult.TotalExecuteTime.Nanoseconds()
	average := float64(totalNanosecond) / float64(aggregateResult.TotalFinished)
	averageDuration := fmt.Sprintf("%.2fns", average)
	aggregateResult.AverageExecuteTime, _ = time.ParseDuration(averageDuration)

	if clientResult.IsSuccessful {
		aggregateResult.NumCorrect++
		aggregateResult.TotalSuccessExecuteTime = aggregateResult.TotalSuccessExecuteTime + clientResult.ExecuteTime
		totalNanosecond = aggregateResult.TotalSuccessExecuteTime.Nanoseconds()
		average = float64(totalNanosecond) / float64(aggregateResult.NumCorrect)
		averageDuration := fmt.Sprintf("%.2fns", average)
		aggregateResult.AverageSuccessExecuteTime, _ = time.ParseDuration(averageDuration)	
	} else {
		if clientResult.Reason == nil {
			clientResult.Reason = errors.New("Unknown")
		} 	
		reg, _ := regexp.Compile("WSARecv tcp \\d{0,3}.\\d{0,3}.\\d{0,3}.\\d{0,3}:\\d{0,5}: resource temporarily unavailable")
		matchResult := reg.FindStringIndex(clientResult.Reason.Error())
		if  matchResult != nil{
			message := clientResult.Reason.Error()[:matchResult[0] - 1]
			fullMessage := message + "WSARecv tcp <address:port>: resource temporarily unavailable"
			clientResult.Reason = errors.New(fullMessage)
		}

		addErrorToTrack(clientResult.Reason)
		if DEBUG_FAILURE {
			fmt.Println("ERROR:", clientResult.Reason.Error())
		}
	}
}

func addErrorToTrack(err error) {
	errorMessage := err.Error()
	if _, exists:=errorAggregation[errorMessage]; !exists {
		errorAggregation[errorMessage] = 0
	}
	errorAggregation[errorMessage]++
}

func reportStatistic(totalConnection int, result AggregateResult) {
	fmt.Println("-------------------------------------------------------------")
	StressTestEngine.TimeEncodedPrintln("")
	fmt.Println("NUMBER OF CONCURRENT CONNECTION:", totalConnection)
	fmt.Printf("NUMBER OF FINISHED CLIENTS: %d\n", result.TotalFinished)
	if (result.TotalFinished != 0) {
	percentCorrect := float64(result.NumCorrect) / float64(result.TotalFinished) * 100
	fmt.Printf("NUMBER OF SUCCESSFUL CLIENTS: %d, ACCOUNT FOR: %3.2f %% \n", result.NumCorrect, percentCorrect)
		totalTimeInMillis := int(result.TotalExecuteTime.Nanoseconds() / 1000000)
		averageTimeInMillis := float64(result.TotalExecuteTime.Nanoseconds()) / 1000000 / float64(result.TotalFinished)
		fmt.Printf("TOTAL RUNNING TIME FOR ALL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
	fmt.Println("---------------------------------------")
		totalTimeInMillis = int(result.TotalSuccessExecuteTime.Nanoseconds() / 1000000)
		averageTimeInMillis = float64(result.TotalSuccessExecuteTime.Nanoseconds()) / 1000000 / float64(result.NumCorrect)
		fmt.Printf("TOTAL RUNNING TIME FOR ALL SUCCESSFUL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL SUCCESSFUL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
	}
	fmt.Println("-------------------------------------------------------------")
}


func reportFinalStatisic(totalConnection int, result AggregateResult, errorAggregation map[string]int) {
	
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("SUMMARY:")
	fmt.Println("NUMBER OF CONCURRENT CONNECTION:", totalConnection)
	fmt.Printf("NUMBER OF FINISHED CLIENTS: %d\n", result.TotalFinished)
	if (result.TotalFinished != 0) {
	percentCorrect := float64(result.NumCorrect) / float64(result.TotalFinished) * 100
	fmt.Printf("NUMBER OF SUCCESSFUL CLIENTS: %d, ACCOUNT FOR: %3.2f %% \n", result.NumCorrect, percentCorrect)
		totalTimeInMillis := int(result.TotalExecuteTime.Nanoseconds() / 1000000)
		averageTimeInMillis := float64(result.TotalExecuteTime.Nanoseconds()) / 1000000 / float64(result.TotalFinished)
		fmt.Printf("TOTAL RUNNING TIME FOR ALL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
		fmt.Println("---------------------------------------")
		totalTimeInMillis = int(result.TotalSuccessExecuteTime.Nanoseconds() / 1000000)
		averageTimeInMillis = float64(result.TotalSuccessExecuteTime.Nanoseconds()) / 1000000 / float64(result.NumCorrect)
		fmt.Printf("TOTAL RUNNING TIME FOR ALL SUCCESSFUL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL SUCCESSFUL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
	}

	// report errors:
	if result.NumCorrect < result.TotalFinished {
		fmt.Println("---------------------------------------")	
		fmt.Println("ERRORS SUMMARY:")
		for reason, occurenceTime := range errorAggregation {
			fmt.Println("OCCUR:", occurenceTime, "- ERROR:", reason)
		}
	}
	fmt.Println("-------------------------------------------------------------")
}