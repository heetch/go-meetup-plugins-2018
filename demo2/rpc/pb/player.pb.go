// Code generated by protoc-gen-go.
// source: player.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	player.proto
	codec.proto

It has these top-level messages:
	RegisterCodecRequest
	RegisterCodecResponse
	AudioFileMetadataRequest
	AudioFileMetadataResponse
	DecodeAudioFileRequest
	DecodeAudioFileResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RegisterCodecRequest struct {
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	Ext  string `protobuf:"bytes,2,opt,name=ext" json:"ext,omitempty"`
}

func (m *RegisterCodecRequest) Reset()                    { *m = RegisterCodecRequest{} }
func (m *RegisterCodecRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterCodecRequest) ProtoMessage()               {}
func (*RegisterCodecRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RegisterCodecResponse struct {
}

func (m *RegisterCodecResponse) Reset()                    { *m = RegisterCodecResponse{} }
func (m *RegisterCodecResponse) String() string            { return proto.CompactTextString(m) }
func (*RegisterCodecResponse) ProtoMessage()               {}
func (*RegisterCodecResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*RegisterCodecRequest)(nil), "pb.RegisterCodecRequest")
	proto.RegisterType((*RegisterCodecResponse)(nil), "pb.RegisterCodecResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Player service

type PlayerClient interface {
	RegisterCodec(ctx context.Context, in *RegisterCodecRequest, opts ...grpc.CallOption) (*RegisterCodecResponse, error)
}

type playerClient struct {
	cc *grpc.ClientConn
}

func NewPlayerClient(cc *grpc.ClientConn) PlayerClient {
	return &playerClient{cc}
}

func (c *playerClient) RegisterCodec(ctx context.Context, in *RegisterCodecRequest, opts ...grpc.CallOption) (*RegisterCodecResponse, error) {
	out := new(RegisterCodecResponse)
	err := grpc.Invoke(ctx, "/pb.Player/RegisterCodec", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Player service

type PlayerServer interface {
	RegisterCodec(context.Context, *RegisterCodecRequest) (*RegisterCodecResponse, error)
}

func RegisterPlayerServer(s *grpc.Server, srv PlayerServer) {
	s.RegisterService(&_Player_serviceDesc, srv)
}

func _Player_RegisterCodec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterCodecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).RegisterCodec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Player/RegisterCodec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).RegisterCodec(ctx, req.(*RegisterCodecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Player_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Player",
	HandlerType: (*PlayerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterCodec",
			Handler:    _Player_RegisterCodec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("player.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xc8, 0x49, 0xac,
	0x4c, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xb2, 0xe1, 0x12,
	0x09, 0x4a, 0x4d, 0xcf, 0x2c, 0x2e, 0x49, 0x2d, 0x72, 0xce, 0x4f, 0x49, 0x4d, 0x0e, 0x4a, 0x2d,
	0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe2, 0x62, 0x49, 0x4c, 0x49, 0x29, 0x92, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85, 0x04, 0xb8, 0x98, 0x53, 0x2b, 0x4a, 0x24, 0x98, 0xc0, 0x42,
	0x20, 0xa6, 0x92, 0x38, 0x97, 0x28, 0x9a, 0xee, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0xa3, 0x00,
	0x2e, 0xb6, 0x00, 0xb0, 0x55, 0x42, 0x6e, 0x5c, 0xbc, 0x28, 0x4a, 0x84, 0x24, 0xf4, 0x0a, 0x92,
	0xf4, 0xb0, 0xd9, 0x29, 0x25, 0x89, 0x45, 0x06, 0x62, 0x9e, 0x12, 0x43, 0x12, 0x1b, 0xd8, 0xcd,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0x85, 0xc7, 0x96, 0xc3, 0x00, 0x00, 0x00,
}
