/**
* @file   FSOpCode.h
* @brief  Header file
* @date   2012-02-10 11:15:03
* @author yanghx <yang.hx@gmail.com>
*/

#ifndef FSOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
#define FSOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
#if _MSC_VER > 1000
	#pragma once
#endif

namespace OpCode
{
	namespace C2S
	{
		enum
		{
			ImageUploadRequest_CMD=0x01,
			ImageDownloadRequest_CMD=0x02,
			ImageUploadRequestHeader_CMD=0x03,
			ImageUploadRequestPart_CMD=0x04,
			//
			BatchUploadRequestHeader_CMD=0x10,
			BatchUploadRequestPart_CMD=0x11,
		};
	}

	//
	namespace S2C
	{
		enum
		{
			ImageUploadResult_CMD=0x01,
			FSGeneralError_CMD=0x02,
			ImageDownloadResult_CMD=0x03,
			ImageUploadRequestPartAck_CMD=0x04,
			ImageUploadRequestHeaderAck_CMD=0x05,

			BatchUploadRequestHeaderAck_CMD=0x10,
			BatchUploadRequestPartAck_CMD=0x11
		};
	}

	namespace S2S
	{
		enum
		{
			QueryFileExistence_CMD = 0x20,
			FileRequest_CMD = 0x21,

			FileExistenceResult_CMD = 0x30,
			FileResponse_CMD = 0x31,
		};
	}
}

#endif //FSOPCODE_H_GUID_1FCE92F0_5451_400F_AB31_F054CB8FC831
