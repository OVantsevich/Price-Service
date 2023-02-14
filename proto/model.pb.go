// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/model.proto

package __

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

type GetPricesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetPricesRequest) Reset() {
	*x = GetPricesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPricesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPricesRequest) ProtoMessage() {}

func (x *GetPricesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPricesRequest.ProtoReflect.Descriptor instead.
func (*GetPricesRequest) Descriptor() ([]byte, []int) {
	return file_proto_model_proto_rawDescGZIP(), []int{0}
}

type GetPricesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prices map[string]float32 `protobuf:"bytes,1,rep,name=prices,proto3" json:"prices,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
}

func (x *GetPricesResponse) Reset() {
	*x = GetPricesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPricesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPricesResponse) ProtoMessage() {}

func (x *GetPricesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPricesResponse.ProtoReflect.Descriptor instead.
func (*GetPricesResponse) Descriptor() ([]byte, []int) {
	return file_proto_model_proto_rawDescGZIP(), []int{1}
}

func (x *GetPricesResponse) GetPrices() map[string]float32 {
	if x != nil {
		return x.Prices
	}
	return nil
}

var File_proto_model_proto protoreflect.FileDescriptor

var file_proto_model_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x86, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a,
	0x06, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x32, 0x44, 0x0a, 0x0c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x34, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x12, 0x11, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_model_proto_rawDescOnce sync.Once
	file_proto_model_proto_rawDescData = file_proto_model_proto_rawDesc
)

func file_proto_model_proto_rawDescGZIP() []byte {
	file_proto_model_proto_rawDescOnce.Do(func() {
		file_proto_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_model_proto_rawDescData)
	})
	return file_proto_model_proto_rawDescData
}

var file_proto_model_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_model_proto_goTypes = []interface{}{
	(*GetPricesRequest)(nil),  // 0: GetPricesRequest
	(*GetPricesResponse)(nil), // 1: GetPricesResponse
	nil,                       // 2: GetPricesResponse.PricesEntry
}
var file_proto_model_proto_depIdxs = []int32{
	2, // 0: GetPricesResponse.prices:type_name -> GetPricesResponse.PricesEntry
	0, // 1: PriceService.GetPrices:input_type -> GetPricesRequest
	1, // 2: PriceService.GetPrices:output_type -> GetPricesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_model_proto_init() }
func file_proto_model_proto_init() {
	if File_proto_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPricesRequest); i {
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
		file_proto_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPricesResponse); i {
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
			RawDescriptor: file_proto_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_model_proto_goTypes,
		DependencyIndexes: file_proto_model_proto_depIdxs,
		MessageInfos:      file_proto_model_proto_msgTypes,
	}.Build()
	File_proto_model_proto = out.File
	file_proto_model_proto_rawDesc = nil
	file_proto_model_proto_goTypes = nil
	file_proto_model_proto_depIdxs = nil
}
