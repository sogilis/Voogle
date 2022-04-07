// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: uploaded_video.proto

package v1

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

type Uploaded_Video struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *Uploaded_Video) Reset() {
	*x = Uploaded_Video{}
	if protoimpl.UnsafeEnabled {
		mi := &file_uploaded_video_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Uploaded_Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Uploaded_Video) ProtoMessage() {}

func (x *Uploaded_Video) ProtoReflect() protoreflect.Message {
	mi := &file_uploaded_video_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Uploaded_Video.ProtoReflect.Descriptor instead.
func (*Uploaded_Video) Descriptor() ([]byte, []int) {
	return file_uploaded_video_proto_rawDescGZIP(), []int{0}
}

func (x *Uploaded_Video) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Uploaded_Video) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

var File_uploaded_video_proto protoreflect.FileDescriptor

var file_uploaded_video_proto_rawDesc = []byte{
	0x0a, 0x14, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x5f, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70, 0x6b, 0x67, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x38, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x65, 0x64, 0x5f, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x53, 0x6f, 0x67, 0x69, 0x6c, 0x69, 0x73, 0x2f, 0x56, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x73, 0x72, 0x63, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x73, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_uploaded_video_proto_rawDescOnce sync.Once
	file_uploaded_video_proto_rawDescData = file_uploaded_video_proto_rawDesc
)

func file_uploaded_video_proto_rawDescGZIP() []byte {
	file_uploaded_video_proto_rawDescOnce.Do(func() {
		file_uploaded_video_proto_rawDescData = protoimpl.X.CompressGZIP(file_uploaded_video_proto_rawDescData)
	})
	return file_uploaded_video_proto_rawDescData
}

var file_uploaded_video_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_uploaded_video_proto_goTypes = []interface{}{
	(*Uploaded_Video)(nil), // 0: pkg.contracts.v1.Uploaded_Video
}
var file_uploaded_video_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_uploaded_video_proto_init() }
func file_uploaded_video_proto_init() {
	if File_uploaded_video_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_uploaded_video_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Uploaded_Video); i {
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
			RawDescriptor: file_uploaded_video_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_uploaded_video_proto_goTypes,
		DependencyIndexes: file_uploaded_video_proto_depIdxs,
		MessageInfos:      file_uploaded_video_proto_msgTypes,
	}.Build()
	File_uploaded_video_proto = out.File
	file_uploaded_video_proto_rawDesc = nil
	file_uploaded_video_proto_goTypes = nil
	file_uploaded_video_proto_depIdxs = nil
}
