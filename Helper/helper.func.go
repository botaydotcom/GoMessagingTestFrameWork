package Helper

import (
	"time"
	"crypto/md5"
	"fmt"
)

var counter int64
var counterTs int64

func RequestId() (string) {
	counter += 1
	result := time.Now().Unix() << 16 + counter
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