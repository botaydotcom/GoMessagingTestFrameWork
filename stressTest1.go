package main

import (
	"flag"
	"fmt"
	"time"
	//"io/ioutil"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"xmltest/btalkTest/StressTestEngine"
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

var DEBUG_WEB_SERVER bool = false

var timeOut int

var inputTestSuites []StressTestEngine.TestSuite
var arrivalRateInMillis float64

var totalConcurrentConnection = 0
var totalNumberConnection = 0

type ClientResult struct {
	TestName           string
	IsSuccessful       bool // success or not
	ExecuteTime        time.Duration
	ResponseTime       time.Duration
	Reason             error
	NumRequest         int
	NumRealTimeRequest int
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
var serverFin chan int
var mutex *sync.Mutex
var rate int
var testDuration time.Duration

var errorAggregation map[string]int
var startServer bool

var flagMap map[string]int

func main() {
	runtime.GOMAXPROCS(6)

	var inputDir string
	var duration int
	mutex = new(sync.Mutex)
	// read the parameters:
	flag.IntVar(&rate, "rate", 5000, "The rate (number) of requests per minute")
	flag.StringVar(&inputDir, "inputDir", "StressTestCase", "The input directory to be parsed")
	flag.IntVar(&duration, "duration", 300, "The duration (s) of the test")
	flag.BoolVar(&startServer, "server", false, "Start webserver to monitor or not?")

	flag.BoolVar(&DEBUG_STEP, "step", false, "Debug showing steps")

	flag.BoolVar(&DEBUG_WEB_SERVER, "web", false, "Debug web server")

	flag.BoolVar(&StressTestEngine.DEBUG_COMPARE_POINTER, "cmp_ptr", false, "Debug compare pointer")
	flag.BoolVar(&StressTestEngine.DEBUG_COMPARE_SLICE, "cmp_slc", false, "Debug compare slice")
	flag.BoolVar(&StressTestEngine.DEBUG_COMPARE_STRUCT, "cmp_str", false, "Debug compare struct")

	flag.BoolVar(&StressTestEngine.DEBUG_CLEANING_UP, "clean", false, "Debug clean up")
	flag.BoolVar(&StressTestEngine.DEBUG_PARSING_MESSAGE, "parse", false, "Debug parse messages")

	flag.BoolVar(&StressTestEngine.DEBUG_READING_MESSAGE, "read", false, "Debug read messages")
	flag.BoolVar(&StressTestEngine.DEBUG_SENDING_MESSAGE, "send", false, "Debug send messages")
	flag.BoolVar(&StressTestEngine.DEBUG_IGNORING_MESSAGE, "ignore", false, "Debug ignore messages")

	flag.BoolVar(&StressTestEngine.DEBUG_WAITING, "wait", false, "Debug waiting")

	flag.StringVar(&StressTestEngine.READ_PENDING_TIMEOUT, "timeout", "10s", "Read time out")

	flag.Parse()

	flagMap = make(map[string]int)
	flagMap["rate"] = rate
	flagMap["duration"] = duration

	updateRate(rate)
	//flag.IntVar(&timeOut, "timeOut", TIME_OUT_FOR_READ, "The time (s) to wait when trying to read from server")

	inputTestSuites = make([]StressTestEngine.TestSuite, 0)
	// read all input test case
	readInputDir(inputDir)

	testDuration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))

	controlChannel = make(chan int, 1)
	resultChannel = make(chan ClientResult, 100000)
	finishChannel = make(chan int, 0)
	serverChan = make(chan OutputData, 100)
	serverFin = make(chan int, 0)
	errorAggregation = make(map[string]int)

	StressTestEngine.SpecialChannel = make(chan int, 1000000)

	// spawn processing routing
	go processingRoutine()

	if startServer {
		go server()
	}

	statisicRoutine()

	if startServer {
		// spawn server

		<-serverFin
	}
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

func updateRate(rate int) {
	arrivalRateInMillis = float64(rate) / (60.0 * 1000.0) // convert rate to rate in millisecond
	fmt.Println(arrivalRateInMillis)
}

func getNextArrivalValue(rate float64) float64 {
	return -math.Log(1.0-rand.Float64()) / rate
}

func getNextArriveTimeInMillisecond(rate float64) time.Duration {
	duration := fmt.Sprintf("%.2fms", getNextArrivalValue(rate))
	result, _ := time.ParseDuration(duration)
	return result
}

func getRandomTestSuite() StressTestEngine.TestSuite {
	return inputTestSuites[rand.Intn(len(inputTestSuites))]
}

func processingRoutine() {
	totalTest := 0
	totalNumberConnection = 0

	startTime := time.Now()
	previousTime := startTime
	for {
		elapsedTime := time.Since(startTime)
		if elapsedTime <= testDuration {
			nextTime := getNextArriveTimeInMillisecond(arrivalRateInMillis)
			nextSpawnTime := previousTime.Add(nextTime)
			timeNow := time.Now()
			if nextSpawnTime.After(timeNow) {
				time.Sleep(nextSpawnTime.Sub(timeNow))
			}
			previousTime = nextSpawnTime
			testSuite := getRandomTestSuite()
			totalTest++
			totalNumberConnection++
			go startAClient(&testSuite, resultChannel) //put channel here)
		} else {
			// stop
			break
		}
	}
	controlChannel <- totalTest
	<-finishChannel
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
	if DEBUG_STEP {
		fmt.Println("START A NEW CLIENT RUNNING:", testSuite.TestSuiteName)
	}
	addr := StressTestEngine.ResolveAddress(testSuite)
	startTestTime := time.Now()
	changeNumConnection(+1)
	testSuiteResult := StressTestEngine.ExecuteTestSuite(addr, testSuite)
	changeNumConnection(-1)
	duration := time.Since(startTestTime)
	defaultDuration, _ := time.ParseDuration("0s")
	result := ClientResult{
		TestName:     testSuite.TestSuiteName,
		IsSuccessful: testSuiteResult.IsCorrect,
		ExecuteTime:  duration,
		ResponseTime: defaultDuration,
		Reason:       testSuiteResult.Reason,
		NumRequest:   testSuiteResult.NumRequest,
	}

	var averageResponse int64 = 0
	// average response time for one request for this client
	if testSuiteResult.NumRequest != 0 {
		averageResponse = testSuiteResult.TimeTaken / int64(testSuiteResult.NumRequest)
	}
	timeStr := fmt.Sprintf("%dns", averageResponse)
	if averageResponse > 1e12 {
		timeStr = fmt.Sprintf("%dms", averageResponse/1e6)
	} else if averageResponse > 1e9 {
		timeStr = fmt.Sprintf("%dus", averageResponse/1e3)
	}
	//fmt.Println(timeStr)

	var err error
	result.ResponseTime, err = time.ParseDuration(timeStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	if DEBUG_STEP {
		fmt.Println("FINISH A CLIENT RUNNING:", testSuite.TestSuiteName)
	}
	resultChannel <- result
}

type AggregateResult struct {
	NumRealTimeRequest        int
	NumRequest                int
	NumCorrect                int
	TotalFinished             int
	TotalExecuteTime          time.Duration
	AverageExecuteTime        time.Duration
	TotalSuccessExecuteTime   time.Duration
	AverageSuccessExecuteTime time.Duration

	TotalResponseTime          time.Duration
	AverageResponseTime        time.Duration
	TotalSuccessResponseTime   time.Duration
	AverageSuccessResponseTime time.Duration
}

var mapUserName map[string]bool

func statisicRoutine() {
	startTestTime := time.Now()
	mapUserName = make(map[string]bool)
	if DEBUG_STEP {
		fmt.Println("START COLLECTING STATISTICS:")
	}
	finishSpawningClient := false
	finishExecuteClient := false

	numClientFinished := 0

	var totalNum int
	defaultDuration, _ := time.ParseDuration("0s")
	aggregateResult := AggregateResult{
		NumRealTimeRequest:         0,
		NumRequest:                 0,
		NumCorrect:                 0,
		TotalFinished:              0,
		TotalExecuteTime:           defaultDuration,
		AverageExecuteTime:         defaultDuration,
		TotalSuccessExecuteTime:    defaultDuration,
		AverageSuccessExecuteTime:  defaultDuration,
		TotalResponseTime:          defaultDuration,
		AverageResponseTime:        defaultDuration,
		TotalSuccessResponseTime:   defaultDuration,
		AverageSuccessResponseTime: defaultDuration,
	}

	var tick int
	for (!finishSpawningClient) || (!finishExecuteClient) {
		durationFromBegining := time.Now().Sub(startTestTime)
		tick := int(durationFromBegining.Seconds())
		fmt.Printf("Tick: %ds\n", tick)
		stop := false
		for !stop {
			select {
			case totalNum = <-controlChannel:
				fmt.Println("CLIENT WILL NOT BE SPAWNED ANY MORE!!! XXXXXXXXXXXXX")
				finishSpawningClient = true

			case <-StressTestEngine.SpecialChannel:
				// special channel is used for sending real-time update
				/*if _, exist := mapUserName[info]; exist {
					log.Fatal("Used user name!!!!!")
					mapUserName[info] = true
				}*/
				aggregateResult.NumRealTimeRequest++

			case clientResult := <-resultChannel:
				//fmt.Println("RECEIVING RESULT SIGNAL")
				//fmt.Println("RESPONSE TIME:", clientResult.ResponseTime)
				numClientFinished++
				mergeResult(clientResult, &aggregateResult)
				if finishSpawningClient {
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
		reportStatistic(totalConcurrentConnection, aggregateResult, tick)
	}

	reportFinalStatisic(totalConcurrentConnection, aggregateResult, errorAggregation, tick)
	finishChannel <- 1
}

func mergeResult(clientResult ClientResult,
	aggregateResult *AggregateResult) {
	aggregateResult.TotalFinished++

	aggregateResult.NumRequest += clientResult.NumRequest

	aggregateResult.TotalExecuteTime = aggregateResult.TotalExecuteTime + clientResult.ExecuteTime
	totalNanosecond := aggregateResult.TotalExecuteTime.Nanoseconds()
	average := float64(totalNanosecond) / float64(aggregateResult.TotalFinished)
	averageDuration := fmt.Sprintf("%.2fns", average)
	aggregateResult.AverageExecuteTime, _ = time.ParseDuration(averageDuration)

	aggregateResult.TotalResponseTime = aggregateResult.TotalResponseTime + clientResult.ResponseTime
	totalNanosecond = aggregateResult.TotalResponseTime.Nanoseconds()
	average = float64(totalNanosecond) / float64(aggregateResult.TotalFinished)
	averageDuration = fmt.Sprintf("%.2fns", average)
	aggregateResult.AverageResponseTime, _ = time.ParseDuration(averageDuration)

	if clientResult.IsSuccessful {
		aggregateResult.NumCorrect++

		aggregateResult.TotalSuccessExecuteTime = aggregateResult.TotalSuccessExecuteTime + clientResult.ExecuteTime
		totalNanosecond = aggregateResult.TotalSuccessExecuteTime.Nanoseconds()
		average = float64(totalNanosecond) / float64(aggregateResult.NumCorrect)
		averageDuration := fmt.Sprintf("%.2fns", average)
		aggregateResult.AverageSuccessExecuteTime, _ = time.ParseDuration(averageDuration)

		aggregateResult.TotalSuccessResponseTime = aggregateResult.TotalSuccessResponseTime + clientResult.ResponseTime
		totalNanosecond = aggregateResult.TotalSuccessResponseTime.Nanoseconds()
		average = float64(totalNanosecond) / float64(aggregateResult.NumCorrect)
		averageDuration = fmt.Sprintf("%.2fns", average)
		aggregateResult.AverageSuccessResponseTime, _ = time.ParseDuration(averageDuration)
	} else {
		if clientResult.Reason == nil {
			clientResult.Reason = errors.New("Unknown")
		}

		reg, _ := regexp.Compile("WSARecv tcp \\d{0,3}.\\d{0,3}.\\d{0,3}.\\d{0,3}:\\d{0,5}: resource temporarily unavailable")
		matchResult := reg.FindStringIndex(clientResult.Reason.Error())
		if matchResult != nil {
			message := clientResult.Reason.Error()[:matchResult[0]]
			fullMessage := message + "WSARecv tcp <address:port>: resource temporarily unavailable"
			clientResult.Reason = errors.New(fullMessage)
		}

		reg1, _ := regexp.Compile("tcp \\d{0,3}.\\d{0,3}.\\d{0,3}.\\d{0,3}:\\d{0,5}: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.")
		matchResult = reg1.FindStringIndex(clientResult.Reason.Error())
		if matchResult != nil {
			message := clientResult.Reason.Error()[:matchResult[0]]
			fullMessage := message + "tcp <address:port>: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond."
			clientResult.Reason = errors.New(fullMessage)
		}

		clientResult.Reason = errors.New(fmt.Sprintf("%s-%s", clientResult.TestName, clientResult.Reason.Error()))

		addErrorToTrack(clientResult.Reason)
	}
}

func addErrorToTrack(err error) {
	errorMessage := err.Error()
	if _, exists := errorAggregation[errorMessage]; !exists {
		errorAggregation[errorMessage] = 0
	}
	errorAggregation[errorMessage]++
}

func reportStatistic(totalConcurrentConnection int, result AggregateResult, tick int) {
	fmt.Println("-------------------------------------------------------------")
	StressTestEngine.TimeEncodedPrintln("")
	fmt.Println("NUMBER OF CONCURRENT CONNECTION:", totalConcurrentConnection)
	fmt.Printf("NUMBER OF FINISHED CLIENTS: %d\n", result.TotalFinished)
	fmt.Printf("NUMBER OF REAL-TIME FINISHED REQUESTS: %d\n", result.NumRealTimeRequest)
	fmt.Printf("NUMBER OF FINISHED REQUESTS: %d\n", result.NumRequest)
	var totalTimeInMillis int64
	var averageTimeInMillis float64
	var percentCorrect float64
	if result.TotalFinished != 0 {
		percentCorrect = float64(result.NumCorrect) / float64(result.TotalFinished) * 100
		fmt.Printf("NUMBER OF SUCCESSFUL CLIENTS: %d, ACCOUNT FOR: %3.2f %% \n", result.NumCorrect, percentCorrect)
		totalTimeInMillis = int64(result.TotalExecuteTime.Nanoseconds() / 1000000)
		averageTimeInMillis = float64(result.TotalExecuteTime.Nanoseconds()) / 1000000 / float64(result.TotalFinished)
		fmt.Println("---------------------------------------")
		fmt.Printf("TOTAL RUNNING TIME FOR ALL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
		fmt.Println("---------------------------------------")

		if result.NumCorrect != 0 {
			totalTimeInMillis = int64(result.TotalSuccessExecuteTime.Nanoseconds() / 1000000)
			averageTimeInMillis = float64(result.TotalSuccessExecuteTime.Nanoseconds()) / 1000000 / float64(result.NumCorrect)
		} else {
			totalTimeInMillis = 0
			averageTimeInMillis = 0
		}
		fmt.Printf("TOTAL RUNNING TIME FOR ALL SUCCESSFUL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RUNNING TIME FOR ALL SUCCESSFUL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
		fmt.Println("---------------------------------------")
		totalTimeInMillis = int64(result.TotalResponseTime.Nanoseconds() / 1000000)
		averageTimeInMillis = float64(result.TotalResponseTime.Nanoseconds()) / 1000000 / float64(result.TotalFinished)
		fmt.Printf("TOTAL RESPONSE TIME FOR ALL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RESPONSE TIME FOR ALL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
		fmt.Println("---------------------------------------")

		if result.NumCorrect != 0 {
			totalTimeInMillis = int64(result.TotalSuccessResponseTime.Nanoseconds() / 1000000)
			averageTimeInMillis = float64(result.TotalSuccessResponseTime.Nanoseconds()) / 1000000 / float64(result.NumCorrect)
		} else {
			totalTimeInMillis = 0
			averageTimeInMillis = 0
		}
		fmt.Printf("TOTAL RESPONSE TIME FOR ALL SUCCESSFUL CLIENTS  : %d MILLISECONDS \n", totalTimeInMillis)
		fmt.Printf("AVERAGE RESPONSE TIME FOR ALL SUCCESSFUL CLIENTS: %3.2f MILLISECONDS \n", averageTimeInMillis)
	}
	output := OutputData{
		TotalConnection:      totalNumberConnection,
		TotalAliveConnection: totalConcurrentConnection,
		NumRealTimeRequest:   result.NumRealTimeRequest,
		TotalFinished:        result.TotalFinished,
		TotalSuccessful:      result.NumCorrect,
		PercentCorrect:       percentCorrect,
		AverageResponse:      averageTimeInMillis,
		ErrorList:            errorAggregation,
		Time:                 tick,
	}

	if startServer {
		if DEBUG_WEB_SERVER {
			fmt.Println("Writing output to server")
		}
		serverChan <- output
	}

	fmt.Println("-------------------------------------------------------------")
}

var serverChan chan OutputData

func reportFinalStatisic(totalConnection int, result AggregateResult, errorAggregation map[string]int, tick int) {

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("SUMMARY:")
	reportStatistic(totalConnection, result, tick)

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

type OutputError struct {
	Name string
	Time int
}

type OutputData struct {
	TotalConnection      int
	TotalAliveConnection int
	TotalFinished        int
	TotalSuccessful      int
	PercentCorrect       float64
	AverageResponse      float64
	NumRealTimeRequest   int
	ErrorList            map[string]int
	Time                 int
}

func handler(w http.ResponseWriter, r *http.Request) {
	if DEBUG_WEB_SERVER {
		fmt.Println("GET A REQUEST")
	}
	result := <-serverChan
	//fmt.Println("GET A RESULT TO RESPONSE", result)
	b, err := json.Marshal(result)
	//fmt.Println(b)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Fprint(w, string(b))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if DEBUG_WEB_SERVER {
		fmt.Println("SERVER RESULT PAGE - ", r.URL.Path[1:])
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

func updateRateHandler(w http.ResponseWriter, r *http.Request) {
	if DEBUG_WEB_SERVER {
		fmt.Println("GET AN UPDATE RATE REQUEST")
	}

	rateValueStr := r.FormValue("rate")
	rateValue, err := strconv.ParseInt(rateValueStr, 0, 0)
	var response http.Response
	if err == nil && rateValue > 0 {
		updateRate(int(rateValue))
		if DEBUG_WEB_SERVER {
			fmt.Println("Update lambda, new value =", rate)
		}
		response.StatusCode = 200
	} else {
		response.StatusCode = 400
	}
	response.Write(w)
}

func getInfoHandler(w http.ResponseWriter, r *http.Request) {
	if DEBUG_WEB_SERVER {
		fmt.Println("GET AN INFO REQUEST")
	}

	b, err := json.Marshal(flagMap)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Fprint(w, string(b))
}

func server() {
	if DEBUG_WEB_SERVER {
		fmt.Println("PREPARE SERVER")
	}
	http.HandleFunc("/getInfo", getInfoHandler)
	http.HandleFunc("/updateRate", updateRateHandler)
	http.HandleFunc("/data/", handler)
	http.HandleFunc("/", viewHandler)
	if DEBUG_WEB_SERVER {
		fmt.Println("SERVER STARTS RUNNING:")
	}
	http.ListenAndServe(":3000", nil)
	serverFin <- 0
}
