// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: gRPC/chittychat.proto

package gRPC

import (
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

// Character limits should be enforced in either client or server.
type ChatReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Text string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *ChatReq) Reset() {
	*x = ChatReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gRPC_chittychat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatReq) ProtoMessage() {}

func (x *ChatReq) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_chittychat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatReq.ProtoReflect.Descriptor instead.
func (*ChatReq) Descriptor() ([]byte, []int) {
	return file_gRPC_chittychat_proto_rawDescGZIP(), []int{0}
}

func (x *ChatReq) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *ChatReq) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type ChatReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User      string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Text      string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Timestamp string `protobuf:"bytes,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *ChatReply) Reset() {
	*x = ChatReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gRPC_chittychat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatReply) ProtoMessage() {}

func (x *ChatReply) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_chittychat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatReply.ProtoReflect.Descriptor instead.
func (*ChatReply) Descriptor() ([]byte, []int) {
	return file_gRPC_chittychat_proto_rawDescGZIP(), []int{1}
}

func (x *ChatReply) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *ChatReply) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *ChatReply) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type Participant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *Participant) Reset() {
	*x = Participant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gRPC_chittychat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Participant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Participant) ProtoMessage() {}

func (x *Participant) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_chittychat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Participant.ProtoReflect.Descriptor instead.
func (*Participant) Descriptor() ([]byte, []int) {
	return file_gRPC_chittychat_proto_rawDescGZIP(), []int{2}
}

func (x *Participant) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

type JoinReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JoinMessage string `protobuf:"bytes,1,opt,name=JoinMessage,proto3" json:"JoinMessage,omitempty"`
	Timestamp   string `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *JoinReply) Reset() {
	*x = JoinReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gRPC_chittychat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinReply) ProtoMessage() {}

func (x *JoinReply) ProtoReflect() protoreflect.Message {
	mi := &file_gRPC_chittychat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinReply.ProtoReflect.Descriptor instead.
func (*JoinReply) Descriptor() ([]byte, []int) {
	return file_gRPC_chittychat_proto_rawDescGZIP(), []int{3}
}

func (x *JoinReply) GetJoinMessage() string {
	if x != nil {
		return x.JoinMessage
	}
	return ""
}

func (x *JoinReply) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

var File_gRPC_chittychat_proto protoreflect.FileDescriptor

var file_gRPC_chittychat_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x52, 0x50, 0x43, 0x2f, 0x63, 0x68, 0x69, 0x74, 0x74, 0x79, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43,
	0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x31, 0x0a, 0x07, 0x43,
	0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x51,
	0x0a, 0x09, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0x21, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x22, 0x4b, 0x0a, 0x09, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x20, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x32, 0xb3, 0x02, 0x0a, 0x0a, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74,
	0x12, 0x45, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x1b, 0x2e, 0x43, 0x68,
	0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74,
	0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68,
	0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4d, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x12, 0x1f, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x6e, 0x74, 0x1a, 0x1d, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68,
	0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x30, 0x01, 0x12, 0x46, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x1f,
	0x2e, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x1a,
	0x1d, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x47,
	0x0a, 0x05, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x12, 0x1f, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79,
	0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x1a, 0x1d, 0x2e, 0x43, 0x68, 0x69, 0x74, 0x74,
	0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4a, 0x6f,
	0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x5d, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x6b, 0x6b, 0x65, 0x6c, 0x42, 0x4b, 0x72, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x73, 0x65, 0x6e, 0x2f, 0x44, 0x53, 0x48, 0x61, 0x6e, 0x64, 0x69, 0x6e,
	0x73, 0x2f, 0x54, 0x72, 0x65, 0x65, 0x2f, 0x4d, 0x61, 0x69, 0x6e, 0x2f, 0x48, 0x61, 0x6e, 0x64,
	0x49, 0x6e, 0x33, 0x5f, 0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x2f, 0x43,
	0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x67, 0x52, 0x50, 0x43, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gRPC_chittychat_proto_rawDescOnce sync.Once
	file_gRPC_chittychat_proto_rawDescData = file_gRPC_chittychat_proto_rawDesc
)

func file_gRPC_chittychat_proto_rawDescGZIP() []byte {
	file_gRPC_chittychat_proto_rawDescOnce.Do(func() {
		file_gRPC_chittychat_proto_rawDescData = protoimpl.X.CompressGZIP(file_gRPC_chittychat_proto_rawDescData)
	})
	return file_gRPC_chittychat_proto_rawDescData
}

var file_gRPC_chittychat_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gRPC_chittychat_proto_goTypes = []interface{}{
	(*ChatReq)(nil),     // 0: ChittyChat_service.ChatReq
	(*ChatReply)(nil),   // 1: ChittyChat_service.ChatReply
	(*Participant)(nil), // 2: ChittyChat_service.Participant
	(*JoinReply)(nil),   // 3: ChittyChat_service.JoinReply
}
var file_gRPC_chittychat_proto_depIdxs = []int32{
	0, // 0: ChittyChat_service.ChittyChat.Publish:input_type -> ChittyChat_service.ChatReq
	2, // 1: ChittyChat_service.ChittyChat.Broadcast:input_type -> ChittyChat_service.Participant
	2, // 2: ChittyChat_service.ChittyChat.Join:input_type -> ChittyChat_service.Participant
	2, // 3: ChittyChat_service.ChittyChat.leave:input_type -> ChittyChat_service.Participant
	1, // 4: ChittyChat_service.ChittyChat.Publish:output_type -> ChittyChat_service.ChatReply
	1, // 5: ChittyChat_service.ChittyChat.Broadcast:output_type -> ChittyChat_service.ChatReply
	3, // 6: ChittyChat_service.ChittyChat.Join:output_type -> ChittyChat_service.JoinReply
	3, // 7: ChittyChat_service.ChittyChat.leave:output_type -> ChittyChat_service.JoinReply
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gRPC_chittychat_proto_init() }
func file_gRPC_chittychat_proto_init() {
	if File_gRPC_chittychat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gRPC_chittychat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatReq); i {
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
		file_gRPC_chittychat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatReply); i {
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
		file_gRPC_chittychat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Participant); i {
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
		file_gRPC_chittychat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinReply); i {
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
			RawDescriptor: file_gRPC_chittychat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gRPC_chittychat_proto_goTypes,
		DependencyIndexes: file_gRPC_chittychat_proto_depIdxs,
		MessageInfos:      file_gRPC_chittychat_proto_msgTypes,
	}.Build()
	File_gRPC_chittychat_proto = out.File
	file_gRPC_chittychat_proto_rawDesc = nil
	file_gRPC_chittychat_proto_goTypes = nil
	file_gRPC_chittychat_proto_depIdxs = nil
}
