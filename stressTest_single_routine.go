package main

import (
	"flag"
	"fmt"
	"time"
	"errors"
	"log"
	"math"
	"math/rand"
	"xmltest/btalkTest/Auth_C2S"
	"xmltest/btalkTest/Auth_S2C"
	"container/list"
	"net"
	"runtime"
	"code.google.com/p/goprotobuf/proto"
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"encoding/json"
	"strconv"
)

const (
	C2S_HelloInfo_CMD                   = 0x01
	C2S_LoginInfo_CMD                   = 0x02
	C2S_RegAccount_CMD                  = 0x03
	C2S_OAuthLogin_CMD                  = 0x04
	C2S_FillFinishReg_CMD               = 0x05
	C2S_KeepAlive_CMD                   = 0x06
	C2S_ChangeNotification_CMD          = 0x07
	C2S_RequestUserNotificationInfo_CMD = 0x08
)

const (
	S2C_HelloInfoResult_CMD       = 0x01
	S2C_LoginUserInfo_CMD         = 0x02
	S2C_CheckVersionErrorInfo_CMD = 0x03
	S2C_NeedFinishReg_CMD         = 0x04
	S2C_KeepAliveAck_CMD          = 0x05
	S2C_NotificationInfo_CMD      = 0x06
)

const (
	ERROR_ErrorInfo_CMD = 0x20
)

const (
	NUM_MESSAGE_TO_SEND = 60
)

var rate int
var	arrivalRateInMillis float64

var duration int
var testDuration time.Duration

var address string
var serverAddr *net.TCPAddr

type Connection struct {
	Conn *net.TCPConn
	ConnId int64
	SendTime int64
	NumMessage int
	ResponseTime int64
}


var waitingConnectionChan chan *Connection
var finishedConnectionChan chan *Connection
var newConnectionChan chan *Connection

var listSendTime []time.Time
var listResponseTime []time.Duration

var numConnection int64
var numFinished int64
var numAlive int64

var finishSpawningChan chan int
var finishSendingChan chan int
var finishReceivingChan chan int

var totalTime int64
var totalRequest int64
var startServer bool

var DEBUG_WEB_SERVER bool = false
var DEBUG_STEP bool = false

var flagMap map[string]int

func main() {
	runtime.GOMAXPROCS(3)

	var duration int
	// read the parameters:
	flag.IntVar(&rate, "rate", 10000, "The rate (number) of requests per minute")
	flag.IntVar(&duration, "duration", 300, "The duration (s) of the test")
	flag.StringVar(&address, "address", "203.117.155.188:9101", "Address of test server")
	flag.BoolVar(&startServer, "server", false, "Start webserver to monitor or not?")
	flag.BoolVar(&DEBUG_WEB_SERVER, "web", false, "Debug web server")

	flag.Parse()
	flagMap = make(map[string]int)
	flagMap["rate"] = rate
	flagMap["duration"] = duration


	testDuration, _ = time.ParseDuration(fmt.Sprintf("%ds", duration))
	updateRate(rate)
	
	if (DEBUG_STEP) {
		TimeEncodedPrintln("Resolve server address:", address)
	}

	var err error
	serverAddr, err = net.ResolveTCPAddr("tcp", address)
	if (err != nil) {
		log.Fatal(err)
	}

	newConnectionChan = make (chan *Connection, 0)
	waitingConnectionChan = make (chan *Connection, 10000)
	finishedConnectionChan = make (chan *Connection, 10000)

	listSendTime = make([] time.Time, 100000)
	listResponseTime = make([] time.Duration, 100000)

	finishSpawningChan = make (chan int)
	finishSendingChan = make (chan int)
	finishReceivingChan = make (chan int)

	// spawn processing routing
	totalTime = 0
	totalRequest = 0
	go spawningRoutine()
	go sendingRoutine()
	go readingRoutine()
	go statisicRoutine()
	if startServer {
		go server()
	}

	<- finishReceivingChan
}

func updateRate(rate int) {
	arrivalRateInMillis = float64(rate) / (60.0 * 1000.0) // convert rate to rate in millisecond
	TimeEncodedPrintln(arrivalRateInMillis)
}

func spawningRoutine() {
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
			conn, err := startAClient(serverAddr) //put channel here)
			// pass the newly created connection to sending routing
			if (err != nil) {
				TimeEncodedPrintln("ERROR", err.Error())
			} else {
				newConnectionChan <- conn
			}
		} else {
			// stop
			break
		}
	}
	finishSpawningChan <- 1
}

func getNextArrivalValue(rate float64) float64 {
	return -math.Log(1.0-rand.Float64()) / rate
}

func getNextArriveTimeInMillisecond(rate float64) time.Duration {
	duration := fmt.Sprintf("%.2fms", getNextArrivalValue(rate))
	result, _ := time.ParseDuration(duration)
	return result
}

func startAClient(addr *net.TCPAddr) (*Connection, error) {
	if (DEBUG_STEP) {
		TimeEncodedPrintln("SPAWNING A CLIENT", addr)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err == nil {
		newConnection := &Connection {
			Conn: conn,
			ConnId : numConnection,
			NumMessage: 0,
		}
		numConnection++
		numAlive++
		return newConnection, nil
	} else {
		return nil, err
	}
	return nil, nil
}

func timer(resultChannel chan int, sleepDuration string) {
	for {
		sleepDuration, _ := time.ParseDuration(sleepDuration)
		time.Sleep(sleepDuration)
		resultChannel <- 1
	}
}

func sendingRoutine() {
	if DEBUG_STEP {
		fmt.Println("SENDING ROUTINE")
	}
	timerChannel := make (chan int)
	go timer(timerChannel, "1s")
	selfPushChannel := make (chan int, 1)

	var listAvailableConnection *list.List
	listAvailableConnection = list.New()
	var currentElement *list.Element
	currentElement = listAvailableConnection.Front()
	finishSpawning := false
	stop := false
	for {
		select {
			case <- finishSpawningChan:
				finishSpawning = true
			case connection:=<-newConnectionChan:
				listAvailableConnection.PushBack(connection)
				//fmt.Println("SEND, Get new connection list size: ", listAvailableConnection.Len())
				if (listAvailableConnection.Len() == 1) {
					selfPushChannel <- 1
				}
			
			case <- selfPushChannel:
				//fmt.Println("SEND: selfpush")
				if listAvailableConnection.Len() != 0 {
					if (currentElement == nil) {
						currentElement = listAvailableConnection.Front()
					}
					connection := currentElement.Value.(*Connection)
					timeNow := time.Now().UnixNano()
					if (timeNow - connection.SendTime > 1000000000) {						
						connection.SendTime = time.Now().UnixNano()
						sendHelloMsg(connection.Conn)
						waitingConnectionChan <- connection

						tempCurrent := currentElement
						currentElement = currentElement.Next()
						listAvailableConnection.Remove(tempCurrent)
						if (currentElement != nil) {
							selfPushChannel <- 1
						}
						if DEBUG_STEP {
							fmt.Println("SEND, Sent a message to connection: ", connection.ConnId, " - send time: ", connection.SendTime)
						}
					} else {
						//fmt.Println("DEFER SENDING MESSAGE FOR CONN: ", connection.ConnId)
						tempCurrent := currentElement
						currentElement = currentElement.Next()
						listAvailableConnection.Remove(tempCurrent)
						listAvailableConnection.PushBack(connection)
						if (currentElement == nil) {
							currentElement = listAvailableConnection.Front()
						}
						selfPushChannel <- 1	
					}
				}
			/*case <- timerChannel:
				//fmt.Println("SEND: timer")
				if listAvailableConnection.Len() != 0 {
					if (currentElement == nil) {
						currentElement = listAvailableConnection.Front()
					}
					connection := currentElement.Value.(*Connection)
					connection.SendTime = time.Now().UnixNano()
					sendHelloMsg(connection.Conn)
					waitingConnectionChan <- connection

					tempCurrent := currentElement
					currentElement = currentElement.Next()
					listAvailableConnection.Remove(tempCurrent)
					if (currentElement != nil) {
						selfPushChannel <- 1
					}
					if DEBUG_STEP {
						fmt.Println("SEND, Sent a message to connection: ", connection.ConnId, " - send time: ", connection.SendTime)
					}
				}*/
			case connection:=<-finishedConnectionChan:
				if (connection.NumMessage < NUM_MESSAGE_TO_SEND) {						
					listAvailableConnection.PushBack(connection)
					if (listAvailableConnection.Len() == 1) {
						selfPushChannel <- 1
					}
					if DEBUG_STEP {
						fmt.Println("SEND, Get a new available connection: ", listAvailableConnection.Len())
					}
				} else {
					if DEBUG_STEP {
						TimeEncodedPrintln("CONNECTION: ", connection.ConnId, connection.NumMessage)
					}
					numFinished++
					numAlive--
					if (finishSpawning && numFinished == numConnection) {
						stop = true
					}
				}
		}
		if stop {
			break
		}
	}
	finishSendingChan <- 1
}

func readingRoutine() {
	var listWaitingConnection *list.List
	listWaitingConnection = list.New()
	var currentElement *list.Element
	currentElement = listWaitingConnection.Front()
	selfPushChannel := make (chan int, 1)
	stop := false
	for {
		select {
			case <- finishSendingChan:
				stop = true
			case connection:=<-waitingConnectionChan:
				listWaitingConnection.PushBack(connection)
				if (listWaitingConnection.Len() == 1) {
					selfPushChannel <- 1
				}
				if DEBUG_STEP {
					fmt.Println("RECEIVE, Receive a WAITING channel: ", listWaitingConnection.Len())
				}
			case <- selfPushChannel:
				if listWaitingConnection.Len() != 0 {
					if (currentElement == nil) {
						currentElement = listWaitingConnection.Front()
					}

					connection := currentElement.Value.(*Connection)
					result, err := readReply(connection.Conn)
					if (result != nil) {
						totalRequest++
						currentTime := time.Now().UnixNano()
						if DEBUG_STEP {
							fmt.Println(currentTime, connection.SendTime, currentTime - connection.SendTime)
						}
						totalTime +=  currentTime - connection.SendTime
						connection.NumMessage++
						finishedConnectionChan <- connection
						tempCurrent := currentElement
						currentElement = currentElement.Next()
						listWaitingConnection.Remove(tempCurrent)
						if (listWaitingConnection.Len() != 0) {
							selfPushChannel <- 1
						} 
						if DEBUG_STEP {
							fmt.Println("RECEIVE, Receive a REPLY FOR: ", connection.ConnId)	
						}
					} else {
						if (err != nil) {
							TimeEncodedPrintln("ERROR:", err)
						} else {
							currentElement = currentElement.Next()
							if (listWaitingConnection.Len() != 0) {
								selfPushChannel <- 1
							}
						}
					}					
				}
		}
		if stop {
			break
		}
	}	
	finishReceivingChan <- 1
}


func sendHelloMsg(conn *net.TCPConn) {
	helloInfo := &Auth_C2S.HelloInfo{
		ClientType: proto.Int32(1),
		Version:    proto.Uint32(2),
	}
	sendMsg(conn, byte(C2S_HelloInfo_CMD), helloInfo)
}

func sendMsg(conn *net.TCPConn, cmd byte, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marchaling error: ", err)
	}
	length := int32(len(data)) + 1

	buf := new(bytes.Buffer)

	if (DEBUG_STEP) {
		fmt.Println("SENDING ROUTING: send a hello message ", length, cmd, data);
	}
	binary.Write(buf, binary.LittleEndian, length)
	binary.Write(buf, binary.LittleEndian, cmd)
	buf.Write(data)
	conn.Write(buf.Bytes())
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(conn *net.TCPConn) (proto.Message, error) {
	duration := time.Millisecond *20
	timeNow := time.Now()
	err := conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		TimeEncodedPrint("Cant set read timeout", err.Error())
		return nil, err
	}
	length := int32(0)
	err = binary.Read(conn, binary.LittleEndian, &length)
	if err != nil {
		return nil, nil
	}

	if (DEBUG_STEP) {
		fmt.Println("RECEIVE MESSAGE LENGTH: ", length)
	}
	// now wait longer for the message data
	duration = time.Millisecond *100
	timeNow = time.Now()
	err = conn.SetReadDeadline(timeNow.Add(duration))

	rbuf := make([]byte, length)
	io.ReadFull(conn, rbuf)
	if (DEBUG_STEP) {
		fmt.Println(rbuf)
	}
	
	cmd := int(rbuf[0])

	switch cmd {
	case S2C_HelloInfoResult_CMD:
		res := &Auth_S2C.HelloInfoResult{}
		err := proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		} else {
			return res, nil
		}
	
	case ERROR_ErrorInfo_CMD:
		res := &Auth_S2C.ErrorInfo{}
		err := proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		}
		msg := fmt.Sprint("Server returns error: ")
		return res, errors.New(msg) 
	default:
		log.Fatal("Unexpected CMD: ", cmd)
	}
	return nil, nil

}


func statisicRoutine() {
	startTestTime := time.Now()
	
	timerChannel := make(chan int, 0)
	go timer(timerChannel, "1s")

	var averageResponseTime float64
	var tick int
	for {
		select {
		case <-timerChannel:
			durationFromBegining := time.Now().Sub(startTestTime)
			tick = int(durationFromBegining.Seconds())
			fmt.Printf("Tick: %ds\n", tick)
			
			averageResponseTime = (float64(totalTime) / float64(totalRequest))
			output := OutputData{
				TotalConnection:      int(numConnection),
				TotalAliveConnection: int(numAlive),
				NumRealTimeRequest:   int(totalRequest),
				TotalFinished:        int(numFinished),
				TotalSuccessful:      0,
				PercentCorrect:       0,
				AverageResponse:      averageResponseTime,
				ErrorList:            nil,
				Time:                 tick,
			}

			fmt.Println("--------------------------------------------")
			fmt.Println("TOTAL TIME:", totalTime)
			fmt.Println("TOTAL REQUEST:", totalRequest)
			fmt.Println("TOTAL ALIVE:", numAlive)
			fmt.Println("AVERAGE RESPONSE TIME:", averageResponseTime / 1000000)
			fmt.Println("--------------------------------------------")
			

			if startServer {
				if DEBUG_WEB_SERVER {
					fmt.Println("Writing output to server")
				}
				serverChan <- output
			}
		}
	}
}

var serverChan chan OutputData

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
	if DEBUG_WEB_SERVER {
		fmt.Println("GET A RESULT TO RESPONSE", result)
	}
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
}


/* TIME ENCODED PRINTS */

func TimeEncodedPrintf(format string, a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Printf(format, a...)
}

func TimeEncodedPrint(a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Print(a...)
}

func TimeEncodedPrintln(a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Println(a...)
}
