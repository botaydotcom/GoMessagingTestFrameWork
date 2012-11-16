package main

import (
	"fmt"
	"log"
	"net"
)

type TestResult struct {
	IsCorrect bool
	Reason    error
	TimeTaken int64
}

type TestCaseResult struct {
	NumTest    int
	NumCorrect int
	Results    []TestResult
}

func main() {
	execute("www.google.com", "80", 20)
}

func execute(targetHost string, targetPort string, timesToTest int) {

	targetAddr := targetHost + ":" + targetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Fatal("Cannot resolve address")
	}

	var testCaseResult TestCaseResult
	testCaseResult.Results = make([]TestResult, timesToTest)
	testCaseResult.NumTest = timesToTest
	testCaseResult.NumCorrect = 0

	chans := make([]chan TestResult, timesToTest)
	for i := 0; i < timesToTest; i++ {
		chans[i] = make(chan TestResult)
	}

	for i := 0; i < timesToTest; i++ {
		go executeOneTest(addr, chans[i])
	}

	for i := 0; i < timesToTest; i++ {
		testCaseResult.Results[i] = <-chans[i]

		if testCaseResult.Results[i].IsCorrect {
			testCaseResult.NumCorrect++
		}
	}
	prettyPrint(testCaseResult)
}

func executeOneTest(addr *net.TCPAddr, resultChan chan TestResult) {
	var result TestResult
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(conn)
		result.IsCorrect = false
		result.Reason = err
	} else {
		fmt.Println("Successfully establish connection")
		result.IsCorrect = true
		result.Reason = nil
	}
	resultChan <- result
}

func prettyPrint(result TestCaseResult) {
	fmt.Println("-------------------------TEST RESULT--------------------------------------------")
	for i := 0; i < result.NumTest; i++ {
		fmt.Println("Case", i+1, ":")
		aResult := result.Results[i]
		if aResult.IsCorrect {
			fmt.Println("	CORRECT.   Time taken:", (aResult.TimeTaken / 1000), "us")
		} else {
			fmt.Println("   INCORRECT. Reason: ", aResult.Reason)
		}
	}
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Total result: ", result.NumCorrect, "/", result.NumTest, " = ", (result.NumCorrect * 100.0 / result.NumTest), "%")
	fmt.Println("-------------------------END OF TEST RESULT-------------------------------------")
}
