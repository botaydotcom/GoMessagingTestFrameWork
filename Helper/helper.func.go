package Helper

import (
	"time"
	"crypto/md5"
	"fmt"
	"code.google.com/p/goprotobuf/proto"
	"xmltest/btalkTest/Auth_C2S"	
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

func GetOAuthRaw(provider string, account string, param string) (string) {
	oauthInfo := &Auth_C2S.OAuthRawInfo{
		Provider: proto.String(provider),
		Account:  proto.String(account),
		Content:  []byte(param),
	}
	data, _ := proto.Marshal(oauthInfo)
	return (string)(data)
}