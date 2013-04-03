package main
import (
    "net"
    "log"
    "code.google.com/p/goprotobuf/proto"
    "garena.com/btalkTest/Auth_C2S"
    "garena.com/btalkTest/Auth_S2C"
    "encoding/binary"
    "encoding/xml"
    "io"
    "bytes"
    "os"
    "io/ioutil"
    "strconv"
    "time"
)
const (
    C2S_DeleteAccount_CMD    = 0x0C
    S2C_DeleteAccountAck_CMD = 0x09
    ERROR_ErrorInfo_CMD      = 0x20
)

type Config struct {
    Addr   string
    Userid string
}

func main() {
    addr := "203.117.155.97:9100"
    userIdStr := ""
    if len(os.Args) > 1 {
        userIdStr = os.Args[1]
    }
    data, err := ioutil.ReadFile("DeleteUser.xml")
    if err == nil {
        var config Config
        err = xml.Unmarshal(data, &config)
        if err == nil {
            addr = config.Addr
            if userIdStr == "" {
                userIdStr = config.Userid
            }
        }
    }
    if err != nil {
        log.Print(err, " use default configuration")
    }
    userId, err := strconv.Atoi(userIdStr)
    if err != nil {
        log.Fatal("wrong userid ", err)
    }
    log.Print("send DeleteUser ", userId, " command to ", addr)
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        log.Fatal("resolve error ", err)
    }
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        log.Fatal("connect error ", err)
    }
    conn.SetReadDeadline(time.Now().Add(time.Second * 5))
    deleteUser := &Auth_C2S.DeleteAccount{Userid : proto.Int(userId)}
    sendMsg(conn, byte(C2S_DeleteAccount_CMD), deleteUser)
    readReply(conn)
}

func sendMsg(conn net.Conn, cmd byte, msg proto.Message) {
    data, err := proto.Marshal(msg)
    if err != nil {
        log.Fatal("marchaling error: ", err)
    }
    length := int32(len(data)) + 1
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, length)
    binary.Write(buf, binary.LittleEndian, cmd)
    buf.Write(data)
    conn.Write(buf.Bytes())
}

func readReply(conn net.Conn) {
    length := int32(0)
    err := binary.Read(conn, binary.LittleEndian, &length)
    if err != nil {
        log.Print("timeout or the connection is broken")
        return
    }
    rbuf := make([]byte, length)
    io.ReadFull(conn, rbuf)
    if len(rbuf) < 1 {
        log.Printf("short message");
        return
    }
    cmd := int(rbuf[0])
    switch cmd {
    case S2C_DeleteAccountAck_CMD:
        res := &Auth_S2C.DeleteAccountAck{}
        err := proto.Unmarshal(rbuf[1:], res)
        if err != nil {
            log.Fatal(err)
        }
        log.Print("User ", res.GetUserid(), " is deleted")
    case ERROR_ErrorInfo_CMD:
        res := &Auth_S2C.ErrorInfo{}
        err := proto.Unmarshal(rbuf[1:], res)
        if err != nil {
            log.Fatal(err)
        }
        errorType := int(res.GetType())
        switch errorType {
        case 64:
            log.Print("Cannot delete normal account without log in")
        case 32:
            log.Print("DB error")
        case 9:
            log.Print("Cannot delete other's account")
        case 1:
            log.Print("User does not exist")
        default:
            log.Print("unknown error: ", res.GetType())
        }
    }
}
