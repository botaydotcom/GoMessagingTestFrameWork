// Code generated by protoc-gen-go.
// source: AuthC2SProtocol.proto
// DO NOT EDIT!

package Auth_C2S

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type RequestChannelList_RequestChannelType int32

const (
	RequestChannelList_DefaultRequest RequestChannelList_RequestChannelType = 0
	RequestChannelList_HotChannel     RequestChannelList_RequestChannelType = 1
	RequestChannelList_GarenaGame     RequestChannelList_RequestChannelType = 2
	RequestChannelList_MyChannel      RequestChannelList_RequestChannelType = 4
	RequestChannelList_AllChannel     RequestChannelList_RequestChannelType = 7
)

var RequestChannelList_RequestChannelType_name = map[int32]string{
	0: "DefaultRequest",
	1: "HotChannel",
	2: "GarenaGame",
	4: "MyChannel",
	7: "AllChannel",
}
var RequestChannelList_RequestChannelType_value = map[string]int32{
	"DefaultRequest": 0,
	"HotChannel":     1,
	"GarenaGame":     2,
	"MyChannel":      4,
	"AllChannel":     7,
}

func (x RequestChannelList_RequestChannelType) Enum() *RequestChannelList_RequestChannelType {
	p := new(RequestChannelList_RequestChannelType)
	*p = x
	return p
}
func (x RequestChannelList_RequestChannelType) String() string {
	return proto.EnumName(RequestChannelList_RequestChannelType_name, int32(x))
}
func (x RequestChannelList_RequestChannelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *RequestChannelList_RequestChannelType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(RequestChannelList_RequestChannelType_value, data, "RequestChannelList_RequestChannelType")
	if err != nil {
		return err
	}
	*x = RequestChannelList_RequestChannelType(value)
	return nil
}

type LoginInfo struct {
	Name             *string `protobuf:"bytes,1,req" json:"Name,omitempty"`
	Password         *string `protobuf:"bytes,2,req" json:"Password,omitempty"`
	ClientType       *int32  `protobuf:"varint,3,req" json:"ClientType,omitempty"`
	MachineId        *string `protobuf:"bytes,4,req" json:"MachineId,omitempty"`
	SoftwareVersion  *int32  `protobuf:"varint,5,req" json:"SoftwareVersion,omitempty"`
	Status           *int32  `protobuf:"varint,6,opt" json:"Status,omitempty"`
	Timestamp        *int64  `protobuf:"varint,7,opt,name=timestamp" json:"timestamp,omitempty"`
	DeviceId         []byte  `protobuf:"bytes,8,opt,name=deviceId" json:"deviceId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *LoginInfo) Reset()         { *this = LoginInfo{} }
func (this *LoginInfo) String() string { return proto.CompactTextString(this) }
func (*LoginInfo) ProtoMessage()       {}

func (this *LoginInfo) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *LoginInfo) GetPassword() string {
	if this != nil && this.Password != nil {
		return *this.Password
	}
	return ""
}

func (this *LoginInfo) GetClientType() int32 {
	if this != nil && this.ClientType != nil {
		return *this.ClientType
	}
	return 0
}

func (this *LoginInfo) GetMachineId() string {
	if this != nil && this.MachineId != nil {
		return *this.MachineId
	}
	return ""
}

func (this *LoginInfo) GetSoftwareVersion() int32 {
	if this != nil && this.SoftwareVersion != nil {
		return *this.SoftwareVersion
	}
	return 0
}

func (this *LoginInfo) GetStatus() int32 {
	if this != nil && this.Status != nil {
		return *this.Status
	}
	return 0
}

func (this *LoginInfo) GetTimestamp() int64 {
	if this != nil && this.Timestamp != nil {
		return *this.Timestamp
	}
	return 0
}

func (this *LoginInfo) GetDeviceId() []byte {
	if this != nil {
		return this.DeviceId
	}
	return nil
}

type HelloInfo struct {
	ClientType       *int32  `protobuf:"varint,1,req" json:"ClientType,omitempty"`
	Version          *uint32 `protobuf:"varint,2,req" json:"Version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *HelloInfo) Reset()         { *this = HelloInfo{} }
func (this *HelloInfo) String() string { return proto.CompactTextString(this) }
func (*HelloInfo) ProtoMessage()       {}

func (this *HelloInfo) GetClientType() int32 {
	if this != nil && this.ClientType != nil {
		return *this.ClientType
	}
	return 0
}

func (this *HelloInfo) GetVersion() uint32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

type RequestToken struct {
	Type             *int32 `protobuf:"varint,1,opt" json:"Type,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RequestToken) Reset()         { *this = RequestToken{} }
func (this *RequestToken) String() string { return proto.CompactTextString(this) }
func (*RequestToken) ProtoMessage()       {}

func (this *RequestToken) GetType() int32 {
	if this != nil && this.Type != nil {
		return *this.Type
	}
	return 0
}

type RequestChannelList struct {
	Type             *RequestChannelList_RequestChannelType `protobuf:"varint,1,req,enum=Auth.C2S.RequestChannelList_RequestChannelType" json:"Type,omitempty"`
	XXX_unrecognized []byte                                 `json:"-"`
}

func (this *RequestChannelList) Reset()         { *this = RequestChannelList{} }
func (this *RequestChannelList) String() string { return proto.CompactTextString(this) }
func (*RequestChannelList) ProtoMessage()       {}

func (this *RequestChannelList) GetType() RequestChannelList_RequestChannelType {
	if this != nil && this.Type != nil {
		return *this.Type
	}
	return 0
}

type RequestChannelInfo struct {
	ChannelId        []int32 `protobuf:"varint,1,rep" json:"ChannelId,omitempty"`
	Option           *int32  `protobuf:"varint,2,opt" json:"Option,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *RequestChannelInfo) Reset()         { *this = RequestChannelInfo{} }
func (this *RequestChannelInfo) String() string { return proto.CompactTextString(this) }
func (*RequestChannelInfo) ProtoMessage()       {}

func (this *RequestChannelInfo) GetOption() int32 {
	if this != nil && this.Option != nil {
		return *this.Option
	}
	return 0
}

type CreateChannelRequest struct {
	MainCateId       *int32  `protobuf:"varint,1,req" json:"MainCateId,omitempty"`
	SubCateId        *int32  `protobuf:"varint,2,req" json:"SubCateId,omitempty"`
	Title            *string `protobuf:"bytes,3,req" json:"Title,omitempty"`
	Memo             *string `protobuf:"bytes,4,opt" json:"Memo,omitempty"`
	Icon             *int64  `protobuf:"varint,5,opt" json:"Icon,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *CreateChannelRequest) Reset()         { *this = CreateChannelRequest{} }
func (this *CreateChannelRequest) String() string { return proto.CompactTextString(this) }
func (*CreateChannelRequest) ProtoMessage()       {}

func (this *CreateChannelRequest) GetMainCateId() int32 {
	if this != nil && this.MainCateId != nil {
		return *this.MainCateId
	}
	return 0
}

func (this *CreateChannelRequest) GetSubCateId() int32 {
	if this != nil && this.SubCateId != nil {
		return *this.SubCateId
	}
	return 0
}

func (this *CreateChannelRequest) GetTitle() string {
	if this != nil && this.Title != nil {
		return *this.Title
	}
	return ""
}

func (this *CreateChannelRequest) GetMemo() string {
	if this != nil && this.Memo != nil {
		return *this.Memo
	}
	return ""
}

func (this *CreateChannelRequest) GetIcon() int64 {
	if this != nil && this.Icon != nil {
		return *this.Icon
	}
	return 0
}

type RequestUserNumList struct {
	ChannelId        []int32 `protobuf:"varint,1,rep" json:"ChannelId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *RequestUserNumList) Reset()         { *this = RequestUserNumList{} }
func (this *RequestUserNumList) String() string { return proto.CompactTextString(this) }
func (*RequestUserNumList) ProtoMessage()       {}

type VMModuleInfoList struct {
	Modules          []*VMModuleInfoList_VMModuleInfo `protobuf:"bytes,1,rep" json:"Modules,omitempty"`
	XXX_unrecognized []byte                           `json:"-"`
}

func (this *VMModuleInfoList) Reset()         { *this = VMModuleInfoList{} }
func (this *VMModuleInfoList) String() string { return proto.CompactTextString(this) }
func (*VMModuleInfoList) ProtoMessage()       {}

type VMModuleInfoList_VMModuleInfo struct {
	Url              *string `protobuf:"bytes,1,req" json:"Url,omitempty"`
	Token            *string `protobuf:"bytes,2,req" json:"Token,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *VMModuleInfoList_VMModuleInfo) Reset()         { *this = VMModuleInfoList_VMModuleInfo{} }
func (this *VMModuleInfoList_VMModuleInfo) String() string { return proto.CompactTextString(this) }
func (*VMModuleInfoList_VMModuleInfo) ProtoMessage()       {}

func (this *VMModuleInfoList_VMModuleInfo) GetUrl() string {
	if this != nil && this.Url != nil {
		return *this.Url
	}
	return ""
}

func (this *VMModuleInfoList_VMModuleInfo) GetToken() string {
	if this != nil && this.Token != nil {
		return *this.Token
	}
	return ""
}

type ChangeMyInfo struct {
	NickName         *string `protobuf:"bytes,1,opt" json:"NickName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ChangeMyInfo) Reset()         { *this = ChangeMyInfo{} }
func (this *ChangeMyInfo) String() string { return proto.CompactTextString(this) }
func (*ChangeMyInfo) ProtoMessage()       {}

func (this *ChangeMyInfo) GetNickName() string {
	if this != nil && this.NickName != nil {
		return *this.NickName
	}
	return ""
}

type OAuthRawInfo struct {
	Provider         *string `protobuf:"bytes,1,req" json:"Provider,omitempty"`
	Account          *string `protobuf:"bytes,2,req" json:"Account,omitempty"`
	Content          []byte  `protobuf:"bytes,3,opt" json:"Content,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *OAuthRawInfo) Reset()         { *this = OAuthRawInfo{} }
func (this *OAuthRawInfo) String() string { return proto.CompactTextString(this) }
func (*OAuthRawInfo) ProtoMessage()       {}

func (this *OAuthRawInfo) GetProvider() string {
	if this != nil && this.Provider != nil {
		return *this.Provider
	}
	return ""
}

func (this *OAuthRawInfo) GetAccount() string {
	if this != nil && this.Account != nil {
		return *this.Account
	}
	return ""
}

func (this *OAuthRawInfo) GetContent() []byte {
	if this != nil {
		return this.Content
	}
	return nil
}

type OAuthLogin struct {
	OAuthInfo        []byte  `protobuf:"bytes,1,req" json:"OAuthInfo,omitempty"`
	ClientType       *int32  `protobuf:"varint,4,req" json:"ClientType,omitempty"`
	MachineId        *string `protobuf:"bytes,5,req" json:"MachineId,omitempty"`
	SoftwareVersion  *int32  `protobuf:"varint,6,req" json:"SoftwareVersion,omitempty"`
	Status           *int32  `protobuf:"varint,7,opt" json:"Status,omitempty"`
	Timestamp        *int64  `protobuf:"varint,8,opt" json:"Timestamp,omitempty"`
	DeviceId         []byte  `protobuf:"bytes,9,opt" json:"DeviceId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *OAuthLogin) Reset()         { *this = OAuthLogin{} }
func (this *OAuthLogin) String() string { return proto.CompactTextString(this) }
func (*OAuthLogin) ProtoMessage()       {}

func (this *OAuthLogin) GetOAuthInfo() []byte {
	if this != nil {
		return this.OAuthInfo
	}
	return nil
}

func (this *OAuthLogin) GetClientType() int32 {
	if this != nil && this.ClientType != nil {
		return *this.ClientType
	}
	return 0
}

func (this *OAuthLogin) GetMachineId() string {
	if this != nil && this.MachineId != nil {
		return *this.MachineId
	}
	return ""
}

func (this *OAuthLogin) GetSoftwareVersion() int32 {
	if this != nil && this.SoftwareVersion != nil {
		return *this.SoftwareVersion
	}
	return 0
}

func (this *OAuthLogin) GetStatus() int32 {
	if this != nil && this.Status != nil {
		return *this.Status
	}
	return 0
}

func (this *OAuthLogin) GetTimestamp() int64 {
	if this != nil && this.Timestamp != nil {
		return *this.Timestamp
	}
	return 0
}

func (this *OAuthLogin) GetDeviceId() []byte {
	if this != nil {
		return this.DeviceId
	}
	return nil
}

type FillFinishReg struct {
	UserId           *int32  `protobuf:"varint,1,opt" json:"UserId,omitempty"`
	NickName         *string `protobuf:"bytes,2,req" json:"NickName,omitempty"`
	Gender           *int32  `protobuf:"varint,3,req" json:"Gender,omitempty"`
	AvatarId         *uint64 `protobuf:"varint,4,req" json:"AvatarId,omitempty"`
	Birthday         *int32  `protobuf:"varint,5,opt" json:"Birthday,omitempty"`
	SoftwareVersion  *int32  `protobuf:"varint,6,opt,name=softwareVersion,def=10103" json:"softwareVersion,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *FillFinishReg) Reset()         { *this = FillFinishReg{} }
func (this *FillFinishReg) String() string { return proto.CompactTextString(this) }
func (*FillFinishReg) ProtoMessage()       {}

const Default_FillFinishReg_SoftwareVersion int32 = 10103

func (this *FillFinishReg) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *FillFinishReg) GetNickName() string {
	if this != nil && this.NickName != nil {
		return *this.NickName
	}
	return ""
}

func (this *FillFinishReg) GetGender() int32 {
	if this != nil && this.Gender != nil {
		return *this.Gender
	}
	return 0
}

func (this *FillFinishReg) GetAvatarId() uint64 {
	if this != nil && this.AvatarId != nil {
		return *this.AvatarId
	}
	return 0
}

func (this *FillFinishReg) GetBirthday() int32 {
	if this != nil && this.Birthday != nil {
		return *this.Birthday
	}
	return 0
}

func (this *FillFinishReg) GetSoftwareVersion() int32 {
	if this != nil && this.SoftwareVersion != nil {
		return *this.SoftwareVersion
	}
	return Default_FillFinishReg_SoftwareVersion
}

type KeepAlive struct {
	RequestId        []byte `protobuf:"bytes,1,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *KeepAlive) Reset()         { *this = KeepAlive{} }
func (this *KeepAlive) String() string { return proto.CompactTextString(this) }
func (*KeepAlive) ProtoMessage()       {}

func (this *KeepAlive) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

type RegAccount struct {
	EMail            *string `protobuf:"bytes,1,req" json:"EMail,omitempty"`
	Password         *string `protobuf:"bytes,2,req" json:"Password,omitempty"`
	Device           *string `protobuf:"bytes,3,opt" json:"Device,omitempty"`
	SoftwareVersion  *int32  `protobuf:"varint,4,opt,name=softwareVersion,def=10103" json:"softwareVersion,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *RegAccount) Reset()         { *this = RegAccount{} }
func (this *RegAccount) String() string { return proto.CompactTextString(this) }
func (*RegAccount) ProtoMessage()       {}

const Default_RegAccount_SoftwareVersion int32 = 10103

func (this *RegAccount) GetEMail() string {
	if this != nil && this.EMail != nil {
		return *this.EMail
	}
	return ""
}

func (this *RegAccount) GetPassword() string {
	if this != nil && this.Password != nil {
		return *this.Password
	}
	return ""
}

func (this *RegAccount) GetDevice() string {
	if this != nil && this.Device != nil {
		return *this.Device
	}
	return ""
}

func (this *RegAccount) GetSoftwareVersion() int32 {
	if this != nil && this.SoftwareVersion != nil {
		return *this.SoftwareVersion
	}
	return Default_RegAccount_SoftwareVersion
}

type RequestUserNotificationInfo struct {
	RequestId        []byte `protobuf:"bytes,1,opt" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RequestUserNotificationInfo) Reset()         { *this = RequestUserNotificationInfo{} }
func (this *RequestUserNotificationInfo) String() string { return proto.CompactTextString(this) }
func (*RequestUserNotificationInfo) ProtoMessage()       {}

func (this *RequestUserNotificationInfo) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

type RequestPluginNotificationInfo struct {
	RequestId        []byte `protobuf:"bytes,1,req" json:"RequestId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RequestPluginNotificationInfo) Reset()         { *this = RequestPluginNotificationInfo{} }
func (this *RequestPluginNotificationInfo) String() string { return proto.CompactTextString(this) }
func (*RequestPluginNotificationInfo) ProtoMessage()       {}

func (this *RequestPluginNotificationInfo) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

type ChangeNotification struct {
	Enabled          *bool  `protobuf:"varint,1,req" json:"Enabled,omitempty"`
	TimeSilence      *bool  `protobuf:"varint,2,opt" json:"TimeSilence,omitempty"`
	StartTime        *int32 `protobuf:"varint,3,opt" json:"StartTime,omitempty"`
	EndTime          *int32 `protobuf:"varint,4,opt" json:"EndTime,omitempty"`
	TimeZone         *int32 `protobuf:"varint,5,opt" json:"TimeZone,omitempty"`
	RequestId        []byte `protobuf:"bytes,6,opt" json:"RequestId,omitempty"`
	Language         *int32 `protobuf:"varint,7,opt" json:"Language,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *ChangeNotification) Reset()         { *this = ChangeNotification{} }
func (this *ChangeNotification) String() string { return proto.CompactTextString(this) }
func (*ChangeNotification) ProtoMessage()       {}

func (this *ChangeNotification) GetEnabled() bool {
	if this != nil && this.Enabled != nil {
		return *this.Enabled
	}
	return false
}

func (this *ChangeNotification) GetTimeSilence() bool {
	if this != nil && this.TimeSilence != nil {
		return *this.TimeSilence
	}
	return false
}

func (this *ChangeNotification) GetStartTime() int32 {
	if this != nil && this.StartTime != nil {
		return *this.StartTime
	}
	return 0
}

func (this *ChangeNotification) GetEndTime() int32 {
	if this != nil && this.EndTime != nil {
		return *this.EndTime
	}
	return 0
}

func (this *ChangeNotification) GetTimeZone() int32 {
	if this != nil && this.TimeZone != nil {
		return *this.TimeZone
	}
	return 0
}

func (this *ChangeNotification) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

func (this *ChangeNotification) GetLanguage() int32 {
	if this != nil && this.Language != nil {
		return *this.Language
	}
	return 0
}

type ChangePluginNotification struct {
	Enabled          *bool   `protobuf:"varint,1,req" json:"Enabled,omitempty"`
	Tag              *string `protobuf:"bytes,2,req" json:"Tag,omitempty"`
	RequestId        []byte  `protobuf:"bytes,3,req" json:"RequestId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ChangePluginNotification) Reset()         { *this = ChangePluginNotification{} }
func (this *ChangePluginNotification) String() string { return proto.CompactTextString(this) }
func (*ChangePluginNotification) ProtoMessage()       {}

func (this *ChangePluginNotification) GetEnabled() bool {
	if this != nil && this.Enabled != nil {
		return *this.Enabled
	}
	return false
}

func (this *ChangePluginNotification) GetTag() string {
	if this != nil && this.Tag != nil {
		return *this.Tag
	}
	return ""
}

func (this *ChangePluginNotification) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

type ChangePassword struct {
	OldPassword      *string `protobuf:"bytes,1,req" json:"OldPassword,omitempty"`
	NewPassword      *string `protobuf:"bytes,2,req" json:"NewPassword,omitempty"`
	RequestId        []byte  `protobuf:"bytes,3,req" json:"RequestId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ChangePassword) Reset()         { *this = ChangePassword{} }
func (this *ChangePassword) String() string { return proto.CompactTextString(this) }
func (*ChangePassword) ProtoMessage()       {}

func (this *ChangePassword) GetOldPassword() string {
	if this != nil && this.OldPassword != nil {
		return *this.OldPassword
	}
	return ""
}

func (this *ChangePassword) GetNewPassword() string {
	if this != nil && this.NewPassword != nil {
		return *this.NewPassword
	}
	return ""
}

func (this *ChangePassword) GetRequestId() []byte {
	if this != nil {
		return this.RequestId
	}
	return nil
}

type MonitorPoll struct {
	Id               *int64 `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Type             *int32 `protobuf:"varint,2,req,name=type" json:"type,omitempty"`
	Parameter        []byte `protobuf:"bytes,3,opt,name=parameter" json:"parameter,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *MonitorPoll) Reset()         { *this = MonitorPoll{} }
func (this *MonitorPoll) String() string { return proto.CompactTextString(this) }
func (*MonitorPoll) ProtoMessage()       {}

func (this *MonitorPoll) GetId() int64 {
	if this != nil && this.Id != nil {
		return *this.Id
	}
	return 0
}

func (this *MonitorPoll) GetType() int32 {
	if this != nil && this.Type != nil {
		return *this.Type
	}
	return 0
}

func (this *MonitorPoll) GetParameter() []byte {
	if this != nil {
		return this.Parameter
	}
	return nil
}

func init() {
	proto.RegisterEnum("Auth.C2S.RequestChannelList_RequestChannelType", RequestChannelList_RequestChannelType_name, RequestChannelList_RequestChannelType_value)
}
