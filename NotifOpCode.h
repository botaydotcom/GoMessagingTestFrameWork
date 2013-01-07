#pragma once
#ifndef NOTIF_OP_CODE_H__
#define NOTIF_OP_CODE_H__
#define NOTIFICATION_PACKET_BASE_COMMAND 0xA5
namespace NotifOpCode
{
	namespace S2S
	{
		enum
		{
			MulticastNotification_CMD = 0x01,
			NotifOperationState_CMD   = 0x02,
		};
	}            

    namespace C2S
    {
        enum
        {
            NotifSubscribe_CMD = 0x10,
            NotifUnsubscribe_CMD = 0x11,
        };
    }

    namespace S2C
    {
        enum
        {
            Notification_CMD        = 0x01,
			NotifOperationState_CMD = 0x02,
        };
    }
}
#endif // NOTIF_OP_CODE_H__
