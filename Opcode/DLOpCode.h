#ifndef DLOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
#define DLOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
#if _MSC_VER > 1000
	#pragma once
#endif

namespace OpCode
{
	namespace C2S
	{
		enum
		{
            DailyLifeRequest_MAINCMD = 161,
			PostDailyItem_CMD=0x01,
			PostComment_CMD=0x02,
			RequestDailyItems_CMD=0x03,
			RequestList_CMD=0x04,
			RequestAllList_CMD=0x05,
			RequestCommentDetail_CMD=0x06,
			RequestCommentContent_CMD=0x07,
			RequestCommentIDList_CMD=0x08,
			ResultAck_CMD = 0x09,
			DeleteDailyItem_CMD = 0x0a,
			RemoveBuddy_CMD = 0x0b,

            RequestRecentPosts_CMD=0x10,
            RequestRecentSummary_CMD=0x11,
            AddCircle_CMD=0x12,
            AddCircleMember_CMD=0x13,
            RemoveCircleMember_CMD=0x14,
            RemoveCircle_CMD=0x15,
            RequestCircles_CMD=0x16,
			RequestItemCircles_CMD=0x17,
		};
	}
    
	//
	namespace S2C
	{
		enum
		{
            DailyLifeResponse_MAINCMD = 161,
			OperationState_CMD=0x01,
			DailyItemResponse_CMD=0x02,
			DailyListResponse_CMD=0x03,
			AllDailyListResponse_CMD=0x04,
			CommentDetailResult_CMD=0x05,
			CommentContentResult_CMD=0x06,
			CommentIDListResponse_CMD=0x07,
            RecentPostsResponse_CMD=0x08,
            RecentSummaryResponse_CMD=0x09,
            CirclesResponse_CMD=0x0a,
			ItemCirclesResponse_CMD=0x0b,
		};
	}
}

#endif //DLOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
