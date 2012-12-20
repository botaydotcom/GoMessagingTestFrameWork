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
	C2S_RequestBuddySimpleInfoList_CMD  = 0x64
	C2S_DeleteBuddy_CMD 				= 0x67

	C2S_RequestMyDiscussion_CMD			= 0x09
	C2S_LeaveDiscussion_CMD				= 0x03
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
	
	S2C_BuddySimpleInfoList_CMD	  = 0x65
	S2C_UserDiscussionList_CMD	  = 0x02

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
	CONN_NUM     = 3
	READ_TIMEOUT = "3s"
	DEFAULT_LINGER_PERIOD = 5
)

// DEBUG FLAGS
const (
	DEBUG_COMPARE_POINTER = false
	DEBUG_COMPARE_SLICE   = false
	DEBUG_COMPARE_STRUCT  = false

	DEBUG_PARSING_MESSAGE = false
	DEBUG_CLEANING_UP	  = true

	DEBUG_SENDING_MESSAGE = true
	DEBUG_READING_MESSAGE = true
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

type CleanUpSuite struct {
	CleanUpSuiteName string
	TargetHost    string
	TargetPort    string
	GlobalVarMap  []Var       `xml:"VarMap>Var"`
	CleanUpSequence []Message `xml:"CleanUpSequence>Message"`
}

var keyMap map[string]string
var valueMap map[string]interface{}

var sliceKeyMap map[string]string
var sliceMap map[string]interface{}

var hasPartial bool
var rawData []byte
var field int
var regVar *regexp.Regexp
var numVar int = 0

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "in", "cleanUp.xml", "The clean up input file")
	flag.Parse()
	// open file
	fmt.Println("Testing file:", inputFile)

	keyMap = make(map[string]string)
	valueMap = make(map[string]interface{})
	sliceKeyMap = make(map[string]string)
	sliceMap = make(map[string]interface{})

	ts, err := readXmlInput(inputFile)

	if err != nil {
		errorMsg := fmt.Sprintf("Cannot read input file because %v", err)
		log.Fatal(errorMsg)
	} else {
		startCleanUp(ts)
	}

}

/******************* steps in test ****************************************/
func startCleanUp(cleanUpSuite *CleanUpSuite) {
	fmt.Println("Start cleaning up now:")

	targetAddr := cleanUpSuite.TargetHost + ":" + cleanUpSuite.TargetPort
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Fatal("Cannot resolve address")
	}

	fmt.Println("Parse var map:")
	// parse test-suite varmap first
	parseVarMap(cleanUpSuite.GlobalVarMap)
	
	fmt.Println("Create connection list to clean up: .....")
	listConnection := createConnection(addr, cleanUpSuite)

	fmt.Println("Ack pending messages for all users: .....")
	ackPendingMessageForAllUsers(listConnection)

	fmt.Println("Remove all discussions: .....")
	removeAllDiscussionsForAllUsers(listConnection)

	fmt.Println("Remove all buddies relations: .....")
	removeAllBuddiesForAllUsers(listConnection)

	for i := range listConnection {
		listConnection[i].Close()
	}
}

func createConnection(addr *net.TCPAddr, cleanUpSuite *CleanUpSuite) (map[int]*net.TCPConn){
	var listConnection = make(map[int]*net.TCPConn)

	// these messages are supposed to be login messages !
	for i, message := range cleanUpSuite.CleanUpSequence {
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
	return listConnection
}

func ackPendingMessageForAllUsers(listConnection map[int]*net.TCPConn) {
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
		}
	}
}

func removeAllBuddiesForAllUsers(listConnection map[int]*net.TCPConn) {
	for {
		stop := true
		for i := range listConnection {
			fmt.Println("Remove buddies for connection ", i)
			conn := listConnection[i]

			queryBuddyMessage := magicVarFunc("Auth_Buddy_C2S_RequestBuddySimpleInfoList")
			protoMessage := queryBuddyMessage.(proto.Message)

			sendMsg(false, byte(0), C2S_RequestBuddySimpleInfoList_CMD, protoMessage, conn)
			expMsg := magicVarFunc("Auth_Buddy_S2C_BuddySimpleInfoList")

			listBuddyMessagePtr, err := readReply(false, 0, S2C_BuddySimpleInfoList_CMD, expMsg.(proto.Message), conn)
			if (err != nil) {
				fmt.Println(err)
				if (err.Error()[0:8] != "Cant read") {
					stop = false
				}
				continue
			}
			listBuddyMessage := *listBuddyMessagePtr

			listBuddy := reflect.Indirect(reflect.ValueOf(listBuddyMessage)).FieldByName("Buddies")

			listBuddyToDelete := make([]int32, listBuddy.Len())
			for index := 0; index < listBuddy.Len(); index++ {
				buddy := listBuddy.Index(index)
				buddyId := reflect.Indirect(buddy).FieldByName("UserId")
				listBuddyToDelete[index] = int32(reflect.Indirect(buddyId).Int())
			}

			deleteMessage := magicVarFunc("Auth_Buddy_C2S_DeleteBuddy")
			structValue := reflect.Indirect(reflect.ValueOf(deleteMessage))
			structValue.FieldByName("UserId").Set(reflect.ValueOf(listBuddyToDelete))			
			sendMsg(false, byte(0), C2S_DeleteBuddy_CMD, deleteMessage.(proto.Message), conn)
		}
		if (stop) {
			break;
		}

	}
}

func removeAllDiscussionsForAllUsers(listConnection map[int]*net.TCPConn) {
	for i := range listConnection {
		conn := listConnection[i]

		queryDiscussionMessage := magicVarFunc("Discussion_C2S_RequestMyDiscussion")
		protoMessage := queryDiscussionMessage.(proto.Message)
		sendMsg(true, byte(DISCUSSION_PACKET_BASE_COMMAND), 
			C2S_RequestMyDiscussion_CMD, protoMessage, conn)
		expMsg := magicVarFunc("Discussion_S2C_UserDiscussionList")

		listDiscussionMessagePtr, _ := readReply(true, DISCUSSION_PACKET_BASE_COMMAND, 
			S2C_UserDiscussionList_CMD, expMsg.(proto.Message), conn)

		listDiscussion := reflect.Indirect(reflect.ValueOf(*listDiscussionMessagePtr)).FieldByName("Discussion")
		for index := 0; index < listDiscussion.Len(); index++ {
			discussion := reflect.Indirect(listDiscussion.Index(index))
			discussionId := discussion.FieldByName("DiscussionId")

			deleteMessage := magicVarFunc("Discussion_C2S_LeaveDiscussion")
			structValue := reflect.Indirect(reflect.ValueOf(deleteMessage))
			structValue.FieldByName("DiscussionId").Set(discussionId)
			sendMsg(true, byte(DISCUSSION_PACKET_BASE_COMMAND), C2S_LeaveDiscussion_CMD, deleteMessage.(proto.Message), conn)
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
	duration, err := time.ParseDuration(READ_TIMEOUT)
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
		log.Print("sending message with length: ", length, " base command / command / message: ", baseCmd, cmd, msg)
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, length)
	if (useBase) {
		binary.Write(buf, binary.LittleEndian, baseCmd)
	}

	binary.Write(buf, binary.LittleEndian, cmd)
	//binary.Write(buf, binary.LittleEndian, int64(0))
	buf.Write(data)
	numSent, err := conn.Write(buf.Bytes())
	if (err != nil) {
		log.Fatal("error while sending data to server ", err, " num bytes sent: ", numSent)
	} else {
		fmt.Println("MESSAGE SUCCESSFULLY SENT", numSent)
	}
}

// read a reply to a buffer based on the expected message type
// return error if reply message has different type of command than expected
func readReply(useBase bool, expBaseCmd byte, expCmd byte, expMsg proto.Message, conn *net.TCPConn) (*proto.Message, error) {
	length := int32(0)
	duration, err := time.ParseDuration(READ_TIMEOUT)
	timeNow := time.Now()

	err = conn.SetReadDeadline(timeNow.Add(duration))
	if err != nil {
		return nil, errors.New("Cant read: Cant set deadline")
	}

	err = binary.Read(conn, binary.LittleEndian, &length)
	fmt.Println("READING:", length, "FROM SOCKET")
	if err != nil {
		fmt.Println("Cant read: Cant read length from socket")
		return nil, errors.New("Cant read: Cant read length from socket")
	}

	if (useBase) {
		length = length - 1
		var baseCmd byte
		err = binary.Read(conn, binary.LittleEndian, &baseCmd)
		if err != nil {
			fmt.Println("Cant read: Cant read base command from socket")
			return nil, errors.New("Cant read: Cant read message from socket")
		}

		if (baseCmd != expBaseCmd) {
			// finish reading the rest of the message, 
			// so that it does not affects other message
			rbuf := make([]byte, length)
			io.ReadFull(conn, rbuf)
			// report error
			fmt.Println("Unexpected BASE CMD")
			errMsg := fmt.Sprintf("Unexpected BASE CMD %d", baseCmd)
			return nil, errors.New(errMsg)	
		}
	}

	rbuf := make([]byte, length)
	_, err = io.ReadFull(conn, rbuf)
	if err != nil {
		fmt.Println("Cant read: Cant read full message from sockets")
		return nil, err
	}

	var res proto.Message
	if len(rbuf) < 1 {
		fmt.Println("Reply message is too short")
		errMes := fmt.Sprintf("Reply message is too short: %d", len(rbuf))
		return nil, errors.New(errMes)
	}

	cmd := rbuf[0]
	// command needs to be equal to expected command
	if cmd != expCmd {
		fmt.Println("Unexpected CMD", cmd)
		errMsg := fmt.Sprintf("Unexpected CMD %d", cmd)
		return nil, errors.New(errMsg)
	}
	
	newValue := reflect.New(reflect.ValueOf(expMsg).Elem().Type()).Interface()
	res = newValue.(proto.Message)
	err = proto.Unmarshal(rbuf[1:], res)

	if err != nil {
		fmt.Println("Cannot unmarshal")
		log.Fatal(err)
	}

	if (DEBUG_READING_MESSAGE) {
		log.Print("Successfully receive message from network: ", reflect.ValueOf(expMsg).Type(), res)
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
	//t := template.Must(template.New("Xml Message").Parse(message))
	//var buffer bytes.Buffer
	//err := t.Execute(&buffer, valueMap)
	//fmt.Println("----------------------------------------")
	//err := t.Execute(os.Stdout, valueMap)
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

func readXmlInput(inputFile string) (*CleanUpSuite, error) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var cleanUpSuite CleanUpSuite
	err = xml.Unmarshal(data, &cleanUpSuite)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}
	return &cleanUpSuite, nil
}