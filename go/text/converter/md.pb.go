// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.1
// source: proto/text/converter/md.proto

package converter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_text_converter_md_proto protoreflect.FileDescriptor

var file_proto_text_converter_md_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x63, 0x6f, 0x6e,
	0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x6d, 0x6f, 0x6e, 0x6f, 0x2e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65,
	0x72, 0x74, 0x65, 0x72, 0x1a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x78, 0x74,
	0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x6d, 0x0a, 0x15, 0x4d, 0x61, 0x72, 0x6b, 0x64,
	0x6f, 0x77, 0x6e, 0x54, 0x65, 0x78, 0x74, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72,
	0x12, 0x54, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x12, 0x23, 0x2e, 0x6d, 0x6f,
	0x6e, 0x6f, 0x2e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74,
	0x1a, 0x24, 0x2e, 0x6d, 0x6f, 0x6e, 0x6f, 0x2e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x63, 0x6f, 0x6e,
	0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x6f, 0x2e, 0x66, 0x31, 0x31,
	0x30, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65,
	0x78, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_text_converter_md_proto_goTypes = []interface{}{
	(*RequestConvert)(nil),  // 0: mono.text.converter.RequestConvert
	(*ResponseConvert)(nil), // 1: mono.text.converter.ResponseConvert
}
var file_proto_text_converter_md_proto_depIdxs = []int32{
	0, // 0: mono.text.converter.MarkdownTextConverter.Convert:input_type -> mono.text.converter.RequestConvert
	1, // 1: mono.text.converter.MarkdownTextConverter.Convert:output_type -> mono.text.converter.ResponseConvert
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_text_converter_md_proto_init() }
func file_proto_text_converter_md_proto_init() {
	if File_proto_text_converter_md_proto != nil {
		return
	}
	file_proto_text_converter_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_text_converter_md_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_text_converter_md_proto_goTypes,
		DependencyIndexes: file_proto_text_converter_md_proto_depIdxs,
	}.Build()
	File_proto_text_converter_md_proto = out.File
	file_proto_text_converter_md_proto_rawDesc = nil
	file_proto_text_converter_md_proto_goTypes = nil
	file_proto_text_converter_md_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MarkdownTextConverterClient is the client API for MarkdownTextConverter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MarkdownTextConverterClient interface {
	Convert(ctx context.Context, in *RequestConvert, opts ...grpc.CallOption) (*ResponseConvert, error)
}

type markdownTextConverterClient struct {
	cc grpc.ClientConnInterface
}

func NewMarkdownTextConverterClient(cc grpc.ClientConnInterface) MarkdownTextConverterClient {
	return &markdownTextConverterClient{cc}
}

func (c *markdownTextConverterClient) Convert(ctx context.Context, in *RequestConvert, opts ...grpc.CallOption) (*ResponseConvert, error) {
	out := new(ResponseConvert)
	err := c.cc.Invoke(ctx, "/mono.text.converter.MarkdownTextConverter/Convert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MarkdownTextConverterServer is the server API for MarkdownTextConverter service.
type MarkdownTextConverterServer interface {
	Convert(context.Context, *RequestConvert) (*ResponseConvert, error)
}

// UnimplementedMarkdownTextConverterServer can be embedded to have forward compatible implementations.
type UnimplementedMarkdownTextConverterServer struct {
}

func (*UnimplementedMarkdownTextConverterServer) Convert(context.Context, *RequestConvert) (*ResponseConvert, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Convert not implemented")
}

func RegisterMarkdownTextConverterServer(s *grpc.Server, srv MarkdownTextConverterServer) {
	s.RegisterService(&_MarkdownTextConverter_serviceDesc, srv)
}

func _MarkdownTextConverter_Convert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestConvert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarkdownTextConverterServer).Convert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mono.text.converter.MarkdownTextConverter/Convert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarkdownTextConverterServer).Convert(ctx, req.(*RequestConvert))
	}
	return interceptor(ctx, in, info, handler)
}

var _MarkdownTextConverter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mono.text.converter.MarkdownTextConverter",
	HandlerType: (*MarkdownTextConverterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Convert",
			Handler:    _MarkdownTextConverter_Convert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/text/converter/md.proto",
}