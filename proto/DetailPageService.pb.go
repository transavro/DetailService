// Code generated by protoc-gen-go. DO NOT EDIT.
// source: DetailPageService.proto

package DetailPageService

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	proto1 "github.com/transavro/ScheduleService/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ButtonType int32

const (
	ButtonType_RemotApi ButtonType = 0
	ButtonType_LocalApi ButtonType = 1
)

var ButtonType_name = map[int32]string{
	0: "RemotApi",
	1: "LocalApi",
}

var ButtonType_value = map[string]int32{
	"RemotApi": 0,
	"LocalApi": 1,
}

func (x ButtonType) String() string {
	return proto.EnumName(ButtonType_name, int32(x))
}

func (ButtonType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2aea33dfb16e7225, []int{0}
}

type TileInfoRequest struct {
	TileId               string   `protobuf:"bytes,1,opt,name=tileId,proto3" json:"tileId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TileInfoRequest) Reset()         { *m = TileInfoRequest{} }
func (m *TileInfoRequest) String() string { return proto.CompactTextString(m) }
func (*TileInfoRequest) ProtoMessage()    {}
func (*TileInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2aea33dfb16e7225, []int{0}
}

func (m *TileInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TileInfoRequest.Unmarshal(m, b)
}
func (m *TileInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TileInfoRequest.Marshal(b, m, deterministic)
}
func (m *TileInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TileInfoRequest.Merge(m, src)
}
func (m *TileInfoRequest) XXX_Size() int {
	return xxx_messageInfo_TileInfoRequest.Size(m)
}
func (m *TileInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TileInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TileInfoRequest proto.InternalMessageInfo

func (m *TileInfoRequest) GetTileId() string {
	if m != nil {
		return m.TileId
	}
	return ""
}

type DetailTileInfo struct {
	Title                string            `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Synopsis             string            `protobuf:"bytes,2,opt,name=synopsis,proto3" json:"synopsis,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Button               []*Button         `protobuf:"bytes,4,rep,name=button,proto3" json:"button,omitempty"`
	ContentTile          []*proto1.Content `protobuf:"bytes,5,rep,name=contentTile,proto3" json:"contentTile,omitempty"`
	BackDrop             []string          `protobuf:"bytes,6,rep,name=backDrop,proto3" json:"backDrop,omitempty"`
	Poster               []string          `protobuf:"bytes,7,rep,name=poster,proto3" json:"poster,omitempty"`
	Portrait             []string          `protobuf:"bytes,8,rep,name=portrait,proto3" json:"portrait,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DetailTileInfo) Reset()         { *m = DetailTileInfo{} }
func (m *DetailTileInfo) String() string { return proto.CompactTextString(m) }
func (*DetailTileInfo) ProtoMessage()    {}
func (*DetailTileInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2aea33dfb16e7225, []int{1}
}

func (m *DetailTileInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetailTileInfo.Unmarshal(m, b)
}
func (m *DetailTileInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetailTileInfo.Marshal(b, m, deterministic)
}
func (m *DetailTileInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetailTileInfo.Merge(m, src)
}
func (m *DetailTileInfo) XXX_Size() int {
	return xxx_messageInfo_DetailTileInfo.Size(m)
}
func (m *DetailTileInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DetailTileInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DetailTileInfo proto.InternalMessageInfo

func (m *DetailTileInfo) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *DetailTileInfo) GetSynopsis() string {
	if m != nil {
		return m.Synopsis
	}
	return ""
}

func (m *DetailTileInfo) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *DetailTileInfo) GetButton() []*Button {
	if m != nil {
		return m.Button
	}
	return nil
}

func (m *DetailTileInfo) GetContentTile() []*proto1.Content {
	if m != nil {
		return m.ContentTile
	}
	return nil
}

func (m *DetailTileInfo) GetBackDrop() []string {
	if m != nil {
		return m.BackDrop
	}
	return nil
}

func (m *DetailTileInfo) GetPoster() []string {
	if m != nil {
		return m.Poster
	}
	return nil
}

func (m *DetailTileInfo) GetPortrait() []string {
	if m != nil {
		return m.Portrait
	}
	return nil
}

type Button struct {
	Title                string     `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Icon                 string     `protobuf:"bytes,2,opt,name=icon,proto3" json:"icon,omitempty"`
	Package              string     `protobuf:"bytes,3,opt,name=package,proto3" json:"package,omitempty"`
	Target               string     `protobuf:"bytes,4,opt,name=target,proto3" json:"target,omitempty"`
	ButtonType           ButtonType `protobuf:"varint,6,opt,name=buttonType,proto3,enum=DetailPageService.ButtonType" json:"buttonType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Button) Reset()         { *m = Button{} }
func (m *Button) String() string { return proto.CompactTextString(m) }
func (*Button) ProtoMessage()    {}
func (*Button) Descriptor() ([]byte, []int) {
	return fileDescriptor_2aea33dfb16e7225, []int{2}
}

func (m *Button) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Button.Unmarshal(m, b)
}
func (m *Button) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Button.Marshal(b, m, deterministic)
}
func (m *Button) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Button.Merge(m, src)
}
func (m *Button) XXX_Size() int {
	return xxx_messageInfo_Button.Size(m)
}
func (m *Button) XXX_DiscardUnknown() {
	xxx_messageInfo_Button.DiscardUnknown(m)
}

var xxx_messageInfo_Button proto.InternalMessageInfo

func (m *Button) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Button) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *Button) GetPackage() string {
	if m != nil {
		return m.Package
	}
	return ""
}

func (m *Button) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *Button) GetButtonType() ButtonType {
	if m != nil {
		return m.ButtonType
	}
	return ButtonType_RemotApi
}

func init() {
	proto.RegisterEnum("DetailPageService.ButtonType", ButtonType_name, ButtonType_value)
	proto.RegisterType((*TileInfoRequest)(nil), "DetailPageService.TileInfoRequest")
	proto.RegisterType((*DetailTileInfo)(nil), "DetailPageService.DetailTileInfo")
	proto.RegisterMapType((map[string]string)(nil), "DetailPageService.DetailTileInfo.MetadataEntry")
	proto.RegisterType((*Button)(nil), "DetailPageService.Button")
}

func init() {
	proto.RegisterFile("DetailPageService.proto", fileDescriptor_2aea33dfb16e7225)
}

var fileDescriptor_2aea33dfb16e7225 = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x71, 0xea, 0x86, 0x29, 0x2d, 0x65, 0xc5, 0xcf, 0x62, 0x81, 0x14, 0x7c, 0x0a, 0x3d,
	0xc4, 0x22, 0x5c, 0x50, 0x2b, 0x0e, 0x85, 0x22, 0x84, 0x00, 0x09, 0xb9, 0x7d, 0x81, 0x8d, 0x33,
	0xb8, 0xab, 0x38, 0xbb, 0x66, 0x3d, 0x89, 0xc8, 0x95, 0x33, 0x9c, 0x78, 0x04, 0x1e, 0x89, 0x57,
	0xe0, 0x41, 0xd0, 0xae, 0xd7, 0x25, 0x21, 0xad, 0x7a, 0xf3, 0x37, 0xf3, 0x7d, 0x33, 0x9e, 0x6f,
	0x66, 0xe1, 0xc1, 0x09, 0x92, 0x90, 0xe5, 0x27, 0x51, 0xe0, 0x29, 0x9a, 0x85, 0xcc, 0x71, 0x58,
	0x19, 0x4d, 0x9a, 0xdd, 0xd9, 0x48, 0xc4, 0x8f, 0x0a, 0xad, 0x8b, 0x12, 0x53, 0x51, 0xc9, 0x54,
	0x28, 0xa5, 0x49, 0x90, 0xd4, 0xaa, 0x6e, 0x04, 0xf1, 0x71, 0x21, 0xe9, 0x7c, 0x3e, 0x1e, 0xe6,
	0x7a, 0x96, 0x92, 0x11, 0xaa, 0x16, 0x0b, 0xa3, 0xd3, 0xd3, 0xfc, 0x1c, 0x27, 0xf3, 0xb2, 0xad,
	0x91, 0x3a, 0x6e, 0x1b, 0x15, 0x66, 0xad, 0x67, 0xf2, 0x14, 0x6e, 0x9f, 0xc9, 0x12, 0xdf, 0xa9,
	0xcf, 0x3a, 0xc3, 0x2f, 0x73, 0xac, 0x89, 0xdd, 0x87, 0x88, 0x6c, 0x68, 0xc2, 0x83, 0x7e, 0x30,
	0xb8, 0x99, 0x79, 0x94, 0x7c, 0x0f, 0x61, 0xaf, 0xf9, 0xc3, 0x56, 0xc1, 0xee, 0xc2, 0x16, 0x49,
	0x2a, 0xd1, 0x33, 0x1b, 0xc0, 0x62, 0xe8, 0xd5, 0x4b, 0xa5, 0xab, 0x5a, 0xd6, 0xbc, 0xe3, 0x12,
	0x17, 0x98, 0xbd, 0x87, 0xde, 0x0c, 0x49, 0x4c, 0x04, 0x09, 0x1e, 0xf6, 0xc3, 0xc1, 0xce, 0x28,
	0x1d, 0x6e, 0xfa, 0xb1, 0xde, 0x66, 0xf8, 0xd1, 0x2b, 0xde, 0x28, 0x32, 0xcb, 0xec, 0xa2, 0x00,
	0x7b, 0x06, 0xd1, 0x78, 0x4e, 0xa4, 0x15, 0xef, 0xba, 0x52, 0x0f, 0x2f, 0x29, 0xf5, 0xca, 0x11,
	0x32, 0x4f, 0x64, 0x47, 0xb0, 0x93, 0x6b, 0x45, 0xa8, 0xc8, 0x56, 0xe7, 0x5b, 0x5e, 0xb7, 0xe1,
	0xce, 0xeb, 0x86, 0x94, 0xad, 0xb2, 0xed, 0x60, 0x63, 0x91, 0x4f, 0x4f, 0x8c, 0xae, 0x78, 0xd4,
	0x0f, 0xed, 0x60, 0x2d, 0xb6, 0xae, 0x55, 0xba, 0x26, 0x34, 0x7c, 0xdb, 0x65, 0x3c, 0xb2, 0x9a,
	0x4a, 0x1b, 0x32, 0x42, 0x12, 0xef, 0x35, 0x9a, 0x16, 0xc7, 0x47, 0xb0, 0xbb, 0x36, 0x1a, 0xdb,
	0x87, 0x70, 0x8a, 0x4b, 0xef, 0xa6, 0xfd, 0xb4, 0x0e, 0x2f, 0x44, 0x39, 0x47, 0x6f, 0x64, 0x03,
	0x0e, 0x3b, 0x2f, 0x82, 0xe4, 0x57, 0x00, 0x51, 0x33, 0xdc, 0x15, 0x6b, 0x60, 0xd0, 0x95, 0xb9,
	0x56, 0x5e, 0xe9, 0xbe, 0x19, 0x87, 0xed, 0x4a, 0xe4, 0x53, 0x51, 0x20, 0x0f, 0x5d, 0xb8, 0x85,
	0x6e, 0xeb, 0xc2, 0x14, 0x48, 0xbc, 0xeb, 0xb7, 0xee, 0x10, 0x7b, 0x09, 0xd0, 0x58, 0x77, 0xb6,
	0xac, 0x90, 0x47, 0xfd, 0x60, 0xb0, 0x37, 0x7a, 0x7c, 0xa5, 0xcf, 0x96, 0x94, 0xad, 0x08, 0x0e,
	0x06, 0x00, 0xff, 0x32, 0xec, 0x16, 0xf4, 0x32, 0x9c, 0x69, 0x3a, 0xae, 0xe4, 0xfe, 0x0d, 0x8b,
	0x3e, 0xe8, 0x5c, 0x94, 0x16, 0x05, 0xa3, 0x1f, 0x01, 0x6c, 0x3e, 0x00, 0xf6, 0x15, 0x76, 0xdf,
	0x22, 0x35, 0x71, 0x77, 0x72, 0xc9, 0x25, 0xbd, 0xff, 0xbb, 0xe0, 0xf8, 0xc9, 0xb5, 0x27, 0x95,
	0xf4, 0xbf, 0xfd, 0xfe, 0xf3, 0xb3, 0x13, 0x27, 0xf7, 0xd2, 0x89, 0x4b, 0xa4, 0xc5, 0x6a, 0x97,
	0xc3, 0xe0, 0x60, 0x1c, 0xb9, 0x07, 0xf2, 0xfc, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8e, 0x27,
	0x22, 0xe4, 0xaf, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DetailPageServiceClient is the client API for DetailPageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DetailPageServiceClient interface {
	GetDetailInfo(ctx context.Context, in *TileInfoRequest, opts ...grpc.CallOption) (*DetailTileInfo, error)
}

type detailPageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDetailPageServiceClient(cc grpc.ClientConnInterface) DetailPageServiceClient {
	return &detailPageServiceClient{cc}
}

func (c *detailPageServiceClient) GetDetailInfo(ctx context.Context, in *TileInfoRequest, opts ...grpc.CallOption) (*DetailTileInfo, error) {
	out := new(DetailTileInfo)
	err := c.cc.Invoke(ctx, "/DetailPageService.DetailPageService/GetDetailInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DetailPageServiceServer is the server API for DetailPageService service.
type DetailPageServiceServer interface {
	GetDetailInfo(context.Context, *TileInfoRequest) (*DetailTileInfo, error)
}

// UnimplementedDetailPageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDetailPageServiceServer struct {
}

func (*UnimplementedDetailPageServiceServer) GetDetailInfo(ctx context.Context, req *TileInfoRequest) (*DetailTileInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDetailInfo not implemented")
}

func RegisterDetailPageServiceServer(s *grpc.Server, srv DetailPageServiceServer) {
	s.RegisterService(&_DetailPageService_serviceDesc, srv)
}

func _DetailPageService_GetDetailInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TileInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DetailPageServiceServer).GetDetailInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DetailPageService.DetailPageService/GetDetailInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DetailPageServiceServer).GetDetailInfo(ctx, req.(*TileInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DetailPageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "DetailPageService.DetailPageService",
	HandlerType: (*DetailPageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDetailInfo",
			Handler:    _DetailPageService_GetDetailInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "DetailPageService.proto",
}
