#pragma once
namespace AccountBindingOpCode
{
	namespace C2S
	{
		enum
		{
			BindRequest_CMD=0x40,
			ContactUpdate_CMD = 0x41,
			RequestBoundAccounts_CMD = 0x42,
			UnBindRequest_CMD = 0x43,
			RequestContactList_CMD = 0x44,
		};
	}

	namespace S2C
	{
		enum
		{
			BindResponse_CMD = 0x40,
			ContactUpdateResponse_CMD = 0x41,
			BoundAccountsResponse_CMD = 0x42,
			UnBindResponse_CMD = 0x43,
			ContactListResponse_CMD = 0x44,
		};
	}
}
