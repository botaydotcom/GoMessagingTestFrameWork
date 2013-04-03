#pragma once
#ifndef DL_ERROR_CODE_H__
#define DL_ERROR_CODE_H__
namespace DL_ERROR_CODE
{
	enum
	{
		InvalidMessage = -1,
		InternalError  = -2,
		ItemSizeTooLarge = -3,
		CommentSizeTooLarge = -4,
		DuplicateRequest    = -5,
		RequestNonBuddyList = -6,
		ItemNotExist        = -7,
		NotAuther           = -8,
	};
}
#endif