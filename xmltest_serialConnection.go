package main

import (
	"errors"

	"code.google.com/p/goprotobuf/proto"
	"log"
	"net"
	"xmltest/btalkTest/Auth_C2S"
	"xmltest/btalkTest/Auth_S2C"
	//"garena.com/btalkTest/TOKEN_S"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
	"time"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"reflect"
	"xmltest/btalkTest/GeneratedDataStructure"
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

type Message struct {
	MessageType string `xml:"type,attr"`
	Data        string `xml:",innerxml"`
}

type TestSuite struct {
	TestSuiteName string
	TestName      string
	TargetHost    string
	TargetPort    string
	TimesToTest   int
	InputType     string
	InputData     []Message `xml:"InputData>Message"`
	OutputData    []Message `xml:"OutputData>Message"`
}

var inputMessages []interface{}
var outputMessages []interface{}

func main() {
	// open file
	data, err := ioutil.ReadFile("test.xml")
	if err != nil {
		panic(err)
	}

	var ts TestSuite
	err = xml.Unmarshal(data, &ts)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("Readed input %v", ts.InputData)

	for i, v := range ts.InputData {
		fmt.Printf("Readed input Number %d \n----------------------------\n%v\n---------------------------\n", i, v.Data)

		input := MagicVarFunc(v.MessageType)

		err = xml.Unmarshal([]byte(v.Data), input)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("\n----------------------------\n Parsed input %v\n---------------------------\n", input)

		s := reflect.ValueOf(input).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			printValue(f)
		}
		inputMessages = append(inputMessages, input)
		fmt.Println("input", input)
	}

	for i, v := range ts.OutputData {
		fmt.Printf("Readed output Number %d \n----------------------------\n%v\n---------------------------\n", i, v.Data)

		message := MagicVarFunc(v.MessageType)

		err = xml.Unmarshal([]byte(v.Data), message)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("\n----------------------------\n Parsed output %v\n---------------------------\n", message)
		s := reflect.ValueOf(message).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			printValue(f)
		}
		outputMessages = append(outputMessages, message)
	}
	fmt.Println(inputMessages, outputMessages)
	execute(ts.TargetHost, ts.TargetPort, ts.TimesToTest, inputMessages, outputMessages)

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

func MagicVarFunc(action string) interface{} {
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
	fmt.Println("-------------------------TEST RESULT--------------------------------------")
	for i := 0; i < result.NumTest; i++ {
		fmt.Println("Case", i, ":")
		aResult := result.Results[i]
		if aResult.IsCorrect {
			fmt.Println("	CORRECT.   Time taken:", (aResult.TimeTaken / 1000), "us")
		} else {
			fmt.Println("   INCORRECT. Reason: ", aResult.Reason)
		}
	}
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("Total result: ", result.NumCorrect, "/", result.NumTest, " = ", (result.NumCorrect * 100.0 / result.NumTest), "%")
	fmt.Println("-------------------------END OF TEST RESULT-------------------------------")
}

func execute(targetHost string, targetPort string, timesToTest int, inputMessages []interface{}, outputMessages []interface{}) {
	var testCaseResult TestCaseResult

	targetAddr := targetHost + ":" + targetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Fatal("Cannot resolve address")
	}

	conns := make([]*net.TCPConn, timesToTest)
	for i := 0; i < timesToTest; i++ {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			log.Fatal("connect error: ", err)
		}
		conns[i] = conn
	}

	testCaseResult.Results = make([]TestResult, timesToTest)
	testCaseResult.NumTest = timesToTest
	testCaseResult.NumCorrect = 0
	chans := make([]chan TestResult, timesToTest)
	for i := 0; i < timesToTest; i++ {
		chans[i] = make(chan TestResult)
	}

	for i := 0; i < timesToTest; i++ {
		go executeOneTest(conns[i], inputMessages, outputMessages, chans[i])
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

func executeOneTest(conn *net.TCPConn, inputMessages []interface{}, outputMessages []interface{}, resultChan chan TestResult) {
	var result TestResult
	// no connection error
	inputMessage := inputMessages[0].(proto.Message)
	outputMessage := outputMessages[0].(proto.Message)
	// testLogin(conns[i], chans[i])
	result = testHello(conn, &inputMessage, &outputMessage)
	resultChan <- result
}

func testHello(conn *net.TCPConn, inMessagePtr *proto.Message, outMessagePtr *proto.Message) TestResult {
	fmt.Println("Testing Log in now ! with message: ", *inMessagePtr, " expecting: ", *outMessagePtr)
	var result TestResult

	begin := time.Now()
	sendMsg(conn, byte(C2S_HelloInfo_CMD), *inMessagePtr)
	end := time.Now()
	replyMessage, err := readReply(conn)

	log.Print("Hello takes ", end.Sub(begin).Nanoseconds()/1000, " us")

	result.TimeTaken = end.Sub(begin).Nanoseconds()

	if err == nil {
		if comparedResult, err := compareReply(*replyMessage, *outMessagePtr); comparedResult {
			fmt.Println("CORRECT reply.")
			result.IsCorrect = true
			result.Reason = nil
		} else {
			fmt.Println("INCORRECT reply! Reason: ", err)
			result.IsCorrect = false
			result.Reason = err
		}
	} else {
		fmt.Println("Encounter error while trying to send: ", err)
		result.IsCorrect = false
		result.Reason = err
	}
	return result
}

func testLogin(conn *net.TCPConn, endSignal chan bool) {
	for i := 0; i < 20; i++ {
		begin := time.Now()
		sendHello(conn)
		readReply(conn)
		//sendLogin(conn)
		//readReply(conn)
		end := time.Now()
		log.Print("login takes ", end.Sub(begin).Nanoseconds()/1000, " us")
	}
	endSignal <- true
}

func testGarena(conn *net.TCPConn, endSignal chan bool) {
	for i := 0; i < 1; i++ {
		begin := time.Now()
		sendHello(conn)
		readReply(conn)
		sendGarenaAuth(conn)
		readReply(conn)
		end := time.Now()
		log.Print("garena login takes ", end.Sub(begin).Nanoseconds()/1000, " us")
	}
	endSignal <- true
}

func testOAuth(conn *net.TCPConn, endSignal chan bool) {
	for i := 0; i < 20; i++ {
		begin := time.Now()
		sendHello(conn)
		readReply(conn)
		sendOAuth(conn)
		readReply(conn)
		end := time.Now()
		log.Print("login takes ", end.Sub(begin).Nanoseconds()/1000, " us")
	}
	endSignal <- true
}

func sendMsg(conn *net.TCPConn, cmd byte, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marchaling error: ", err)
	}
	length := int32(len(data)) + 1
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, length)
	binary.Write(buf, binary.LittleEndian, cmd)
	//binary.Write(buf, binary.LittleEndian, int64(0))
	buf.Write(data)
	conn.Write(buf.Bytes())
}

func sendLogin(conn *net.TCPConn) {
	m := md5.New()
	io.WriteString(m, "123456")
	password := string(m.Sum(nil))
	password = hex.EncodeToString([]byte(password))
	loginInfo := &Auth_C2S.LoginInfo{
		Name: proto.String("	"),
		Passworld:       proto.String(password),
		ClientType:      proto.Int32(1),
		MachineId:       proto.String("ddd"),
		SoftwareVersion: proto.Int32(1),
	}
	sendMsg(conn, byte(C2S_LoginInfo_CMD), loginInfo)
}

func sendHello(conn *net.TCPConn) {
	helloInfo := &Auth_C2S.HelloInfo{
		ClientType: proto.Int32(1),
		Version:    proto.Uint32(1),
	}
	sendMsg(conn, byte(C2S_HelloInfo_CMD), helloInfo)
}

func sendOAuth(conn *net.TCPConn) {
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
	sendMsg(conn, byte(C2S_OAuthLogin_CMD), oauth)
}

func sendGarenaAuth(conn *net.TCPConn) {
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
	sendMsg(conn, byte(C2S_OAuthLogin_CMD), oauth)
}

// read reply to a buffer
func readReply(conn *net.TCPConn) (*proto.Message, error) {
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

	cmd := int(rbuf[0])
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
		fmt.Println("cmd = unexpected", cmd)
	}
	return &res, err
}

func compareReply(rep proto.Message, out proto.Message) (bool, error) {
	var errorMsg string
	fmt.Println(rep, "|", out)

	replyMessage := reflect.ValueOf(rep).Elem()
	outputMessage := reflect.ValueOf(out).Elem()

	if reflect.ValueOf(rep).Type() != reflect.ValueOf(out).Type() {
		errorMsg = fmt.Sprintf("Reply has different type from expected. Expect: %d, get: %d", reflect.ValueOf(rep).Type(), reflect.ValueOf(out).Type())
		return false, errors.New(errorMsg)
	}

	for i := 0; i < outputMessage.NumField(); i++ {
		fmt.Println("Comparing field ", i)
		var dif bool = false
		ptrOut := outputMessage.Field(i) // pointer
		ptrRep := replyMessage.Field(i)  // pointer
		fmt.Println(ptrOut, ptrRep, ptrOut.Kind())

		if ptrOut.Kind() == reflect.Ptr {
			outValue := reflect.Indirect(reflect.ValueOf(ptrOut.Interface()))
			repValue := reflect.Indirect(reflect.ValueOf(ptrRep.Interface()))
			compare := (outValue.Interface() == repValue.Interface())
			fmt.Println("Comparing: ", outValue.Interface(), repValue.Interface(), " result ", compare)
			if !compare {
				dif = true
				fmt.Println(outValue, repValue)
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
					fmt.Println("Comparing: ", ptrOut.Index(index), ptrRep.Index(index), " result ", compare)

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
			fmt.Println("They are the same")
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
		//fmt.Println(value.Kind(), value.Elem())
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			return value
		}
	}

	return ptrValue
}
