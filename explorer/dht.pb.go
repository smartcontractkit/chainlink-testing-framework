package explorer

// original pb https://github.com/smartcontractkit/libocr/blob/master/networking/dht-router/serialization/cl_dht_addr_announcement.pb.go

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SignedAnnouncement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addrs      [][]byte `protobuf:"bytes,1,rep,name=addrs,proto3" json:"addrs,omitempty"`
	UserPrefix uint32   `protobuf:"varint,2,opt,name=userPrefix,proto3" json:"userPrefix,omitempty"`
	Counter    uint64   `protobuf:"varint,3,opt,name=counter,proto3" json:"counter,omitempty"`
	PublicKey  []byte   `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Sig        []byte   `protobuf:"bytes,5,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *SignedAnnouncement) Reset() {
	*x = SignedAnnouncement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_serialization_cl_dht_addr_announcement_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedAnnouncement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedAnnouncement) ProtoMessage() {}

func (x *SignedAnnouncement) ProtoReflect() protoreflect.Message {
	mi := &file_serialization_cl_dht_addr_announcement_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedAnnouncement.ProtoReflect.Descriptor instead.
func (*SignedAnnouncement) Descriptor() ([]byte, []int) {
	return file_serialization_cl_dht_addr_announcement_proto_rawDescGZIP(), []int{0}
}

func (x *SignedAnnouncement) GetAddrs() [][]byte {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *SignedAnnouncement) GetUserPrefix() uint32 {
	if x != nil {
		return x.UserPrefix
	}
	return 0
}

func (x *SignedAnnouncement) GetCounter() uint64 {
	if x != nil {
		return x.Counter
	}
	return 0
}

func (x *SignedAnnouncement) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *SignedAnnouncement) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

var File_serialization_cl_dht_addr_announcement_proto protoreflect.FileDescriptor

var file_serialization_cl_dht_addr_announcement_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x63, 0x6c, 0x5f, 0x64, 0x68, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x5f, 0x61, 0x6e, 0x6e, 0x6f,
	0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x64, 0x68, 0x74, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x22, 0x95, 0x01, 0x0a, 0x12, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52,
	0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x50, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72,
	0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69,
	0x67, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x3b, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_serialization_cl_dht_addr_announcement_proto_rawDescOnce sync.Once
	file_serialization_cl_dht_addr_announcement_proto_rawDescData = file_serialization_cl_dht_addr_announcement_proto_rawDesc
)

func file_serialization_cl_dht_addr_announcement_proto_rawDescGZIP() []byte {
	file_serialization_cl_dht_addr_announcement_proto_rawDescOnce.Do(func() {
		file_serialization_cl_dht_addr_announcement_proto_rawDescData = protoimpl.X.CompressGZIP(file_serialization_cl_dht_addr_announcement_proto_rawDescData)
	})
	return file_serialization_cl_dht_addr_announcement_proto_rawDescData
}

var file_serialization_cl_dht_addr_announcement_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_serialization_cl_dht_addr_announcement_proto_goTypes = []interface{}{
	(*SignedAnnouncement)(nil), // 0: dhtrouter.SignedAnnouncement
}
var file_serialization_cl_dht_addr_announcement_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_serialization_cl_dht_addr_announcement_proto_init() }
func file_serialization_cl_dht_addr_announcement_proto_init() {
	if File_serialization_cl_dht_addr_announcement_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_serialization_cl_dht_addr_announcement_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedAnnouncement); i {
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
			RawDescriptor: file_serialization_cl_dht_addr_announcement_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_serialization_cl_dht_addr_announcement_proto_goTypes,
		DependencyIndexes: file_serialization_cl_dht_addr_announcement_proto_depIdxs,
		MessageInfos:      file_serialization_cl_dht_addr_announcement_proto_msgTypes,
	}.Build()
	File_serialization_cl_dht_addr_announcement_proto = out.File
	file_serialization_cl_dht_addr_announcement_proto_rawDesc = nil
	file_serialization_cl_dht_addr_announcement_proto_goTypes = nil
	file_serialization_cl_dht_addr_announcement_proto_depIdxs = nil
}
