package main

import (
	"code.google.com/p/goprotobuf/proto"
	"garena.com/btalkTest/Auth_C2S"
	"garena.com/btalkTest/Auth_S2C"
	"log"
	"net"
	//"garena.com/btalkTest/TOKEN_S"
	"encoding/binary"
	"io"
	//"fmt"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"time"
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

func Execute(inputMessages []interface{}, outputMessages []interface{}) {

	addr, err := net.ResolveTCPAddr("tcp", "203.117.155.188:9100")
	if err != nil {
		log.Fatal("wrong address")
	}
	conns := make([]*net.TCPConn, CONN_NUM)
	for i := 0; i < CONN_NUM; i++ {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			log.Fatal("connect error: ", err)
		}
		conns[i] = conn
	}
	chans := make([]chan bool, CONN_NUM)

	for i := 0; i < CONN_NUM; i++ {
		chans[i] = make(chan bool)
	}

	for i := 0; i < CONN_NUM; i += 3 {
		go testLogin(conns[i], chans[i])
		//go testOAuth(conns[i+1], chans[i+1])
		//go testGarena(conns[i+2], chans[i+2])
	}

	for i := 0; i < CONN_NUM; i++ {
		//    <-chans[i]
		<-chans[2]
	}
}

func testLogin(conn *net.TCPConn, endSignal chan bool) {
	for i := 0; i < 20; i++ {
		begin := time.Now()
		sendHello(conn)
		readReply(conn)
		sendLogin(conn)
		readReply(conn)
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
	log.Print("length: ", length)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, length)
	binary.Write(buf, binary.LittleEndian, cmd)
	//binary.Write(buf, binary.LittleEndian, int64(0))
	buf.Write(data)
	conn.Write(buf.Bytes())
}

func sendHello(conn *net.TCPConn) {
	helloInfo := &Auth_C2S.HelloInfo{
		ClientType: proto.Int32(1),
		Version:    proto.Uint32(1),
	}
	sendMsg(conn, byte(C2S_HelloInfo_CMD), helloInfo)
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

func readReply(conn *net.TCPConn) {
	length := int32(0)
	binary.Read(conn, binary.LittleEndian, &length)
	rbuf := make([]byte, length)
	io.ReadFull(conn, rbuf)
	if len(rbuf) < 9 {
		return
	}
	cmd := int(rbuf[8])
	switch cmd {
	case S2C_HelloInfoResult_CMD:
		res := &Auth_S2C.HelloInfoResult{}
		err := proto.Unmarshal(rbuf[9:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res.GetToken())
	case S2C_LoginUserInfo_CMD:
		res := &Auth_S2C.LoginUserInfo{}
		err := proto.Unmarshal(rbuf[9:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(res.GetMyInfo(), res.GetServerTime())
	case ERROR_ErrorInfo_CMD:
		res := &Auth_S2C.ErrorInfo{}
		err := proto.Unmarshal(rbuf[9:], res)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("error type: ", res.GetType())
	default:
		fmt.Println("Unexpected CMD: ", cmd)
	}
}
