// Code generated by protoc-gen-go.
// source: inapi/operator.proto
// DO NOT EDIT!

package inapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PbOpLogEntry struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Status  string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	Updated uint64 `protobuf:"varint,3,opt,name=updated" json:"updated,omitempty"`
	Message string `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
}

func (m *PbOpLogEntry) Reset()                    { *m = PbOpLogEntry{} }
func (m *PbOpLogEntry) String() string            { return proto.CompactTextString(m) }
func (*PbOpLogEntry) ProtoMessage()               {}
func (*PbOpLogEntry) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *PbOpLogEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbOpLogEntry) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PbOpLogEntry) GetUpdated() uint64 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func (m *PbOpLogEntry) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type PbOpLogSets struct {
	Name    string          `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Version uint32          `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
	Items   []*PbOpLogEntry `protobuf:"bytes,4,rep,name=items" json:"items,omitempty"`
}

func (m *PbOpLogSets) Reset()                    { *m = PbOpLogSets{} }
func (m *PbOpLogSets) String() string            { return proto.CompactTextString(m) }
func (*PbOpLogSets) ProtoMessage()               {}
func (*PbOpLogSets) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *PbOpLogSets) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbOpLogSets) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PbOpLogSets) GetItems() []*PbOpLogEntry {
	if m != nil {
		return m.Items
	}
	return nil
}

func init() {
	proto.RegisterType((*PbOpLogEntry)(nil), "inapi.PbOpLogEntry")
	proto.RegisterType((*PbOpLogSets)(nil), "inapi.PbOpLogSets")
}

func init() { proto.RegisterFile("inapi/operator.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xc1, 0x8a, 0xc2, 0x30,
	0x10, 0x86, 0xe9, 0x36, 0x6d, 0xd9, 0xe9, 0xee, 0x25, 0x8a, 0xe4, 0x58, 0x7a, 0xaa, 0x97, 0x0a,
	0xfa, 0x0c, 0xde, 0x04, 0x25, 0x3e, 0x41, 0x4a, 0xc7, 0x92, 0x43, 0x93, 0x90, 0x4c, 0x05, 0xdf,
	0x5e, 0x8c, 0x2d, 0x78, 0xf0, 0x96, 0xef, 0xcf, 0xcc, 0xfc, 0x7c, 0xb0, 0xd6, 0x46, 0x39, 0xbd,
	0xb3, 0x0e, 0xbd, 0x22, 0xeb, 0x5b, 0xe7, 0x2d, 0x59, 0x9e, 0xc5, 0xb4, 0x36, 0xf0, 0x77, 0xe9,
	0xce, 0xee, 0x64, 0x87, 0xa3, 0x21, 0xff, 0xe0, 0x1c, 0x98, 0x51, 0x23, 0x8a, 0xa4, 0x4a, 0x9a,
	0x5f, 0x19, 0xdf, 0x7c, 0x03, 0x79, 0x20, 0x45, 0x53, 0x10, 0x3f, 0x31, 0x9d, 0x89, 0x0b, 0x28,
	0x26, 0xd7, 0x2b, 0xc2, 0x5e, 0xa4, 0x55, 0xd2, 0x30, 0xb9, 0xe0, 0xeb, 0x67, 0xc4, 0x10, 0xd4,
	0x80, 0x82, 0xc5, 0x95, 0x05, 0xeb, 0x1b, 0x94, 0x73, 0xdf, 0x15, 0x29, 0x7c, 0xad, 0x13, 0x50,
	0xdc, 0xd1, 0x07, 0x6d, 0x4d, 0x3c, 0xfb, 0x2f, 0x17, 0xe4, 0x5b, 0xc8, 0x34, 0xe1, 0x18, 0x04,
	0xab, 0xd2, 0xa6, 0xdc, 0xaf, 0xda, 0xe8, 0xd0, 0x7e, 0x0a, 0xc8, 0xf7, 0x44, 0x97, 0x47, 0xcb,
	0xc3, 0x33, 0x00, 0x00, 0xff, 0xff, 0xd8, 0x5a, 0x4d, 0x52, 0xfd, 0x00, 0x00, 0x00,
}
