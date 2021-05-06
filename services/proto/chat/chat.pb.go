// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: chat.proto

package chat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InfoUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Role string `protobuf:"bytes,2,opt,name=Role,proto3" json:"Role,omitempty"`
}

func (x *InfoUser) Reset() {
	*x = InfoUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoUser) ProtoMessage() {}

func (x *InfoUser) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoUser.ProtoReflect.Descriptor instead.
func (*InfoUser) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *InfoUser) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *InfoUser) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type Participants struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender    *InfoUser `protobuf:"bytes,1,opt,name=Sender,proto3" json:"Sender,omitempty"`
	Recipient *InfoUser `protobuf:"bytes,2,opt,name=Recipient,proto3" json:"Recipient,omitempty"`
}

func (x *Participants) Reset() {
	*x = Participants{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Participants) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Participants) ProtoMessage() {}

func (x *Participants) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Participants.ProtoReflect.Descriptor instead.
func (*Participants) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *Participants) GetSender() *InfoUser {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Participants) GetRecipient() *InfoUser {
	if x != nil {
		return x.Recipient
	}
	return nil
}

type Speakers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Speaker1 *InfoUser `protobuf:"bytes,1,opt,name=speaker1,proto3" json:"speaker1,omitempty"`
	Speaker2 *InfoUser `protobuf:"bytes,2,opt,name=speaker2,proto3" json:"speaker2,omitempty"`
}

func (x *Speakers) Reset() {
	*x = Speakers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Speakers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Speakers) ProtoMessage() {}

func (x *Speakers) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Speakers.ProtoReflect.Descriptor instead.
func (*Speakers) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *Speakers) GetSpeaker1() *InfoUser {
	if x != nil {
		return x.Speaker1
	}
	return nil
}

func (x *Speakers) GetSpeaker2() *InfoUser {
	if x != nil {
		return x.Speaker2
	}
	return nil
}

type InfoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32         `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Date         string        `protobuf:"bytes,2,opt,name=Date,proto3" json:"Date,omitempty"`
	Text         string        `protobuf:"bytes,3,opt,name=Text,proto3" json:"Text,omitempty"`
	Participants *Participants `protobuf:"bytes,4,opt,name=Participants,proto3" json:"Participants,omitempty"`
}

func (x *InfoMessage) Reset() {
	*x = InfoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoMessage) ProtoMessage() {}

func (x *InfoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoMessage.ProtoReflect.Descriptor instead.
func (*InfoMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{3}
}

func (x *InfoMessage) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *InfoMessage) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *InfoMessage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *InfoMessage) GetParticipants() *Participants {
	if x != nil {
		return x.Participants
	}
	return nil
}

type MoreInfoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	More []*InfoMessage `protobuf:"bytes,1,rep,name=More,proto3" json:"More,omitempty"`
}

func (x *MoreInfoMessage) Reset() {
	*x = MoreInfoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoreInfoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoreInfoMessage) ProtoMessage() {}

func (x *MoreInfoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoreInfoMessage.ProtoReflect.Descriptor instead.
func (*MoreInfoMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{4}
}

func (x *MoreInfoMessage) GetMore() []*InfoMessage {
	if x != nil {
		return x.More
	}
	return nil
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x22, 0x2e, 0x0a, 0x08, 0x49, 0x6e, 0x66, 0x6f, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x6e, 0x0a, 0x0c, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43,
	0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x06, 0x53, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43,
	0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x52, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x6c, 0x0a, 0x08, 0x53, 0x70, 0x65, 0x61, 0x6b,
	0x65, 0x72, 0x73, 0x12, 0x2f, 0x0a, 0x08, 0x73, 0x70, 0x65, 0x61, 0x6b, 0x65, 0x72, 0x31, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61,
	0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x08, 0x73, 0x70, 0x65, 0x61,
	0x6b, 0x65, 0x72, 0x31, 0x12, 0x2f, 0x0a, 0x08, 0x73, 0x70, 0x65, 0x61, 0x6b, 0x65, 0x72, 0x32,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68,
	0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x08, 0x73, 0x70, 0x65,
	0x61, 0x6b, 0x65, 0x72, 0x32, 0x22, 0x82, 0x01, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x78,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x3b, 0x0a,
	0x0c, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x2e,
	0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x0c, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x22, 0x3d, 0x0a, 0x0f, 0x4d, 0x6f,
	0x72, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a,
	0x04, 0x4d, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x04, 0x4d, 0x6f, 0x72, 0x65, 0x32, 0xcb, 0x01, 0x0a, 0x04, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x3e, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x43, 0x68, 0x61, 0x74,
	0x73, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e,
	0x66, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68,
	0x61, 0x74, 0x2e, 0x4d, 0x6f, 0x72, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74,
	0x2e, 0x53, 0x70, 0x65, 0x61, 0x6b, 0x65, 0x72, 0x73, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x43, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x6f, 0x72, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x40, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43,
	0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x2e, 0x49, 0x6e, 0x66, 0x6f,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6f, 0x72, 0x73, 0x63, 0x68, 0x74, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData = file_chat_proto_rawDesc
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_chat_proto_goTypes = []interface{}{
	(*InfoUser)(nil),        // 0: protoChat.InfoUser
	(*Participants)(nil),    // 1: protoChat.Participants
	(*Speakers)(nil),        // 2: protoChat.Speakers
	(*InfoMessage)(nil),     // 3: protoChat.InfoMessage
	(*MoreInfoMessage)(nil), // 4: protoChat.MoreInfoMessage
}
var file_chat_proto_depIdxs = []int32{
	0, // 0: protoChat.Participants.Sender:type_name -> protoChat.InfoUser
	0, // 1: protoChat.Participants.Recipient:type_name -> protoChat.InfoUser
	0, // 2: protoChat.Speakers.speaker1:type_name -> protoChat.InfoUser
	0, // 3: protoChat.Speakers.speaker2:type_name -> protoChat.InfoUser
	1, // 4: protoChat.InfoMessage.Participants:type_name -> protoChat.Participants
	3, // 5: protoChat.MoreInfoMessage.More:type_name -> protoChat.InfoMessage
	0, // 6: protoChat.Chat.GetAllChats:input_type -> protoChat.InfoUser
	2, // 7: protoChat.Chat.GetAllMessages:input_type -> protoChat.Speakers
	3, // 8: protoChat.Chat.ProcessMessage:input_type -> protoChat.InfoMessage
	4, // 9: protoChat.Chat.GetAllChats:output_type -> protoChat.MoreInfoMessage
	4, // 10: protoChat.Chat.GetAllMessages:output_type -> protoChat.MoreInfoMessage
	3, // 11: protoChat.Chat.ProcessMessage:output_type -> protoChat.InfoMessage
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoUser); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Participants); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Speakers); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MoreInfoMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatClient interface {
	GetAllChats(ctx context.Context, in *InfoUser, opts ...grpc.CallOption) (*MoreInfoMessage, error)
	GetAllMessages(ctx context.Context, in *Speakers, opts ...grpc.CallOption) (*MoreInfoMessage, error)
	ProcessMessage(ctx context.Context, in *InfoMessage, opts ...grpc.CallOption) (*InfoMessage, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) GetAllChats(ctx context.Context, in *InfoUser, opts ...grpc.CallOption) (*MoreInfoMessage, error) {
	out := new(MoreInfoMessage)
	err := c.cc.Invoke(ctx, "/protoChat.Chat/GetAllChats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetAllMessages(ctx context.Context, in *Speakers, opts ...grpc.CallOption) (*MoreInfoMessage, error) {
	out := new(MoreInfoMessage)
	err := c.cc.Invoke(ctx, "/protoChat.Chat/GetAllMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) ProcessMessage(ctx context.Context, in *InfoMessage, opts ...grpc.CallOption) (*InfoMessage, error) {
	out := new(InfoMessage)
	err := c.cc.Invoke(ctx, "/protoChat.Chat/ProcessMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
type ChatServer interface {
	GetAllChats(context.Context, *InfoUser) (*MoreInfoMessage, error)
	GetAllMessages(context.Context, *Speakers) (*MoreInfoMessage, error)
	ProcessMessage(context.Context, *InfoMessage) (*InfoMessage, error)
}

// UnimplementedChatServer can be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (*UnimplementedChatServer) GetAllChats(context.Context, *InfoUser) (*MoreInfoMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllChats not implemented")
}
func (*UnimplementedChatServer) GetAllMessages(context.Context, *Speakers) (*MoreInfoMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMessages not implemented")
}
func (*UnimplementedChatServer) ProcessMessage(context.Context, *InfoMessage) (*InfoMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessMessage not implemented")
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_GetAllChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetAllChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoChat.Chat/GetAllChats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetAllChats(ctx, req.(*InfoUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetAllMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Speakers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetAllMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoChat.Chat/GetAllMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetAllMessages(ctx, req.(*Speakers))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_ProcessMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).ProcessMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoChat.Chat/ProcessMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).ProcessMessage(ctx, req.(*InfoMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protoChat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllChats",
			Handler:    _Chat_GetAllChats_Handler,
		},
		{
			MethodName: "GetAllMessages",
			Handler:    _Chat_GetAllMessages_Handler,
		},
		{
			MethodName: "ProcessMessage",
			Handler:    _Chat_ProcessMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}
