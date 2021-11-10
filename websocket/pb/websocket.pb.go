// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: websocket.proto

package pb

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

// Method
type Method int32

const (
	// UNKNOWN
	Method_UNKNOWN Method = 0
	// PING 心跳
	Method_PING Method = 1
	// SUBSCRIBE 消息订阅
	Method_SUBSCRIBE Method = 2
	// UNSUBSCRIBE 取消订阅
	Method_UNSUBSCRIBE Method = 3
	// BROADCAST 广播
	Method_BROADCAST Method = 4
	// DISCONNECT 断开
	Method_DISCONNECT Method = 5
)

// Enum value maps for Method.
var (
	Method_name = map[int32]string{
		0: "UNKNOWN",
		1: "PING",
		2: "SUBSCRIBE",
		3: "UNSUBSCRIBE",
		4: "BROADCAST",
		5: "DISCONNECT",
	}
	Method_value = map[string]int32{
		"UNKNOWN":     0,
		"PING":        1,
		"SUBSCRIBE":   2,
		"UNSUBSCRIBE": 3,
		"BROADCAST":   4,
		"DISCONNECT":  5,
	}
)

func (x Method) Enum() *Method {
	p := new(Method)
	*p = x
	return p
}

func (x Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Method) Descriptor() protoreflect.EnumDescriptor {
	return file_websocket_proto_enumTypes[0].Descriptor()
}

func (Method) Type() protoreflect.EnumType {
	return &file_websocket_proto_enumTypes[0]
}

func (x Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Method.Descriptor instead.
func (Method) EnumDescriptor() ([]byte, []int) {
	return file_websocket_proto_rawDescGZIP(), []int{0}
}

var File_websocket_proto protoreflect.FileDescriptor

var file_websocket_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x77, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x2a, 0x5e, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x50, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52,
	0x49, 0x42, 0x45, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x55, 0x42, 0x53, 0x43,
	0x52, 0x49, 0x42, 0x45, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43,
	0x41, 0x53, 0x54, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e,
	0x45, 0x43, 0x54, 0x10, 0x05, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_websocket_proto_rawDescOnce sync.Once
	file_websocket_proto_rawDescData = file_websocket_proto_rawDesc
)

func file_websocket_proto_rawDescGZIP() []byte {
	file_websocket_proto_rawDescOnce.Do(func() {
		file_websocket_proto_rawDescData = protoimpl.X.CompressGZIP(file_websocket_proto_rawDescData)
	})
	return file_websocket_proto_rawDescData
}

var file_websocket_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_websocket_proto_goTypes = []interface{}{
	(Method)(0), // 0: pb.Method
}
var file_websocket_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_websocket_proto_init() }
func file_websocket_proto_init() {
	if File_websocket_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_websocket_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_websocket_proto_goTypes,
		DependencyIndexes: file_websocket_proto_depIdxs,
		EnumInfos:         file_websocket_proto_enumTypes,
	}.Build()
	File_websocket_proto = out.File
	file_websocket_proto_rawDesc = nil
	file_websocket_proto_goTypes = nil
	file_websocket_proto_depIdxs = nil
}