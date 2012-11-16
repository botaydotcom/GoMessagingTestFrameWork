package main

import (
	"strconv"
	"errors"
	"regexp"
	"code.google.com/p/goprotobuf/proto"
	"log"
	"net"
	//"xmltest/btalkTest/Auth_C2S"
	"xmltest/btalkTest/Auth_S2C"
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
	"reflect"
	"time"
	"text/template"
	"xmltest/btalkTest/GeneratedDataStructure"
	"os"
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
	CONN_NUM = 3 //has to be even
)

type InnerXml struct {
	Xml 		string `xml:",innerxml"`
}

type Message struct {
	MessageType  string `xml:"type,attr"`
	FromClient   bool   `xml:"fromClient,attr"`
	IsPartial    bool   `xml:"isPartial,attr"`
	PartialField int    `xml:"partialField,attr"`
	Command      string 
	Data         InnerXml 
}

type TestSuite struct {
	TestSuiteName string
	TestName      string
	TargetHost    string
	TargetPort    string
	TimesToTest   int
	InputType     string
	MessageSequence     []Message `xml:"MessageSequence>Message"`
}


var keyMap 		  map[string] string

var valueMap 	  map[string] interface{}

var defaultMap    map[string] string


var ts TestSuite

var messageSequence []interface{}

var hasPartial bool
var rawData []byte
var field int
	
var regVar *regexp.Regexp

var numVar int = 0
	
/*
func parseAMessage(v Message) (interface{}, int, error) {
	message := magicVarFunc(v.MessageType)
	err = xml.Unmarshal([]byte(v.Data.Xml), message)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	s := reflect.ValueOf(message).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		printValue(f)
	}
		
	cmd, _ := strconv.ParseInt(v.Command, 0, 0)
	outputMessages = append(outputMessages, message)
	outputCommands = append(outputCommands, int(cmd))
	return message, cmd, err
}*/

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

func main() {	
	var inputFile string
	flag.StringVar(&inputFile, "in", "test.xml", "The test input file")
	flag.Parse()
	// open file
	fmt.Println("Testing file:", inputFile)

	keyMap = make(map[string] string)
	valueMap = make(map[string] interface{})
	defaultMap = make(map[string] string)
	
	readXmlInput(inputFile)

	executeTestSuite(ts.TargetHost, ts.TargetPort, ts.TimesToTest)
}

func magicVarFunc(action string) interface{} {
	return reflect.New(GeneratedDataStructure.Map[action]).Interface()
}

/******************************************************************************/

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

func plugValue(message string) (string, error) {
	fmt.Println("Plug in value to template now!")
	t := template.Must(template.New("Xml Message").Parse(message))
	var buffer bytes.Buffer 
	err := t.Execute(os.Stdout, valueMap)
	err = t.Execute(&buffer, valueMap)
	fmt.Println("----------------------------------------")
	return buffer.String(), err
}

func plugDefaultValue(message string) (string, error) {
	fmt.Println("Plug in value to template now!")
	t := template.Must(template.New("Xml Message").Parse(message))
	var buffer bytes.Buffer 
	err := t.Execute(os.Stdout, defaultMap)
	err = t.Execute(&buffer, defaultMap)
	fmt.Println("----------------------------------------")
	return buffer.String(), err
}


func parseAMessage(v Message) (interface{}, int, error) {
	message := magicVarFunc(v.MessageType)
	addedXmlMessage, err := plugValue(v.Data.Xml)
	err =  xml.Unmarshal([]byte(addedXmlMessage), message)

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

//(interface{}, int, error)
func preProcess(v Message)  {
	regVar = regexp.MustCompile("{{.([a-zA-Z0-9]*)}}")
	var rawData = v.Data.Xml
	var processedData = rawData
	fmt.Println("Matching ", rawData)
	listVar := regVar.FindAllStringSubmatch(rawData, -1)
	for _, matchedValue := range listVar {
		fmt.Print(matchedValue[1])
		key := matchedValue[1]
		strMarkValue := ""

		_, present := defaultMap[key]
		fmt.Println(" key present: ", present)
		if !present {
			newMarkValue := START_MARK_VALUE - numVar
			numVar++
			strMarkValue = strconv.Itoa(newMarkValue)
			keyMap[strMarkValue] = key
			defaultMap[key] = strMarkValue
			valueMap[key] = strMarkValue
			// now replace
			reg, _ := regexp.Compile("{{." + key + "}}")
			processedData = reg.ReplaceAllString(processedData, defaultMap[key])
		}
	}	
	fmt.Println("After matching to value", processedData)
}

func compareGetValueForPointer(expPtr reflect.Value, repPtr reflect.Value) (bool, error){
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
		if (expValue.Kind() != reflect.Struct) {
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

func compareGetValueForSlice(expSlice reflect.Value, repSlice reflect.Value) (bool, error){
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

func compareGetValueForStruct(expStruct reflect.Value, repStruct reflect.Value) (bool, error){
	fmt.Println("Comparing structs")
	result := true
	var err error
	for i := 0; i < expStruct.NumField(); i++ {
		fmt.Println("Comparing field", i)
		expStructField := expStruct.Field(i) 	// pointer
		repStructField := repStruct.Field(i)  	// pointer
		
		fmt.Println(expStructField, repStructField, expStructField.Kind())
		var err1 error
		var isEqual bool
		if (expStructField.Kind() == reflect.Ptr) {
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
			break;
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

func test(conn *net.TCPConn) TestResult {
	fmt.Println("Start test now!!!")
	var result = TestResult{
		IsCorrect: true,
		Reason:    nil}
	var totalTime int64 = 0

	for i, message := range ts.MessageSequence {
		preProcess(message)
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
				fmt.Println("Comparing received message now! ");
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

// read reply to a buffer
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
	
	
	
	
	
	
	
	switch cmd {
	case S2C_HelloInfoResult_CMD:
		fmt.Println("cmd = helloinforresult")
		res = &Auth_S2C.HelloInfoResult{}
		err = proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res.(*Auth_S2C.HelloInfoResult).GetToken())
	case S2C_LoginUserInfo_CMD:
		fmt.Println("cmd = login user info")
		res = &Auth_S2C.LoginUserInfo{}
		err = proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res.(*Auth_S2C.LoginUserInfo).GetMyInfo(), res.(*Auth_S2C.LoginUserInfo).GetServerTime())
	case ERROR_ErrorInfo_CMD:
		fmt.Println("cmd = error info")
		res = &Auth_S2C.ErrorInfo{}
		err = proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("error type: ", res.(*Auth_S2C.ErrorInfo).GetType())
	default:
		errMsg := fmt.Sprintf("Unexpected CMD %d", cmd)
		return nil, errors.New(errMsg)
	}
	return &res, err
}

func compareReply(rep proto.Message, out proto.Message) (bool, error) {
	var errorMsg string
	//fmt.Println(rep, "|", out)

	replyMessage := reflect.ValueOf(rep).Elem()
	outputMessage := reflect.ValueOf(out).Elem()

	if reflect.ValueOf(rep).Type() != reflect.ValueOf(out).Type() {
		errorMsg = fmt.Sprintf("Reply has different type from expected. Expect: %d, get: %d", reflect.ValueOf(rep).Type(), reflect.ValueOf(out).Type())
		return false, errors.New(errorMsg)
	}

	for i := 0; i < outputMessage.NumField(); i++ {
		//fmt.Println("Comparing field ", i)
		var dif bool = false
		ptrOut := outputMessage.Field(i) // pointer
		ptrRep := replyMessage.Field(i)  // pointer
		//fmt.Println(ptrOut, ptrRep, ptrOut.Kind())

		if ptrOut.Kind() == reflect.Ptr {
			outValue := reflect.Indirect(reflect.ValueOf(ptrOut.Interface()))
			repValue := reflect.Indirect(reflect.ValueOf(ptrRep.Interface()))
			compare := (outValue.Interface() == repValue.Interface())
			//fmt.Println("Comparing: ", outValue.Interface(), repValue.Interface(), " result ", compare)
			if !compare {
				dif = true
				//fmt.Println(outValue, repValue)
			}
			if dif {
				errorMsg = fmt.Sprintf("Reply field value is different from expected field value. Expect: %v, get: %v",
					outValue, repValue)
			}
		} else if ptrOut.Kind() == reflect.Slice {
			if ptrOut.Cap() != ptrRep.Cap() {
				dif = true
			} else {
				for index := 0; index < ptrOut.Cap(); index++ {
					compare := (ptrOut.Index(index) == ptrRep.Index(index))
					//fmt.Println("Comparing: ", ptrOut.Index(index), ptrRep.Index(index), " result ", compare)

					if !compare {
						dif = true
						break
					}
				}
			}
			if dif {
				errorMsg = fmt.Sprintf("Reply field value is different from expected field value. Expect: %v, get: %v",
					reflect.ValueOf(ptrOut).Interface(), reflect.ValueOf(ptrRep).Interface())
			}
		}

		if dif {
			return false, errors.New(errorMsg)
		} else {
			//fmt.Println("They are the same")
		}
	}

	return true, nil
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
