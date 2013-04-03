// Code generated by protoc-gen-go.
// source: proto files/ServiceC2S.proto
// DO NOT EDIT!

package Service_C2S

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Login struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Token            *string `protobuf:"bytes,2,req,name=token" json:"token,omitempty"`
	Version          *int32  `protobuf:"varint,3,req,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Login) Reset()         { *this = Login{} }
func (this *Login) String() string { return proto.CompactTextString(this) }
func (*Login) ProtoMessage()       {}

func (this *Login) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *Login) GetToken() string {
	if this != nil && this.Token != nil {
		return *this.Token
	}
	return ""
}

func (this *Login) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

type Chat struct {
	UserId           *int32  `protobuf:"varint,1,req" json:"UserId,omitempty"`
	MsgId            *int64  `protobuf:"varint,2,req" json:"MsgId,omitempty"`
	MetaTag          *string `protobuf:"bytes,3,opt" json:"MetaTag,omitempty"`
	Message          *string `protobuf:"bytes,4,req" json:"Message,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Chat) Reset()         { *this = Chat{} }
func (this *Chat) String() string { return proto.CompactTextString(this) }
func (*Chat) ProtoMessage()       {}

func (this *Chat) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *Chat) GetMsgId() int64 {
	if this != nil && this.MsgId != nil {
		return *this.MsgId
	}
	return 0
}

func (this *Chat) GetMetaTag() string {
	if this != nil && this.MetaTag != nil {
		return *this.MetaTag
	}
	return ""
}

func (this *Chat) GetMessage() string {
	if this != nil && this.Message != nil {
		return *this.Message
	}
	return ""
}

type MultiChat struct {
	UserIds          []int32 `protobuf:"varint,1,rep" json:"UserIds,omitempty"`
	MsgId            *int64  `protobuf:"varint,2,req" json:"MsgId,omitempty"`
	MetaTag          *string `protobuf:"bytes,3,opt" json:"MetaTag,omitempty"`
	Message          *string `protobuf:"bytes,4,req" json:"Message,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *MultiChat) Reset()         { *this = MultiChat{} }
func (this *MultiChat) String() string { return proto.CompactTextString(this) }
func (*MultiChat) ProtoMessage()       {}

func (this *MultiChat) GetMsgId() int64 {
	if this != nil && this.MsgId != nil {
		return *this.MsgId
	}
	return 0
}

func (this *MultiChat) GetMetaTag() string {
	if this != nil && this.MetaTag != nil {
		return *this.MetaTag
	}
	return ""
}

func (this *MultiChat) GetMessage() string {
	if this != nil && this.Message != nil {
		return *this.Message
	}
	return ""
}

type AddResult struct {
	FromId           *int32 `protobuf:"varint,1,req" json:"FromId,omitempty"`
	RequestId        *int64 `protobuf:"varint,2,req" json:"RequestId,omitempty"`
	Action           *int32 `protobuf:"varint,3,opt" json:"Action,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *AddResult) Reset()         { *this = AddResult{} }
func (this *AddResult) String() string { return proto.CompactTextString(this) }
func (*AddResult) ProtoMessage()       {}

func (this *AddResult) GetFromId() int32 {
	if this != nil && this.FromId != nil {
		return *this.FromId
	}
	return 0
}

func (this *AddResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *AddResult) GetAction() int32 {
	if this != nil && this.Action != nil {
		return *this.Action
	}
	return 0
}

type RemoveResult struct {
	FromId           *int32 `protobuf:"varint,1,req" json:"FromId,omitempty"`
	ErrorCode        *int32 `protobuf:"varint,2,opt" json:"ErrorCode,omitempty"`
	RequestId        *int64 `protobuf:"varint,3,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RemoveResult) Reset()         { *this = RemoveResult{} }
func (this *RemoveResult) String() string { return proto.CompactTextString(this) }
func (*RemoveResult) ProtoMessage()       {}

func (this *RemoveResult) GetFromId() int32 {
	if this != nil && this.FromId != nil {
		return *this.FromId
	}
	return 0
}

func (this *RemoveResult) GetErrorCode() int32 {
	if this != nil && this.ErrorCode != nil {
		return *this.ErrorCode
	}
	return 0
}

func (this *RemoveResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

type RequestUserInfo struct {
	UserId           *int32 `protobuf:"varint,2,opt" json:"UserId,omitempty"`
	RequestId        *int64 `protobuf:"varint,3,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RequestUserInfo) Reset()         { *this = RequestUserInfo{} }
func (this *RequestUserInfo) String() string { return proto.CompactTextString(this) }
func (*RequestUserInfo) ProtoMessage()       {}

func (this *RequestUserInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *RequestUserInfo) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

type AddRequest struct {
	ServiceName      *string `protobuf:"bytes,1,req" json:"ServiceName,omitempty"`
	RequestId        *int64  `protobuf:"varint,2,req" json:"RequestId,omitempty"`
	Message          *string `protobuf:"bytes,3,opt" json:"Message,omitempty"`
	Language         *int32  `protobuf:"varint,4,opt" json:"Language,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *AddRequest) Reset()         { *this = AddRequest{} }
func (this *AddRequest) String() string { return proto.CompactTextString(this) }
func (*AddRequest) ProtoMessage()       {}

func (this *AddRequest) GetServiceName() string {
	if this != nil && this.ServiceName != nil {
		return *this.ServiceName
	}
	return ""
}

func (this *AddRequest) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *AddRequest) GetMessage() string {
	if this != nil && this.Message != nil {
		return *this.Message
	}
	return ""
}

func (this *AddRequest) GetLanguage() int32 {
	if this != nil && this.Language != nil {
		return *this.Language
	}
	return 0
}

type RemoveRequest struct {
	ServiceId        *int32 `protobuf:"varint,1,req" json:"ServiceId,omitempty"`
	RequestId        *int64 `protobuf:"varint,2,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RemoveRequest) Reset()         { *this = RemoveRequest{} }
func (this *RemoveRequest) String() string { return proto.CompactTextString(this) }
func (*RemoveRequest) ProtoMessage()       {}

func (this *RemoveRequest) GetServiceId() int32 {
	if this != nil && this.ServiceId != nil {
		return *this.ServiceId
	}
	return 0
}

func (this *RemoveRequest) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

type MyServicesQuery struct {
	RequestId        *int64 `protobuf:"varint,1,opt" json:"RequestId,omitempty"`
	Version          *int32 `protobuf:"varint,2,opt" json:"Version,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *MyServicesQuery) Reset()         { *this = MyServicesQuery{} }
func (this *MyServicesQuery) String() string { return proto.CompactTextString(this) }
func (*MyServicesQuery) ProtoMessage()       {}

func (this *MyServicesQuery) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *MyServicesQuery) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func init() {
}
