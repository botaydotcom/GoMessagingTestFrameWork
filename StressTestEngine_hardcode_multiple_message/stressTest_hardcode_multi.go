package StressTestEngine_hardcode_multiple_message

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	"net"
	"regexp"
	//"strconv"
	"xmltest/btalkTest/Auth_C2S"
	"xmltest/btalkTest/Auth_S2C"
	"xmltest/btalkTest/Helper"
	//"garena.com/btalkTest/TOKEN_S"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	//"os"
	//"reflect"
	//"text/template"
	"time"
	//"xmltest/btalkTest/GeneratedDataStructure"
)

const (
	START_MARK_VALUE = -732235123
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

	C2S_AddBuddyResult_CMD     = 0x66
	C2S_Chat2AckRemote_CMD     = 0x82
	C2S_InviteMemberResult_CMD = 0x08
	C2S_ChatInfoRecvedAck_CMD  = 0x0B // discussion chat ack
)

const (
	S2C_HelloInfoResult_CMD       = 0x01
	S2C_LoginUserInfo_CMD         = 0x02
	S2C_CheckVersionErrorInfo_CMD = 0x03
	S2C_NeedFinishReg_CMD         = 0x04
	S2C_KeepAliveAck_CMD          = 0x05
	S2C_NotificationInfo_CMD      = 0x06

	S2C_ChatInfo2_CMD             = 0x82
	S2C_RemoteRequestAddBuddy_CMD = 0x72
	S2C_InviteMember_CMD          = 0x09
	S2C_ChatInfo_CMD              = 0x08 // discussion chat ack

	DISCUSSION_PACKET_BASE_COMMAND = 0xA0
	DailyLifeRequest_MAINCMD       = 0xA1
)

const (
	ERROR_ErrorInfo_CMD = 0x20
)

const (
	EX_FB_AuthRequest = 0x01
	EX_TW_AuthRequest = 0x02
	EX_GA_AuthRequest = 0x03
)

const (
	EX_FB_AuthResponse = 0x01
	EX_TW_AuthResponse = 0x02
	EX_GA_AuthResponse = 0x03
)

const (
	CONN_NUM             = 3
)

// DEBUG FLAGS
var READ_PENDING_TIMEOUT string = "10s"

var DEBUG_COMPARE_POINTER bool = false
var DEBUG_COMPARE_SLICE bool = false
var DEBUG_COMPARE_STRUCT bool = false

var DEBUG_CLEANING_UP bool = false
var DEBUG_PARSING_MESSAGE bool = false

var DEBUG_SENDING_MESSAGE bool = false
var DEBUG_READING_MESSAGE bool = true
var DEBUG_IGNORING_MESSAGE bool = false

var BYPASS_CONNECTION_SERVER bool = false

var DEBUG_WAITING bool = true

var USE_UNBOUND bool = false
var USE_FORCE_CHECK bool = false
var USE_BYTE_ENCODE bool = false
var USE_PREPROCESS bool = false
var USE_VARIABLE bool = true
var USE_CHECK_REPLY bool = false
var USE_FILL_SLICE bool = false

type InnerXml struct {
	Xml string `xml:",innerxml"`
}

type UnboundVarDesc struct {
	VarName string `xml:"name,attr"`
	Field   string
}

type ByteEncodedVarDesc struct {
	MessageType string `xml:"type,attr"`
	Field   	string
	Data 		InnerXml
}

/*
	IsPartial    bool   `xml:"isPartial,attr"`
	PartialField int    `xml:"partialField,attr"`
*/
type Message struct {
	Repeat      int    `xml:"repeat,attr"`
	MessageType string `xml:"type,attr"`
	FromClient  bool   `xml:"fromClient,attr"`
	Connection  int    `xml:"connection,attr"`
	Close       bool   `xml:"close,attr"`
	Wait        string `xml:"wait,attr"`
	BaseCommand string
	Command     string
	Unbound     []UnboundVarDesc `xml:"Unbound>Var"`
	ForceCheck  []string 	 `xml:"ForceCheck>FieldName"`
	ByteEncoded []ByteEncodedVarDesc `xml:"ByteEncoded>Var"`
	Data        InnerXml
}

type Var struct {
	Name   string `xml:"name,attr"`
	IsFunc bool   `xml:"isFunc,attr"`
	SendBack bool  `xml:"sendBack,attr"`
	Value  string
	Params []string
}

type SpecialMessage struct {
	Name            string    `xml:"name,attr"`
	MessageSequence []Message `xml:"MessageSequence>Message"`
}

type IgnoreMessageSignature struct {
	BaseCommand string
	Command     string
	MessageType string
}

type TestInfo struct {
	Skip           bool                     `xml:"skip,attr"`
	Repeat         int                      `xml:"repeat,attr"`
	IgnoreMessages []IgnoreMessageSignature `xml:"IgnoreMessages>Message"`
	Name           string
}

type TestCase struct {
	VarMap          []Var     `xml:"VarMap>Var"`
	MessageSequence []Message `xml:"MessageSequence>Message"`
	CleanUpSequence []Message `xml:"CleanUpSequence>Message"`
}

type TestSuite struct {
	TestSuiteName        string
	TargetHost           string
	TargetPort           string
	GlobalSpecicalMessages []SpecialMessage         `xml:"MessageMap>Var"`
	GlobalIgnoreMessages []IgnoreMessageSignature `xml:"IgnoreMessages>Message"`
	GlobalVarMap         []Var                    `xml:"VarMap>Var"`
	ListTestInfos        []TestInfo               `xml:"ListTest>TestInfo"`
	ListTestCases        []TestCase               `xml:"Tests>Test"`
}

type TestResult struct {
	IsCorrect bool
	Reason    error
	TimeTaken int64
	NumRequest int
}

type TestCaseResult struct {
	NumTest    int
	NumCorrect int
	Results    []TestResult
}

type TestSuiteResult struct {
	Skip            []bool
	TestCaseResults []TestCaseResult
}

type UniqueConnection struct {
	Connection *net.TCPConn
	Identifier int64
}

var hasPartial bool
var rawData []byte
var field int
var regVar *regexp.Regexp
var numVar int = 0
var readTimeOut string

func ResolveAddress(testSuite *TestSuite) (addr *net.TCPAddr){
	targetAddr := testSuite.TargetHost + ":" + testSuite.TargetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		TimeEncodedPrintln("Cannot resolve address")
		log.Fatal("Cannot resolve address")
	}
	return addr
}

func Initialize() {
	regVar = regexp.MustCompile("{{.([a-zA-Z0-9_]*)}}")
}

type Data struct {
	KeyMap map[string]string
	ValueMap map[string]interface{}
	SliceKeyMap map[string]string
	SliceMap map[string]interface{}
	ForceCheckMap map[string]bool
}
var SpecialChannel chan int
func parseAndAugmentSpecialMessage(testSuite *TestSuite) {
	//SPECIAL MESSAGE PART
	specialMessageMap := make(map[string]SpecialMessage)
	// parse global special messages
	parseGlobalSpecialMessages(testSuite.GlobalSpecicalMessages, specialMessageMap)
	for i, _ := range testSuite.ListTestCases {
		augmentTestCaseWithSpecialMessages(&testSuite.ListTestCases[i], specialMessageMap)
	}
}

func parseGlobalSpecialMessages(specialMessages []SpecialMessage, 
	specialMessageMap map[string]SpecialMessage) {
	for _, specialMessage := range specialMessages {
		messageName := specialMessage.Name
		specialMessageMap[messageName] = specialMessage
	}
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("SPECIAL MESSAGE MAP:")
		fmt.Println(specialMessageMap)
	}
}

func augmentTestCaseWithSpecialMessages(testCase *TestCase, 
		specialMessageMap map[string] SpecialMessage) {
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("AUGMENT TEST CASE MESSAGE SEQUENCE WITH SPECIAL MESSAGES")
	}
	messageSequences := testCase.MessageSequence
	augmentedMessageSequence := make([]Message, 0)
	for _, message := range messageSequences {
		messageType := message.MessageType
		if specialMessage, exist := specialMessageMap[messageType]; exist {
			repeat := message.Repeat
			if repeat == 0 {
				repeat = 1
			}
			for i := 0; i < repeat; i++ {
				augmentedMessageSequence = append(augmentedMessageSequence, specialMessage.MessageSequence...)
			}
		} else {
			augmentedMessageSequence = append(augmentedMessageSequence, message)
		}
	}
	testCase.MessageSequence = augmentedMessageSequence
	for _, message := range testCase.MessageSequence {
		if DEBUG_PARSING_MESSAGE {
			fmt.Println(message.MessageType)
		}
	}
}

/******************* steps in test ****************************************/
func ExecuteTestSuite(addr *net.TCPAddr, testSuite *TestSuite) (TestResult){
	var data Data
	data.KeyMap = make(map[string]string)
	data.ValueMap = make(map[string]interface{})

	data.SliceKeyMap = make(map[string]string)
	data.SliceMap = make(map[string]interface{})

	parseAndAugmentSpecialMessage(testSuite)
	readTimeOut = READ_PENDING_TIMEOUT
	var result TestResult 
	result.IsCorrect = true
	skipped := false
	for i := range testSuite.ListTestInfos {
		if !skipped && !testSuite.ListTestInfos[i].Skip {
			testCaseResult := executeTestCase(addr, &(testSuite.ListTestInfos[i]), &(testSuite.ListTestCases[i]), &data)
			result.TimeTaken = result.TimeTaken + testCaseResult.TimeTaken
			result.NumRequest = result.NumRequest + testCaseResult.NumRequest
			if !testCaseResult.IsCorrect {
				result.IsCorrect = false
				result.Reason = errors.New(fmt.Sprintf("Case: %d-%s", i, testCaseResult.Reason.Error()))				
				break
			} 
			
		}
	}
	
	return result
}

func executeTestCase(addr *net.TCPAddr, testInfo *TestInfo, testCase *TestCase, data *Data) TestResult {
	result := PerformTestCaseOnce(addr, testCase, data)	
	return result
}

func PerformTestCaseOnce(addr *net.TCPAddr, testCase *TestCase, data *Data) (TestResult){
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil,
		TimeTaken: 0,
		NumRequest: 0,
	}

	var listConnection = make(map[int]*UniqueConnection)
	var connectionState = make(map[int]bool)
	for i := range connectionState {
		connectionState[i] = false // all closed
	}

	var numRequest = 0
	var totalSendTime, totalReceiveTime int64
	totalSendTime = 0
	totalReceiveTime = 0

	beginTest := time.Now()
	var totalTime int64 = 0
	userName := Helper.GetNextUserEmail()

	for i, message := range testCase.MessageSequence {
		if message.Wait != "" {
			duration, err := time.ParseDuration(message.Wait)
			if err == nil {
				if DEBUG_WAITING {
					TimeEncodedPrintln("WAITING FOR: ", duration)
				}
				//var beforeConnection time.Time = time.Now()
				time.Sleep(duration)
				//fmt.Println("SLEEP TIME:", time.Since(beforeConnection).Seconds())
			}
		} else {
			if DEBUG_WAITING {
				fmt.Println("NO WAITING FOUND!")
			}
		}
		
		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if isOpened, exist := connectionState[connectionId]; !exist || !isOpened {
			//Connection does not exist or closed, create new connection for connection 
			
			conn, err := net.DialTCP("tcp", nil, addr)
			if (DEBUG_SENDING_MESSAGE) {
				fmt.Println("TRYING TO CONNECT TO: ", addr, " PROBLEM: ", err)
			}
			//conn.SetLinger(0)
			if err != nil {
				result.IsCorrect = false
				result.Reason = err
				break
			} else {
				newConnection := UniqueConnection{
					Connection: conn,
					Identifier: Helper.GetNextUniqueNumber(),
				}
				listConnection[connectionId] = &newConnection
				connectionState[connectionId] = true
			}
		} else {
			// Connection  exists => reuses
		}
		conn := listConnection[connectionId]
		var err error
		if DEBUG_SENDING_MESSAGE && DEBUG_READING_MESSAGE && BYPASS_CONNECTION_SERVER {
			fmt.Println("User name - identifier: ", userName, conn.Identifier)
		}
		switch {
		case message.MessageType == "Auth_C2S_HelloInfo":
			//fmt.Println("HELLO INFO")			
			totalSendTime += time.Now().Sub(beginTest).Nanoseconds()
			sendHello(conn)
		case message.MessageType == "Auth_S2C_HelloInfoResult": 
			//fmt.Println("HELLO INFO RESULT")
			err = readHelloReply(conn)
			totalReceiveTime += time.Now().Sub(beginTest).Nanoseconds()
			numRequest ++
			SpecialChannel <- 1 // 1 more request
		case message.MessageType == "Auth_C2S_LoginInfo":
			//fmt.Println("LOGIN INFO")
			totalSendTime += time.Now().Sub(beginTest).Nanoseconds()			
			sendLogin(conn, userName)
		case message.MessageType == "Auth_S2C_LoginUserInfo":
			//fmt.Println("LOGIN INFO RESULT")
			_, err = readLoginReply(conn)
			totalReceiveTime += time.Now().Sub(beginTest).Nanoseconds()
			numRequest ++
			SpecialChannel <- 1 // 1 more request
		}
		if err != nil {
			// INCORRECT Reply message 
			result.IsCorrect = false
			errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
			result.Reason = errors.New(errMsg)
			break
		} 
		if message.Close {
			connectionState[connectionId] = false
			conn.Connection.Close()
		}
	}

	for i := range listConnection {
		if connectionState[i] {
			listConnection[i].Connection.Close()
		}
	}

	totalTime = totalReceiveTime - totalSendTime
	result.TimeTaken = totalTime
	result.NumRequest = numRequest

	//cleanUpAfterTest(addr, testCase, data)

	return result
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(useBase bool, baseCmd byte, cmd byte, msg proto.Message, conn *UniqueConnection) {
	data, err := proto.Marshal(msg)
	if err != nil {
		TimeEncodedPrintln("marshaling error: ", err)
		log.Fatal("marshaling error: ", err)
	}
	length := int32(len(data)) + 1

	if useBase {
		length = length + 1
	}

	if BYPASS_CONNECTION_SERVER {
		length = length + 8
	}

	if DEBUG_SENDING_MESSAGE {
		TimeEncodedPrintln("sending message base length: ", length, " command / command / message: ", baseCmd, cmd, msg)
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, length)
	if useBase {
		binary.Write(buf, binary.LittleEndian, baseCmd)
	}

	binary.Write(buf, binary.LittleEndian, cmd)
	if BYPASS_CONNECTION_SERVER {
		if DEBUG_SENDING_MESSAGE {
			TimeEncodedPrintln("Send extra 8 bytes: ", conn.Identifier)
		}
		binary.Write(buf, binary.LittleEndian, int64(conn.Identifier))
	}
	buf.Write(data)
	n, err := conn.Connection.Write(buf.Bytes())
	if DEBUG_SENDING_MESSAGE {
		fmt.Println("Sent", n, "bytes", "encounter error:", err)
	}
}


func sendHello(conn *UniqueConnection) {
	helloInfo := &Auth_C2S.HelloInfo{
		ClientType: proto.Int32(1),
		Version:    proto.Uint32(1),
	}
	sendMsg(false, 0, byte(C2S_HelloInfo_CMD), helloInfo, conn)
}

func sendLogin(conn *UniqueConnection, userName string) {
	m := md5.New()
	io.WriteString(m, "123456")
	password := string(m.Sum(nil))
	password = hex.EncodeToString([]byte(password))
	loginInfo := &Auth_C2S.LoginInfo{
		Name: proto.String(userName),
		Password:       proto.String(password),
		ClientType:      proto.Int32(1),
		MachineId:       proto.String("1"),
		SoftwareVersion: proto.Int32(10200),
	}
	sendMsg(false, 0, byte(C2S_LoginInfo_CMD), loginInfo, conn)
}

func sendOAuth(conn *UniqueConnection) {
	oauthInfo := &Auth_C2S.OAuthRawInfo{
		Provider: proto.String("facebook"),
		Account:  proto.String("579766697@facebook"),
		Content:  []byte("BAAFUf1RW0coBAHDc4d3ZCmJD86miGqJd6egnmDqJUEsJkVYlQZC30CP0TqbCa7ZAWZA4lYZAnQATvpwbc1LNWtfMJW6HeoJQAgDiUESS2ttZCmasMDGrIqDzVHbTNPCaoZD"),
	}
	data, err := proto.Marshal(oauthInfo)
	if err != nil {
		log.Fatal("marshal error: ", err)
	}
	oauth := &Auth_C2S.OAuthLogin{
		OAuthInfo:       data,
		ClientType:      proto.Int32(1),
		MachineId:       proto.String("ddd"),
		SoftwareVersion: proto.Int32(1),
	}
	sendMsg(false, 0, byte(C2S_OAuthLogin_CMD), oauth, conn)
}

func sendGarenaAuth(conn *UniqueConnection) {
	oauthInfo := &Auth_C2S.OAuthRawInfo{
		Provider: proto.String("garena"),
		Account:  proto.String("Cheng.Wei@garena"),
		Content:  []byte("bb334276b3d83e7dcf64632ee027e0e1"),
	}
	data, err := proto.Marshal(oauthInfo)
	if err != nil {
		log.Fatal("marshal error: ", err)
	}
	oauth := &Auth_C2S.OAuthLogin{
		OAuthInfo:       data,
		ClientType:      proto.Int32(1),
		MachineId:       proto.String("ddd"),
		SoftwareVersion: proto.Int32(1),
	}
	sendMsg(false, 0, byte(C2S_OAuthLogin_CMD), oauth, conn)
}

func sendChatMsg(conn *UniqueConnection){

}

func readHelloReply(conn *UniqueConnection) (error){
	var receivedMessage interface{}
	var err error
	for {
		receivedMessage, err = readReply(S2C_HelloInfoResult_CMD, conn)
		// nil, nil means message is ignored
		if (receivedMessage != nil || err != nil) {
			break
		}
	}
	return err
}

func readLoginReply(conn *UniqueConnection) (int32, error){
	var receivedMessage interface{}
	var err error
	for {
		receivedMessage, err = readReply(S2C_LoginUserInfo_CMD, conn)
		// nil, nil means message is ignored
		if (receivedMessage != nil || err != nil) {
			break
		}
	}

	if err != nil {
		return 0, err
	} else {
		loginMessage := receivedMessage.(*Auth_S2C.LoginUserInfo)
		return *loginMessage.MyInfo.UserId, nil
	}
}


// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(command int, conn *UniqueConnection) (interface{}, error) {
	duration := time.Second *10
	timeNow := time.Now()
	err := conn.Connection.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		TimeEncodedPrint("Cant set read timeout", err.Error())
		return nil, err
	}
	length := int32(0)
	err = binary.Read(conn.Connection, binary.LittleEndian, &length)
	if DEBUG_READING_MESSAGE {
		fmt.Println("TRYING TO READ MESSAGE LENGTH => ", length, " ERROR: ", err)
	}

	if err != nil {
		return nil, err
	}

	if BYPASS_CONNECTION_SERVER {
		fromValue := int64(0)
		err = binary.Read(conn.Connection, binary.LittleEndian, &fromValue)
		if DEBUG_READING_MESSAGE {
			fmt.Println("Trying to read extra 8 bytes:" , fromValue, " PROBLEM: ", err)
		}
		if (fromValue != conn.Identifier) {
			if DEBUG_READING_MESSAGE {
				fmt.Println("Ignore this message")
			}		
			// ignore this message
			rbuf := make([]byte, length)
			io.ReadFull(conn.Connection, rbuf)
			return nil, nil
		}
		
		if err != nil {
			return nil, err
		}
		length = length - 8
	}

	rbuf := make([]byte, length)
	n, err := io.ReadFull(conn.Connection, rbuf)
	if int32(n) < length {
		return nil, errors.New("Corrupted message")
	}
	cmd := int(rbuf[0])
	if (cmd != command) {
		fmt.Println("Message not match", command, rbuf)
		message := fmt.Sprintf("Message not match: %v", rbuf)
		return nil, errors.New(message)
	}
	switch cmd {
	case S2C_HelloInfoResult_CMD:
		res := &Auth_S2C.HelloInfoResult{}
		err := proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		} else {
			return res, nil
		}

	case S2C_LoginUserInfo_CMD:
		res := &Auth_S2C.LoginUserInfo{}
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

/********************* Input - output helper *********************************/

func ReadXmlInput(inputFile string) (*TestSuite, error) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var testSuite TestSuite
	err = xml.Unmarshal(data, &testSuite)
	if err != nil {
		TimeEncodedPrintf("error: %v\n", err)
		return nil, err
	}
	return &testSuite, nil
}

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
