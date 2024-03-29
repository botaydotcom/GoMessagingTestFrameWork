// Code generated by protoc-gen-go.
// source: AuthCacheProtocol.proto
// DO NOT EDIT!

package C_A

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import P_Common "CommonProtocol.pb"
// discarding unused import Auth_Buddy_S2C "S2CBuddyProtocol.pb"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type UserLoginInfo struct {
	ServerInfo       *uint64 `protobuf:"varint,1,req" json:"ServerInfo,omitempty"`
	SocketInfo       *uint64 `protobuf:"varint,2,req" json:"SocketInfo,omitempty"`
	UserId           *int32  `protobuf:"varint,3,req" json:"UserId,omitempty"`
	Status           *int32  `protobuf:"varint,4,req" json:"Status,omitempty"`
	ClientType       *int32  `protobuf:"varint,5,opt" json:"ClientType,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserLoginInfo) Reset()         { *this = UserLoginInfo{} }
func (this *UserLoginInfo) String() string { return proto.CompactTextString(this) }
func (*UserLoginInfo) ProtoMessage()       {}

func (this *UserLoginInfo) GetServerInfo() uint64 {
	if this != nil && this.ServerInfo != nil {
		return *this.ServerInfo
	}
	return 0
}

func (this *UserLoginInfo) GetSocketInfo() uint64 {
	if this != nil && this.SocketInfo != nil {
		return *this.SocketInfo
	}
	return 0
}

func (this *UserLoginInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *UserLoginInfo) GetStatus() int32 {
	if this != nil && this.Status != nil {
		return *this.Status
	}
	return 0
}

func (this *UserLoginInfo) GetClientType() int32 {
	if this != nil && this.ClientType != nil {
		return *this.ClientType
	}
	return 0
}

type ClientVersionInfo struct {
	ClientType       *int32 `protobuf:"varint,1,req" json:"ClientType,omitempty"`
	SoftwareVersion  *int32 `protobuf:"varint,2,req" json:"SoftwareVersion,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *ClientVersionInfo) Reset()         { *this = ClientVersionInfo{} }
func (this *ClientVersionInfo) String() string { return proto.CompactTextString(this) }
func (*ClientVersionInfo) ProtoMessage()       {}

func (this *ClientVersionInfo) GetClientType() int32 {
	if this != nil && this.ClientType != nil {
		return *this.ClientType
	}
	return 0
}

func (this *ClientVersionInfo) GetSoftwareVersion() int32 {
	if this != nil && this.SoftwareVersion != nil {
		return *this.SoftwareVersion
	}
	return 0
}

type UserIdInfo struct {
	UserId           *int32 `protobuf:"varint,1,req" json:"UserId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *UserIdInfo) Reset()         { *this = UserIdInfo{} }
func (this *UserIdInfo) String() string { return proto.CompactTextString(this) }
func (*UserIdInfo) ProtoMessage()       {}

func (this *UserIdInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

type UserInfo struct {
	UserId           *int32  `protobuf:"varint,1,req" json:"UserId,omitempty"`
	Account          *string `protobuf:"bytes,2,req" json:"Account,omitempty"`
	Password         *string `protobuf:"bytes,3,req" json:"Password,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserInfo) Reset()         { *this = UserInfo{} }
func (this *UserInfo) String() string { return proto.CompactTextString(this) }
func (*UserInfo) ProtoMessage()       {}

func (this *UserInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *UserInfo) GetAccount() string {
	if this != nil && this.Account != nil {
		return *this.Account
	}
	return ""
}

func (this *UserInfo) GetPassword() string {
	if this != nil && this.Password != nil {
		return *this.Password
	}
	return ""
}

type UserInfoEx struct {
	Level            *int32  `protobuf:"varint,1,opt,def=0" json:"Level,omitempty"`
	UserType         *int32  `protobuf:"varint,2,opt,def=0" json:"UserType,omitempty"`
	Icon             *uint64 `protobuf:"varint,3,opt,def=0" json:"Icon,omitempty"`
	NickName         *string `protobuf:"bytes,4,opt" json:"NickName,omitempty"`
	Version          *int32  `protobuf:"varint,5,opt,def=1" json:"Version,omitempty"`
	Gender           *int32  `protobuf:"varint,6,opt" json:"Gender,omitempty"`
	Birthday         *int32  `protobuf:"varint,7,opt" json:"Birthday,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserInfoEx) Reset()         { *this = UserInfoEx{} }
func (this *UserInfoEx) String() string { return proto.CompactTextString(this) }
func (*UserInfoEx) ProtoMessage()       {}

const Default_UserInfoEx_Level int32 = 0
const Default_UserInfoEx_UserType int32 = 0
const Default_UserInfoEx_Icon uint64 = 0
const Default_UserInfoEx_Version int32 = 1

func (this *UserInfoEx) GetLevel() int32 {
	if this != nil && this.Level != nil {
		return *this.Level
	}
	return Default_UserInfoEx_Level
}

func (this *UserInfoEx) GetUserType() int32 {
	if this != nil && this.UserType != nil {
		return *this.UserType
	}
	return Default_UserInfoEx_UserType
}

func (this *UserInfoEx) GetIcon() uint64 {
	if this != nil && this.Icon != nil {
		return *this.Icon
	}
	return Default_UserInfoEx_Icon
}

func (this *UserInfoEx) GetNickName() string {
	if this != nil && this.NickName != nil {
		return *this.NickName
	}
	return ""
}

func (this *UserInfoEx) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return Default_UserInfoEx_Version
}

func (this *UserInfoEx) GetGender() int32 {
	if this != nil && this.Gender != nil {
		return *this.Gender
	}
	return 0
}

func (this *UserInfoEx) GetBirthday() int32 {
	if this != nil && this.Birthday != nil {
		return *this.Birthday
	}
	return 0
}

type UserInfoToken struct {
	Token            *string `protobuf:"bytes,1,req" json:"Token,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserInfoToken) Reset()         { *this = UserInfoToken{} }
func (this *UserInfoToken) String() string { return proto.CompactTextString(this) }
func (*UserInfoToken) ProtoMessage()       {}

func (this *UserInfoToken) GetToken() string {
	if this != nil && this.Token != nil {
		return *this.Token
	}
	return ""
}

type CateChannelList struct {
	ChannelId        []int32 `protobuf:"varint,1,rep" json:"ChannelId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *CateChannelList) Reset()         { *this = CateChannelList{} }
func (this *CateChannelList) String() string { return proto.CompactTextString(this) }
func (*CateChannelList) ProtoMessage()       {}

type UserChannelList struct {
	UserChannel      []*UserChannelList_UserChannelInfo `protobuf:"bytes,1,rep" json:"UserChannel,omitempty"`
	XXX_unrecognized []byte                             `json:"-"`
}

func (this *UserChannelList) Reset()         { *this = UserChannelList{} }
func (this *UserChannelList) String() string { return proto.CompactTextString(this) }
func (*UserChannelList) ProtoMessage()       {}

type UserChannelList_UserChannelInfo struct {
	ChannelId        *int32 `protobuf:"varint,1,req" json:"ChannelId,omitempty"`
	Membership       *int32 `protobuf:"varint,2,req" json:"Membership,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *UserChannelList_UserChannelInfo) Reset()         { *this = UserChannelList_UserChannelInfo{} }
func (this *UserChannelList_UserChannelInfo) String() string { return proto.CompactTextString(this) }
func (*UserChannelList_UserChannelInfo) ProtoMessage()       {}

func (this *UserChannelList_UserChannelInfo) GetChannelId() int32 {
	if this != nil && this.ChannelId != nil {
		return *this.ChannelId
	}
	return 0
}

func (this *UserChannelList_UserChannelInfo) GetMembership() int32 {
	if this != nil && this.Membership != nil {
		return *this.Membership
	}
	return 0
}

type BuddyRelationInfo struct {
	UserId           *int32 `protobuf:"varint,1,req" json:"UserId,omitempty"`
	CateId           *int32 `protobuf:"varint,2,req" json:"CateId,omitempty"`
	Relationship     *int32 `protobuf:"varint,3,req" json:"Relationship,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *BuddyRelationInfo) Reset()         { *this = BuddyRelationInfo{} }
func (this *BuddyRelationInfo) String() string { return proto.CompactTextString(this) }
func (*BuddyRelationInfo) ProtoMessage()       {}

func (this *BuddyRelationInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *BuddyRelationInfo) GetCateId() int32 {
	if this != nil && this.CateId != nil {
		return *this.CateId
	}
	return 0
}

func (this *BuddyRelationInfo) GetRelationship() int32 {
	if this != nil && this.Relationship != nil {
		return *this.Relationship
	}
	return 0
}

type UserBuddyList struct {
	Buddy            []*BuddyRelationInfo `protobuf:"bytes,1,rep" json:"Buddy,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (this *UserBuddyList) Reset()         { *this = UserBuddyList{} }
func (this *UserBuddyList) String() string { return proto.CompactTextString(this) }
func (*UserBuddyList) ProtoMessage()       {}

type UserSignatureInfo struct {
	Signature        []byte `protobuf:"bytes,1,opt" json:"Signature,omitempty"`
	Version          *int32 `protobuf:"varint,2,opt" json:"Version,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *UserSignatureInfo) Reset()         { *this = UserSignatureInfo{} }
func (this *UserSignatureInfo) String() string { return proto.CompactTextString(this) }
func (*UserSignatureInfo) ProtoMessage()       {}

func (this *UserSignatureInfo) GetSignature() []byte {
	if this != nil {
		return this.Signature
	}
	return nil
}

func (this *UserSignatureInfo) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

type UserStatusInfo struct {
	Status           *int32 `protobuf:"varint,1,req" json:"Status,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *UserStatusInfo) Reset()         { *this = UserStatusInfo{} }
func (this *UserStatusInfo) String() string { return proto.CompactTextString(this) }
func (*UserStatusInfo) ProtoMessage()       {}

func (this *UserStatusInfo) GetStatus() int32 {
	if this != nil && this.Status != nil {
		return *this.Status
	}
	return 0
}

type UserChatInfoList struct {
	Chat             [][]byte `protobuf:"bytes,1,rep" json:"Chat,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (this *UserChatInfoList) Reset()         { *this = UserChatInfoList{} }
func (this *UserChatInfoList) String() string { return proto.CompactTextString(this) }
func (*UserChatInfoList) ProtoMessage()       {}

type BuddyInfoHolder struct {
	UserId           *int32             `protobuf:"varint,1,req" json:"UserId,omitempty"`
	BasicInfo        *UserInfo          `protobuf:"bytes,2,opt" json:"BasicInfo,omitempty"`
	UserInfo         *UserInfoEx        `protobuf:"bytes,3,opt" json:"UserInfo,omitempty"`
	StatusInfo       *UserStatusInfo    `protobuf:"bytes,4,opt" json:"StatusInfo,omitempty"`
	RelationInfo     *BuddyRelationInfo `protobuf:"bytes,5,opt" json:"RelationInfo,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (this *BuddyInfoHolder) Reset()         { *this = BuddyInfoHolder{} }
func (this *BuddyInfoHolder) String() string { return proto.CompactTextString(this) }
func (*BuddyInfoHolder) ProtoMessage()       {}

func (this *BuddyInfoHolder) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *BuddyInfoHolder) GetBasicInfo() *UserInfo {
	if this != nil {
		return this.BasicInfo
	}
	return nil
}

func (this *BuddyInfoHolder) GetUserInfo() *UserInfoEx {
	if this != nil {
		return this.UserInfo
	}
	return nil
}

func (this *BuddyInfoHolder) GetStatusInfo() *UserStatusInfo {
	if this != nil {
		return this.StatusInfo
	}
	return nil
}

func (this *BuddyInfoHolder) GetRelationInfo() *BuddyRelationInfo {
	if this != nil {
		return this.RelationInfo
	}
	return nil
}

type BuddyInfoHolderList struct {
	Buddy            []*BuddyInfoHolder `protobuf:"bytes,1,rep" json:"Buddy,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (this *BuddyInfoHolderList) Reset()         { *this = BuddyInfoHolderList{} }
func (this *BuddyInfoHolderList) String() string { return proto.CompactTextString(this) }
func (*BuddyInfoHolderList) ProtoMessage()       {}

type UserInfoList struct {
	Info             []*UserInfoEx `protobuf:"bytes,1,rep" json:"Info,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (this *UserInfoList) Reset()         { *this = UserInfoList{} }
func (this *UserInfoList) String() string { return proto.CompactTextString(this) }
func (*UserInfoList) ProtoMessage()       {}

type UserImportInfo struct {
	Imported         *int32 `protobuf:"varint,1,req" json:"Imported,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *UserImportInfo) Reset()         { *this = UserImportInfo{} }
func (this *UserImportInfo) String() string { return proto.CompactTextString(this) }
func (*UserImportInfo) ProtoMessage()       {}

func (this *UserImportInfo) GetImported() int32 {
	if this != nil && this.Imported != nil {
		return *this.Imported
	}
	return 0
}

type UserMemDataInfo struct {
	Tag              *string `protobuf:"bytes,1,req" json:"Tag,omitempty"`
	Content          []byte  `protobuf:"bytes,2,req" json:"Content,omitempty"`
	Version          *int32  `protobuf:"varint,3,req" json:"Version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserMemDataInfo) Reset()         { *this = UserMemDataInfo{} }
func (this *UserMemDataInfo) String() string { return proto.CompactTextString(this) }
func (*UserMemDataInfo) ProtoMessage()       {}

func (this *UserMemDataInfo) GetTag() string {
	if this != nil && this.Tag != nil {
		return *this.Tag
	}
	return ""
}

func (this *UserMemDataInfo) GetContent() []byte {
	if this != nil {
		return this.Content
	}
	return nil
}

func (this *UserMemDataInfo) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

type UserMemDataList struct {
	MemData          []*UserMemDataInfo `protobuf:"bytes,1,rep" json:"MemData,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (this *UserMemDataList) Reset()         { *this = UserMemDataList{} }
func (this *UserMemDataList) String() string { return proto.CompactTextString(this) }
func (*UserMemDataList) ProtoMessage()       {}

type UserDiscussionList struct {
	DiscussionId     []int64 `protobuf:"varint,1,rep" json:"DiscussionId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *UserDiscussionList) Reset()         { *this = UserDiscussionList{} }
func (this *UserDiscussionList) String() string { return proto.CompactTextString(this) }
func (*UserDiscussionList) ProtoMessage()       {}

type DiscussionInfo struct {
	DiscussionId     *int64  `protobuf:"varint,1,req" json:"DiscussionId,omitempty"`
	Title            *string `protobuf:"bytes,2,req" json:"Title,omitempty"`
	Kind             *int32  `protobuf:"varint,3,opt" json:"Kind,omitempty"`
	Option           *int32  `protobuf:"varint,4,opt" json:"Option,omitempty"`
	Icon             *int64  `protobuf:"varint,5,opt" json:"Icon,omitempty"`
	Version          *int32  `protobuf:"varint,6,req" json:"Version,omitempty"`
	MemberVersion    *int32  `protobuf:"varint,7,opt" json:"MemberVersion,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *DiscussionInfo) Reset()         { *this = DiscussionInfo{} }
func (this *DiscussionInfo) String() string { return proto.CompactTextString(this) }
func (*DiscussionInfo) ProtoMessage()       {}

func (this *DiscussionInfo) GetDiscussionId() int64 {
	if this != nil && this.DiscussionId != nil {
		return *this.DiscussionId
	}
	return 0
}

func (this *DiscussionInfo) GetTitle() string {
	if this != nil && this.Title != nil {
		return *this.Title
	}
	return ""
}

func (this *DiscussionInfo) GetKind() int32 {
	if this != nil && this.Kind != nil {
		return *this.Kind
	}
	return 0
}

func (this *DiscussionInfo) GetOption() int32 {
	if this != nil && this.Option != nil {
		return *this.Option
	}
	return 0
}

func (this *DiscussionInfo) GetIcon() int64 {
	if this != nil && this.Icon != nil {
		return *this.Icon
	}
	return 0
}

func (this *DiscussionInfo) GetVersion() int32 {
	if this != nil && this.Version != nil {
		return *this.Version
	}
	return 0
}

func (this *DiscussionInfo) GetMemberVersion() int32 {
	if this != nil && this.MemberVersion != nil {
		return *this.MemberVersion
	}
	return 0
}

type DiscussionMemberInfo struct {
	UserId           *int32 `protobuf:"varint,1,req" json:"UserId,omitempty"`
	Position         *int32 `protobuf:"varint,2,opt" json:"Position,omitempty"`
	Option           *int32 `protobuf:"varint,3,opt" json:"Option,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *DiscussionMemberInfo) Reset()         { *this = DiscussionMemberInfo{} }
func (this *DiscussionMemberInfo) String() string { return proto.CompactTextString(this) }
func (*DiscussionMemberInfo) ProtoMessage()       {}

func (this *DiscussionMemberInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func (this *DiscussionMemberInfo) GetPosition() int32 {
	if this != nil && this.Position != nil {
		return *this.Position
	}
	return 0
}

func (this *DiscussionMemberInfo) GetOption() int32 {
	if this != nil && this.Option != nil {
		return *this.Option
	}
	return 0
}

type DiscussionMemberList struct {
	Member           []*DiscussionMemberInfo `protobuf:"bytes,1,rep" json:"Member,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (this *DiscussionMemberList) Reset()         { *this = DiscussionMemberList{} }
func (this *DiscussionMemberList) String() string { return proto.CompactTextString(this) }
func (*DiscussionMemberList) ProtoMessage()       {}

type DiscussionInviteMemberList struct {
	Member           []*DiscussionInviteMemberList_InviteMemberInfo `protobuf:"bytes,1,rep" json:"Member,omitempty"`
	XXX_unrecognized []byte                                         `json:"-"`
}

func (this *DiscussionInviteMemberList) Reset()         { *this = DiscussionInviteMemberList{} }
func (this *DiscussionInviteMemberList) String() string { return proto.CompactTextString(this) }
func (*DiscussionInviteMemberList) ProtoMessage()       {}

type DiscussionInviteMemberList_InviteMemberInfo struct {
	InviteId         []byte `protobuf:"bytes,1,req" json:"InviteId,omitempty"`
	UserId           *int32 `protobuf:"varint,2,req" json:"UserId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *DiscussionInviteMemberList_InviteMemberInfo) Reset() {
	*this = DiscussionInviteMemberList_InviteMemberInfo{}
}
func (this *DiscussionInviteMemberList_InviteMemberInfo) String() string {
	return proto.CompactTextString(this)
}
func (*DiscussionInviteMemberList_InviteMemberInfo) ProtoMessage() {}

func (this *DiscussionInviteMemberList_InviteMemberInfo) GetInviteId() []byte {
	if this != nil {
		return this.InviteId
	}
	return nil
}

func (this *DiscussionInviteMemberList_InviteMemberInfo) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

type CommonMessageItemHolder struct {
	Command          *int32 `protobuf:"varint,1,req" json:"Command,omitempty"`
	Body             []byte `protobuf:"bytes,2,req" json:"Body,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *CommonMessageItemHolder) Reset()         { *this = CommonMessageItemHolder{} }
func (this *CommonMessageItemHolder) String() string { return proto.CompactTextString(this) }
func (*CommonMessageItemHolder) ProtoMessage()       {}

func (this *CommonMessageItemHolder) GetCommand() int32 {
	if this != nil && this.Command != nil {
		return *this.Command
	}
	return 0
}

func (this *CommonMessageItemHolder) GetBody() []byte {
	if this != nil {
		return this.Body
	}
	return nil
}

type CommonMessageListHolder struct {
	UserId           *int32                     `protobuf:"varint,1,req" json:"UserId,omitempty"`
	Item             []*CommonMessageItemHolder `protobuf:"bytes,2,rep" json:"Item,omitempty"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (this *CommonMessageListHolder) Reset()         { *this = CommonMessageListHolder{} }
func (this *CommonMessageListHolder) String() string { return proto.CompactTextString(this) }
func (*CommonMessageListHolder) ProtoMessage()       {}

func (this *CommonMessageListHolder) GetUserId() int32 {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return 0
}

func init() {
}
