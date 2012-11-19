
package GeneratedDataStructure

import (
	"reflect"
	
	"xmltest/btalkTest/Auth_Buddy_C2S"
	
	"xmltest/btalkTest/Auth_Buddy_S2C"
	
	"xmltest/btalkTest/Auth_C2S"
	
	"xmltest/btalkTest/Auth_S2C"
	
	"xmltest/btalkTest/TOKEN_S"
	
	
)
// exported Map
var Map map[string] reflect.Type

func init() {
	Map = make(map[string] reflect.Type)
	
	var aAuth_Buddy_C2S_AddBuddyResult Auth_Buddy_C2S.AddBuddyResult
	Map["Auth_Buddy_C2S_AddBuddyResult"] = reflect.TypeOf(aAuth_Buddy_C2S_AddBuddyResult)
	
	var aAuth_Buddy_C2S_ChangeBuddyCategory Auth_Buddy_C2S.ChangeBuddyCategory
	Map["Auth_Buddy_C2S_ChangeBuddyCategory"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeBuddyCategory)
	
	var aAuth_Buddy_C2S_ChangeMeMemData Auth_Buddy_C2S.ChangeMeMemData
	Map["Auth_Buddy_C2S_ChangeMeMemData"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeMeMemData)
	
	var aAuth_Buddy_C2S_ChangeMyCurrChannel Auth_Buddy_C2S.ChangeMyCurrChannel
	Map["Auth_Buddy_C2S_ChangeMyCurrChannel"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeMyCurrChannel)
	
	var aAuth_Buddy_C2S_ChangeUserInfo Auth_Buddy_C2S.ChangeUserInfo
	Map["Auth_Buddy_C2S_ChangeUserInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeUserInfo)
	
	var aAuth_Buddy_C2S_ChangeUserInfo_ChangeInfo Auth_Buddy_C2S.ChangeUserInfo_ChangeInfo
	Map["Auth_Buddy_C2S_ChangeUserInfo_ChangeInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeUserInfo_ChangeInfo)
	
	var aAuth_Buddy_C2S_ChangeUserInfo_ChangeOnlineStatus Auth_Buddy_C2S.ChangeUserInfo_ChangeOnlineStatus
	Map["Auth_Buddy_C2S_ChangeUserInfo_ChangeOnlineStatus"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeUserInfo_ChangeOnlineStatus)
	
	var aAuth_Buddy_C2S_ChangeUserInfo_ChangeSignature Auth_Buddy_C2S.ChangeUserInfo_ChangeSignature
	Map["Auth_Buddy_C2S_ChangeUserInfo_ChangeSignature"] = reflect.TypeOf(aAuth_Buddy_C2S_ChangeUserInfo_ChangeSignature)
	
	var aAuth_Buddy_C2S_Chat2AckRemote Auth_Buddy_C2S.Chat2AckRemote
	Map["Auth_Buddy_C2S_Chat2AckRemote"] = reflect.TypeOf(aAuth_Buddy_C2S_Chat2AckRemote)
	
	var aAuth_Buddy_C2S_ChatAckRemote Auth_Buddy_C2S.ChatAckRemote
	Map["Auth_Buddy_C2S_ChatAckRemote"] = reflect.TypeOf(aAuth_Buddy_C2S_ChatAckRemote)
	
	var aAuth_Buddy_C2S_ChatInfo Auth_Buddy_C2S.ChatInfo
	Map["Auth_Buddy_C2S_ChatInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_ChatInfo)
	
	var aAuth_Buddy_C2S_ChatInfo2 Auth_Buddy_C2S.ChatInfo2
	Map["Auth_Buddy_C2S_ChatInfo2"] = reflect.TypeOf(aAuth_Buddy_C2S_ChatInfo2)
	
	var aAuth_Buddy_C2S_ChatMediaInfo Auth_Buddy_C2S.ChatMediaInfo
	Map["Auth_Buddy_C2S_ChatMediaInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_ChatMediaInfo)
	
	var aAuth_Buddy_C2S_DeleteBuddy Auth_Buddy_C2S.DeleteBuddy
	Map["Auth_Buddy_C2S_DeleteBuddy"] = reflect.TypeOf(aAuth_Buddy_C2S_DeleteBuddy)
	
	var aAuth_Buddy_C2S_DeleteCategory Auth_Buddy_C2S.DeleteCategory
	Map["Auth_Buddy_C2S_DeleteCategory"] = reflect.TypeOf(aAuth_Buddy_C2S_DeleteCategory)
	
	var aAuth_Buddy_C2S_DoActionCode Auth_Buddy_C2S.DoActionCode
	Map["Auth_Buddy_C2S_DoActionCode"] = reflect.TypeOf(aAuth_Buddy_C2S_DoActionCode)
	
	var aAuth_Buddy_C2S_EditCategory Auth_Buddy_C2S.EditCategory
	Map["Auth_Buddy_C2S_EditCategory"] = reflect.TypeOf(aAuth_Buddy_C2S_EditCategory)
	
	var aAuth_Buddy_C2S_MobileClearAction Auth_Buddy_C2S.MobileClearAction
	Map["Auth_Buddy_C2S_MobileClearAction"] = reflect.TypeOf(aAuth_Buddy_C2S_MobileClearAction)
	
	var aAuth_Buddy_C2S_MobileFeatureAction Auth_Buddy_C2S.MobileFeatureAction
	Map["Auth_Buddy_C2S_MobileFeatureAction"] = reflect.TypeOf(aAuth_Buddy_C2S_MobileFeatureAction)
	
	var aAuth_Buddy_C2S_P2PRequest Auth_Buddy_C2S.P2PRequest
	Map["Auth_Buddy_C2S_P2PRequest"] = reflect.TypeOf(aAuth_Buddy_C2S_P2PRequest)
	
	var aAuth_Buddy_C2S_P2PResponse Auth_Buddy_C2S.P2PResponse
	Map["Auth_Buddy_C2S_P2PResponse"] = reflect.TypeOf(aAuth_Buddy_C2S_P2PResponse)
	
	var aAuth_Buddy_C2S_RequestAddBuddy Auth_Buddy_C2S.RequestAddBuddy
	Map["Auth_Buddy_C2S_RequestAddBuddy"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestAddBuddy)
	
	var aAuth_Buddy_C2S_RequestAllBuddyMemData Auth_Buddy_C2S.RequestAllBuddyMemData
	Map["Auth_Buddy_C2S_RequestAllBuddyMemData"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestAllBuddyMemData)
	
	var aAuth_Buddy_C2S_RequestBuddyCurrChannelList Auth_Buddy_C2S.RequestBuddyCurrChannelList
	Map["Auth_Buddy_C2S_RequestBuddyCurrChannelList"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddyCurrChannelList)
	
	var aAuth_Buddy_C2S_RequestBuddyInfo Auth_Buddy_C2S.RequestBuddyInfo
	Map["Auth_Buddy_C2S_RequestBuddyInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddyInfo)
	
	var aAuth_Buddy_C2S_RequestBuddyInfo_Request Auth_Buddy_C2S.RequestBuddyInfo_Request
	Map["Auth_Buddy_C2S_RequestBuddyInfo_Request"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddyInfo_Request)
	
	var aAuth_Buddy_C2S_RequestBuddyOnlineList Auth_Buddy_C2S.RequestBuddyOnlineList
	Map["Auth_Buddy_C2S_RequestBuddyOnlineList"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddyOnlineList)
	
	var aAuth_Buddy_C2S_RequestBuddySignature Auth_Buddy_C2S.RequestBuddySignature
	Map["Auth_Buddy_C2S_RequestBuddySignature"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddySignature)
	
	var aAuth_Buddy_C2S_RequestBuddySimpleInfoList Auth_Buddy_C2S.RequestBuddySimpleInfoList
	Map["Auth_Buddy_C2S_RequestBuddySimpleInfoList"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestBuddySimpleInfoList)
	
	var aAuth_Buddy_C2S_RequestCategoryList Auth_Buddy_C2S.RequestCategoryList
	Map["Auth_Buddy_C2S_RequestCategoryList"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestCategoryList)
	
	var aAuth_Buddy_C2S_RequestOneBuddyMemData Auth_Buddy_C2S.RequestOneBuddyMemData
	Map["Auth_Buddy_C2S_RequestOneBuddyMemData"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestOneBuddyMemData)
	
	var aAuth_Buddy_C2S_RequestSignature Auth_Buddy_C2S.RequestSignature
	Map["Auth_Buddy_C2S_RequestSignature"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestSignature)
	
	var aAuth_Buddy_C2S_RequestSyncBuddy Auth_Buddy_C2S.RequestSyncBuddy
	Map["Auth_Buddy_C2S_RequestSyncBuddy"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestSyncBuddy)
	
	var aAuth_Buddy_C2S_RequestUserInfo Auth_Buddy_C2S.RequestUserInfo
	Map["Auth_Buddy_C2S_RequestUserInfo"] = reflect.TypeOf(aAuth_Buddy_C2S_RequestUserInfo)
	
	var aAuth_Buddy_S2C_BuddyCurrChannelInfo Auth_Buddy_S2C.BuddyCurrChannelInfo
	Map["Auth_Buddy_S2C_BuddyCurrChannelInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyCurrChannelInfo)
	
	var aAuth_Buddy_S2C_BuddyCurrChannelList Auth_Buddy_S2C.BuddyCurrChannelList
	Map["Auth_Buddy_S2C_BuddyCurrChannelList"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyCurrChannelList)
	
	var aAuth_Buddy_S2C_BuddyDeleted Auth_Buddy_S2C.BuddyDeleted
	Map["Auth_Buddy_S2C_BuddyDeleted"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyDeleted)
	
	var aAuth_Buddy_S2C_BuddyInfo Auth_Buddy_S2C.BuddyInfo
	Map["Auth_Buddy_S2C_BuddyInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyInfo)
	
	var aAuth_Buddy_S2C_BuddyInfoList Auth_Buddy_S2C.BuddyInfoList
	Map["Auth_Buddy_S2C_BuddyInfoList"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyInfoList)
	
	var aAuth_Buddy_S2C_BuddyMemDataInfo Auth_Buddy_S2C.BuddyMemDataInfo
	Map["Auth_Buddy_S2C_BuddyMemDataInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyMemDataInfo)
	
	var aAuth_Buddy_S2C_BuddyMemDataItem Auth_Buddy_S2C.BuddyMemDataItem
	Map["Auth_Buddy_S2C_BuddyMemDataItem"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyMemDataItem)
	
	var aAuth_Buddy_S2C_BuddyOnline Auth_Buddy_S2C.BuddyOnline
	Map["Auth_Buddy_S2C_BuddyOnline"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyOnline)
	
	var aAuth_Buddy_S2C_BuddyOnlineList Auth_Buddy_S2C.BuddyOnlineList
	Map["Auth_Buddy_S2C_BuddyOnlineList"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyOnlineList)
	
	var aAuth_Buddy_S2C_BuddyRelationInfo Auth_Buddy_S2C.BuddyRelationInfo
	Map["Auth_Buddy_S2C_BuddyRelationInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddyRelationInfo)
	
	var aAuth_Buddy_S2C_BuddySimpleInfoList Auth_Buddy_S2C.BuddySimpleInfoList
	Map["Auth_Buddy_S2C_BuddySimpleInfoList"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddySimpleInfoList)
	
	var aAuth_Buddy_S2C_BuddySimpleInfoList_BuddySimpleInfo Auth_Buddy_S2C.BuddySimpleInfoList_BuddySimpleInfo
	Map["Auth_Buddy_S2C_BuddySimpleInfoList_BuddySimpleInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_BuddySimpleInfoList_BuddySimpleInfo)
	
	var aAuth_Buddy_S2C_CategoryItem Auth_Buddy_S2C.CategoryItem
	Map["Auth_Buddy_S2C_CategoryItem"] = reflect.TypeOf(aAuth_Buddy_S2C_CategoryItem)
	
	var aAuth_Buddy_S2C_CategoryList Auth_Buddy_S2C.CategoryList
	Map["Auth_Buddy_S2C_CategoryList"] = reflect.TypeOf(aAuth_Buddy_S2C_CategoryList)
	
	var aAuth_Buddy_S2C_ChangeBuddyCategoryResult Auth_Buddy_S2C.ChangeBuddyCategoryResult
	Map["Auth_Buddy_S2C_ChangeBuddyCategoryResult"] = reflect.TypeOf(aAuth_Buddy_S2C_ChangeBuddyCategoryResult)
	
	var aAuth_Buddy_S2C_ChangeBuddyCurrChannel Auth_Buddy_S2C.ChangeBuddyCurrChannel
	Map["Auth_Buddy_S2C_ChangeBuddyCurrChannel"] = reflect.TypeOf(aAuth_Buddy_S2C_ChangeBuddyCurrChannel)
	
	var aAuth_Buddy_S2C_ChangeMeMemDataResult Auth_Buddy_S2C.ChangeMeMemDataResult
	Map["Auth_Buddy_S2C_ChangeMeMemDataResult"] = reflect.TypeOf(aAuth_Buddy_S2C_ChangeMeMemDataResult)
	
	var aAuth_Buddy_S2C_Chat2Ack Auth_Buddy_S2C.Chat2Ack
	Map["Auth_Buddy_S2C_Chat2Ack"] = reflect.TypeOf(aAuth_Buddy_S2C_Chat2Ack)
	
	var aAuth_Buddy_S2C_ChatAck Auth_Buddy_S2C.ChatAck
	Map["Auth_Buddy_S2C_ChatAck"] = reflect.TypeOf(aAuth_Buddy_S2C_ChatAck)
	
	var aAuth_Buddy_S2C_ChatInfo Auth_Buddy_S2C.ChatInfo
	Map["Auth_Buddy_S2C_ChatInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_ChatInfo)
	
	var aAuth_Buddy_S2C_ChatInfo2 Auth_Buddy_S2C.ChatInfo2
	Map["Auth_Buddy_S2C_ChatInfo2"] = reflect.TypeOf(aAuth_Buddy_S2C_ChatInfo2)
	
	var aAuth_Buddy_S2C_ChatMediaInfo Auth_Buddy_S2C.ChatMediaInfo
	Map["Auth_Buddy_S2C_ChatMediaInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_ChatMediaInfo)
	
	var aAuth_Buddy_S2C_ChatMediaInfoAck Auth_Buddy_S2C.ChatMediaInfoAck
	Map["Auth_Buddy_S2C_ChatMediaInfoAck"] = reflect.TypeOf(aAuth_Buddy_S2C_ChatMediaInfoAck)
	
	var aAuth_Buddy_S2C_DeleteBuddyError Auth_Buddy_S2C.DeleteBuddyError
	Map["Auth_Buddy_S2C_DeleteBuddyError"] = reflect.TypeOf(aAuth_Buddy_S2C_DeleteBuddyError)
	
	var aAuth_Buddy_S2C_DeleteBuddyError_DeleteBuddyErrorType Auth_Buddy_S2C.DeleteBuddyError_DeleteBuddyErrorType
	Map["Auth_Buddy_S2C_DeleteBuddyError_DeleteBuddyErrorType"] = reflect.TypeOf(aAuth_Buddy_S2C_DeleteBuddyError_DeleteBuddyErrorType)
	
	var aAuth_Buddy_S2C_DeleteCategoryError Auth_Buddy_S2C.DeleteCategoryError
	Map["Auth_Buddy_S2C_DeleteCategoryError"] = reflect.TypeOf(aAuth_Buddy_S2C_DeleteCategoryError)
	
	var aAuth_Buddy_S2C_DeleteCategoryError_DeleteCategoryErrorType Auth_Buddy_S2C.DeleteCategoryError_DeleteCategoryErrorType
	Map["Auth_Buddy_S2C_DeleteCategoryError_DeleteCategoryErrorType"] = reflect.TypeOf(aAuth_Buddy_S2C_DeleteCategoryError_DeleteCategoryErrorType)
	
	var aAuth_Buddy_S2C_DeleteCategoryResponse Auth_Buddy_S2C.DeleteCategoryResponse
	Map["Auth_Buddy_S2C_DeleteCategoryResponse"] = reflect.TypeOf(aAuth_Buddy_S2C_DeleteCategoryResponse)
	
	var aAuth_Buddy_S2C_EditCategoryError Auth_Buddy_S2C.EditCategoryError
	Map["Auth_Buddy_S2C_EditCategoryError"] = reflect.TypeOf(aAuth_Buddy_S2C_EditCategoryError)
	
	var aAuth_Buddy_S2C_EditCategoryError_EditCategoryErrorType Auth_Buddy_S2C.EditCategoryError_EditCategoryErrorType
	Map["Auth_Buddy_S2C_EditCategoryError_EditCategoryErrorType"] = reflect.TypeOf(aAuth_Buddy_S2C_EditCategoryError_EditCategoryErrorType)
	
	var aAuth_Buddy_S2C_EditCategoryResponse Auth_Buddy_S2C.EditCategoryResponse
	Map["Auth_Buddy_S2C_EditCategoryResponse"] = reflect.TypeOf(aAuth_Buddy_S2C_EditCategoryResponse)
	
	var aAuth_Buddy_S2C_MobileUserLocationResponse Auth_Buddy_S2C.MobileUserLocationResponse
	Map["Auth_Buddy_S2C_MobileUserLocationResponse"] = reflect.TypeOf(aAuth_Buddy_S2C_MobileUserLocationResponse)
	
	var aAuth_Buddy_S2C_MobileUserLocationResponse_MobileUserLocationInfo Auth_Buddy_S2C.MobileUserLocationResponse_MobileUserLocationInfo
	Map["Auth_Buddy_S2C_MobileUserLocationResponse_MobileUserLocationInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_MobileUserLocationResponse_MobileUserLocationInfo)
	
	var aAuth_Buddy_S2C_MyInfoChanged Auth_Buddy_S2C.MyInfoChanged
	Map["Auth_Buddy_S2C_MyInfoChanged"] = reflect.TypeOf(aAuth_Buddy_S2C_MyInfoChanged)
	
	var aAuth_Buddy_S2C_MyInfoChanged_MyInfo Auth_Buddy_S2C.MyInfoChanged_MyInfo
	Map["Auth_Buddy_S2C_MyInfoChanged_MyInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_MyInfoChanged_MyInfo)
	
	var aAuth_Buddy_S2C_RemoteAddBuddyResult Auth_Buddy_S2C.RemoteAddBuddyResult
	Map["Auth_Buddy_S2C_RemoteAddBuddyResult"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteAddBuddyResult)
	
	var aAuth_Buddy_S2C_RemoteChat2Ack Auth_Buddy_S2C.RemoteChat2Ack
	Map["Auth_Buddy_S2C_RemoteChat2Ack"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteChat2Ack)
	
	var aAuth_Buddy_S2C_RemoteChatAck Auth_Buddy_S2C.RemoteChatAck
	Map["Auth_Buddy_S2C_RemoteChatAck"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteChatAck)
	
	var aAuth_Buddy_S2C_RemoteDoActionCode Auth_Buddy_S2C.RemoteDoActionCode
	Map["Auth_Buddy_S2C_RemoteDoActionCode"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteDoActionCode)
	
	var aAuth_Buddy_S2C_RemoteP2PRequest Auth_Buddy_S2C.RemoteP2PRequest
	Map["Auth_Buddy_S2C_RemoteP2PRequest"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteP2PRequest)
	
	var aAuth_Buddy_S2C_RemoteP2PResponse Auth_Buddy_S2C.RemoteP2PResponse
	Map["Auth_Buddy_S2C_RemoteP2PResponse"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteP2PResponse)
	
	var aAuth_Buddy_S2C_RemoteP2PResponseAck Auth_Buddy_S2C.RemoteP2PResponseAck
	Map["Auth_Buddy_S2C_RemoteP2PResponseAck"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteP2PResponseAck)
	
	var aAuth_Buddy_S2C_RemoteRequestAddBuddy Auth_Buddy_S2C.RemoteRequestAddBuddy
	Map["Auth_Buddy_S2C_RemoteRequestAddBuddy"] = reflect.TypeOf(aAuth_Buddy_S2C_RemoteRequestAddBuddy)
	
	var aAuth_Buddy_S2C_RequestAddBuddyError Auth_Buddy_S2C.RequestAddBuddyError
	Map["Auth_Buddy_S2C_RequestAddBuddyError"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestAddBuddyError)
	
	var aAuth_Buddy_S2C_RequestAddBuddyError_RequestAddBuddyErrorType Auth_Buddy_S2C.RequestAddBuddyError_RequestAddBuddyErrorType
	Map["Auth_Buddy_S2C_RequestAddBuddyError_RequestAddBuddyErrorType"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestAddBuddyError_RequestAddBuddyErrorType)
	
	var aAuth_Buddy_S2C_RequestAllBuddyMemDataResult Auth_Buddy_S2C.RequestAllBuddyMemDataResult
	Map["Auth_Buddy_S2C_RequestAllBuddyMemDataResult"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestAllBuddyMemDataResult)
	
	var aAuth_Buddy_S2C_RequestUserInfoError Auth_Buddy_S2C.RequestUserInfoError
	Map["Auth_Buddy_S2C_RequestUserInfoError"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestUserInfoError)
	
	var aAuth_Buddy_S2C_RequestUserInfoError_RequestUserInfoErrorType Auth_Buddy_S2C.RequestUserInfoError_RequestUserInfoErrorType
	Map["Auth_Buddy_S2C_RequestUserInfoError_RequestUserInfoErrorType"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestUserInfoError_RequestUserInfoErrorType)
	
	var aAuth_Buddy_S2C_RequestUserInfoResult Auth_Buddy_S2C.RequestUserInfoResult
	Map["Auth_Buddy_S2C_RequestUserInfoResult"] = reflect.TypeOf(aAuth_Buddy_S2C_RequestUserInfoResult)
	
	var aAuth_Buddy_S2C_UserInfo Auth_Buddy_S2C.UserInfo
	Map["Auth_Buddy_S2C_UserInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_UserInfo)
	
	var aAuth_Buddy_S2C_UserSignatureInfo Auth_Buddy_S2C.UserSignatureInfo
	Map["Auth_Buddy_S2C_UserSignatureInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_UserSignatureInfo)
	
	var aAuth_Buddy_S2C_UserSignatureList Auth_Buddy_S2C.UserSignatureList
	Map["Auth_Buddy_S2C_UserSignatureList"] = reflect.TypeOf(aAuth_Buddy_S2C_UserSignatureList)
	
	var aAuth_Buddy_S2C_UserStatusInfo Auth_Buddy_S2C.UserStatusInfo
	Map["Auth_Buddy_S2C_UserStatusInfo"] = reflect.TypeOf(aAuth_Buddy_S2C_UserStatusInfo)
	
	var aAuth_C2S_ChangeMyInfo Auth_C2S.ChangeMyInfo
	Map["Auth_C2S_ChangeMyInfo"] = reflect.TypeOf(aAuth_C2S_ChangeMyInfo)
	
	var aAuth_C2S_ChangeNotification Auth_C2S.ChangeNotification
	Map["Auth_C2S_ChangeNotification"] = reflect.TypeOf(aAuth_C2S_ChangeNotification)
	
	var aAuth_C2S_ChangePassword Auth_C2S.ChangePassword
	Map["Auth_C2S_ChangePassword"] = reflect.TypeOf(aAuth_C2S_ChangePassword)
	
	var aAuth_C2S_CreateChannelRequest Auth_C2S.CreateChannelRequest
	Map["Auth_C2S_CreateChannelRequest"] = reflect.TypeOf(aAuth_C2S_CreateChannelRequest)
	
	var aAuth_C2S_FillFinishReg Auth_C2S.FillFinishReg
	Map["Auth_C2S_FillFinishReg"] = reflect.TypeOf(aAuth_C2S_FillFinishReg)
	
	var aAuth_C2S_HelloInfo Auth_C2S.HelloInfo
	Map["Auth_C2S_HelloInfo"] = reflect.TypeOf(aAuth_C2S_HelloInfo)
	
	var aAuth_C2S_KeepAlive Auth_C2S.KeepAlive
	Map["Auth_C2S_KeepAlive"] = reflect.TypeOf(aAuth_C2S_KeepAlive)
	
	var aAuth_C2S_LoginInfo Auth_C2S.LoginInfo
	Map["Auth_C2S_LoginInfo"] = reflect.TypeOf(aAuth_C2S_LoginInfo)
	
	var aAuth_C2S_MonitorPoll Auth_C2S.MonitorPoll
	Map["Auth_C2S_MonitorPoll"] = reflect.TypeOf(aAuth_C2S_MonitorPoll)
	
	var aAuth_C2S_OAuthLogin Auth_C2S.OAuthLogin
	Map["Auth_C2S_OAuthLogin"] = reflect.TypeOf(aAuth_C2S_OAuthLogin)
	
	var aAuth_C2S_OAuthRawInfo Auth_C2S.OAuthRawInfo
	Map["Auth_C2S_OAuthRawInfo"] = reflect.TypeOf(aAuth_C2S_OAuthRawInfo)
	
	var aAuth_C2S_RegAccount Auth_C2S.RegAccount
	Map["Auth_C2S_RegAccount"] = reflect.TypeOf(aAuth_C2S_RegAccount)
	
	var aAuth_C2S_RequestChannelInfo Auth_C2S.RequestChannelInfo
	Map["Auth_C2S_RequestChannelInfo"] = reflect.TypeOf(aAuth_C2S_RequestChannelInfo)
	
	var aAuth_C2S_RequestChannelList Auth_C2S.RequestChannelList
	Map["Auth_C2S_RequestChannelList"] = reflect.TypeOf(aAuth_C2S_RequestChannelList)
	
	var aAuth_C2S_RequestChannelList_RequestChannelType Auth_C2S.RequestChannelList_RequestChannelType
	Map["Auth_C2S_RequestChannelList_RequestChannelType"] = reflect.TypeOf(aAuth_C2S_RequestChannelList_RequestChannelType)
	
	var aAuth_C2S_RequestToken Auth_C2S.RequestToken
	Map["Auth_C2S_RequestToken"] = reflect.TypeOf(aAuth_C2S_RequestToken)
	
	var aAuth_C2S_RequestUserNotificationInfo Auth_C2S.RequestUserNotificationInfo
	Map["Auth_C2S_RequestUserNotificationInfo"] = reflect.TypeOf(aAuth_C2S_RequestUserNotificationInfo)
	
	var aAuth_C2S_RequestUserNumList Auth_C2S.RequestUserNumList
	Map["Auth_C2S_RequestUserNumList"] = reflect.TypeOf(aAuth_C2S_RequestUserNumList)
	
	var aAuth_C2S_VMModuleInfoList Auth_C2S.VMModuleInfoList
	Map["Auth_C2S_VMModuleInfoList"] = reflect.TypeOf(aAuth_C2S_VMModuleInfoList)
	
	var aAuth_C2S_VMModuleInfoList_VMModuleInfo Auth_C2S.VMModuleInfoList_VMModuleInfo
	Map["Auth_C2S_VMModuleInfoList_VMModuleInfo"] = reflect.TypeOf(aAuth_C2S_VMModuleInfoList_VMModuleInfo)
	
	var aAuth_S2C_ChangeMyInfoResult Auth_S2C.ChangeMyInfoResult
	Map["Auth_S2C_ChangeMyInfoResult"] = reflect.TypeOf(aAuth_S2C_ChangeMyInfoResult)
	
	var aAuth_S2C_ChangeNotificationResult Auth_S2C.ChangeNotificationResult
	Map["Auth_S2C_ChangeNotificationResult"] = reflect.TypeOf(aAuth_S2C_ChangeNotificationResult)
	
	var aAuth_S2C_ChangePasswordResult Auth_S2C.ChangePasswordResult
	Map["Auth_S2C_ChangePasswordResult"] = reflect.TypeOf(aAuth_S2C_ChangePasswordResult)
	
	var aAuth_S2C_ChannelIdInfo Auth_S2C.ChannelIdInfo
	Map["Auth_S2C_ChannelIdInfo"] = reflect.TypeOf(aAuth_S2C_ChannelIdInfo)
	
	var aAuth_S2C_ChannelIdInfoList Auth_S2C.ChannelIdInfoList
	Map["Auth_S2C_ChannelIdInfoList"] = reflect.TypeOf(aAuth_S2C_ChannelIdInfoList)
	
	var aAuth_S2C_ChannelSimpleInfoList Auth_S2C.ChannelSimpleInfoList
	Map["Auth_S2C_ChannelSimpleInfoList"] = reflect.TypeOf(aAuth_S2C_ChannelSimpleInfoList)
	
	var aAuth_S2C_ChannelUserInfo Auth_S2C.ChannelUserInfo
	Map["Auth_S2C_ChannelUserInfo"] = reflect.TypeOf(aAuth_S2C_ChannelUserInfo)
	
	var aAuth_S2C_ChannelUserNum Auth_S2C.ChannelUserNum
	Map["Auth_S2C_ChannelUserNum"] = reflect.TypeOf(aAuth_S2C_ChannelUserNum)
	
	var aAuth_S2C_ChannelUserNumList Auth_S2C.ChannelUserNumList
	Map["Auth_S2C_ChannelUserNumList"] = reflect.TypeOf(aAuth_S2C_ChannelUserNumList)
	
	var aAuth_S2C_CheckVersionErrorInfo Auth_S2C.CheckVersionErrorInfo
	Map["Auth_S2C_CheckVersionErrorInfo"] = reflect.TypeOf(aAuth_S2C_CheckVersionErrorInfo)
	
	var aAuth_S2C_CreateChannelFailed Auth_S2C.CreateChannelFailed
	Map["Auth_S2C_CreateChannelFailed"] = reflect.TypeOf(aAuth_S2C_CreateChannelFailed)
	
	var aAuth_S2C_CreateChannelFailed_ErrorCodeType Auth_S2C.CreateChannelFailed_ErrorCodeType
	Map["Auth_S2C_CreateChannelFailed_ErrorCodeType"] = reflect.TypeOf(aAuth_S2C_CreateChannelFailed_ErrorCodeType)
	
	var aAuth_S2C_CreateChannelResult Auth_S2C.CreateChannelResult
	Map["Auth_S2C_CreateChannelResult"] = reflect.TypeOf(aAuth_S2C_CreateChannelResult)
	
	var aAuth_S2C_ErrorInfo Auth_S2C.ErrorInfo
	Map["Auth_S2C_ErrorInfo"] = reflect.TypeOf(aAuth_S2C_ErrorInfo)
	
	var aAuth_S2C_HelloInfoResult Auth_S2C.HelloInfoResult
	Map["Auth_S2C_HelloInfoResult"] = reflect.TypeOf(aAuth_S2C_HelloInfoResult)
	
	var aAuth_S2C_KeepAliveAck Auth_S2C.KeepAliveAck
	Map["Auth_S2C_KeepAliveAck"] = reflect.TypeOf(aAuth_S2C_KeepAliveAck)
	
	var aAuth_S2C_LoginTokenInfo Auth_S2C.LoginTokenInfo
	Map["Auth_S2C_LoginTokenInfo"] = reflect.TypeOf(aAuth_S2C_LoginTokenInfo)
	
	var aAuth_S2C_LoginUserInfo Auth_S2C.LoginUserInfo
	Map["Auth_S2C_LoginUserInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo)
	
	var aAuth_S2C_LoginUserInfo_LastLoginInfo Auth_S2C.LoginUserInfo_LastLoginInfo
	Map["Auth_S2C_LoginUserInfo_LastLoginInfo"] = reflect.TypeOf(aAuth_S2C_LoginUserInfo_LastLoginInfo)
	
	var aAuth_S2C_MonitorResult Auth_S2C.MonitorResult
	Map["Auth_S2C_MonitorResult"] = reflect.TypeOf(aAuth_S2C_MonitorResult)
	
	var aAuth_S2C_NeedFinishReg Auth_S2C.NeedFinishReg
	Map["Auth_S2C_NeedFinishReg"] = reflect.TypeOf(aAuth_S2C_NeedFinishReg)
	
	var aAuth_S2C_RequestAckInfo Auth_S2C.RequestAckInfo
	Map["Auth_S2C_RequestAckInfo"] = reflect.TypeOf(aAuth_S2C_RequestAckInfo)
	
	var aAuth_S2C_RequestTokenResult Auth_S2C.RequestTokenResult
	Map["Auth_S2C_RequestTokenResult"] = reflect.TypeOf(aAuth_S2C_RequestTokenResult)
	
	var aAuth_S2C_UserNotificationInfo Auth_S2C.UserNotificationInfo
	Map["Auth_S2C_UserNotificationInfo"] = reflect.TypeOf(aAuth_S2C_UserNotificationInfo)
	
	var aTOKEN_S_UserInfo TOKEN_S.UserInfo
	Map["TOKEN_S_UserInfo"] = reflect.TypeOf(aTOKEN_S_UserInfo)
	
}
