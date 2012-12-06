//
//  DiscussionOpCode.h
//  BTalkApp
//
//  Created by yang hx on 04/20/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#ifndef _DiscussionOpCode____FILEEXTENSION___
#define _DiscussionOpCode____FILEEXTENSION___


#define DISCUSSION_PACKET_BASE_COMMAND 0xA0
namespace DiscussionNS
{
    namespace C2S
    {
        enum
        {
            CreateDiscussion_CMD = 0x01,
            InviteMember_CMD,
            LeaveDiscussion_CMD,
            ChatInfo_CMD,
            ChangeMyOption_CMD,
            ChangeDiscussionInfo_CMD,
            RequestDiscussionInfo_CMD,
			InviteMemberResult_CMD,
			RequestMyDiscussion_CMD=0x09,
			RequestMyOption_CMD,
			ChatInfoRecvedAck_CMD
        };
    }
    namespace S2C
    {
        enum
        {
            DiscussionInfo_CMD = 0x01,
            UserDiscussionList_CMD,
            DiscussionMemberList_CMD,
            DiscussionError_CMD,
            NewDiscussionArrival_CMD,
            MemberLeave_CMD,
            MemberJoin_CMD,
            ChatInfo_CMD,
			InviteMember_CMD,
			ChatInfoAck_CMD,
			MyDiscussionOption_CMD
        };
    }
}

#endif
