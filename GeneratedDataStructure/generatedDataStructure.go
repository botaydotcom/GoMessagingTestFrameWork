
package GeneratedDataStructure

import (
	"reflect"
	
	"xmltest/btalkTest/Auth_C2S"
	
	"xmltest/btalkTest/Auth_S2C"
	
	"xmltest/btalkTest/TOKEN_S"
	
	
)
// exported Map
var Map map[string] reflect.Type

func init() {
	Map = make(map[string] reflect.Type)
	
	var aAuth_C2S_HelloInfo Auth_C2S.HelloInfo
	Map["HelloInfo"] = reflect.TypeOf(aAuth_C2S_HelloInfo)
	
	var aAuth_C2S_LoginInfo Auth_C2S.LoginInfo
	Map["LoginInfo"] = reflect.TypeOf(aAuth_C2S_LoginInfo)
	
	var aAuth_C2S_OAuthLogin Auth_C2S.OAuthLogin
	Map["OAuthLogin"] = reflect.TypeOf(aAuth_C2S_OAuthLogin)
	
	var aAuth_C2S_OAuthRawInfo Auth_C2S.OAuthRawInfo
	Map["OAuthRawInfo"] = reflect.TypeOf(aAuth_C2S_OAuthRawInfo)
	
	var aAuth_S2C_CheckVersionErrorInfo Auth_S2C.CheckVersionErrorInfo
	Map["CheckVersionErrorInfo"] = reflect.TypeOf(aAuth_S2C_CheckVersionErrorInfo)
	
	var aAuth_S2C_ErrorInfo Auth_S2C.ErrorInfo
	Map["ErrorInfo"] = reflect.TypeOf(aAuth_S2C_ErrorInfo)
	
	var aAuth_S2C_HelloInfoResult Auth_S2C.HelloInfoResult
	Map["HelloInfoResult"] = reflect.TypeOf(aAuth_S2C_HelloInfoResult)
	
	var aAuth_S2C_LoginTokenInfo Auth_S2C.LoginTokenInfo
	Map["LoginTokenInfo"] = reflect.TypeOf(aAuth_S2C_LoginTokenInfo)
	
	var aAuth_S2C_LoginUserInfo Auth_S2C.LoginUserInfo
	Map["LoginUserInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo)
	
	var aAuth_S2C_LoginUserInfo_LastLoginInfo Auth_S2C.LoginUserInfo_LastLoginInfo
	Map["LoginUserInfo_LastLoginInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo_LastLoginInfo)
	
	var aTOKEN_S_UserInfo TOKEN_S.UserInfo
	Map["UserInfo"] = reflect.TypeOf(aTOKEN_S_UserInfo)
	
}
