// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package ibcmiddlewarev1

import (
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: ibcmiddleware/v1/tx.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_ibcmiddleware_v1_tx_proto protoreflect.FileDescriptor

var file_ibcmiddleware_v1_tx_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x62, 0x63, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x69, 0x62, 0x63,
	0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x67,
	0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x42, 0xc3, 0x01, 0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x62, 0x63, 0x6d,
	0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x07, 0x54, 0x78,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x6c, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x2f, 0x73,
	0x69, 0x6d, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x62, 0x63, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x62, 0x63, 0x6d, 0x69,
	0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x58, 0x58,
	0xaa, 0x02, 0x10, 0x49, 0x62, 0x63, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x10, 0x49, 0x62, 0x63, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1c, 0x49, 0x62, 0x63, 0x6d, 0x69, 0x64, 0x64,
	0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x11, 0x49, 0x62, 0x63, 0x6d, 0x69, 0x64, 0x64, 0x6c,
	0x65, 0x77, 0x61, 0x72, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_ibcmiddleware_v1_tx_proto_goTypes = []interface{}{}
var file_ibcmiddleware_v1_tx_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ibcmiddleware_v1_tx_proto_init() }
func file_ibcmiddleware_v1_tx_proto_init() {
	if File_ibcmiddleware_v1_tx_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ibcmiddleware_v1_tx_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ibcmiddleware_v1_tx_proto_goTypes,
		DependencyIndexes: file_ibcmiddleware_v1_tx_proto_depIdxs,
	}.Build()
	File_ibcmiddleware_v1_tx_proto = out.File
	file_ibcmiddleware_v1_tx_proto_rawDesc = nil
	file_ibcmiddleware_v1_tx_proto_goTypes = nil
	file_ibcmiddleware_v1_tx_proto_depIdxs = nil
}
