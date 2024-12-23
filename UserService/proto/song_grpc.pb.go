// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.3
// source: proto/song.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SongServiceClient is the client API for SongService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SongServiceClient interface {
	GetSongById(ctx context.Context, in *SongRequest, opts ...grpc.CallOption) (*SongResponse, error)
}

type songServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSongServiceClient(cc grpc.ClientConnInterface) SongServiceClient {
	return &songServiceClient{cc}
}

func (c *songServiceClient) GetSongById(ctx context.Context, in *SongRequest, opts ...grpc.CallOption) (*SongResponse, error) {
	out := new(SongResponse)
	err := c.cc.Invoke(ctx, "/song.SongService/GetSongById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SongServiceServer is the server API for SongService service.
// All implementations must embed UnimplementedSongServiceServer
// for forward compatibility
type SongServiceServer interface {
	GetSongById(context.Context, *SongRequest) (*SongResponse, error)
	mustEmbedUnimplementedSongServiceServer()
}

// UnimplementedSongServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSongServiceServer struct {
}

func (UnimplementedSongServiceServer) GetSongById(context.Context, *SongRequest) (*SongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSongById not implemented")
}
func (UnimplementedSongServiceServer) mustEmbedUnimplementedSongServiceServer() {}

// UnsafeSongServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SongServiceServer will
// result in compilation errors.
type UnsafeSongServiceServer interface {
	mustEmbedUnimplementedSongServiceServer()
}

func RegisterSongServiceServer(s grpc.ServiceRegistrar, srv SongServiceServer) {
	s.RegisterService(&SongService_ServiceDesc, srv)
}

func _SongService_GetSongById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SongRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SongServiceServer).GetSongById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/song.SongService/GetSongById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SongServiceServer).GetSongById(ctx, req.(*SongRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SongService_ServiceDesc is the grpc.ServiceDesc for SongService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SongService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "song.SongService",
	HandlerType: (*SongServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSongById",
			Handler:    _SongService_GetSongById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/song.proto",
}

// ArtistServiceClient is the client API for ArtistService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArtistServiceClient interface {
	GetArtistById(ctx context.Context, in *ArtistRequest, opts ...grpc.CallOption) (*ArtistResponse, error)
}

type artistServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArtistServiceClient(cc grpc.ClientConnInterface) ArtistServiceClient {
	return &artistServiceClient{cc}
}

func (c *artistServiceClient) GetArtistById(ctx context.Context, in *ArtistRequest, opts ...grpc.CallOption) (*ArtistResponse, error) {
	out := new(ArtistResponse)
	err := c.cc.Invoke(ctx, "/song.ArtistService/GetArtistById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArtistServiceServer is the server API for ArtistService service.
// All implementations must embed UnimplementedArtistServiceServer
// for forward compatibility
type ArtistServiceServer interface {
	GetArtistById(context.Context, *ArtistRequest) (*ArtistResponse, error)
	mustEmbedUnimplementedArtistServiceServer()
}

// UnimplementedArtistServiceServer must be embedded to have forward compatible implementations.
type UnimplementedArtistServiceServer struct {
}

func (UnimplementedArtistServiceServer) GetArtistById(context.Context, *ArtistRequest) (*ArtistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArtistById not implemented")
}
func (UnimplementedArtistServiceServer) mustEmbedUnimplementedArtistServiceServer() {}

// UnsafeArtistServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArtistServiceServer will
// result in compilation errors.
type UnsafeArtistServiceServer interface {
	mustEmbedUnimplementedArtistServiceServer()
}

func RegisterArtistServiceServer(s grpc.ServiceRegistrar, srv ArtistServiceServer) {
	s.RegisterService(&ArtistService_ServiceDesc, srv)
}

func _ArtistService_GetArtistById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArtistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArtistServiceServer).GetArtistById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/song.ArtistService/GetArtistById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArtistServiceServer).GetArtistById(ctx, req.(*ArtistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ArtistService_ServiceDesc is the grpc.ServiceDesc for ArtistService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArtistService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "song.ArtistService",
	HandlerType: (*ArtistServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArtistById",
			Handler:    _ArtistService_GetArtistById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/song.proto",
}

// AlbumServiceClient is the client API for AlbumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlbumServiceClient interface {
	GetAlbumById(ctx context.Context, in *AlbumRequest, opts ...grpc.CallOption) (*AlbumResponse, error)
}

type albumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAlbumServiceClient(cc grpc.ClientConnInterface) AlbumServiceClient {
	return &albumServiceClient{cc}
}

func (c *albumServiceClient) GetAlbumById(ctx context.Context, in *AlbumRequest, opts ...grpc.CallOption) (*AlbumResponse, error) {
	out := new(AlbumResponse)
	err := c.cc.Invoke(ctx, "/song.AlbumService/GetAlbumById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlbumServiceServer is the server API for AlbumService service.
// All implementations must embed UnimplementedAlbumServiceServer
// for forward compatibility
type AlbumServiceServer interface {
	GetAlbumById(context.Context, *AlbumRequest) (*AlbumResponse, error)
	mustEmbedUnimplementedAlbumServiceServer()
}

// UnimplementedAlbumServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAlbumServiceServer struct {
}

func (UnimplementedAlbumServiceServer) GetAlbumById(context.Context, *AlbumRequest) (*AlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbumById not implemented")
}
func (UnimplementedAlbumServiceServer) mustEmbedUnimplementedAlbumServiceServer() {}

// UnsafeAlbumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlbumServiceServer will
// result in compilation errors.
type UnsafeAlbumServiceServer interface {
	mustEmbedUnimplementedAlbumServiceServer()
}

func RegisterAlbumServiceServer(s grpc.ServiceRegistrar, srv AlbumServiceServer) {
	s.RegisterService(&AlbumService_ServiceDesc, srv)
}

func _AlbumService_GetAlbumById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).GetAlbumById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/song.AlbumService/GetAlbumById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).GetAlbumById(ctx, req.(*AlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AlbumService_ServiceDesc is the grpc.ServiceDesc for AlbumService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlbumService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "song.AlbumService",
	HandlerType: (*AlbumServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAlbumById",
			Handler:    _AlbumService_GetAlbumById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/song.proto",
}