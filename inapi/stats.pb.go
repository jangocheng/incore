// Code generated by protoc-gen-go. DO NOT EDIT.
// source: inapi/stats.proto

package inapi

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type PbStatsSampleFeed struct {
	Kind                 string                `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	Cycle                uint32                `protobuf:"varint,2,opt,name=cycle,proto3" json:"cycle,omitempty"`
	Items                []*PbStatsSampleEntry `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PbStatsSampleFeed) Reset()         { *m = PbStatsSampleFeed{} }
func (m *PbStatsSampleFeed) String() string { return proto.CompactTextString(m) }
func (*PbStatsSampleFeed) ProtoMessage()    {}
func (*PbStatsSampleFeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_46f1f394ff918ee5, []int{0}
}

func (m *PbStatsSampleFeed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbStatsSampleFeed.Unmarshal(m, b)
}
func (m *PbStatsSampleFeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbStatsSampleFeed.Marshal(b, m, deterministic)
}
func (m *PbStatsSampleFeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbStatsSampleFeed.Merge(m, src)
}
func (m *PbStatsSampleFeed) XXX_Size() int {
	return xxx_messageInfo_PbStatsSampleFeed.Size(m)
}
func (m *PbStatsSampleFeed) XXX_DiscardUnknown() {
	xxx_messageInfo_PbStatsSampleFeed.DiscardUnknown(m)
}

var xxx_messageInfo_PbStatsSampleFeed proto.InternalMessageInfo

func (m *PbStatsSampleFeed) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *PbStatsSampleFeed) GetCycle() uint32 {
	if m != nil {
		return m.Cycle
	}
	return 0
}

func (m *PbStatsSampleFeed) GetItems() []*PbStatsSampleEntry {
	if m != nil {
		return m.Items
	}
	return nil
}

type PbStatsSampleEntry struct {
	Name                 string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Items                []*PbStatsSampleValue `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PbStatsSampleEntry) Reset()         { *m = PbStatsSampleEntry{} }
func (m *PbStatsSampleEntry) String() string { return proto.CompactTextString(m) }
func (*PbStatsSampleEntry) ProtoMessage()    {}
func (*PbStatsSampleEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_46f1f394ff918ee5, []int{1}
}

func (m *PbStatsSampleEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbStatsSampleEntry.Unmarshal(m, b)
}
func (m *PbStatsSampleEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbStatsSampleEntry.Marshal(b, m, deterministic)
}
func (m *PbStatsSampleEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbStatsSampleEntry.Merge(m, src)
}
func (m *PbStatsSampleEntry) XXX_Size() int {
	return xxx_messageInfo_PbStatsSampleEntry.Size(m)
}
func (m *PbStatsSampleEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_PbStatsSampleEntry.DiscardUnknown(m)
}

var xxx_messageInfo_PbStatsSampleEntry proto.InternalMessageInfo

func (m *PbStatsSampleEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbStatsSampleEntry) GetItems() []*PbStatsSampleValue {
	if m != nil {
		return m.Items
	}
	return nil
}

type PbStatsSampleValue struct {
	Time                 uint32   `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Value                int64    `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbStatsSampleValue) Reset()         { *m = PbStatsSampleValue{} }
func (m *PbStatsSampleValue) String() string { return proto.CompactTextString(m) }
func (*PbStatsSampleValue) ProtoMessage()    {}
func (*PbStatsSampleValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_46f1f394ff918ee5, []int{2}
}

func (m *PbStatsSampleValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbStatsSampleValue.Unmarshal(m, b)
}
func (m *PbStatsSampleValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbStatsSampleValue.Marshal(b, m, deterministic)
}
func (m *PbStatsSampleValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbStatsSampleValue.Merge(m, src)
}
func (m *PbStatsSampleValue) XXX_Size() int {
	return xxx_messageInfo_PbStatsSampleValue.Size(m)
}
func (m *PbStatsSampleValue) XXX_DiscardUnknown() {
	xxx_messageInfo_PbStatsSampleValue.DiscardUnknown(m)
}

var xxx_messageInfo_PbStatsSampleValue proto.InternalMessageInfo

func (m *PbStatsSampleValue) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *PbStatsSampleValue) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type PbStatsIndexList struct {
	IndexCycle           uint32              `protobuf:"varint,1,opt,name=index_cycle,json=indexCycle,proto3" json:"index_cycle,omitempty"`
	SampleCycle          uint32              `protobuf:"varint,2,opt,name=sample_cycle,json=sampleCycle,proto3" json:"sample_cycle,omitempty"`
	Items                []*PbStatsIndexFeed `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *PbStatsIndexList) Reset()         { *m = PbStatsIndexList{} }
func (m *PbStatsIndexList) String() string { return proto.CompactTextString(m) }
func (*PbStatsIndexList) ProtoMessage()    {}
func (*PbStatsIndexList) Descriptor() ([]byte, []int) {
	return fileDescriptor_46f1f394ff918ee5, []int{3}
}

func (m *PbStatsIndexList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbStatsIndexList.Unmarshal(m, b)
}
func (m *PbStatsIndexList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbStatsIndexList.Marshal(b, m, deterministic)
}
func (m *PbStatsIndexList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbStatsIndexList.Merge(m, src)
}
func (m *PbStatsIndexList) XXX_Size() int {
	return xxx_messageInfo_PbStatsIndexList.Size(m)
}
func (m *PbStatsIndexList) XXX_DiscardUnknown() {
	xxx_messageInfo_PbStatsIndexList.DiscardUnknown(m)
}

var xxx_messageInfo_PbStatsIndexList proto.InternalMessageInfo

func (m *PbStatsIndexList) GetIndexCycle() uint32 {
	if m != nil {
		return m.IndexCycle
	}
	return 0
}

func (m *PbStatsIndexList) GetSampleCycle() uint32 {
	if m != nil {
		return m.SampleCycle
	}
	return 0
}

func (m *PbStatsIndexList) GetItems() []*PbStatsIndexFeed {
	if m != nil {
		return m.Items
	}
	return nil
}

type PbStatsIndexFeed struct {
	Time                 uint32                `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Items                []*PbStatsSampleEntry `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PbStatsIndexFeed) Reset()         { *m = PbStatsIndexFeed{} }
func (m *PbStatsIndexFeed) String() string { return proto.CompactTextString(m) }
func (*PbStatsIndexFeed) ProtoMessage()    {}
func (*PbStatsIndexFeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_46f1f394ff918ee5, []int{4}
}

func (m *PbStatsIndexFeed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbStatsIndexFeed.Unmarshal(m, b)
}
func (m *PbStatsIndexFeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbStatsIndexFeed.Marshal(b, m, deterministic)
}
func (m *PbStatsIndexFeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbStatsIndexFeed.Merge(m, src)
}
func (m *PbStatsIndexFeed) XXX_Size() int {
	return xxx_messageInfo_PbStatsIndexFeed.Size(m)
}
func (m *PbStatsIndexFeed) XXX_DiscardUnknown() {
	xxx_messageInfo_PbStatsIndexFeed.DiscardUnknown(m)
}

var xxx_messageInfo_PbStatsIndexFeed proto.InternalMessageInfo

func (m *PbStatsIndexFeed) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *PbStatsIndexFeed) GetItems() []*PbStatsSampleEntry {
	if m != nil {
		return m.Items
	}
	return nil
}

func init() {
	proto.RegisterType((*PbStatsSampleFeed)(nil), "inapi.PbStatsSampleFeed")
	proto.RegisterType((*PbStatsSampleEntry)(nil), "inapi.PbStatsSampleEntry")
	proto.RegisterType((*PbStatsSampleValue)(nil), "inapi.PbStatsSampleValue")
	proto.RegisterType((*PbStatsIndexList)(nil), "inapi.PbStatsIndexList")
	proto.RegisterType((*PbStatsIndexFeed)(nil), "inapi.PbStatsIndexFeed")
}

func init() { proto.RegisterFile("inapi/stats.proto", fileDescriptor_46f1f394ff918ee5) }

var fileDescriptor_46f1f394ff918ee5 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcc, 0xcc, 0x4b, 0x2c,
	0xc8, 0xd4, 0x2f, 0x2e, 0x49, 0x2c, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05,
	0x0b, 0x29, 0xe5, 0x71, 0x09, 0x06, 0x24, 0x05, 0x83, 0xc4, 0x83, 0x13, 0x73, 0x0b, 0x72, 0x52,
	0xdd, 0x52, 0x53, 0x53, 0x84, 0x84, 0xb8, 0x58, 0xb2, 0x33, 0xf3, 0x52, 0x24, 0x18, 0x15, 0x18,
	0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x11, 0x2e, 0xd6, 0xe4, 0xca, 0xe4, 0x9c, 0x54, 0x09, 0x26,
	0x05, 0x46, 0x0d, 0xde, 0x20, 0x08, 0x47, 0x48, 0x9f, 0x8b, 0x35, 0xb3, 0x24, 0x35, 0xb7, 0x58,
	0x82, 0x59, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x52, 0x0f, 0x6c, 0xaa, 0x1e, 0x8a, 0x91, 0xae, 0x79,
	0x25, 0x45, 0x95, 0x41, 0x10, 0x75, 0x4a, 0x91, 0x5c, 0x42, 0x98, 0x92, 0x20, 0x0b, 0xf3, 0x12,
	0x73, 0x53, 0x61, 0x16, 0x82, 0xd8, 0x08, 0xa3, 0x99, 0x70, 0x1b, 0x1d, 0x96, 0x98, 0x53, 0x9a,
	0x0a, 0x33, 0xda, 0x0e, 0xcd, 0x68, 0xb0, 0x24, 0xc8, 0xe8, 0x92, 0x4c, 0xa8, 0xd1, 0xbc, 0x41,
	0x60, 0x36, 0xc8, 0x2f, 0x65, 0x20, 0x49, 0xb0, 0x5f, 0x98, 0x83, 0x20, 0x1c, 0xa5, 0x56, 0x46,
	0x2e, 0x01, 0xa8, 0x01, 0x9e, 0x79, 0x29, 0xa9, 0x15, 0x3e, 0x99, 0xc5, 0x25, 0x42, 0xf2, 0x5c,
	0xdc, 0x99, 0x20, 0x4e, 0x3c, 0xc4, 0xf3, 0x10, 0x53, 0xb8, 0xc0, 0x42, 0xce, 0xe0, 0x10, 0x50,
	0xe4, 0xe2, 0x29, 0x06, 0x5b, 0x17, 0x8f, 0x1c, 0x3c, 0xdc, 0x10, 0x31, 0x88, 0x12, 0x5d, 0xd4,
	0x40, 0x12, 0x47, 0xf5, 0x09, 0xd8, 0x2e, 0x50, 0xb0, 0xc3, 0xfc, 0x11, 0x8e, 0xea, 0x0c, 0x58,
	0x8c, 0x60, 0xf8, 0x82, 0x98, 0x00, 0x42, 0x0e, 0xfb, 0x24, 0x36, 0x70, 0xcc, 0x1b, 0x03, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x46, 0xfc, 0x5c, 0xfc, 0x0e, 0x02, 0x00, 0x00,
}
