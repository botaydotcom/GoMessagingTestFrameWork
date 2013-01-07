#ifndef _D_SYS_OP_CODE_
#define _D_SYS_OP_CODE_
#if _MSC_VER > 1000
	#pragma once
#endif
namespace AUTHS2CNS
{
	namespace LOGIN
	{
		enum
		{
			HelloInfoResult_CMD=0x01,
			LoginUserInfo_CMD=0x02,
			CheckVersionErrorInfo_CMD=0x03,
            NeedFinishReg_CMD = 0x04,
            KeepAliveAck_CMD = 0x05,
			NotificationInfo_CMD = 0x06,
            ChangePasswordAck_CMD = 0x07,
            PluginNotificationInfo_CMD = 0x08,

			DeviceTokenInvalid_CMD = 0x92,
		};
	}

	namespace COMMON
	{
		enum
		{
			RequestAckInfo_CMD=0x0A
		};
	}
	namespace CHANNEL
	{
		enum
		{
			ChannelIdInfoList_CMD=0x10,
			ChannelSimpleInfoList_CMD,
			CreateChannelResult_CMD,
			CreateChannelFailed_CMD,
			ChannelUserNumList_CMD
		};
	}
	
	namespace ERRORNS
	{
		enum
		{
			ErrorInfo_CMD=0x20
		};
	}
	
	namespace TOKEN
	{
		enum
		{
			RequestTokenResult_CMD=0x30
		};
	}
	
	namespace CORE
	{
		enum
		{
			Connected_CMD=235,
			Disconnect_CMD=236,
            MonitorResponse_CMD=237,

		};
	}
	
	namespace USER
	{
		enum
		{
			ChangeMyInfoResult_CMD=0x50
		};
	}
}

namespace AUTHC2SNS
{
	namespace LOGIN
	{
		enum
		{
			HelloInfo_CMD=0x01,
			LoginInfo_CMD=0x02,
			RegAccount_CMD=0x03,
            OAuthLogin_CMD = 0x04,
            FillFinishReg_CMD = 0x05,
            KeepAlive_CMD = 0x06,
			ChangeNotification_CMD = 0x07,
			RequestUserNotificationInfo_CMD = 0x08, 
            ChangePassword_CMD = 0x09,
			RequestUserPluginNotificationInfo_CMD = 0x0a,
            ChangePluginNotification_CMD = 0x0b,

            DeviceTokenSubmit_CMD = 0x90,
			DeviceTokenRemove_CMD = 0x91,
		};
	}
	
	namespace CHANNEL
	{
		enum
		{
			RequestChannelList_CMD=0x10,
			RequestChannelInfo_CMD,
			CreateChannelRequest_CMD,
			RequestUserNumList_CMD
		};
	}
	namespace TOKEN
	{
		enum
		{
			RequestToken_CMD=0x30
		};
	}
	
	namespace CORE
	{
		enum
		{
			Connected_CMD=235,
			Disconnect_CMD=236,
            Monitor_CMD = 237,
		};
	}
	
	namespace USER
	{
		enum
		{
			ChangeMyInfo_CMD=0x50
		};
	}
}
#endif //_D_SYS_OP_CODE_
