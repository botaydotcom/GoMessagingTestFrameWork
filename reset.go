package main

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	"net"
	"strconv"
	"xmltest/btalkTest/Auth_C2S"
	"xmltest/btalkTest/Auth_S2C"
	//"garena.com/btalkTest/TOKEN_S"
	"bytes"
	//"crypto/md5"
	"encoding/binary"
	"flag"
	//"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	//"os"
	//"text/template"
	"time"
)

const (
	C2S_DeleteAccount_CMD = 0x0c
)

const (
	ERROR_ErrorInfo_CMD = 0x20
)

const (
	READ_PENDING_TIMEOUT = "3s"
	CONN_NUM             = 3
)

type TestSuite struct {
	TargetHost string
	TargetPort string
	UserId     string
}

var userIdToDelete string
var userId int32

func main() {
	var inputFile string
	inputFile = "reset.xml"
	//flag.StringVar(&inputFile, "in", "test.xml", "The test input file")
	flag.StringVar(&userIdToDelete, "user", "", "The user id to delete")

	flag.Parse()

	ts, err := readXmlInput(inputFile)

	if userIdToDelete == "" {
		userIdToDelete = ts.UserId
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Cannot read input file because %v", err)
		timeEncodedPrintln(errorMsg)
		log.Fatal(errorMsg)
	}
	value, err := strconv.ParseInt(userIdToDelete, 0, 0)

	if err != nil {
		errorMsg := fmt.Sprintf("Cannot parse user id because %v", err)
		timeEncodedPrintln(errorMsg)
		log.Fatal(errorMsg)
	} else {
		userId = int32(value)
		executeTestSuite(ts)
	}
}

/******************* steps in test ****************************************/
func executeTestSuite(testSuite *TestSuite) {
	targetAddr := testSuite.TargetHost + ":" + testSuite.TargetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		timeEncodedPrintln("Cannot resolve address")
		log.Fatal("Cannot resolve address")
	}
	conn, err := net.DialTCP("tcp", nil, addr)

	sendDeleteMessageToServer(userId, conn)
	err, protoMessage := readReply(conn)
	if (err == nil) {
		prettyPrintTestSuite(testSuite, true, nil, nil)
	} else {
		prettyPrintTestSuite(testSuite, false, err, protoMessage)
	}
	conn.Close()
}

/***************Send and Receive Messages********************************/
// send a message with command and message using given connection
func sendMsg(conn *net.TCPConn, cmd byte, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		timeEncodedPrintln("marshaling error: ", err)
		log.Fatal("marshaling error: ", err)
	}
	length := int32(len(data)) + 1

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, length)
	
	binary.Write(buf, binary.LittleEndian, cmd)
	buf.Write(data)
	fmt.Println("writing:", cmd, msg, "to output socket")
	conn.Write(buf.Bytes())
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(conn *net.TCPConn) (error, proto.Message) {
	duration, _ := time.ParseDuration("10s")
	timeNow := time.Now()
	err := conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		fmt.Print("Error: Cannot set read timeout", err.Error())
		return err, nil
	}
	length := int32(0)
	err = binary.Read(conn, binary.LittleEndian, &length)
	if err != nil {
		return nil, nil
	}

	rbuf := make([]byte, length)
	io.ReadFull(conn, rbuf)
	fmt.Println(rbuf)
	
	cmd := int(rbuf[0])

	switch cmd {
	case ERROR_ErrorInfo_CMD:
		res := &Auth_S2C.ErrorInfo{}
		err := proto.Unmarshal(rbuf[1:], res)
		if err != nil {
			log.Fatal(err)
		}
		msg := fmt.Sprint("Server returns error: ")
		return errors.New(msg), res
	default:
		msg := fmt.Sprint("Unexpected CMD: ", cmd)
		return errors.New(msg), nil
	}
	return nil, nil

}

/************************* Pointer value helper functions *************************/
func readXmlInput(inputFile string) (*TestSuite, error) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var testSuite TestSuite
	err = xml.Unmarshal(data, &testSuite)
	if err != nil {
		timeEncodedPrintf("error: %v\n", err)
		return nil, err
	}
	return &testSuite, nil
}

func prettyPrintTestSuite(testSuite *TestSuite, isSuccessful bool, err error, msg proto.Message) {
	fmt.Println("--------------------------------------------------------------------------------")
	if isSuccessful {
		fmt.Println("SUCCEEDED!")
	} else {
		fmt.Println("FAILED")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if msg != nil {
			fmt.Println("Server error message :", msg)
		}
	}
	fmt.Println("END OF TEST SUITE RESULT")
	fmt.Println("--------------------------------------------------------------------------------")

	
}


func sendDeleteMessageToServer(userId int32, conn *net.TCPConn) {
	message := &Auth_C2S.DeleteAccount{
		Userid: proto.Int32(userId),
	}
	sendMsg(conn, byte(C2S_DeleteAccount_CMD), message)
}

func timeEncodedPrintf(format string, a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Printf(format, a...)
}

func timeEncodedPrint(a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Print(a...)
}

func timeEncodedPrintln(a ...interface{}) {
	fmt.Print(time.Now().Format("Jan _2 15:04:05"), ": ")
	fmt.Println(a...)
}
