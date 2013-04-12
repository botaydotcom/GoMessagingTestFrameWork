package StressTestEngine

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	"net"
	"regexp"
	"strconv"
	//"xmltest/btalkTest/Auth_C2S"
	//"xmltest/btalkTest/Auth_S2C"
	//"garena.com/btalkTest/TOKEN_S"
	"bytes"
	//"crypto/md5"
	"encoding/binary"
	//"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	//"os"
	"reflect"
	//"text/template"
	"time"
	"xmltest/btalkTest/GeneratedDataStructure"
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

var DEBUG_WAITING bool = false
var BYPASS_CONNECTION_SERVER bool = false

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

type Data struct {
	KeyMap map[string]string
	ValueMap map[string]interface{}
	SliceKeyMap map[string]string
	SliceMap map[string]interface{}
	GlobalIgnoreMessages map[int]string
	CurrentIgnoreMessages map[int]string
	ForceCheckMap map[string]bool
	SpecialMessageMap map[string]SpecialMessage
}
var SpecialChannel chan int

// do any initialize if needed
func Initialize() {
}


/******************* steps in test ****************************************/
func ExecuteTestSuite(addr *net.TCPAddr, testSuite *TestSuite) (TestResult){
	var data Data
	data.KeyMap = make(map[string]string)
	data.ValueMap = make(map[string]interface{})

	data.SliceKeyMap = make(map[string]string)
	data.SliceMap = make(map[string]interface{})

	data.GlobalIgnoreMessages = make(map[int] string)
	data.SpecialMessageMap = make(map[string]SpecialMessage)

	// parse global special messages
	parseGlobalSpecialMessages(testSuite.GlobalSpecicalMessages, &data)
	for i, _ := range testSuite.ListTestCases {
		augmentTestCaseWithSpecialMessages(&testSuite.ListTestCases[i], &data)
	}

	// parse global ignore messages
	parseGlobalIgnoreMessages(testSuite.GlobalIgnoreMessages, &data)

	// parse test-suite varmap first
	parseVarMap(testSuite.GlobalVarMap, &data)
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
	prepareCurrentIgnoreMap(testInfo.IgnoreMessages, data)

	result := PerformTestCaseOnce(addr, testCase, data)	

	//fmt.Println("TEST CASE RESULT:", result)
	return result
}

func PerformTestCaseOnce(addr *net.TCPAddr, testCase *TestCase, data *Data) (TestResult){
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil,
		TimeTaken: 0,
		NumRequest: 0,
	}

	var listConnection = make(map[int]*net.TCPConn)
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
	for i, message := range testCase.MessageSequence {
		if (DEBUG_PARSING_MESSAGE) {
			fmt.Println("Input message:", message)
		}
		// parse message (plug in values if necessary)
		parsedMessage, comm, useBase, baseCmd, err := parseAMessage(message, data)
		if err != nil {
			result.IsCorrect = false
			errMsg := fmt.Sprintf("Cannot parse message %d", i)
			result.Reason = errors.New(errMsg)
			break
		}
		protoParsedMessage := parsedMessage.(proto.Message)

		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if isOpened, exist := connectionState[connectionId]; !exist || !isOpened {
			//Connection does not exist or closed, create new connection for connection 
			
			conn, err := net.DialTCP("tcp", nil, addr)
			//conn.SetLinger(0)
			if err != nil {
				result.IsCorrect = false
				result.Reason = err
				break
			} else {
				listConnection[connectionId] = conn
				connectionState[connectionId] = true
			}
		} else {
			// Connection  exists => reuses
		}
		conn := listConnection[connectionId]

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

		if message.FromClient { // message from client, so send to server now
			if DEBUG_SENDING_MESSAGE {
				fmt.Println("Preparing to send a message to server ", message.MessageType)
			}
			begin := time.Now()
			sendMsg(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
			end := time.Now()
			totalTime += end.Sub(begin).Nanoseconds()
			totalSendTime += begin.Sub(beginTest).Nanoseconds()
		} else { // message from server, so read from buffer now
			// try to see if we can ignore this message
			if !useBase {
				baseCmd = 0
			}
			if tryIgnoreExpectedMessage(int(baseCmd), int(comm), protoParsedMessage, data) {
				// ignore this message
				continue
			}
			if DEBUG_READING_MESSAGE {
				fmt.Println("Expecting to receive a message from server ", message.MessageType, protoParsedMessage)
			}

			var replyMessage *proto.Message
			var err error
			for {
				begin := time.Now()
				replyMessage, err = readReply(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn, data)
				end := time.Now()
				totalReceiveTime += end.Sub(beginTest).Nanoseconds()
				totalTime += end.Sub(begin).Nanoseconds()
				numRequest ++
				SpecialChannel <- 1 // 1 more request
				if DEBUG_READING_MESSAGE {
					TimeEncodedPrintln("reply message, err: ", replyMessage, err)
				}
				if replyMessage != nil || err != nil {
					break
					// when both reply message / err == nil signals that the message is ignored
				}
			}

			if err == nil {
				for _, unboundDesc := range message.Unbound {
					data.KeyMap[unboundDesc.Field] = unboundDesc.VarName
				}

				// mark force check map for this message
				prepareForceCheckMap(message.ForceCheck, data)
				
				if comparedResult, err := compareGetValueForProtoMessage(protoParsedMessage, *replyMessage, data); comparedResult {
					// CORRECT Reply message.
				} else {
					// INCORRECT Reply message 
					result.IsCorrect = false
					errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
					result.Reason = errors.New(errMsg)
					break
				}
			} else {
				// Encounter error while trying to parse reply message
				result.IsCorrect = false
				errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
				result.Reason = errors.New(errMsg)
				break
			}
		}
		if message.Close {
			connectionState[connectionId] = false
			conn.Close()
		}
	}

	for i := range listConnection {
		if connectionState[i] {
			listConnection[i].Close()
		}
	}

	totalTime = totalReceiveTime - totalSendTime
	result.TimeTaken = totalTime
	result.NumRequest = numRequest

	cleanUpAfterTest(addr, testCase, data)

	return result
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(useBase bool, baseCmd byte, cmd byte, msg proto.Message, conn *net.TCPConn) {
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
		binary.Write(buf, binary.LittleEndian, int64(0))
	}
	buf.Write(data)
	n, err := conn.Write(buf.Bytes())
	if DEBUG_SENDING_MESSAGE {
		fmt.Println("Sent", n, "bytes", "encounter error:", err)
	}
}

func tryIgnoreReceivedMessage(buffer []byte, data *Data) bool {
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Trying to ignore received message buffer: ", buffer)
		fmt.Println("Global ignore list:", data.GlobalIgnoreMessages)
		fmt.Println("Current ignore list:", data.CurrentIgnoreMessages)
	}

	baseCmd := buffer[0]
	// report error
	signature := calculateSignature(0, int(baseCmd))
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Signature 1: ", signature)
	}
	if ignoreType, exist := data.CurrentIgnoreMessages[signature]; exist {
		newValue := magicVarFunc(ignoreType)
		res := newValue.(proto.Message)
		err := proto.Unmarshal(buffer[1:], res)
		if err == nil {
			// can ignore this message
			if DEBUG_IGNORING_MESSAGE {
				fmt.Println("Ignored!")
			}
			return true
		}
	}

	if len(buffer) >= 2 {
		signature = calculateSignature(int(baseCmd), int(buffer[1]))
		if DEBUG_IGNORING_MESSAGE {
			fmt.Println("Signature 2: ", signature)
		}
		if ignoreType, exist := data.CurrentIgnoreMessages[signature]; exist {
			newValue := magicVarFunc(ignoreType)
			res := newValue.(proto.Message)
			err := proto.Unmarshal(buffer[1:], res)
			if err == nil {
				// can ignore this message
				if DEBUG_IGNORING_MESSAGE {
					fmt.Println("Ignored!")
				}
				return true
			}
		}
	}
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Not ignored!")
	}
	return false
}

// when we expect some messages, and get error while receiving
// check if the expected messages can be ignored, and ignore it.
func tryIgnoreExpectedMessage(baseCmd int, cmd int, expMsg proto.Message, data *Data) bool {
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Trying to ignore expected message: ", baseCmd, cmd, expMsg)
		fmt.Println("Global ignore list:", data.GlobalIgnoreMessages)
		fmt.Println("Current ignore list:", data.CurrentIgnoreMessages)
	}
	signature := calculateSignature(baseCmd, cmd)
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Signature: ", signature)
	}
	if ignoreType, exist := data.CurrentIgnoreMessages[signature]; exist {
		fmt.Println(ignoreType, expMsg)
		if DEBUG_IGNORING_MESSAGE {
			fmt.Println("Ignored!")
		}
		return true
	}
	if DEBUG_IGNORING_MESSAGE {
		fmt.Println("Not Ignored!")
	}
	return false
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(useBase bool, expBaseCmd byte, expCmd byte, expMsg proto.Message, conn *net.TCPConn, data *Data) (*proto.Message, error) {
	length := int32(0)
	duration, err := time.ParseDuration(readTimeOut)
	timeNow := time.Now()
	if DEBUG_READING_MESSAGE {
		fmt.Println("set read deadline")
	}
	err = conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		return nil, err
	}

	err = binary.Read(conn, binary.LittleEndian, &length)
	if DEBUG_READING_MESSAGE {
		fmt.Println("read length", length)
	}
	if err != nil {
		return nil, err
	}

	if BYPASS_CONNECTION_SERVER {
		tempBuf := make([]byte, 8)
		_, err = io.ReadFull(conn, tempBuf)
		if err != nil {
			return nil, err
		}
		length = length - 8
	}

	var baseCmd byte
	if useBase {
		length = length - 1
		err = binary.Read(conn, binary.LittleEndian, &baseCmd)
		if baseCmd != expBaseCmd {
			// finish reading the rest of the message, 
			// so that it does not affects other message
			rbuf := make([]byte, length)
			io.ReadFull(conn, rbuf)

			// try to ignore this message before report error
			fullBuffer := make([]byte, 0)
			fullBuffer = append(fullBuffer, baseCmd)
			fullBuffer = append(fullBuffer, rbuf...)

			if tryIgnoreReceivedMessage(fullBuffer, data) {
				return nil, nil // signal ignored
			}

			errMsg := fmt.Sprintf("Unexpected BASE CMD %d", baseCmd)
			return nil, errors.New(errMsg)
		}
	}

	rbuf := make([]byte, length)
	_, err = io.ReadFull(conn, rbuf)
	if err != nil {
		return nil, err
	}

	var res proto.Message
	if len(rbuf) < 1 {
		errMes := fmt.Sprintf("Reply message is too short: %d", len(rbuf))
		return nil, errors.New(errMes)
	}

	cmd := rbuf[0]
	// command needs to be equal to expected command
	if cmd != expCmd {
		// try to ignore this message before report error
		fullBuffer := make([]byte, 0)
		if useBase {
			fullBuffer = append(fullBuffer, baseCmd)
		}
		fullBuffer = append(fullBuffer, rbuf...)
		if tryIgnoreReceivedMessage(fullBuffer, data) {
			return nil, nil // signal ignored
		}
		if DEBUG_READING_MESSAGE {
			fmt.Println("buffer, baseCmd", rbuf, baseCmd)
			fmt.Println("fullBuffer", fullBuffer)
		}

		errMsg := fmt.Sprintf("Unexpected CMD %d", cmd)
		return nil, errors.New(errMsg)
	}

	newValue := reflect.New(reflect.ValueOf(expMsg).Elem().Type()).Interface()
	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
		fmt.Println("Error when receiving message from server:")
		TimeEncodedPrintln(err)
		log.Fatal(err)
	}

	if DEBUG_READING_MESSAGE {
		TimeEncodedPrintln("Successfully receive message from network: ", res)
	}
	return &res, err
}

/**************Compare Values and Bind value from reply message to unbound var in value map*****/
func compareGetValueForPointer(expPtr reflect.Value, repPtr reflect.Value,
	domainName string, data *Data) (bool, error) {
	if DEBUG_COMPARE_POINTER {
		fmt.Println("Comparing pointers", expPtr.Elem().IsValid(), repPtr.Elem().IsValid())
		fmt.Println("Domain name:", domainName)
	}
	// Now check the current domain name of this pointer
	// if it is marked as: expecting value, and the variable is not bound
	// bind the value of the key
	key, present := data.KeyMap[domainName]
	if present {
		if DEBUG_COMPARE_POINTER {
			fmt.Println("Unbound pointer variable, bind value to map now")
		}
		// exp pointer != null. If pointer of reply value is null, return false
		if !repPtr.Elem().IsValid() {
			errorMsg := fmt.Sprintf("Reply has no value while expected to get a value for binding")
			return false, errors.New(errorMsg)
		}
		repValue := reflect.Indirect(reflect.ValueOf(repPtr.Interface()))
		// bound variable to value, return true
		data.ValueMap[key] = repValue.Interface()
		delete(data.KeyMap, domainName)
		// no error
		return true, nil
	}

	// if pointer of expected value is null, just dont compare, except when it's force check
	if !expPtr.Elem().IsValid(){
		if DEBUG_COMPARE_POINTER {
			fmt.Println("No expected value, check force check for:", domainName, "forceCheckMap:", data.ForceCheckMap)
		}
		if _, present = data.ForceCheckMap[domainName]; !present {
			return true, nil
		} else {
			if repPtr.Elem().IsValid() {
				repValue := reflect.Indirect(reflect.ValueOf(repPtr.Interface()))
				return false, errors.New(fmt.Sprint("Reply has value while expect no value for forced check field", domainName, " ", repValue.Interface()))
			} else {
				return true, nil
			}
		}
	}
	expValue := reflect.Indirect(reflect.ValueOf(expPtr.Interface()))

	// exp pointer != null. If pointer of reply value is null, return false
	if !repPtr.Elem().IsValid() {
		errorMsg := fmt.Sprintf("Reply has no value while expected value is %v", expValue.Interface())
		return false, errors.New(errorMsg)
	}
	repValue := reflect.Indirect(reflect.ValueOf(repPtr.Interface()))

	strVar := fmt.Sprintf("%v", expValue.Interface())

	if DEBUG_COMPARE_POINTER {
		fmt.Println("Expected pointer in message has value: ", strVar)
	}

	result := true
	var err error

	// if the expected field is a variable field that's not bound, bind the value of the key now
	key, present = data.KeyMap[strVar]
	if present {
		if DEBUG_COMPARE_POINTER {
			fmt.Println("Unbound pointer variable, bind value to map now")
		}
		// bound variable to value, return true
		data.ValueMap[key] = repValue.Interface()
	} else {
		// compare the two value of pointer
		// if they are not struct, just compare
		if expValue.Kind() == reflect.Struct {
			return compareGetValueForStruct(expValue, repValue, domainName, data)
		} else {
			isEqual := (expValue.Interface() == repValue.Interface())

			if DEBUG_COMPARE_POINTER {
				fmt.Println("Comparing: ", expValue.Interface(), repValue.Interface(), " equal? ", isEqual)
			}

			if !isEqual {
				errorMsg := fmt.Sprintf("Reply field value is different from expected field value. Expect: |%v|, get: |%v|",
					expValue.Interface(), repValue.Interface())
				result = false
				err = errors.New(errorMsg)
			}
		}
	}
	return result, err
}

func compareGetValueForSlice(expSlice reflect.Value, repSlice reflect.Value,
	domainName string, data *Data) (bool, error) {
	if DEBUG_COMPARE_SLICE {
		fmt.Println("Comparing slices")
		fmt.Println("Domain name:", domainName)
	}

	result := true
	var err error

	// Now check the current domain name of this slice
	// if it is marked as: expecting value, and the variable is not bound
	// bind the value of the key
	key, present := data.KeyMap[domainName]
	if present {
		if DEBUG_COMPARE_SLICE {
			fmt.Println("Unbound slice variable, bind value to map now")
		}
		// bound variable to value, return true
		data.ValueMap[key] = repSlice.Interface()
		delete(data.KeyMap, domainName)
		// no error
		return true, nil
	}

	var errorMsg string
	dif := false

	strVar := ""
	// try to convert the whole slice to a string value to see if it's unbound? (default value)
	byteArray, canConvert := expSlice.Interface().([]byte)
	if canConvert {
		strVar = string(byteArray)
	}

	// if the expected field is a variable field that's not bound, (its value exists in key map)
	// bind the value of the key now
	key, present = data.KeyMap[strVar]
	if present {
		if DEBUG_COMPARE_SLICE {
			fmt.Println("Unbound slice variable, bind value to map now")
		}

		// bound value to variable, return true		
		data.SliceKeyMap[strVar] = key
		data.SliceMap[key] = repSlice.Interface()
	} else {
		if expSlice.Len() != repSlice.Len() {
			errorMsg = fmt.Sprintf("Slices have different len, expect: %d, get: %d", expSlice.Len(), repSlice.Len())
			dif = true
		} else {
			for index := 0; index < expSlice.Len(); index++ {
				value1 := expSlice.Index(index)
				value2 := repSlice.Index(index)
				currentDomainName := fmt.Sprintf("%s[%d]", domainName, index)

				var isEqual bool
				if value1.Kind() == reflect.Ptr {
					isEqual, err = compareGetValueForPointer(value1, value2, currentDomainName, data)
				} else if value1.Kind() == reflect.Slice {
					isEqual, err = compareGetValueForSlice(value1, value2, currentDomainName, data)
				} else if value1.Kind() == reflect.Struct {
					isEqual, err = compareGetValueForStruct(value1, value2, currentDomainName, data)
				} else {
					// Now check the current domain name of this slice element
					// if it is marked as: expecting value, and the variable is not bound
					// bind the value of the key					
					key, present = data.KeyMap[currentDomainName]
					if present {
						if DEBUG_COMPARE_SLICE {
							fmt.Println("Unbound slice element-variable, bind value to map now")
						}
						// bound variable to value, return true
						data.ValueMap[key] = value2.Interface()
						delete(data.KeyMap, currentDomainName)
						// no error
						continue
					}
					isEqual = (value1.Interface() == value2.Interface())
				}
				if !isEqual {
					errorMsg = fmt.Sprintf("Reply field in slice has different value from expected field value. Expect: %v, get: %v",
						value1.Interface(), value2.Interface())
					dif = true
					break
				}
			}
		}
	}
	if dif {
		err = errors.New(errorMsg)
		result = false
	}
	return result, err
}

func compareGetValueForStruct(expStruct reflect.Value, repStruct reflect.Value,
	domainName string, data *Data) (bool, error) {
	if DEBUG_COMPARE_STRUCT {
		fmt.Println("Comparing structs")
		fmt.Println("Domain name:", domainName)
	}

	// Now check the current domain name of this struct
	// if it is marked as: expecting value, and the variable is not bound
	// bind the value of the key
	key, present := data.KeyMap[domainName]
	if present {
		if DEBUG_COMPARE_STRUCT {
			fmt.Println("Unbound struct variable, bind value to map now")
		}
		// bound variable to value, return true
		data.ValueMap[key] = repStruct.Interface()
		delete(data.KeyMap, domainName)
		// no error
		return true, nil
	}

	result := true
	var err error
	for i := 0; i < expStruct.NumField(); i++ {
		expStructField := expStruct.Field(i) // pointer
		repStructField := repStruct.Field(i) // pointer
		f := expStruct.Type().Field(i)
		currentDomainName := fmt.Sprintf("%s.%s", domainName, f.Name)

		if DEBUG_COMPARE_STRUCT {
			fmt.Println(expStructField, repStructField, expStructField.Kind())
		}

		var err1 error
		var isEqual bool
		if expStructField.Kind() == reflect.Ptr {
			isEqual, err1 = compareGetValueForPointer(expStructField,
				repStructField, currentDomainName, data)
		} else if expStructField.Kind() == reflect.Slice {
			isEqual, err1 = compareGetValueForSlice(expStructField,
				repStructField, currentDomainName, data)
		} else if expStructField.Kind() == reflect.Struct {
			isEqual, err1 = compareGetValueForStruct(expStructField,
				repStructField, currentDomainName, data)
		}
		if !isEqual {
			result = false
			errorMsg := fmt.Sprintf(" Field %d: %s", i, err1.Error())
			err = errors.New(errorMsg)
			break
		}
	}
	return result, err
}

func compareGetValueForProtoMessage(protoExp proto.Message, protoRep proto.Message, data *Data) (bool, error) {
	result := true
	var err error

	if reflect.ValueOf(protoExp).Type() != reflect.ValueOf(protoRep).Type() {
		errorMsg := fmt.Sprintf("Reply has different type from expected. Expect: %d, get: %d",
			reflect.ValueOf(protoExp).Type(), reflect.ValueOf(protoRep).Type())
		return false, errors.New(errorMsg)
	}

	expMessage := reflect.ValueOf(protoExp).Elem()
	repMessage := reflect.ValueOf(protoRep).Elem()

	result, err = compareGetValueForStruct(expMessage, repMessage, "", data)
	return result, err

}

/************************* Parse message ****************************/
func augmentTestCaseWithSpecialMessages(testCase *TestCase, data *Data) {
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("AUGMENT TEST CASE MESSAGE SEQUENCE WITH SPECIAL MESSAGES")
	}
	messageSequences := testCase.MessageSequence
	augmentedMessageSequence := make([]Message, 0)
	for _, message := range messageSequences {
		messageType := message.MessageType
		if specialMessage, exist := data.SpecialMessageMap[messageType]; exist {
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

func parseGlobalSpecialMessages(specialMessages []SpecialMessage, data *Data) {
	for _, specialMessage := range specialMessages {
		messageName := specialMessage.Name
		data.SpecialMessageMap[messageName] = specialMessage
	}
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("SPECIAL MESSAGE MAP:")
		fmt.Println(data.SpecialMessageMap)
	}
}


// parse varMap from input and put value to the value map
func parseVarMap(varMap []Var, data *Data) {
	for _, newVar := range varMap {
		if !newVar.IsFunc {
			data.ValueMap[newVar.Name] = newVar.Value
		} else {
			params:= make([]reflect.Value, 0)
			for _, value:= range newVar.Params {
				params = append(params, reflect.ValueOf(value))
			}
			data.ValueMap[newVar.Name] = magicCallFunc(newVar.Value, params)[0].Interface()
		}
		if DEBUG_PARSING_MESSAGE {
			fmt.Println("Paring global var: ", newVar.Name, " ==> ", data.ValueMap[newVar.Name])
		}
		/*if (newVar.SendBack) {
			result := fmt.Sprint(data.ValueMap[newVar.Name])
			SpecialChannel <- result
		}*/
	}
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("Global var map:", data.ValueMap)
	}
}

func calculateSignature(baseCommand int, command int) int {
	return baseCommand*500 + command
}

// parse global ignore messages
func parseGlobalIgnoreMessages(ignoreMessages []IgnoreMessageSignature, data *Data) {
	for _, messageType := range ignoreMessages {
		if DEBUG_PARSING_MESSAGE {
			fmt.Println("Parsing: ", messageType)
		}
		baseCmd, _ := strconv.ParseInt(messageType.BaseCommand, 0, 0)
		cmd, _ := strconv.ParseInt(messageType.Command, 0, 0)
		signature := calculateSignature(int(baseCmd), int(cmd))
		data.GlobalIgnoreMessages[signature] = messageType.MessageType
	}
}

func prepareForceCheckMap(forceCheckList []string, data *Data) {
	data.ForceCheckMap = make(map[string]bool)
	for _, value := range forceCheckList {
		data.ForceCheckMap[value] = true
	}
}

func prepareCurrentIgnoreMap(ignoreMessages []IgnoreMessageSignature, data *Data) {
	// clean the current ignore message, copy the whole global ignore map here
	data.CurrentIgnoreMessages = make(map[int]string)
	for key, value := range data.GlobalIgnoreMessages {
		data.CurrentIgnoreMessages[key] = value
	}
	for _, messageType := range ignoreMessages {
		baseCmd, _ := strconv.ParseInt(messageType.BaseCommand, 0, 0)
		cmd, _ := strconv.ParseInt(messageType.Command, 0, 0)
		signature := calculateSignature(int(baseCmd), int(cmd))
		data.CurrentIgnoreMessages[signature] = messageType.MessageType
	}
}

func plugValueForVar(message string, varName string, value interface{}) string {
	valueToReplace := ""
	kind := reflect.ValueOf(value).Kind()
	if kind == reflect.Slice {
		// dont plug value here	
	} else {
		reg, _ := regexp.Compile("{{." + varName + "}}")
		valueToReplace = fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(value)).Interface())
		message = reg.ReplaceAllString(message, valueToReplace)
	}
	return message
}

// plug value from value map to the raw xml message
func plugValue(message string, data *Data) (string, error) {
	for key, value := range data.ValueMap {
		message = plugValueForVar(message, key, value)
	}
	return message, nil
}

// pre-process a message and put unbound variable to value map
func preProcess(rawData string, data *Data) {
	regVar = regexp.MustCompile("{{.([a-zA-Z0-9_]*)}}")
	var processedData = rawData
	listVar := regVar.FindAllStringSubmatch(rawData, -1)
	for _, matchedValue := range listVar {
		key := matchedValue[1]
		if DEBUG_PARSING_MESSAGE {
			fmt.Println("Parsing variable key: ", key)
		}
		strMarkValue := ""

		_, present := data.ValueMap[key]
		if !present {
			if DEBUG_PARSING_MESSAGE {
				fmt.Println("Unbound key: ", key)
			}
			newMarkValue := START_MARK_VALUE - numVar
			numVar++
			strMarkValue = strconv.Itoa(newMarkValue)
			data.KeyMap[strMarkValue] = key
			data.ValueMap[key] = strMarkValue
			// now replace
			reg, _ := regexp.Compile("{{." + key + "}}")
			processedData = reg.ReplaceAllString(processedData, strMarkValue)
		}
	}
}

func fillSliceValueToMessage(rawMessage *interface{}, data *Data) {
	protoMsg := (*rawMessage).(proto.Message)
	message := reflect.ValueOf(protoMsg).Elem()

	for i := 0; i < message.NumField(); i++ {
		expStructField := message.Field(i) // pointer

		if expStructField.Kind() == reflect.Slice {
			strVar := ""
			byteArray, canConvert := expStructField.Interface().([]byte)

			if canConvert {
				strVar = string(byteArray)
				// if the expected field is a variable field that's not bound, bind the value of the key now
				key, present := data.KeyMap[strVar]

				/*if present {
					if sliceValue, valuePresent := data.SliceMap[key]; valuePresent {
						message.Field(i).Set(reflect.ValueOf(sliceValue))
					}
				}*/
				if present {
					// this is to fill a byte value receive from server to the same var
					if sliceValue, valuePresent := data.SliceMap[key]; valuePresent {
						message.Field(i).Set(reflect.ValueOf(sliceValue))
					}
				} else {
					// this is to fill a byte value from a function value to the same var
					matchValue := regVar.FindStringSubmatch(strVar)
					if (DEBUG_PARSING_MESSAGE) {
						fmt.Println("Match value:", matchValue)
					}
					if (matchValue != nil) {
						key = matchValue[1]	
						if sliceValue, valuePresent := data.ValueMap[key]; valuePresent {
							message.Field(i).Set(reflect.ValueOf(sliceValue))
						}						
					}
				}
			}
		}
	}

}



/**************Try to fill in byte-encoded value*****/

func visitPointerPluginByteEncoded(refPtrMessage *reflect.Value, value reflect.Value,
	domainName string, expectedDomainName string) (bool, error){
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("Try to fill byte-encoded value to pointer:")
		fmt.Println("Domain name:", domainName)
	}

	ptrMessage := *refPtrMessage

	message := reflect.Indirect(ptrMessage)
	if (domainName == expectedDomainName) {
		message.Set(value)
		return true, nil
	} else {
		if message.Kind() == reflect.Struct {
			return visitStructPluginByteEncoded(&message, value, domainName, expectedDomainName)
		} 
	}
	return false, nil
}

func visitSlicePluginByteEncoded(refSliceMessage *reflect.Value, value reflect.Value,
	domainName string, expectedDomainName string) (bool, error){
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("Try to fill in byte encoded value to slice")
		fmt.Println("Domain name:", domainName)
	}


	stop := false
	var err error

	sliceMessage := *refSliceMessage
	for index := 0; index < sliceMessage.Len(); index++ {
		currentDomainName := fmt.Sprintf("%s[%d]", domainName, index)
		field := sliceMessage.Index(index)		

		if field.Kind() == reflect.Ptr {
			stop, err = visitPointerPluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
		} else if field.Kind() == reflect.Struct {
			stop, err = visitStructPluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
		} else {
			if (currentDomainName == expectedDomainName) {
				sliceMessage.Field(index).Set(value)
				return true, nil
			} else if field.Kind() == reflect.Slice {
				stop, err = visitSlicePluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
			} 
		}
		if (stop) {
			break
		}
	}
		
	return stop, err
}

func visitStructPluginByteEncoded(refStructMessage *reflect.Value, value reflect.Value,
	domainName string, expectedDomainName string) (bool, error) {
	if DEBUG_PARSING_MESSAGE {
		fmt.Println("Try to fill in byte encoded value")
		fmt.Println("Domain name:", domainName)
	}

	stop := true
	var err error

	structMessage := *refStructMessage
	for i := 0; i < structMessage.NumField(); i++ {
		f := structMessage.Type().Field(i)
		currentDomainName := fmt.Sprintf("%s.%s", domainName, f.Name)
		field := structMessage.Field(i)

		if field.Kind() == reflect.Ptr {
			stop, err = visitPointerPluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
		} else if field.Kind() == reflect.Struct {
			stop, err = visitStructPluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
		} else {
			if (currentDomainName == expectedDomainName) {
				structMessage.Field(i).Set(value)
				return true, nil
			} else if field.Kind() == reflect.Slice {
				stop, err = visitSlicePluginByteEncoded(&field,
				value, currentDomainName, expectedDomainName)
			} 
		}
		if (stop) {
			return stop, err
		}
	}
	return stop, err
}


// parse a message from message sequence
// return: message / command / use base command? / base command / error
func parseAMessage(v Message, data *Data) (interface{}, int, bool, int, error) {
	// first pass: parse the message and put new var to value map
	preProcess(v.Data.Xml, data)

	// second pass: plug all value to the var in the raw string 
	addedXmlMessage, err := plugValue(v.Data.Xml, data)

	if DEBUG_PARSING_MESSAGE {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Processing message: ", addedXmlMessage)
	}

	// parse message
	message := magicVarFunc(v.MessageType)
	err = xml.Unmarshal([]byte(addedXmlMessage), message)

	if (len(v.ByteEncoded) != 0) {
		for _, byteEncoded := range v.ByteEncoded {
			preProcess(byteEncoded.Data.Xml, data)
			addedXmlMessage, err := plugValue(byteEncoded.Data.Xml, data)
			if DEBUG_PARSING_MESSAGE {				
				fmt.Println("Byte encoded message: ", addedXmlMessage)
			}
			if err != nil {
				TimeEncodedPrintf("error: %v\n", err)
				log.Fatal(err)
			}
			value := magicVarFunc(byteEncoded.MessageType)
			err = xml.Unmarshal([]byte(addedXmlMessage), value)
			byteEncodedData, _ := proto.Marshal(value.(proto.Message))
			reflectMessage := reflect.Indirect(reflect.ValueOf(message))
			visitStructPluginByteEncoded(&reflectMessage, reflect.ValueOf(byteEncodedData), "", byteEncoded.Field)
		}
	}

	if DEBUG_PARSING_MESSAGE {
		TimeEncodedPrintln("after unmarshal parsing, error = ", err, " message = ", message)
	}

	// run through the slice map to fill in missing value
	fillSliceValueToMessage(&message, data)

	if DEBUG_PARSING_MESSAGE {
		fmt.Println("After filling in slice: ", message)
	}

	if err != nil {
		TimeEncodedPrintf("error: %v\n", err)
	}

	if DEBUG_PARSING_MESSAGE {
		/*s := reflect.ValueOf(message).Elem()
		TimeEncodedPrintln("Message type: ", s.Type())
		for i := 0; i < s.NumField(); i++ {
			TimeEncodedPrint("Print value field ", i, " : ")
			f := s.Field(i)
			printValue(f)
		}*/
		//s := reflect.ValueOf(message).Elem()
		//printValueForStruct(s, 0)
	}

	cmd, _ := strconv.ParseInt(v.Command, 0, 0)

	useBase := true
	baseCmd, err1 := strconv.ParseInt(v.BaseCommand, 0, 0)
	if err1 != nil {
		useBase = false
		baseCmd = 0
	}

	return message, int(cmd), useBase, int(baseCmd), err
}

/************************* Pointer value helper functions *************************/
// call a function with a specific name 
func magicCallFunc(typeName string, in []reflect.Value) []reflect.Value {
	return GeneratedDataStructure.FuncMap[typeName].Call(in)
}

// create a zero pointer to a value of the correponding type in the map
func magicVarFunc(typeName string) interface{} {
	return reflect.New(GeneratedDataStructure.Map[typeName]).Interface()
}

func getKindForPointer(ptrValue reflect.Value) reflect.Kind {
	if ptrValue.Kind() == reflect.Ptr {
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			return value.Kind()
		}
	}

	return reflect.Invalid
}

func getValueForPointer(ptrValue reflect.Value) reflect.Value {
	if ptrValue.Kind() == reflect.Ptr {
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			return value
		}
	}

	return ptrValue
}

// use to get the value that a pointer points to
func getValue(ptrValue reflect.Value) (interface{}, error) {
	if ptrValue.Kind() == reflect.Ptr {
		//TimeEncodedPrintln(value.Kind(), value.Elem())
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			if _, ok := value.Interface().(interface {
				Int()
			}); ok {
				return value.Int(), nil
			}
			if _, ok := value.Interface().(interface {
				String()
			}); ok {
				return value.String(), nil
			}
			if _, ok := value.Interface().(interface {
				Float()
			}); ok {
				return value.Float(), nil
			}
		}
	}
	return ptrValue, errors.New("Not a pointer")
}

func getOffsetForLevel(level int) string {
	offset := ""
	for i := 0; i < level; i++ {
		offset = offset + "   "
	}
	return offset
}

func printValueForPointer(ptrValue reflect.Value, level int) {
	offset := getOffsetForLevel(level)
	if ptrValue.Elem().IsValid() {
		value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))

		fmt.Print(offset)
		fmt.Println("value:", value.Interface())
		if value.Kind() == reflect.Slice {
			printValueForSlice(value, level+1)
		}
		if value.Kind() == reflect.Struct {
			printValueForStruct(value, level+1)
		}
	} else {
		fmt.Print(offset)
		fmt.Println("Uninitialized")
	}
}

func printValueForSlice(sliceValue reflect.Value, level int) {
	offset := getOffsetForLevel(level)
	fmt.Print(offset)
	fmt.Println("Slice: ")
	for index := 0; index < sliceValue.Len(); index++ {
		fmt.Print(offset)
		fmt.Println("Index", index, ":")
		f := sliceValue.Index(index)
		if f.Kind() == reflect.Ptr {
			printValueForPointer(f, level+1)
		} else if f.Kind() == reflect.Slice {
			printValueForSlice(f, level+1)
		} else if f.Kind() == reflect.Struct {
			printValueForStruct(f, level+1)
		} else {
			fmt.Print(offset)
			fmt.Println(f.Interface())
		}
	}
}

func printValueForStruct(structValue reflect.Value, level int) {
	offset := getOffsetForLevel(level)
	fmt.Print(offset)
	fmt.Println("Struct: ", structValue.Type())
	for i := 0; i < structValue.NumField(); i++ {
		fmt.Print(offset)
		fmt.Println("Value field", i, ":")
		f := structValue.Field(i)
		if f.Kind() == reflect.Ptr {
			printValueForPointer(f, level+1)
		} else if f.Kind() == reflect.Slice {
			printValueForSlice(f, level+1)
		} else if f.Kind() == reflect.Struct {
			printValueForStruct(f, level+1)
		}
	}
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

func cleanUpAfterTest(addr *net.TCPAddr, testCase *TestCase, data *Data) {
	var listConnection = make(map[int]*net.TCPConn)

	// these messages are supposed to be login messages !
	for i, message := range testCase.CleanUpSequence {
		// parse message (plug in values if necessary)
		parsedMessage, comm, useBase, baseCmd, err := parseAMessage(message, data)
		if err != nil {
			TimeEncodedPrintln("Error in parsing clean up message at message:", i)
			continue
		}
		protoParsedMessage := parsedMessage.(proto.Message)

		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if _, exist := listConnection[connectionId]; !exist {
			// Connection does not exist or closed, create new connection for connection
			conn, err := net.DialTCP("tcp", nil, addr)
			conn.SetLinger(0)
			if err != nil {
				TimeEncodedPrintln("Cannot establish connection for clean up:", connectionId)
				continue
			} else {
				listConnection[connectionId] = conn
			}
		} else {
			// Connection exists => reuses
		}

		conn := listConnection[connectionId]
		if message.FromClient {
			// message from client, so send to server now
			sendMsg(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
		} else {
			// message from server, so read from buffer now
			readReply(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn, data)
		}
		if err != nil {
			continue
		}
	}

	pending := true
	for pending {
		pending = false
		for i, conn := range listConnection {
			if DEBUG_CLEANING_UP {
				TimeEncodedPrintln("Trying to read from connection", i, "and ack all pending messages")
			}
			// now all the established connections are the one to be clean up
			for {
				pendingMessage, err, comm := readPendingMessage(conn)
				if err != nil {
					break
				}
				if pendingMessage != nil {
					if DEBUG_CLEANING_UP {
						TimeEncodedPrintln("Pending message value: ", *pendingMessage)
					}
					// ack server
					ackPendingMessageToServer(pendingMessage, comm, conn)
					pending = true
				}
			}
			conn.Close()
		}
	}
}

// read a reply to a buffer WITHOUT knowing the message type before hand
// the message can only be in one of the few types: 
// AddBuddyRequest
// ChatInfo
// this is only used for cleaning up phase
// return error if cannot read a message in these type
func readPendingMessage(conn *net.TCPConn) (*proto.Message, error, byte) {
	length := int32(0)
	duration, err := time.ParseDuration(READ_PENDING_TIMEOUT)
	timeNow := time.Now()

	err = conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		return nil, err, 0
	}

	err = binary.Read(conn, binary.LittleEndian, &length)
	if err != nil {
		return nil, err, 0
	}

	rbuf := make([]byte, length)
	_, err = io.ReadFull(conn, rbuf)
	if err != nil {
		return nil, err, 0
	}

	var res proto.Message
	if len(rbuf) < 1 {
		errMes := fmt.Sprintf("Reply message is too short: %d", len(rbuf))
		return nil, errors.New(errMes), 0
	}

	baseCmd := rbuf[0]
	useBase := false
	if baseCmd == DISCUSSION_PACKET_BASE_COMMAND {
		// Found a discussion message with command 

		useBase = true
		rbuf = rbuf[1:]
	} else {
		baseCmd = 0
	}

	//TimeEncodedPrintln("Buffer read from network: ", rbuf)

	var newValue interface{}
	cmd := rbuf[0]
	switch cmd {
	case S2C_RemoteRequestAddBuddy_CMD:
		if !useBase {
			if DEBUG_CLEANING_UP {
				TimeEncodedPrintln("Found a request add buddy message, respond now: ")
			}
			newValue = magicVarFunc("Auth_Buddy_S2C_RemoteRequestAddBuddy")
		}

	case S2C_ChatInfo2_CMD:
		if !useBase {
			if DEBUG_CLEANING_UP {
				TimeEncodedPrintln("Found a chat message, respond now: ")
			}
			newValue = magicVarFunc("Auth_Buddy_S2C_ChatInfo2")
		}

	case S2C_InviteMember_CMD:
		if useBase {
			if DEBUG_CLEANING_UP {
				TimeEncodedPrintln("Found an invite member message, respond now: ")
			}
			newValue = magicVarFunc("Discussion_S2C_InviteMember")
		}

	case S2C_ChatInfo_CMD:
		if useBase {
			if DEBUG_CLEANING_UP {
				TimeEncodedPrintln("Found a discussion chat message, respond now: ")
			}
			newValue = magicVarFunc("Discussion_S2C_ChatInfo")
		}
	default:
		return nil, nil, 0 // receive a non-offline message
	}

	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
		TimeEncodedPrintln(err)
		log.Fatal(err)
	}
	return &res, err, cmd
}

func ackPendingMessageToServer(pendingMessage *proto.Message, comm byte, conn *net.TCPConn) {
	var replyMessage interface{}
	var res proto.Message
	var structValue reflect.Value
	var replyComm byte

	useBase := false
	baseCmd := 0

	switch comm {
	case S2C_RemoteRequestAddBuddy_CMD:
		replyMessage = magicVarFunc("Auth_Buddy_C2S_AddBuddyResult")
		res = replyMessage.(proto.Message)
		structValue = reflect.Indirect(reflect.ValueOf(replyMessage))
		userId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("FromId")
		structValue.FieldByName("UserId").Set(userId)
		var val int32 = 2
		structValue.FieldByName("Action").Set(reflect.ValueOf(&val))
		replyComm = C2S_AddBuddyResult_CMD

	case S2C_ChatInfo2_CMD:
		replyMessage = magicVarFunc("Auth_Buddy_C2S_Chat2AckRemote")
		res = replyMessage.(proto.Message)
		structValue = reflect.Indirect(reflect.ValueOf(replyMessage))
		userId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("FromId")
		msgId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("MsgId")
		structValue.FieldByName("UserId").Set(userId)
		structValue.FieldByName("MsgId").Set(msgId)
		replyComm = C2S_Chat2AckRemote_CMD

	case S2C_InviteMember_CMD:
		replyMessage = magicVarFunc("Discussion_C2S_InviteMemberResult")
		res = replyMessage.(proto.Message)
		structValue = reflect.Indirect(reflect.ValueOf(replyMessage))
		inviteId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("InviteId")
		structValue.FieldByName("InviteId").Set(inviteId)

		var val int32 = 0
		structValue.FieldByName("Agree").Set(reflect.ValueOf(&val))
		replyComm = C2S_InviteMemberResult_CMD
		useBase = true
		baseCmd = DISCUSSION_PACKET_BASE_COMMAND

	case S2C_ChatInfo_CMD:
		replyMessage = magicVarFunc("Discussion_C2S_ChatInfoRecvedAck")
		res = replyMessage.(proto.Message)
		structValue = reflect.Indirect(reflect.ValueOf(replyMessage))
		discussionId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("DiscussionId")
		structValue.FieldByName("DiscussionId").Set(discussionId)
		messageId := reflect.Indirect(reflect.ValueOf(*pendingMessage)).FieldByName("MsgId")
		structValue.FieldByName("MsgId").Set(messageId)

		replyComm = C2S_ChatInfoRecvedAck_CMD
		useBase = true
		baseCmd = DISCUSSION_PACKET_BASE_COMMAND
	}
	if DEBUG_CLEANING_UP {
		TimeEncodedPrintln("Reply with message", res)
	}
	sendMsg(useBase, byte(baseCmd), replyComm, res, conn)
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
