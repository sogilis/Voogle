// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: video.proto

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

type Video_VideoStatus int32

const (
	Video_VIDEO_STATUS_UNSPECIFIED Video_VideoStatus = 0
	Video_VIDEO_STATUS_UPLOADING   Video_VideoStatus = 1
	Video_VIDEO_STATUS_UPLOADED    Video_VideoStatus = 2
	Video_VIDEO_STATUS_ENCODING    Video_VideoStatus = 3
	Video_VIDEO_STATUS_COMPLETE    Video_VideoStatus = 4
	Video_VIDEO_STATUS_UNKNOWN     Video_VideoStatus = 5
	Video_VIDEO_STATUS_FAIL_UPLOAD Video_VideoStatus = 6
	Video_VIDEO_STATUS_FAIL_ENCODE Video_VideoStatus = 7
)

// Enum value maps for Video_VideoStatus.
var (
	Video_VideoStatus_name = map[int32]string{
		0: "VIDEO_STATUS_UNSPECIFIED",
		1: "VIDEO_STATUS_UPLOADING",
		2: "VIDEO_STATUS_UPLOADED",
		3: "VIDEO_STATUS_ENCODING",
		4: "VIDEO_STATUS_COMPLETE",
		5: "VIDEO_STATUS_UNKNOWN",
		6: "VIDEO_STATUS_FAIL_UPLOAD",
		7: "VIDEO_STATUS_FAIL_ENCODE",
	}
	Video_VideoStatus_value = map[string]int32{
		"VIDEO_STATUS_UNSPECIFIED": 0,
		"VIDEO_STATUS_UPLOADING":   1,
		"VIDEO_STATUS_UPLOADED":    2,
		"VIDEO_STATUS_ENCODING":    3,
		"VIDEO_STATUS_COMPLETE":    4,
		"VIDEO_STATUS_UNKNOWN":     5,
		"VIDEO_STATUS_FAIL_UPLOAD": 6,
		"VIDEO_STATUS_FAIL_ENCODE": 7,
	}
)

func (x Video_VideoStatus) Enum() *Video_VideoStatus {
	p := new(Video_VideoStatus)
	*p = x
	return p
}

func (x Video_VideoStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Video_VideoStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_video_proto_enumTypes[0].Descriptor()
}

func (Video_VideoStatus) Type() protoreflect.EnumType {
	return &file_video_proto_enumTypes[0]
}

func (x Video_VideoStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Video_VideoStatus.Descriptor instead.
func (Video_VideoStatus) EnumDescriptor() ([]byte, []int) {
	return file_video_proto_rawDescGZIP(), []int{0, 0}
}

type Video struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status Video_VideoStatus `protobuf:"varint,2,opt,name=status,proto3,enum=pkg.contracts.v1.Video_VideoStatus" json:"status,omitempty"`
	Source string            `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *Video) Reset() {
	*x = Video{}
	if protoimpl.UnsafeEnabled {
		mi := &file_video_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Video) ProtoMessage() {}

func (x *Video) ProtoReflect() protoreflect.Message {
	mi := &file_video_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Video.ProtoReflect.Descriptor instead.
func (*Video) Descriptor() ([]byte, []int) {
	return file_video_proto_rawDescGZIP(), []int{0}
}

func (x *Video) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Video) GetStatus() Video_VideoStatus {
	if x != nil {
		return x.Status
	}
	return Video_VIDEO_STATUS_UNSPECIFIED
}

func (x *Video) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

var File_video_proto protoreflect.FileDescriptor

var file_video_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70,
	0x6b, 0x67, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x22,
	0xdd, 0x02, 0x0a, 0x05, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3b, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x70, 0x6b, 0x67, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x69, 0x64,
	0x65, 0x6f, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0xee,
	0x01, 0x0a, 0x0b, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c,
	0x0a, 0x18, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16,
	0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x50, 0x4c,
	0x4f, 0x41, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x56, 0x49, 0x44, 0x45,
	0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x50, 0x4c, 0x4f, 0x41, 0x44, 0x45,
	0x44, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x45, 0x4e, 0x43, 0x4f, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x19,
	0x0a, 0x15, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43,
	0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x04, 0x12, 0x18, 0x0a, 0x14, 0x56, 0x49, 0x44,
	0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x05, 0x12, 0x1c, 0x0a, 0x18, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x5f, 0x55, 0x50, 0x4c, 0x4f, 0x41, 0x44, 0x10,
	0x06, 0x12, 0x1c, 0x0a, 0x18, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x5f, 0x45, 0x4e, 0x43, 0x4f, 0x44, 0x45, 0x10, 0x07, 0x42,
	0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x6f,
	0x67, 0x69, 0x6c, 0x69, 0x73, 0x2f, 0x56, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x73, 0x72, 0x63,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_video_proto_rawDescOnce sync.Once
	file_video_proto_rawDescData = file_video_proto_rawDesc
)

func file_video_proto_rawDescGZIP() []byte {
	file_video_proto_rawDescOnce.Do(func() {
		file_video_proto_rawDescData = protoimpl.X.CompressGZIP(file_video_proto_rawDescData)
	})
	return file_video_proto_rawDescData
}

var file_video_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_video_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_video_proto_goTypes = []interface{}{
	(Video_VideoStatus)(0), // 0: pkg.contracts.v1.Video.VideoStatus
	(*Video)(nil),          // 1: pkg.contracts.v1.Video
}
var file_video_proto_depIdxs = []int32{
	0, // 0: pkg.contracts.v1.Video.status:type_name -> pkg.contracts.v1.Video.VideoStatus
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_video_proto_init() }
func file_video_proto_init() {
	if File_video_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_video_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Video); i {
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
			RawDescriptor: file_video_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_video_proto_goTypes,
		DependencyIndexes: file_video_proto_depIdxs,
		EnumInfos:         file_video_proto_enumTypes,
		MessageInfos:      file_video_proto_msgTypes,
	}.Build()
	File_video_proto = out.File
	file_video_proto_rawDesc = nil
	file_video_proto_goTypes = nil
	file_video_proto_depIdxs = nil
}
