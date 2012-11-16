
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
	
	var aAuth_C2S_ChangeMyInfo Auth_C2S.ChangeMyInfo
	Map["ChangeMyInfo"] = reflect.TypeOf(aAuth_C2S_ChangeMyInfo)
	
	var aAuth_C2S_ChangeNotification Auth_C2S.ChangeNotification
	Map["ChangeNotification"] = reflect.TypeOf(aAuth_C2S_ChangeNotification)
	
	var aAuth_C2S_ChangePassword Auth_C2S.ChangePassword
	Map["ChangePassword"] = reflect.TypeOf(aAuth_C2S_ChangePassword)
	
	var aAuth_C2S_CreateChannelRequest Auth_C2S.CreateChannelRequest
	Map["CreateChannelRequest"] = reflect.TypeOf(aAuth_C2S_CreateChannelRequest)
	
	var aAuth_C2S_FillFinishReg Auth_C2S.FillFinishReg
	Map["FillFinishReg"] = reflect.TypeOf(aAuth_C2S_FillFinishReg)
	
	var aAuth_C2S_HelloInfo Auth_C2S.HelloInfo
	Map["HelloInfo"] = reflect.TypeOf(aAuth_C2S_HelloInfo)
	
	var aAuth_C2S_KeepAlive Auth_C2S.KeepAlive
	Map["KeepAlive"] = reflect.TypeOf(aAuth_C2S_KeepAlive)
	
	var aAuth_C2S_LoginInfo Auth_C2S.LoginInfo
	Map["LoginInfo"] = reflect.TypeOf(aAuth_C2S_LoginInfo)
	
	var aAuth_C2S_MonitorPoll Auth_C2S.MonitorPoll
	Map["MonitorPoll"] = reflect.TypeOf(aAuth_C2S_MonitorPoll)
	
	var aAuth_C2S_OAuthLogin Auth_C2S.OAuthLogin
	Map["OAuthLogin"] = reflect.TypeOf(aAuth_C2S_OAuthLogin)
	
	var aAuth_C2S_OAuthRawInfo Auth_C2S.OAuthRawInfo
	Map["OAuthRawInfo"] = reflect.TypeOf(aAuth_C2S_OAuthRawInfo)
	
	var aAuth_C2S_RegAccount Auth_C2S.RegAccount
	Map["RegAccount"] = reflect.TypeOf(aAuth_C2S_RegAccount)
	
	var aAuth_C2S_RequestChannelInfo Auth_C2S.RequestChannelInfo
	Map["RequestChannelInfo"] = reflect.TypeOf(aAuth_C2S_RequestChannelInfo)
	
	var aAuth_C2S_RequestChannelList Auth_C2S.RequestChannelList
	Map["RequestChannelList"] = reflect.TypeOf(aAuth_C2S_RequestChannelList)
	
	var aAuth_C2S_RequestChannelList_RequestChannelType Auth_C2S.RequestChannelList_RequestChannelType
	Map["RequestChannelList_RequestChannelType"] = reflect.TypeOf(aAuth_C2S_RequestChannelList_RequestChannelType)
	
	var aAuth_C2S_RequestToken Auth_C2S.RequestToken
	Map["RequestToken"] = reflect.TypeOf(aAuth_C2S_RequestToken)
	
	var aAuth_C2S_RequestUserNotificationInfo Auth_C2S.RequestUserNotificationInfo
	Map["RequestUserNotificationInfo"] = reflect.TypeOf(aAuth_C2S_RequestUserNotificationInfo)
	
	var aAuth_C2S_RequestUserNumList Auth_C2S.RequestUserNumList
	Map["RequestUserNumList"] = reflect.TypeOf(aAuth_C2S_RequestUserNumList)
	
	var aAuth_C2S_VMModuleInfoList Auth_C2S.VMModuleInfoList
	Map["VMModuleInfoList"] = reflect.TypeOf(aAuth_C2S_VMModuleInfoList)
	
	var aAuth_C2S_VMModuleInfoList_VMModuleInfo Auth_C2S.VMModuleInfoList_VMModuleInfo
	Map["VMModuleInfoList_VMModuleInfo"] = reflect.TypeOf(aAuth_C2S_VMModuleInfoList_VMModuleInfo)
	
	var aAuth_S2C_ChangeMyInfoResult Auth_S2C.ChangeMyInfoResult
	Map["ChangeMyInfoResult"] = reflect.TypeOf(aAuth_S2C_ChangeMyInfoResult)
	
	var aAuth_S2C_ChangeNotificationResult Auth_S2C.ChangeNotificationResult
	Map["ChangeNotificationResult"] = reflect.TypeOf(aAuth_S2C_ChangeNotificationResult)
	
	var aAuth_S2C_ChangePasswordResult Auth_S2C.ChangePasswordResult
	Map["ChangePasswordResult"] = reflect.TypeOf(aAuth_S2C_ChangePasswordResult)
	
	var aAuth_S2C_ChannelIdInfo Auth_S2C.ChannelIdInfo
	Map["ChannelIdInfo"] = reflect.TypeOf(aAuth_S2C_ChannelIdInfo)
	
	var aAuth_S2C_ChannelIdInfoList Auth_S2C.ChannelIdInfoList
	Map["ChannelIdInfoList"] = reflect.TypeOf(aAuth_S2C_ChannelIdInfoList)
	
	var aAuth_S2C_ChannelSimpleInfoList Auth_S2C.ChannelSimpleInfoList
	Map["ChannelSimpleInfoList"] = reflect.TypeOf(aAuth_S2C_ChannelSimpleInfoList)
	
	var aAuth_S2C_ChannelUserInfo Auth_S2C.ChannelUserInfo
	Map["ChannelUserInfo"] = reflect.TypeOf(aAuth_S2C_ChannelUserInfo)
	
	var aAuth_S2C_ChannelUserNum Auth_S2C.ChannelUserNum
	Map["ChannelUserNum"] = reflect.TypeOf(aAuth_S2C_ChannelUserNum)
	
	var aAuth_S2C_ChannelUserNumList Auth_S2C.ChannelUserNumList
	Map["ChannelUserNumList"] = reflect.TypeOf(aAuth_S2C_ChannelUserNumList)
	
	var aAuth_S2C_CheckVersionErrorInfo Auth_S2C.CheckVersionErrorInfo
	Map["CheckVersionErrorInfo"] = reflect.TypeOf(aAuth_S2C_CheckVersionErrorInfo)
	
	var aAuth_S2C_CreateChannelFailed Auth_S2C.CreateChannelFailed
	Map["CreateChannelFailed"] = reflect.TypeOf(aAuth_S2C_CreateChannelFailed)
	
	var aAuth_S2C_CreateChannelFailed_ErrorCodeType Auth_S2C.CreateChannelFailed_ErrorCodeType
	Map["CreateChannelFailed_ErrorCodeType"] = reflect.TypeOf(aAuth_S2C_CreateChannelFailed_ErrorCodeType)
	
	var aAuth_S2C_CreateChannelResult Auth_S2C.CreateChannelResult
	Map["CreateChannelResult"] = reflect.TypeOf(aAuth_S2C_CreateChannelResult)
	
	var aAuth_S2C_ErrorInfo Auth_S2C.ErrorInfo
	Map["ErrorInfo"] = reflect.TypeOf(aAuth_S2C_ErrorInfo)
	
	var aAuth_S2C_HelloInfoResult Auth_S2C.HelloInfoResult
	Map["HelloInfoResult"] = reflect.TypeOf(aAuth_S2C_HelloInfoResult)
	
	var aAuth_S2C_KeepAliveAck Auth_S2C.KeepAliveAck
	Map["KeepAliveAck"] = reflect.TypeOf(aAuth_S2C_KeepAliveAck)
	
	var aAuth_S2C_LoginTokenInfo Auth_S2C.LoginTokenInfo
	Map["LoginTokenInfo"] = reflect.TypeOf(aAuth_S2C_LoginTokenInfo)
	
	var aAuth_S2C_LoginUserInfo Auth_S2C.LoginUserInfo
	Map["LoginUserInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo)
	
	var aAuth_S2C_LoginUserInfo_LastLoginInfo Auth_S2C.LoginUserInfo_LastLoginInfo
	Map["LoginUserInfo_LastLoginInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo_LastLoginInfo)
	
	var aAuth_S2C_MonitorResult Auth_S2C.MonitorResult
	Map["MonitorResult"] = reflect.TypeOf(aAuth_S2C_MonitorResult)
	
	var aAuth_S2C_NeedFinishReg Auth_S2C.NeedFinishReg
	Map["NeedFinishReg"] = reflect.TypeOf(aAuth_S2C_NeedFinishReg)
	
	var aAuth_S2C_RequestAckInfo Auth_S2C.RequestAckInfo
	Map["RequestAckInfo"] = reflect.TypeOf(aAuth_S2C_RequestAckInfo)
	
	var aAuth_S2C_RequestTokenResult Auth_S2C.RequestTokenResult
	Map["RequestTokenResult"] = reflect.TypeOf(aAuth_S2C_RequestTokenResult)
	
	var aAuth_S2C_UserNotificationInfo Auth_S2C.UserNotificationInfo
	Map["UserNotificationInfo"] = reflect.TypeOf(aAuth_S2C_UserNotificationInfo)
	
	var aTOKEN_S_UserInfo TOKEN_S.UserInfo
	Map["UserInfo"] = reflect.TypeOf(aTOKEN_S_UserInfo)
	
}
