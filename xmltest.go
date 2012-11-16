package main

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
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"
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
	CONN_NUM = 3 
)

type InnerXml struct {
	Xml string `xml:",innerxml"`
}

type Message struct {
	MessageType  string `xml:"type,attr"`
	FromClient   bool   `xml:"fromClient,attr"`
	IsPartial    bool   `xml:"isPartial,attr"`
	PartialField int    `xml:"partialField,attr"`
	Command      string
	Data         InnerXml
}

type Var struct {
	Name  string `xml:"name,attr"`
	Value string
}

type TestSuite struct {
	TestSuiteName   string
	TestName        string
	TargetHost      string
	TargetPort      string
	TimesToTest     int
	VarMap          []Var     `xml:"VarMap>Var"`
	MessageSequence []Message `xml:"MessageSequence>Message"`
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


var keyMap map[string]string
var valueMap map[string]interface{}

var ts TestSuite

var messageSequence []interface{}

var hasPartial bool
var rawData []byte
var field int
var regVar *regexp.Regexp
var numVar int = 0

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "in", "test.xml", "The test input file")
	flag.Parse()
	// open file
	fmt.Println("Testing file:", inputFile)

	keyMap = make(map[string]string)
	valueMap = make(map[string]interface{})

	readXmlInput(inputFile)
	executeTestSuite(ts.TargetHost, ts.TargetPort, ts.TimesToTest)
}

/******************* steps in test ****************************************/
func executeTestSuite(targetHost string, targetPort string, timesToTest int) {
	var testCaseResult TestCaseResult

	targetAddr := targetHost + ":" + targetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Fatal("Cannot resolve address")
	}

	testCaseResult.Results = make([]TestResult, timesToTest)
	testCaseResult.NumTest = timesToTest
	testCaseResult.NumCorrect = 0
	chans := make([]chan TestResult, timesToTest)
	for i := 0; i < timesToTest; i++ {
		chans[i] = make(chan TestResult)
	}

	// parse test-suite varmap first
	parseVarMap()
	for i := 0; i < timesToTest; i++ {
		go executeOneTest(addr, chans[i])
	}

	for i := 0; i < timesToTest; i++ {
		testCaseResult.Results[i] = <-chans[i]
		if testCaseResult.Results[i].IsCorrect {
			testCaseResult.NumCorrect++
		}
	}

	fmt.Println(testCaseResult)
	prettyPrint(testCaseResult)
}

func executeOneTest(addr *net.TCPAddr, resultChan chan TestResult) {
	var result TestResult
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		//log.Fatal("connect error: ", err)
		fmt.Println(err)
		log.Print("connect error: ", err)
		result.IsCorrect = false
		result.Reason = err
	} else {
		fmt.Println("Successfully establish connection")
		// no connection error
		// testLogin(conns[i], chans[i])
		result = test(conn)
	}
	resultChan <- result
}

func test(conn *net.TCPConn) TestResult {
	fmt.Println("Start test now!!!")
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil}
	var totalTime int64 = 0

	for i, message := range ts.MessageSequence {
		parsedMessage, comm, err := parseAMessage(message)

		if err != nil {
			result.IsCorrect = false
			errMsg := fmt.Sprintf("Cannot parse message %d", i)
			result.Reason = errors.New(errMsg)
			break
		}

		protoParsedMessage := parsedMessage.(proto.Message)

		if message.FromClient { // message from client, so send to server now
			fmt.Println("Message to be sent: ", protoParsedMessage)
			begin := time.Now()
			sendMsg(byte(comm), protoParsedMessage, conn)
			end := time.Now()
			totalTime += end.Sub(begin).Nanoseconds()
		} else { // message from server, so read from buffer now
			replyMessage, err := readReply(byte(comm), protoParsedMessage, conn)
			if err == nil {
				fmt.Println("Correctly received: ", replyMessage)
				fmt.Println("Comparing received message now! ")
				if comparedResult, err := compareGetValueForProtoMessage(protoParsedMessage, *replyMessage); comparedResult {
					fmt.Println("CORRECT Reply message.")
				} else {
					fmt.Println("INCORRECT Reply message ", err)
					result.IsCorrect = false
					errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
					result.Reason = errors.New(errMsg)
					break
				}
			} else {
				fmt.Println("Encounter error while trying to parse reply message: ", err)
				result.IsCorrect = false
				errMsg := fmt.Sprintf("Message %d: %s", i+1, err.Error())
				result.Reason = errors.New(errMsg)
				break
			}
		}
	}

	logMsg := fmt.Sprintf("Test suite %s - Test %s takes %d us", ts.TestSuiteName, ts.TestName, totalTime/1000)
	log.Print(logMsg)
	result.TimeTaken = totalTime

	return result
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(cmd byte, msg proto.Message, conn *net.TCPConn) {
	log.Print("sending message command: ", cmd)
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marchaling error: ", err)
	}
	length := int32(len(data)) + 1
	log.Print("sending message length: ", length)

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, length)
	binary.Write(buf, binary.LittleEndian, cmd)
	//binary.Write(buf, binary.LittleEndian, int64(0))
	buf.Write(data)
	conn.Write(buf.Bytes())
	log.Print("sending message buffer: ", buf.Bytes())
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(expCmd byte, expMsg proto.Message, conn *net.TCPConn) (*proto.Message, error) {
	length := int32(0)
	binary.Read(conn, binary.LittleEndian, &length)
	rbuf := make([]byte, length)
	io.ReadFull(conn, rbuf)
	var err error
	var res proto.Message
	if len(rbuf) < 1 {
		errMes := fmt.Sprintf("Reply message is too short: %d", len(rbuf))
		return nil, errors.New(errMes)
	}

	log.Print("Buffer read from network: ", rbuf)

	cmd := rbuf[0]
	// command needs to be equal to expected command
	if cmd != expCmd {
		errMsg := fmt.Sprintf("Unexpected CMD %d", cmd)
		return nil, errors.New(errMsg)
	}
	fmt.Println("CMD is CORRECT, expected message type: ", reflect.ValueOf(expMsg).Type())
	newValue := reflect.New(reflect.ValueOf(expMsg).Elem().Type()).Interface()
	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully receive message", res)
	return &res, err
}

/**************Compare Values and Bind value from reply message to unbound var in value map*****/
func compareGetValueForPointer(expPtr reflect.Value, repPtr reflect.Value) (bool, error) {
	fmt.Println("Comparing pointers", expPtr.Elem().IsValid(), repPtr.Elem().IsValid())
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
	fmt.Println("Expected pointer in message has value: ", strVar)

	result := true
	var err error

	// if the expected field is a variable field that's not bound, bind the value of the key now
	key, present := keyMap[strVar]
	if present {
		fmt.Println("Bound value to map")
		// bound variable to value, return true
		valueMap[key] = repValue.Interface()
	} else {
		// compare the two value of pointer
		// if they are not struct, just compare
		if expValue.Kind() != reflect.Struct {
			isEqual := (expValue.Interface() == repValue.Interface())
			fmt.Println("Comparing: ", expValue.Interface(), repValue.Interface(), " result ", isEqual)

			if !isEqual {
				errorMsg := fmt.Sprintf("Reply field value is different from expected field value. Expect: %v, get: %v",
					expValue, repValue)
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
	fmt.Println("Comparing slices")
	result := true
	var err error
	dif := false
	if expSlice.Cap() != repSlice.Cap() {
		dif = true
	} else {
		for index := 0; index < expSlice.Cap(); index++ {
			isEqual := (expSlice.Index(index) == repSlice.Index(index))
			if !isEqual {
				dif = true
				break
			}
		}
	}
	if dif {
		errorMsg := fmt.Sprintf("Reply field value is different from expected field value. Expect: %v, get: %v",
			reflect.ValueOf(expSlice).Interface(), reflect.ValueOf(repSlice).Interface())
		err = errors.New(errorMsg)
		result = false
	}
	return result, err
}

func compareGetValueForStruct(expStruct reflect.Value, repStruct reflect.Value) (bool, error) {
	fmt.Println("Comparing structs")
	result := true
	var err error
	for i := 0; i < expStruct.NumField(); i++ {
		fmt.Println("Comparing field", i)
		expStructField := expStruct.Field(i) // pointer
		repStructField := repStruct.Field(i) // pointer

		fmt.Println(expStructField, repStructField, expStructField.Kind())
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
func parseVarMap() {
	for _, newVar := range ts.VarMap {		
		valueMap[newVar.Name] = newVar.Value
	}
}

// plug value from value map to the raw xml message
func plugValue(message string) (string, error) {
	fmt.Println("Plug value from value map to message")
	t := template.Must(template.New("Xml Message").Parse(message))
	var buffer bytes.Buffer
	err := t.Execute(&buffer, valueMap)
	fmt.Println("----------------------------------------")
	err = t.Execute(os.Stdout, valueMap)
	return buffer.String(), err
}

// pre-process a message and put unbound variable to value map
func preProcess(rawData string) {
	fmt.Println("Preprocess ", rawData)
	regVar = regexp.MustCompile("{{.([a-zA-Z0-9]*)}}")
	var processedData = rawData
	listVar := regVar.FindAllStringSubmatch(rawData, -1)
	for _, matchedValue := range listVar {
		fmt.Print(matchedValue[1])
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

// parse a message from message sequence
func parseAMessage(v Message) (interface{}, int, error) {
	// first pass: parse the message and put new var to value map
	preProcess(v.Data.Xml)

	// second pass: plug all value to the var in the raw string 
	addedXmlMessage, err := plugValue(v.Data.Xml)

	fmt.Println("After second pass, message = ", addedXmlMessage)
	message := magicVarFunc(v.MessageType)	
	err = xml.Unmarshal([]byte(addedXmlMessage), message)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	s := reflect.ValueOf(message).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		printValue(f)
	}

	cmd, _ := strconv.ParseInt(v.Command, 0, 0)
	return message, int(cmd), err
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
			fmt.Println("Invalid value")
		}
	} else {
		fmt.Println("Not a pointer")
	}
}


/********************* Input - output helper *********************************/

func readXmlInput(inputFile string) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, &ts)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
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