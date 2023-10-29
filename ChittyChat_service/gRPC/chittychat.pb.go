// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: grpc/chittychat.proto

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username  string `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
	Timestamp int32  `protobuf:"varint,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_chittychat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_chittychat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_grpc_chittychat_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Message) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Message) GetTimestamp() int32 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_grpc_chittychat_proto protoreflect.FileDescriptor

var file_grpc_chittychat_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x68, 0x69, 0x74, 0x74, 0x79, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x52, 0x50, 0x43, 0x22, 0x5d, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x3d, 0x0a, 0x0a,
	0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x12, 0x2f, 0x0a, 0x0b, 0x43, 0x68,
	0x61, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0d, 0x2e, 0x67, 0x52, 0x50, 0x43,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x67, 0x52, 0x50, 0x43, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x4a, 0x5a, 0x48, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x6b, 0x6b, 0x65, 0x6c,
	0x42, 0x4b, 0x72, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x73, 0x65, 0x6e, 0x2f, 0x44, 0x53, 0x48, 0x61,
	0x6e, 0x64, 0x69, 0x6e, 0x73, 0x2f, 0x54, 0x72, 0x65, 0x65, 0x2f, 0x4d, 0x61, 0x69, 0x6e, 0x2f,
	0x43, 0x68, 0x69, 0x74, 0x74, 0x79, 0x43, 0x68, 0x61, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x67, 0x52, 0x50, 0x43, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_chittychat_proto_rawDescOnce sync.Once
	file_grpc_chittychat_proto_rawDescData = file_grpc_chittychat_proto_rawDesc
)

func file_grpc_chittychat_proto_rawDescGZIP() []byte {
	file_grpc_chittychat_proto_rawDescOnce.Do(func() {
		file_grpc_chittychat_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_chittychat_proto_rawDescData)
	})
	return file_grpc_chittychat_proto_rawDescData
}

var file_grpc_chittychat_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_grpc_chittychat_proto_goTypes = []interface{}{
	(*Message)(nil), // 0: gRPC.Message
}
var file_grpc_chittychat_proto_depIdxs = []int32{
	0, // 0: gRPC.ChittyChat.ChatService:input_type -> gRPC.Message
	0, // 1: gRPC.ChittyChat.ChatService:output_type -> gRPC.Message
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_chittychat_proto_init() }
func file_grpc_chittychat_proto_init() {
	if File_grpc_chittychat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_chittychat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_grpc_chittychat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_chittychat_proto_goTypes,
		DependencyIndexes: file_grpc_chittychat_proto_depIdxs,
		MessageInfos:      file_grpc_chittychat_proto_msgTypes,
	}.Build()
	File_grpc_chittychat_proto = out.File
	file_grpc_chittychat_proto_rawDesc = nil
	file_grpc_chittychat_proto_goTypes = nil
	file_grpc_chittychat_proto_depIdxs = nil
}
