package Helper

import (
	"time"
	"crypto/md5"
	"fmt"
	"code.google.com/p/goprotobuf/proto"
	"xmltest/btalkTest/Auth_C2S"
	"sync"
	"io/ioutil"
	"os"
	"log"
	"strconv"
)

const (
	FILE_CHUNK_SIZE = 10000
)

var counter int64
var counterTs int64

func RequestId() (string) {	
	counter += 1
	result := time.Now().Unix() << 16 + counter
	fmt.Println("GETTING NEW REQUEST ID - COUNTER: ", counter, " RESULT: ", result)
	return fmt.Sprintf("%v", result)	
}

func Timestamp() (string) {
	result := time.Now().Unix();
	return fmt.Sprintf("%v", result)
}

func IncTimestamp() (string) {
	counterTs += 1
	result := time.Now().Unix() + counterTs
	return fmt.Sprintf("%v", result)
}

func Md5For(plainText string) (string) {
    h := md5.New()
    h.Write([]byte(plainText))
    return fmt.Sprintf("%x", h.Sum(nil))
}

func GetOAuthRaw(provider string, account string, param string) (string) {
	oauthInfo := &Auth_C2S.OAuthRawInfo{
		Provider: proto.String(provider),
		Account:  proto.String(account),
		Content:  []byte(param),
	}
	data, _ := proto.Marshal(oauthInfo)
	return (string)(data)
}

var userMin int = 0
var userMax int = 20000
var counterUser int = 0
var mutex sync.Mutex

func getNextUserNumber() int{
	(&mutex).Lock()
	counterUser++
	result := (userMin + counterUser) % userMax
	(&mutex).Unlock()
	return result
}

func getCurrentUserNumber() int{
	(&mutex).Lock()
	result := (userMin + counterUser) % userMax
	(&mutex).Unlock()
	return result
}

func GetNextUserEmail() string {
	user := fmt.Sprintf("indotrial_user_%d@test.com", getNextUserNumber())
	return user
}

func GetCurrentUserName() string {
	user := fmt.Sprintf("indotrial_user_%d", getCurrentUserNumber())
	return user
}

func FileContent(fileName string) ([]byte) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}

func FileContentMultiPart(fileName string, partIdStr string) ([]byte) {
	partId, _ := strconv.ParseInt(partIdStr, 0, 0)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	fileSize := getFileSize(fileName)
	//fmt.Println("Trying to get offset for file", fileName, "part:", partId)
	offsetStart, offsetEnd := getOffsetForPart(fileSize, partId)
	//fmt.Println("result:", offsetStart, offsetEnd)
	return data[offsetStart:offsetEnd]
}

func getOffsetForPart(fileSize int64, partId int64)(int64, int64) {
	if int64(partId) * FILE_CHUNK_SIZE > fileSize {
		return -1, -1
	}
	start := int64(partId) * FILE_CHUNK_SIZE
	end := int64(partId + 1) * FILE_CHUNK_SIZE

	if end > fileSize {
		end = fileSize
	}
	return start, end
}

func getFileSize(fileName string) (int64) {
	info, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	return info.Size()
}

func FileSize(fileName string) (string) {
	fileSize := getFileSize(fileName)
	result := fmt.Sprintf("%d", fileSize)
	return result
}

func numFilePart(fileName string) (int64) {
	fileSize := getFileSize(fileName)
	return (fileSize - 1) / 10000 + 1
}

func NumFilePart(fileName string) (string) {
	numPart := numFilePart(fileName)
	result := fmt.Sprintf("%d", numPart)
	return result
}