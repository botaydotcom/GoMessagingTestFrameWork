package TestEngine

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

	C2S_AddBuddyResult_CMD              = 0x66
	C2S_Chat2AckRemote_CMD              = 0x82
	C2S_InviteMemberResult_CMD   		= 0x08
	C2S_ChatInfoRecvedAck_CMD   		= 0x0B // discussion chat ack
)

const (
	S2C_HelloInfoResult_CMD       = 0x01
	S2C_LoginUserInfo_CMD         = 0x02
	S2C_CheckVersionErrorInfo_CMD = 0x03
	S2C_NeedFinishReg_CMD         = 0x04
	S2C_KeepAliveAck_CMD          = 0x05
	S2C_NotificationInfo_CMD      = 0x06

	S2C_ChatInfo2_CMD       	  = 0x82
	S2C_RemoteRequestAddBuddy_CMD = 0x72
	S2C_InviteMember_CMD		  = 0x09
	S2C_ChatInfo_CMD    		  = 0x08 // discussion chat ack

	DISCUSSION_PACKET_BASE_COMMAND = 0xA0

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
	READ_PENDING_TIMEOUT = "3s"
	CONN_NUM     		 = 3
)

// DEBUG FLAGS
const (
	DEBUG_COMPARE_POINTER = false
	DEBUG_COMPARE_SLICE   = false
	DEBUG_COMPARE_STRUCT  = false

	DEBUG_CLEANING_UP	  = false
	DEBUG_PARSING_MESSAGE = false

	DEBUG_IGNORING_MESSAGE = false
	DEBUG_SENDING_MESSAGE = false
	DEBUG_READING_MESSAGE = false
)

type InnerXml struct {
	Xml string `xml:",innerxml"`
}

/*
	IsPartial    bool   `xml:"isPartial,attr"`
	PartialField int    `xml:"partialField,attr"`
*/

type Message struct {
	MessageType string `xml:"type,attr"`
	FromClient  bool   `xml:"fromClient,attr"`
	Connection  int    `xml:"connection,attr"`
	Close       bool   `xml:"close,attr"`
	BaseCommand string
	Command     string
	Data        InnerXml
}

type Var struct {
	Name  string `xml:"name,attr"`
	Value string
}

type IgnoreMessageSignature struct {
	BaseCommand string
	Command 	string
	MessageType string
}

type TestInfo struct {
	Skip   bool `xml:"skip,attr"`
	Repeat int  `xml:"repeat,attr"`
	IgnoreMessages []IgnoreMessageSignature `xml:"IgnoreMessages>Message"`
	Name   string
}

type TestCase struct {
	VarMap          []Var     `xml:"VarMap>Var"`
	MessageSequence []Message `xml:"MessageSequence>Message"`
	CleanUpSequence []Message `xml:"CleanUpSequence>Message"`
}

type TestSuite struct {
	TestSuiteName string
	TargetHost    string
	TargetPort    string
	GlobalIgnoreMessages []IgnoreMessageSignature `xml:"IgnoreMessages>Message"`
	GlobalVarMap  []Var      `xml:"VarMap>Var"`
	ListTestInfos []TestInfo `xml:"ListTest>TestInfo"`
	ListTestCases []TestCase `xml:"Tests>Test"`
}


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

type TestSuiteResult struct {
	Skip            []bool
	TestCaseResults []TestCaseResult
}

var keyMap map[string]string
var valueMap map[string]interface{}

var sliceKeyMap map[string]string
var sliceMap map[string]interface{}

var globalIgnoreMessages map[int]string
var currentIgnoreMessages map[int]string

var hasPartial bool
var rawData []byte
var field int
var regVar *regexp.Regexp
var numVar int = 0
var readTimeOut string

func ExecuteTest(inputFile string, timeOut int) {
// open file
	readTimeOut = fmt.Sprintf("%ds", timeOut)
	fmt.Println("Testing file:", inputFile)

	keyMap = make(map[string]string)
	valueMap = make(map[string]interface{})
	sliceKeyMap = make(map[string]string)
	sliceMap = make(map[string]interface{})
	globalIgnoreMessages = make(map[int] string)

	ts, err := readXmlInput(inputFile)

	if err != nil {
		errorMsg := fmt.Sprintf("Cannot read input file because %v", err)
		log.Fatal(errorMsg)
	} else {
		executeTestSuite(ts)
	}

}

/******************* steps in test ****************************************/
func executeTestSuite(testSuite *TestSuite) {
	targetAddr := testSuite.TargetHost + ":" + testSuite.TargetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Fatal("Cannot resolve address")
	}

	// parse global ignore messages
	parseGlobalIgnoreMessages(testSuite.GlobalIgnoreMessages)
	
	// parse test-suite varmap first
	parseVarMap(testSuite.GlobalVarMap)

	var testSuiteResult TestSuiteResult
	testSuiteResult.Skip = make([]bool, len(testSuite.ListTestInfos))
	testSuiteResult.TestCaseResults = make([]TestCaseResult, len(testSuite.ListTestInfos))

	skipped := false
	for i := range testSuite.ListTestInfos {
		if !skipped && !testSuite.ListTestInfos[i].Skip {
			testSuiteResult.Skip[i] = false
			testSuiteResult.TestCaseResults[i] = executeTestCase(addr, &(testSuite.ListTestInfos[i]), &(testSuite.ListTestCases[i]))
			if testSuiteResult.TestCaseResults[i].NumCorrect !=  testSuiteResult.TestCaseResults[i].NumTest {
				skipped = true
				// skip the rest of the test suite
			}
		} else {
			testSuiteResult.Skip[i] = true
		}
	}
	prettyPrintTestSuite(testSuite, &testSuiteResult)
}

func executeTestCase(addr *net.TCPAddr, testInfo *TestInfo, testCase *TestCase) TestCaseResult {
	fmt.Println("Running test case:", testInfo.Repeat, "times")
	fmt.Println(testInfo.Name)

	prepareCurrentIgnoreMap(testInfo.IgnoreMessages)

	var testCaseResult TestCaseResult

	testCaseResult.Results = make([]TestResult, testInfo.Repeat)
	testCaseResult.NumTest = testInfo.Repeat
	testCaseResult.NumCorrect = 0

	chans := make([]chan TestResult, testInfo.Repeat)
	for i := 0; i < testInfo.Repeat; i++ {
		chans[i] = make(chan TestResult)
		go performTestCaseOnce(addr, testCase, chans[i])
	}

	for i := 0; i < testInfo.Repeat; i++ {
		testCaseResult.Results[i] = <-chans[i]
		if testCaseResult.Results[i].IsCorrect {
			testCaseResult.NumCorrect++
		} 
	}

	return testCaseResult
}

func performTestCaseOnce(addr *net.TCPAddr, testCase *TestCase, resultChan chan TestResult) {
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil}

	var listConnection = make(map[int]*net.TCPConn)
	var connectionState = make(map[int]bool)
	for i := range connectionState {
		connectionState[i] = false // all closed
	}

	var totalTime int64 = 0
	for i, message := range testCase.MessageSequence {
		// parse message (plug in values if necessary)
		parsedMessage, comm, useBase, baseCmd, err := parseAMessage(message)
		if err != nil {
			result.IsCorrect = false
			errMsg := fmt.Sprintf("Cannot parse message %d", i)
			result.Reason = errors.New(errMsg)
			break
		}
		protoParsedMessage := parsedMessage.(proto.Message)

		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if isOpened, exist := connectionState[connectionId]; (!exist ||  !isOpened) {
			//Connection does not exist or closed, create new connection for connection 
			conn, err := net.DialTCP("tcp", nil, addr)
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

		if message.FromClient { // message from client, so send to server now
			begin := time.Now()
			sendMsg(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
			end := time.Now()
			totalTime += end.Sub(begin).Nanoseconds()
		} else { // message from server, so read from buffer now
			// try to see if we can ignore this message
			if !useBase {
				baseCmd = 0
			}
			if (tryIgnoreExpectedMessage(int(baseCmd), int(comm), protoParsedMessage)) {
				// ignore this message
				continue
			}
			if (DEBUG_READING_MESSAGE) {
				fmt.Println("Expecting to receive a message from server ", message.MessageType)
			}

			var replyMessage *proto.Message 
			var err error
			for {
				replyMessage, err = readReply(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)	
				if (DEBUG_READING_MESSAGE) {
					fmt.Println("reply message, err: ", replyMessage, err)
				}
				if (replyMessage != nil || err != nil) {
					break;
					// when both reply message / err == nil signals that the message is ignored
				}

			}
			if err == nil {
				if comparedResult, err := compareGetValueForProtoMessage(protoParsedMessage, *replyMessage); comparedResult {
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

	result.TimeTaken = totalTime

	cleanUpAfterTest(addr, testCase)

	resultChan <- result
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(useBase bool, baseCmd byte, cmd byte, msg proto.Message, conn *net.TCPConn) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	length := int32(len(data)) + 1

	if (useBase) {
		length = length + 1
	}

	if (DEBUG_SENDING_MESSAGE) {
		log.Print("sending message base length: ", length, " command / command / message: ", baseCmd, cmd, msg)
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, length)
	if (useBase) {
		binary.Write(buf, binary.LittleEndian, baseCmd)
	}

	binary.Write(buf, binary.LittleEndian, cmd)
	//binary.Write(buf, binary.LittleEndian, int64(0))
	buf.Write(data)
	conn.Write(buf.Bytes())
}

func tryIgnoreReceivedMessage(buffer []byte) (bool){
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Trying to ignore received message buffer: ", buffer)
		fmt.Println("Global ignore list:", globalIgnoreMessages)
		fmt.Println("Current ignore list:", currentIgnoreMessages)
	}

	baseCmd := buffer[0]
	// report error
	signature := calculateSignature(0, int(baseCmd))
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Signature 1: ", signature)
	}	
	if ignoreType, exist := currentIgnoreMessages[signature]; exist {
		newValue := magicVarFunc(ignoreType)
		res := newValue.(proto.Message)
		err := proto.Unmarshal(buffer[1:], res)
		if err == nil {
			// can ignore this message
			if (DEBUG_IGNORING_MESSAGE) {
				fmt.Println("Ignored!")
			}	
			return true
		}
	}

	if (len(buffer) >= 2) {
		signature = calculateSignature(int(baseCmd), int(buffer[1]))
		if (DEBUG_IGNORING_MESSAGE) {
			fmt.Println("Signature 2: ", signature)
		}	
		if ignoreType, exist:= currentIgnoreMessages[signature]; exist {
			newValue := magicVarFunc(ignoreType)
			res := newValue.(proto.Message)
			err := proto.Unmarshal(buffer[1:], res)
			if err == nil {
				// can ignore this message
				if (DEBUG_IGNORING_MESSAGE) {
					fmt.Println("Ignored!")
				}
				return true
			}
		}
	}
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Not ignored!")
	}	
	return false
}

// when we expect some messages, and get error while receiving
// check if the expected messages can be ignored, and ignore it.
func tryIgnoreExpectedMessage(baseCmd int, cmd int, expMsg proto.Message) (bool){
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Trying to ignore expected message: ", baseCmd, cmd, expMsg)
		fmt.Println("Global ignore list:", globalIgnoreMessages)
		fmt.Println("Current ignore list:", currentIgnoreMessages)
	}
	signature := calculateSignature(baseCmd, cmd)
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Signature: ", signature)
	}
	if ignoreType, exist := currentIgnoreMessages[signature]; exist {
		fmt.Println(ignoreType, expMsg)
		if (DEBUG_IGNORING_MESSAGE) {
			fmt.Println("Ignored!")
		}
		return true
	}
	if (DEBUG_IGNORING_MESSAGE) {
		fmt.Println("Not Ignored!")
	}
	return false
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(useBase bool, expBaseCmd byte, expCmd byte, expMsg proto.Message, conn *net.TCPConn) (*proto.Message, error) {
	length := int32(0)
	duration, err := time.ParseDuration(readTimeOut)
	timeNow := time.Now()
	if (DEBUG_READING_MESSAGE) {
		fmt.Println("set read deadline")
	}
	err = conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		return nil, err
	}

	err = binary.Read(conn, binary.LittleEndian, &length)
	if (DEBUG_READING_MESSAGE) {
		fmt.Println("read length", length)
	}
	if err != nil {
		return nil, err
	}

	var baseCmd byte
	if (useBase) {
		length = length - 1
		err = binary.Read(conn, binary.LittleEndian, &baseCmd)
		if (baseCmd != expBaseCmd) {
			// finish reading the rest of the message, 
			// so that it does not affects other message
			rbuf := make([]byte, length)
			io.ReadFull(conn, rbuf)

			// try to ignore this message before report error
			fullBuffer := make([]byte, 0)
			fullBuffer = append(fullBuffer, baseCmd)
			fullBuffer = append(fullBuffer, rbuf...)

			if tryIgnoreReceivedMessage(fullBuffer) {
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
		if (useBase) {
			fullBuffer = append(fullBuffer, baseCmd)
		}
		fullBuffer = append(fullBuffer, rbuf...)
		if tryIgnoreReceivedMessage(fullBuffer) {
			return nil, nil // signal ignored
		}
		if (DEBUG_READING_MESSAGE) {
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
		log.Fatal(err)
	}

	if (DEBUG_READING_MESSAGE) {
		log.Print("Successfully receive message from network: ", res)
	}
	return &res, err
}

/**************Compare Values and Bind value from reply message to unbound var in value map*****/
func compareGetValueForPointer(expPtr reflect.Value, repPtr reflect.Value) (bool, error) {
	if (DEBUG_COMPARE_POINTER) {
		fmt.Println("Comparing pointers", expPtr.Elem().IsValid(), repPtr.Elem().IsValid())
	}
	// if pointer of expected value is null, just dont compare
	if !expPtr.Elem().IsValid() {
		return true, nil
	}
	expValue := reflect.Indirect(reflect.ValueOf(expPtr.Interface()))

	// if pointer of reply value is null, return false
	if !repPtr.Elem().IsValid() {
		errorMsg := fmt.Sprintf("Reply has no value while expected value is %v", expValue.Interface())
		return false, errors.New(errorMsg)
	}
	repValue := reflect.Indirect(reflect.ValueOf(repPtr.Interface()))

	strVar := fmt.Sprintf("%v", expValue.Interface())

	if (DEBUG_COMPARE_POINTER) {
		fmt.Println("Expected pointer in message has value: ", strVar)
	}

	result := true
	var err error

	// if the expected field is a variable field that's not bound, bind the value of the key now
	key, present := keyMap[strVar]
	if present {
		if (DEBUG_COMPARE_POINTER) {
			fmt.Println("Unbound pointer variable, bind value to map now")
		}
		// bound variable to value, return true
		valueMap[key] = repValue.Interface()
	} else {
		// compare the two value of pointer
		// if they are not struct, just compare
		if expValue.Kind() != reflect.Struct {
			isEqual := (expValue.Interface() == repValue.Interface())

			if (DEBUG_COMPARE_POINTER) {
				fmt.Println("Comparing: ", expValue.Interface(), repValue.Interface(), " equal? ", isEqual)
			}

			if !isEqual {
				errorMsg := fmt.Sprintf("Reply field value is different from expected field value. Expect: |%v|, get: |%v|",
					expValue.Interface(), repValue.Interface())
				result = false
				err = errors.New(errorMsg)
			}
		} else { // compare struct
			return compareGetValueForStruct(expValue, repValue)
		}
	}
	return result, err
}

func compareGetValueForSlice(expSlice reflect.Value, repSlice reflect.Value) (bool, error) {
	if (DEBUG_COMPARE_SLICE) {
		fmt.Println("Comparing slices")
	}

	result := true
	var err error
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
	key, present := keyMap[strVar]
	if present {
		if (DEBUG_COMPARE_SLICE) {
			fmt.Println("Unbound slice variable, bind value to map now")
		}

		// bound value to variable, return true		
		sliceKeyMap[strVar] = key
		sliceMap[key] = repSlice.Interface()
	} else {
		if expSlice.Len() != repSlice.Len() {
			errorMsg = fmt.Sprintf("Slices have different len, expect: %d, get: %d", expSlice.Len(), repSlice.Len())
			dif = true
		} else {
			for index := 0; index < expSlice.Len(); index++ {
				value1 := expSlice.Index(index)
				value2 := repSlice.Index(index)
				var isEqual bool
				if value1.Kind() == reflect.Ptr {
					isEqual, err = compareGetValueForPointer(value1, value2)
				} else if value1.Kind() == reflect.Slice {
					isEqual, err = compareGetValueForSlice(value1, value2)
				} else if value1.Kind() == reflect.Struct {
					isEqual, err = compareGetValueForStruct(value1, value2)
				} else {
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

func compareGetValueForStruct(expStruct reflect.Value, repStruct reflect.Value) (bool, error) {
	if (DEBUG_COMPARE_STRUCT) {
		fmt.Println("Comparing structs")
	}
	result := true
	var err error
	for i := 0; i < expStruct.NumField(); i++ {
		expStructField := expStruct.Field(i) // pointer
		repStructField := repStruct.Field(i) // pointer

		if (DEBUG_COMPARE_STRUCT) {
			fmt.Println(expStructField, repStructField, expStructField.Kind())
		}

		var err1 error
		var isEqual bool
		if expStructField.Kind() == reflect.Ptr {
			isEqual, err1 = compareGetValueForPointer(expStructField, repStructField)
		} else if expStructField.Kind() == reflect.Slice {
			isEqual, err1 = compareGetValueForSlice(expStructField, repStructField)
		} else if expStructField.Kind() == reflect.Struct {
			isEqual, err1 = compareGetValueForStruct(expStructField, repStructField)
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

func compareGetValueForProtoMessage(protoExp proto.Message, protoRep proto.Message) (bool, error) {
	result := true
	var err error

	if reflect.ValueOf(protoExp).Type() != reflect.ValueOf(protoRep).Type() {
		errorMsg := fmt.Sprintf("Reply has different type from expected. Expect: %d, get: %d",
			reflect.ValueOf(protoExp).Type(), reflect.ValueOf(protoRep).Type())
		return false, errors.New(errorMsg)
	}

	expMessage := reflect.ValueOf(protoExp).Elem()
	repMessage := reflect.ValueOf(protoRep).Elem()

	result, err = compareGetValueForStruct(expMessage, repMessage)
	return result, err

}

/************************* Parse message ****************************/
// parse varMap from input and put value to the value map
func parseVarMap(varMap []Var) {
	for _, newVar := range varMap {
		valueMap[newVar.Name] = newVar.Value
	}
}

func calculateSignature(baseCommand int, command int) (int) {
	return baseCommand * 500 + command;
}

// parse global ignore messages
func parseGlobalIgnoreMessages(ignoreMessages []IgnoreMessageSignature) {
	for _, messageType := range ignoreMessages {
		fmt.Println("Parsing: ", messageType)
		baseCmd, _ := strconv.ParseInt(messageType.BaseCommand, 0, 0)
		cmd, _ := strconv.ParseInt(messageType.Command, 0, 0)
		signature := calculateSignature(int(baseCmd), int(cmd))
		globalIgnoreMessages[signature] = messageType.MessageType
	}
}

func prepareCurrentIgnoreMap(ignoreMessages []IgnoreMessageSignature) {
	// clean the current ignore message, copy the whole global ignore map here
	currentIgnoreMessages = make(map[int]string)
	for key, value := range globalIgnoreMessages {
		currentIgnoreMessages[key] = value
	}
	for _, messageType := range ignoreMessages {
		baseCmd, _ := strconv.ParseInt(messageType.BaseCommand, 0, 0)
		cmd, _ := strconv.ParseInt(messageType.Command, 0, 0)
		signature := calculateSignature(int(baseCmd), int(cmd))
		currentIgnoreMessages[signature] = messageType.MessageType
	}
}

func plugValueForVar(message string, varName string, value interface{}) (string) {
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
func plugValue(message string) (string, error) {
	for key, value := range(valueMap) {
		message = plugValueForVar(message, key, value)
	}
	return message, nil
}

// pre-process a message and put unbound variable to value map
func preProcess(rawData string) {
	regVar = regexp.MustCompile("{{.([a-zA-Z0-9]*)}}")
	var processedData = rawData
	listVar := regVar.FindAllStringSubmatch(rawData, -1)
	for _, matchedValue := range listVar {
		key := matchedValue[1]
		strMarkValue := ""

		_, present := valueMap[key]
		if !present {
			newMarkValue := START_MARK_VALUE - numVar
			numVar++
			strMarkValue = strconv.Itoa(newMarkValue)
			keyMap[strMarkValue] = key
			valueMap[key] = strMarkValue
			// now replace
			reg, _ := regexp.Compile("{{." + key + "}}")
			processedData = reg.ReplaceAllString(processedData, strMarkValue)
		}
	}
}

func fillSliceValueToMessage(rawMessage *interface{}) {
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
				key, present := keyMap[strVar]
				
				if present {
					if sliceValue, valuePresent := sliceMap[key]; valuePresent {
						message.Field(i).Set(reflect.ValueOf(sliceValue))
					}
				}
			}
		}
	}
 
}

// parse a message from message sequence
// return: message / command / use base command? / base command / error
func parseAMessage(v Message) (interface{}, int, bool, int, error) {
	// first pass: parse the message and put new var to value map
	preProcess(v.Data.Xml)

	// second pass: plug all value to the var in the raw string 
	addedXmlMessage, err := plugValue(v.Data.Xml)

	if (DEBUG_PARSING_MESSAGE) {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Processing message: ", addedXmlMessage)
	}

	message := magicVarFunc(v.MessageType)
	err = xml.Unmarshal([]byte(addedXmlMessage), message)

	// run through the slice map to fill in missing value
	fillSliceValueToMessage(&message)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	
	if (DEBUG_PARSING_MESSAGE) {
		s := reflect.ValueOf(message).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			printValue(f)
		}
	}

	cmd, _ := strconv.ParseInt(v.Command, 0, 0)

	useBase := true
	baseCmd, err1 := strconv.ParseInt(v.BaseCommand, 0, 0)
	if (err1 != nil) {
		useBase = false
		baseCmd = 0
	}
	
	return message, int(cmd), useBase, int(baseCmd), err
}

/************************* Pointer value helper functions *************************/
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
		//fmt.Println(value.Kind(), value.Elem())
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

func printValue(ptrValue reflect.Value) {
	if ptrValue.Kind() == reflect.Ptr {
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			//fmt.Println("value:", value, value.String(), value.Interface().(string))

			if _, ok := value.Interface().(interface {
				Int()
			}); ok {
				fmt.Println(value.Int())
			}
			if _, ok := value.Interface().(interface {
				Float()
			}); ok {
				fmt.Println(value.Float())
			}
			if value.Kind() == reflect.String {
				fmt.Println(value.String())
			}
		} else {
			fmt.Println("Uninitialized")
		}
	} else {
		fmt.Println("Not a pointer")
	}
}

/********************* Input - output helper *********************************/

func readXmlInput(inputFile string) (*TestSuite, error) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var testSuite TestSuite
	err = xml.Unmarshal(data, &testSuite)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	return &testSuite, nil
}

func prettyPrintTestCase(testCaseInfo *TestInfo, isSkipped bool, result *TestCaseResult) (bool){
	incorrect := false
	var offset = "     "
	fmt.Printf(" %sTEST CASE: %s\n", offset, testCaseInfo.Name)
	if isSkipped {
		fmt.Println(offset, "    Skipped!")
	} else {
		totalTime := 0
		for i := 0; i < result.NumTest; i++ {
			fmt.Print(offset, "    Repeat time ", i+1, ": ")
			aResult := result.Results[i]
			if aResult.IsCorrect {
				fmt.Printf("	CORRECT.   Time taken: %5d us \n", (aResult.TimeTaken / 1000))
				totalTime += int(aResult.TimeTaken / 1000)
			} else {
				incorrect = true
				fmt.Println("   INCORRECT. Reason: ", aResult.Reason)
				totalTime += int(aResult.TimeTaken / 1000)
			}
		}
		fmt.Println(offset, "---------------------------")
		fmt.Printf("%s    Total result:      %2d/%2d = %3d%%   Time taken: %5d us\n", offset,
			result.NumCorrect, result.NumTest,
			(result.NumCorrect * 100.0 / result.NumTest), totalTime)
	}
	fmt.Println(offset, "END OF TEST CASE RESULT")
	return incorrect
}

func prettyPrintTestSuite(testSuite *TestSuite, result *TestSuiteResult) {
	incorrect := false
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("TEST SUITE: %s\n", testSuite.TestSuiteName)
	for i := 0; i < len(testSuite.ListTestInfos); i++ {
		fmt.Println()
		if (prettyPrintTestCase(&(testSuite.ListTestInfos[i]), result.Skip[i], &(result.TestCaseResults[i]))) {
			incorrect = true
		}
		fmt.Println()
	}
	fmt.Println("END OF TEST SUITE RESULT")
	fmt.Println("--------------------------------------------------------------------------------")
	if (incorrect) {
		log.Fatal("TEST CASE FAILED, STOP!")
	}
}

func cleanUpAfterTest(addr *net.TCPAddr, testCase *TestCase) {
	var listConnection = make(map[int]*net.TCPConn)

	// these messages are supposed to be login messages !
	for i, message := range testCase.CleanUpSequence {
		// parse message (plug in values if necessary)
		parsedMessage, comm, useBase, baseCmd, err := parseAMessage(message)
		if err != nil {
			log.Println("Error in parsing clean up message at message:", i)
			continue
		}
		protoParsedMessage := parsedMessage.(proto.Message)

		// get connection from list or make one if it does not exist yet
		connectionId := message.Connection
		if _, exist := listConnection[connectionId]; !exist {
			// Connection does not exist or closed, create new connection for connection
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				log.Print("Cannot establish connection for clean up:", connectionId)
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
			readReply(useBase, byte(baseCmd), byte(comm), protoParsedMessage, conn)
		}
		if err != nil {
			continue
		}
	}

	pending := true
	for pending {
		pending = false
		for i, conn := range listConnection {
			if (DEBUG_CLEANING_UP) {
				fmt.Println("Trying to read from connection", i, "and ack all pending messages")
			}
			// now all the established connections are the one to be clean up
			for {
				pendingMessage, err, comm := readPendingMessage(conn)
				if err != nil {
					break
				}
				if pendingMessage != nil {
					if (DEBUG_CLEANING_UP) {
						fmt.Println("Pending message value: ", *pendingMessage)
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

	//log.Print("Buffer read from network: ", rbuf)

	var newValue interface{}
	cmd := rbuf[0]
	switch cmd {
	case S2C_RemoteRequestAddBuddy_CMD:
		if !useBase {
			if (DEBUG_CLEANING_UP) {
				fmt.Println("Found a request add buddy message, respond now: ")
			}
			newValue = magicVarFunc("Auth_Buddy_S2C_RemoteRequestAddBuddy")
		}

	case S2C_ChatInfo2_CMD:
		if !useBase {
			if (DEBUG_CLEANING_UP) {
				fmt.Println("Found a chat message, respond now: ")
			}			
			newValue = magicVarFunc("Auth_Buddy_S2C_ChatInfo2")
		}

	case S2C_InviteMember_CMD:
		if useBase {
			if (DEBUG_CLEANING_UP) {
				fmt.Println("Found an invite member message, respond now: ")
			}
			newValue = magicVarFunc("Discussion_S2C_InviteMember")
		}

	case S2C_ChatInfo_CMD:
		if useBase {
			if (DEBUG_CLEANING_UP) {
				fmt.Println("Found a discussion chat message, respond now: ")
			}
			newValue = magicVarFunc("Discussion_S2C_ChatInfo")
		}
	default:
		return nil, nil, 0 // receive a non-offline message
	}

	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
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
	if (DEBUG_CLEANING_UP) {
		fmt.Println("Reply with message", res)
	}
	sendMsg(useBase, byte(baseCmd), replyComm, res, conn)
}