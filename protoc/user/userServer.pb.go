// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v5.27.1
// source: userServer.proto

package user

import (
	_ "github.com/gogo/protobuf/gogoproto"
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

type RequestHello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *RequestHello) Reset() {
	*x = RequestHello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userServer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestHello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestHello) ProtoMessage() {}

func (x *RequestHello) ProtoReflect() protoreflect.Message {
	mi := &file_userServer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestHello.ProtoReflect.Descriptor instead.
func (*RequestHello) Descriptor() ([]byte, []int) {
	return file_userServer_proto_rawDescGZIP(), []int{0}
}

func (x *RequestHello) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ResponseHello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *ResponseHello) Reset() {
	*x = ResponseHello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userServer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseHello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseHello) ProtoMessage() {}

func (x *ResponseHello) ProtoReflect() protoreflect.Message {
	mi := &file_userServer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseHello.ProtoReflect.Descriptor instead.
func (*ResponseHello) Descriptor() ([]byte, []int) {
	return file_userServer_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseHello) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

var File_userServer_proto protoreflect.FileDescriptor

var file_userServer_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xea, 0xde, 0x1f, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0xf2,
	0xde, 0x1f, 0x0b, 0x66, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x21, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x32, 0x3b, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x33, 0x0a, 0x08, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x12, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x1a,
	0x13, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x42, 0x07, 0x5a, 0x05, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_userServer_proto_rawDescOnce sync.Once
	file_userServer_proto_rawDescData = file_userServer_proto_rawDesc
)

func file_userServer_proto_rawDescGZIP() []byte {
	file_userServer_proto_rawDescOnce.Do(func() {
		file_userServer_proto_rawDescData = protoimpl.X.CompressGZIP(file_userServer_proto_rawDescData)
	})
	return file_userServer_proto_rawDescData
}

var file_userServer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_userServer_proto_goTypes = []interface{}{
	(*RequestHello)(nil),  // 0: user.RequestHello
	(*ResponseHello)(nil), // 1: user.ResponseHello
}
var file_userServer_proto_depIdxs = []int32{
	0, // 0: user.User.SayHello:input_type -> user.RequestHello
	1, // 1: user.User.SayHello:output_type -> user.ResponseHello
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_userServer_proto_init() }
func file_userServer_proto_init() {
	if File_userServer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_userServer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestHello); i {
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
		file_userServer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseHello); i {
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
			RawDescriptor: file_userServer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userServer_proto_goTypes,
		DependencyIndexes: file_userServer_proto_depIdxs,
		MessageInfos:      file_userServer_proto_msgTypes,
	}.Build()
	File_userServer_proto = out.File
	file_userServer_proto_rawDesc = nil
	file_userServer_proto_goTypes = nil
	file_userServer_proto_depIdxs = nil
}
