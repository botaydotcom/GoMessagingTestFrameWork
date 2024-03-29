/**
* @file   ConstValue.h
* @brief  Header file for class CConstValue For AuthServer
* @date   Mon Sep 12 16:17:03 UTC+0800 2011
* @author yanghx<yanghx@garena.com>
*/

#ifndef _D_1832B0DB_B628_4900_8457_91259BDB3D98_H_
#define _D_1832B0DB_B628_4900_8457_91259BDB3D98_H_
#if _MSC_VER > 1000
	#pragma once
#endif
namespace ERROR_CODE
{
	enum
	{
		NO_SUCH_USER=0x01,
		LOGIN_CLIENT_VERSION_ERROR,
		LOGIN_PWD_NOT_MATCH,
		LOGIN_NAME_CONFUSE,
		LOGIN_DUPLICATE_USER,
		CHANGE_PASSWORD_OLD_PWD_NOT_MATCH,
		INVALID_MESSAGE,
		LOGIN_EXPIRED,
		VERIFICATION_ERROR,
		NEED_VERIFICATION,
		CONTACT_LIST_VERSION_ERROR,
		REG_EXISTS_ACCOUNT=0x10,
		FILL_FINISH_REG_NOT_EXIST = 0x11,
		SRY_DB_ERROR=0x20,
		SERVER_TOO_BUSY=0x30,
		NOT_LOGINED=0x40,		
	};

	namespace SYSTEM_LIMIT_NS
	{
		enum
		{
			USER_DISCUSSION_MAX_NUM=0x11
		};
	}
}

namespace DB_CONST
{
	enum
	{
		MIN_BUDDY_CATE_ID=100
	};
}

namespace SYSTEM_LIMIT
{
	enum
	{
		MAX_BUDDY_REQUEST=500,
		MAX_DISCUSSION_FOR_USER=10
	};

	namespace BUDDY
	{
		enum
		{
			CHAT_MAX_LENGTH=6000,
			SIGN_MAX_LENGTH=600,
			MEDIA_MAX_LENGTH=10000
		};
	}
}

namespace USER_EVENT_MASK
{
	enum
	{
		WHEN_DC_CHANGE_STATUS=0x01
	};
}


namespace CLIENT_TYPE_ENUM
{
	enum
	{
		PC_VERSION=1,
		MAC_VERSION=0x05,
		HAS_APN_SERVER=0x10,
		IPHONE_VERSION=HAS_APN_SERVER+1,
		ANDROID_VERSION=HAS_APN_SERVER+2
	};
}

#define DISCUSSION_MUTE_OPTION(v) (((v)&0x01)==0x01)
#define DISCUSSION_NO_NOTIFICATION_OPTION(v) (((v)&0x02)==0x02)
#endif //_D_1832B0DB_B628_4900_8457_91259BDB3D98_H_