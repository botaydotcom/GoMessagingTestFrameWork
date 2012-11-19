#ifndef _D_BUDDY_OPCODE_
#define _D_BUDDY_OPCODE_
#if _MSC_VER > 1000
#pragma once
#endif

/*
public static class BuddyOpCode
{
public static final byte CategoryList_CMD = 0x60;
public static final byte EditCategoryError_CMD = 0x61;
public static final byte NewCategoryResponse_CMD = 0x62;
public static final byte DeleteCategoryError_CMD = 0x63;
public static final byte DeleteCategoryResponse_CMD = 0x64;
public static final byte BuddySimpleInfoList_CMD = 0x65;
public static final byte RequestAddBuddyError_CMD = 0x66;
public static final byte RemoteAddBuddyResult_CMD = 0x67;
public static final byte BuddyDeleted_CMD = 0x68;
public static final byte DeleteBuddyError_CMD = 0x69;
public static final byte BuddyInfoList_CMD = 0x6A;
public static final byte UserSignatureList_CMD = 0x6B;
public static final byte BuddyOnlineList_CMD = 0x6C;
public static final byte MyInfoChange_CMD = 0x6D;
}
 */
namespace AUTHS2C
{
	namespace BuddyOpCode
	{
		enum
		{
			CategoryList_CMD=0x60,
			EditCategoryError_CMD,
			EditCategoryResponse_CMD,
			DeleteCategoryError_CMD,
			DeleteCategoryResponse_CMD,
			BuddySimpleInfoList_CMD,
			RequestAddBuddyError_CMD,
			RemoteAddBuddyResult_CMD,
			BuddyDeleted_CMD,
			DeleteBuddyError_CMD,
			BuddyInfoList_CMD,
			UserSignatureList_CMD,
			BuddyOnlineList_CMD,
			DynamicModuleList_CMD,
			ChatInfo_CMD,
			ChatAck_CMD,
			RequestUserInfoResult_CMD=0x70,
			RequestUserInfoError_CMD,
			RemoteRequestAddBuddy_CMD,
			BuddyOnline_CMD,
			ChangeBuddyCategoryResult_CMD,
			MyInfoChanged_CMD,
			//New Command
			ChatMediaInfoAck_CMD,
			ChatMediaInfo_CMD,
			RemoteP2PRequest_CMD,
			RemoteP2PResponse_CMD,
			RemoteP2PResponseAck_CMD,
			RemoteDoActionCode_CMD,
			RemoteChatAck_CMD=0x7c,
			UserSignatureInfo_CMD=0x7d,
			//
			ChangeMeMemDataResult_CMD,
			BuddyMemDataInfo_CMD,
			RequestAllBuddyMemDataResult_CMD=0x80,
            MobileUserLocationResponse_CMD=0x81,

			ChatInfo2_CMD,
			Chat2Ack_CMD,
			RemoteChat2Ack_CMD
		};
	}
}

/*
public static class BuddyOpCode
{
public static final byte RequestCategoryList_CMD = 0x60;
public static final byte EditCategory_CMD = 0x61;
public static final byte DeleteCategory_CMD = 0x62;
public static final byte RequestBuddyOnlineList_CMD = 0x63;
public static final byte RequestBuddySimpleInfoList_CMD = 0x64;
public static final byte RequestAddBuddy_CMD = 0x65;
public static final byte AddBuddyResult_CMD = 0x66;
public static final byte DeleteBuddy_CMD = 0x67;
public static final byte RequestBuddyInfo_CMD = 0x68;
public static final byte RequestSignature_CMD = 0x69;
public static final byte RequestBuddySignature_CMD = 0x6A;
public static final byte RequestMyOnlineStatusChange_CMD = 0x6B;
}
}
 */
namespace AUTHC2S
{
	namespace BuddyOpCode
	{
		enum
		{
			RequestCategoryList_CMD=0x60,
			EditCategory_CMD,
			DeleteCategory_CMD,
			RequestBuddyOnlineList_CMD,
			RequestBuddySimpleInfoList_CMD,
			RequestAddBuddy_CMD,
			AddBuddyResult_CMD,
			DeleteBuddy_CMD,
			RequestBuddyInfo_CMD,
			RequestSignature_CMD,
			RequestBuddySignature_CMD,
			RequestMyOnlineStatusChange_CMD,
			ChatInfo_CMD,
			RequestUserInfo_CMD,
			ChangeBuddyCategory_CMD,
			ChangeUserInfo_CMD,
			//New Command
			ChatMediaInfo_CMD=0x70,
			DoActionCode_CMD,
			P2PRequest_CMD,
			P2PResponse_CMD,
			//MemData 
			ChangeMeMemData_CMD,
			RequestAllBuddyMemData_CMD,
			RequestOneBuddyMemData_CMD=0x76,
            ChatAckRemote_CMD=0x77,
            RequestSyncBuddy_CMD,
            MobileFeatureAction_CMD=0x79,
            MobileClearAction_CMD=0x80,

			ChatInfo2_CMD,
			Chat2AckRemote_CMD
		};
	}
}
#endif