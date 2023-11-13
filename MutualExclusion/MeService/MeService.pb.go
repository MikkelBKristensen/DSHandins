// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: MeService/MeService.proto

package MeService

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

	Timestamp int64  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	NodeId    string `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_MeService_MeService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_MeService_MeService_proto_msgTypes[0]
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
	return file_MeService_MeService_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Message) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

type ConnectionMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port      string `protobuf:"bytes,1,opt,name=port,proto3" json:"port,omitempty"`
	IsJoin    bool   `protobuf:"varint,2,opt,name=isJoin,proto3" json:"isJoin,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ConnectionMsg) Reset() {
	*x = ConnectionMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_MeService_MeService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionMsg) ProtoMessage() {}

func (x *ConnectionMsg) ProtoReflect() protoreflect.Message {
	mi := &file_MeService_MeService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionMsg.ProtoReflect.Descriptor instead.
func (*ConnectionMsg) Descriptor() ([]byte, []int) {
	return file_MeService_MeService_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectionMsg) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *ConnectionMsg) GetIsJoin() bool {
	if x != nil {
		return x.IsJoin
	}
	return false
}

func (x *ConnectionMsg) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_MeService_MeService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_MeService_MeService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_MeService_MeService_proto_rawDescGZIP(), []int{2}
}

var File_MeService_MeService_proto protoreflect.FileDescriptor

var file_MeService_MeService_proto_rawDesc = []byte{
	0x0a, 0x19, 0x4d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x4d, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x4d, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x3f, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x59, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x69, 0x73, 0x4a, 0x6f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73,
	0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x85, 0x01, 0x0a, 0x09,
	0x4d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x10, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x2e,
	0x4d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x1a, 0x12, 0x2e, 0x4d, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x0c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x2e, 0x4d, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x12, 0x2e, 0x4d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x4d, 0x69, 0x6b, 0x6b, 0x65, 0x6c, 0x42, 0x4b, 0x72, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x73, 0x65, 0x6e, 0x2f, 0x44, 0x53, 0x48, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x73, 0x2f, 0x54, 0x72,
	0x65, 0x65, 0x2f, 0x4d, 0x61, 0x69, 0x6e, 0x2f, 0x4d, 0x75, 0x74, 0x75, 0x61, 0x6c, 0x45, 0x78,
	0x63, 0x6c, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x4d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_MeService_MeService_proto_rawDescOnce sync.Once
	file_MeService_MeService_proto_rawDescData = file_MeService_MeService_proto_rawDesc
)

func file_MeService_MeService_proto_rawDescGZIP() []byte {
	file_MeService_MeService_proto_rawDescOnce.Do(func() {
		file_MeService_MeService_proto_rawDescData = protoimpl.X.CompressGZIP(file_MeService_MeService_proto_rawDescData)
	})
	return file_MeService_MeService_proto_rawDescData
}

var file_MeService_MeService_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_MeService_MeService_proto_goTypes = []interface{}{
	(*Message)(nil),       // 0: MeService.Message
	(*ConnectionMsg)(nil), // 1: MeService.ConnectionMsg
	(*Empty)(nil),         // 2: MeService.Empty
}
var file_MeService_MeService_proto_depIdxs = []int32{
	1, // 0: MeService.MeService.ConnectionStatus:input_type -> MeService.ConnectionMsg
	0, // 1: MeService.MeService.RequestEntry:input_type -> MeService.Message
	0, // 2: MeService.MeService.ConnectionStatus:output_type -> MeService.Message
	0, // 3: MeService.MeService.RequestEntry:output_type -> MeService.Message
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_MeService_MeService_proto_init() }
func file_MeService_MeService_proto_init() {
	if File_MeService_MeService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_MeService_MeService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_MeService_MeService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionMsg); i {
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
		file_MeService_MeService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_MeService_MeService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_MeService_MeService_proto_goTypes,
		DependencyIndexes: file_MeService_MeService_proto_depIdxs,
		MessageInfos:      file_MeService_MeService_proto_msgTypes,
	}.Build()
	File_MeService_MeService_proto = out.File
	file_MeService_MeService_proto_rawDesc = nil
	file_MeService_MeService_proto_goTypes = nil
	file_MeService_MeService_proto_depIdxs = nil
}
