// Code generated by protoc-gen-go.
// source: proto files/ServiceS2C.proto
// DO NOT EDIT!

package Service_S2C

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type ServiceInfo struct {
	ServiceId        *int32  `protobuf:"varint,1,req" json:"ServiceId,omitempty"`
	ServiceName      *string `protobuf:"bytes,2,req" json:"ServiceName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ServiceInfo) Reset()         { *this = ServiceInfo{} }
func (this *ServiceInfo) String() string { return proto.CompactTextString(this) }
func (*ServiceInfo) ProtoMessage()       {}

func (this *ServiceInfo) GetServiceId() int32 {
	if this != nil && this.ServiceId != nil {
		return *this.ServiceId
	}
	return 0
}

func (this *ServiceInfo) GetServiceName() string {
	if this != nil && this.ServiceName != nil {
		return *this.ServiceName
	}
	return ""
}

type LoginInfo struct {
	MyInfo           *ServiceInfo             `protobuf:"bytes,1,req" json:"MyInfo,omitempty"`
	LastLogin        *LoginInfo_LastLoginInfo `protobuf:"bytes,2,opt" json:"LastLogin,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (this *LoginInfo) Reset()         { *this = LoginInfo{} }
func (this *LoginInfo) String() string { return proto.CompactTextString(this) }
func (*LoginInfo) ProtoMessage()       {}

func (this *LoginInfo) GetMyInfo() *ServiceInfo {
	if this != nil {
		return this.MyInfo
	}
	return nil
}

func (this *LoginInfo) GetLastLogin() *LoginInfo_LastLoginInfo {
	if this != nil {
		return this.LastLogin
	}
	return nil
}

type LoginInfo_LastLoginInfo struct {
	IP               *int32  `protobuf:"varint,1,opt" json:"IP,omitempty"`
	Time             *int32  `protobuf:"varint,2,opt" json:"Time,omitempty"`
	Country          *string `protobuf:"bytes,3,opt" json:"Country,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *LoginInfo_LastLoginInfo) Reset()         { *this = LoginInfo_LastLoginInfo{} }
func (this *LoginInfo_LastLoginInfo) String() string { return proto.CompactTextString(this) }
func (*LoginInfo_LastLoginInfo) ProtoMessage()       {}

func (this *LoginInfo_LastLoginInfo) GetIP() int32 {
	if this != nil && this.IP != nil {
		return *this.IP
	}
	return 0
}

func (this *LoginInfo_LastLoginInfo) GetTime() int32 {
	if this != nil && this.Time != nil {
		return *this.Time
	}
	return 0
}

func (this *LoginInfo_LastLoginInfo) GetCountry() string {
	if this != nil && this.Country != nil {
		return *this.Country
	}
	return ""
}

type Chat struct {
	UserId           *int32  `protobuf:"varint,1,req" json:"UserId,omitempty"`
	MetaTag          *string `protobuf:"bytes,2,opt" json:"MetaTag,omitempty"`
	Message          *string `protobuf:"bytes,3,req" json:"Message,omitempty"`
	Timestamp        *int32  `protobuf:"varint,4,opt" json:"Timestamp,omitempty"`
	MsgId            *int64  `protobuf:"varint,5,req" json:"MsgId,omitempty"`
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

func (this *Chat) GetTimestamp() int32 {
	if this != nil && this.Timestamp != nil {
		return *this.Timestamp
	}
	return 0
}

func (this *Chat) GetMsgId() int64 {
	if this != nil && this.MsgId != nil {
		return *this.MsgId
	}
	return 0
}

type RemoteAddRequest struct {
	UserId           *int32  `protobuf:"varint,1,req" json:"UserId,omitempty"`
	RequestId        *int64  `protobuf:"varint,2,req" json:"RequestId,omitempty"`
	Message          *string `protobuf:"bytes,3,opt" json:"Message,omitempty"`
	Language         *int32  `protobuf:"varint,4,opt" json:"Language,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *RemoteAddRequest) Reset()         { *this = RemoteAddRequest{} }
func (this *RemoteAddRequest) String() string { return proto.CompactTextString(this) }
func (*RemoteAddRequest) ProtoMessage()       {}

func (this *RemoteAddRequest) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *RemoteAddRequest) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *RemoteAddRequest) GetMessage() string {
	if this != nil && this.Message != nil {
		return *this.Message
	}
	return ""
}

func (this *RemoteAddRequest) GetLanguage() int32 {
	if this != nil && this.Language != nil {
		return *this.Language
	}
	return 0
}

type RemoteRemoveRequest struct {
	UserId           *int32 `protobuf:"varint,1,req" json:"UserId,omitempty"`
	RequestId        *int64 `protobuf:"varint,2,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RemoteRemoveRequest) Reset()         { *this = RemoteRemoveRequest{} }
func (this *RemoteRemoveRequest) String() string { return proto.CompactTextString(this) }
func (*RemoteRemoveRequest) ProtoMessage()       {}

func (this *RemoteRemoveRequest) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *RemoteRemoveRequest) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

type RemoteAddResult struct {
	Sid              *int32 `protobuf:"varint,1,req" json:"Sid,omitempty"`
	RequestId        *int64 `protobuf:"varint,2,req" json:"RequestId,omitempty"`
	Action           *int32 `protobuf:"varint,3,opt" json:"Action,omitempty"`
	Version          *int32 `protobuf:"varint,4,opt" json:"Version,omitempty"`
	OldVersion       *int32 `protobuf:"varint,5,opt" json:"OldVersion,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RemoteAddResult) Reset()         { *this = RemoteAddResult{} }
func (this *RemoteAddResult) String() string { return proto.CompactTextString(this) }
func (*RemoteAddResult) ProtoMessage()       {}

func (this *RemoteAddResult) GetSid() int32 {
	if this != nil && this.Sid != nil {
		return *this.Sid
	}
	return 0
}

func (this *RemoteAddResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *RemoteAddResult) GetAction() int32 {
	if this != nil && this.Action != nil {
		return *this.Action
	}
	return 0
}

func (this *RemoteAddResult) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func (this *RemoteAddResult) GetOldVersion() int32 {
	if this != nil && this.OldVersion != nil {
		return *this.OldVersion
	}
	return 0
}

type RemoteRemoveResult struct {
	Sid              *int32 `protobuf:"varint,1,req" json:"Sid,omitempty"`
	ErrorCode        *int32 `protobuf:"varint,2,opt" json:"ErrorCode,omitempty"`
	RequestId        *int64 `protobuf:"varint,3,opt" json:"RequestId,omitempty"`
	Version          *int32 `protobuf:"varint,4,opt" json:"Version,omitempty"`
	OldVersion       *int32 `protobuf:"varint,5,opt" json:"OldVersion,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RemoteRemoveResult) Reset()         { *this = RemoteRemoveResult{} }
func (this *RemoteRemoveResult) String() string { return proto.CompactTextString(this) }
func (*RemoteRemoveResult) ProtoMessage()       {}

func (this *RemoteRemoveResult) GetSid() int32 {
	if this != nil && this.Sid != nil {
		return *this.Sid
	}
	return 0
}

func (this *RemoteRemoveResult) GetErrorCode() int32 {
	if this != nil && this.ErrorCode != nil {
		return *this.ErrorCode
	}
	return 0
}

func (this *RemoteRemoveResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func (this *RemoteRemoveResult) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func (this *RemoteRemoveResult) GetOldVersion() int32 {
	if this != nil && this.OldVersion != nil {
		return *this.OldVersion
	}
	return 0
}

type MyServicesResult struct {
	Services         []*ServiceInfo `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
	Version          *int32         `protobuf:"varint,2,req" json:"Version,omitempty"`
	RequestId        *int64         `protobuf:"varint,3,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (this *MyServicesResult) Reset()         { *this = MyServicesResult{} }
func (this *MyServicesResult) String() string { return proto.CompactTextString(this) }
func (*MyServicesResult) ProtoMessage()       {}

func (this *MyServicesResult) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func (this *MyServicesResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

type UserInfo struct {
	Nickname         *string `protobuf:"bytes,1,req" json:"Nickname,omitempty"`
	Icon             *int64  `protobuf:"varint,2,req" json:"Icon,omitempty"`
	Exp              *int32  `protobuf:"varint,3,req" json:"Exp,omitempty"`
	Version          *int32  `protobuf:"varint,4,req" json:"Version,omitempty"`
	Gender           *int32  `protobuf:"varint,5,req" json:"Gender,omitempty"`
	Birthday         *int32  `protobuf:"varint,6,req" json:"Birthday,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserInfo) Reset()         { *this = UserInfo{} }
func (this *UserInfo) String() string { return proto.CompactTextString(this) }
func (*UserInfo) ProtoMessage()       {}

func (this *UserInfo) GetNickname() string {
	if this != nil && this.Nickname != nil {
		return *this.Nickname
	}
	return ""
}

func (this *UserInfo) GetIcon() int64 {
	if this != nil && this.Icon != nil {
		return *this.Icon
	}
	return 0
}

func (this *UserInfo) GetExp() int32 {
	if this != nil && this.Exp != nil {
		return *this.Exp
	}
	return 0
}

func (this *UserInfo) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func (this *UserInfo) GetGender() int32 {
	if this != nil && this.Gender != nil {
		return *this.Gender
	}
	return 0
}

func (this *UserInfo) GetBirthday() int32 {
	if this != nil && this.Birthday != nil {
		return *this.Birthday
	}
	return 0
}

type UserInfoResult struct {
	UserId           *int32    `protobuf:"varint,1,req" json:"UserId,omitempty"`
	Signature        []byte    `protobuf:"bytes,2,opt" json:"Signature,omitempty"`
	UserInfo         *UserInfo `protobuf:"bytes,3,opt" json:"UserInfo,omitempty"`
	RequestId        *int64    `protobuf:"varint,4,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (this *UserInfoResult) Reset()         { *this = UserInfoResult{} }
func (this *UserInfoResult) String() string { return proto.CompactTextString(this) }
func (*UserInfoResult) ProtoMessage()       {}

func (this *UserInfoResult) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *UserInfoResult) GetSignature() []byte {
	if this != nil {
		return this.Signature
	}
	return nil
}

func (this *UserInfoResult) GetUserInfo() *UserInfo {
	if this != nil {
		return this.UserInfo
	}
	return nil
}

func (this *UserInfoResult) GetRequestId() int64 {
	if this != nil && this.RequestId != nil {
		return *this.RequestId
	}
	return 0
}

func init() {
}
