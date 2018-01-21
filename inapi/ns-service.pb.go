// Code generated by protoc-gen-go.
// source: inapi/ns-service.proto
// DO NOT EDIT!

package inapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type NsPodServiceHost struct {
	Rep  uint32 `protobuf:"varint,1,opt,name=rep" json:"rep,omitempty"`
	Ip   string `protobuf:"bytes,2,opt,name=ip" json:"ip,omitempty"`
	Port uint32 `protobuf:"varint,3,opt,name=port" json:"port,omitempty"`
}

func (m *NsPodServiceHost) Reset()                    { *m = NsPodServiceHost{} }
func (m *NsPodServiceHost) String() string            { return proto.CompactTextString(m) }
func (*NsPodServiceHost) ProtoMessage()               {}
func (*NsPodServiceHost) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *NsPodServiceHost) GetRep() uint32 {
	if m != nil {
		return m.Rep
	}
	return 0
}

func (m *NsPodServiceHost) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *NsPodServiceHost) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type NsPodServiceEntry struct {
	Port  uint32              `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	Items []*NsPodServiceHost `protobuf:"bytes,2,rep,name=items" json:"items,omitempty"`
}

func (m *NsPodServiceEntry) Reset()                    { *m = NsPodServiceEntry{} }
func (m *NsPodServiceEntry) String() string            { return proto.CompactTextString(m) }
func (*NsPodServiceEntry) ProtoMessage()               {}
func (*NsPodServiceEntry) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *NsPodServiceEntry) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *NsPodServiceEntry) GetItems() []*NsPodServiceHost {
	if m != nil {
		return m.Items
	}
	return nil
}

type NsPodServiceMap struct {
	Id       string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User     string               `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Services []*NsPodServiceEntry `protobuf:"bytes,3,rep,name=services" json:"services,omitempty"`
	Updated  uint64               `protobuf:"varint,4,opt,name=updated" json:"updated,omitempty"`
}

func (m *NsPodServiceMap) Reset()                    { *m = NsPodServiceMap{} }
func (m *NsPodServiceMap) String() string            { return proto.CompactTextString(m) }
func (*NsPodServiceMap) ProtoMessage()               {}
func (*NsPodServiceMap) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *NsPodServiceMap) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NsPodServiceMap) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *NsPodServiceMap) GetServices() []*NsPodServiceEntry {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *NsPodServiceMap) GetUpdated() uint64 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func init() {
	proto.RegisterType((*NsPodServiceHost)(nil), "inapi.NsPodServiceHost")
	proto.RegisterType((*NsPodServiceEntry)(nil), "inapi.NsPodServiceEntry")
	proto.RegisterType((*NsPodServiceMap)(nil), "inapi.NsPodServiceMap")
}

func init() { proto.RegisterFile("inapi/ns-service.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xbf, 0x4e, 0x04, 0x21,
	0x10, 0xc6, 0x03, 0xec, 0xf9, 0x67, 0x8c, 0x7a, 0x4e, 0xa1, 0x53, 0x92, 0xad, 0x68, 0x6e, 0x4d,
	0xd4, 0x57, 0x30, 0xb9, 0x46, 0x63, 0x30, 0xb1, 0x5f, 0x85, 0x82, 0xc2, 0x85, 0x00, 0x67, 0xe2,
	0x03, 0xf8, 0xde, 0x66, 0x67, 0xcf, 0x75, 0x63, 0xec, 0x3e, 0x98, 0x5f, 0xe6, 0xfb, 0x01, 0x5c,
	0x86, 0xa1, 0x4f, 0xe1, 0x7a, 0x28, 0x9b, 0xe2, 0xf3, 0x47, 0x78, 0xf3, 0x5d, 0xca, 0xb1, 0x46,
	0x5c, 0xf1, 0x7d, 0xbb, 0x85, 0xf5, 0x63, 0x79, 0x8a, 0xee, 0x79, 0x1a, 0x6e, 0x63, 0xa9, 0xb8,
	0x06, 0x95, 0x7d, 0x22, 0xa1, 0x85, 0x39, 0xb5, 0x63, 0xc4, 0x33, 0x90, 0x21, 0x91, 0xd4, 0xc2,
	0x1c, 0x5b, 0x19, 0x12, 0x22, 0x34, 0x29, 0xe6, 0x4a, 0x8a, 0x11, 0xce, 0xed, 0x0b, 0x5c, 0x2c,
	0x37, 0xdd, 0x0f, 0x35, 0x7f, 0xce, 0xa0, 0xf8, 0x05, 0x71, 0x03, 0xab, 0x50, 0xfd, 0x7b, 0x21,
	0xa9, 0x95, 0x39, 0xb9, 0xb9, 0xea, 0xd8, 0xa4, 0xfb, 0xab, 0x61, 0x27, 0xaa, 0xfd, 0x12, 0x70,
	0xbe, 0x9c, 0x3d, 0xf4, 0x93, 0x8f, 0xe3, 0xa5, 0xa3, 0x8f, 0x1b, 0x6b, 0x76, 0xc5, 0xe7, 0xbd,
	0x21, 0x67, 0xbc, 0x83, 0xa3, 0xfd, 0x8b, 0x0b, 0x29, 0x6e, 0xa2, 0x7f, 0x9a, 0x58, 0xd3, 0xce,
	0x24, 0x12, 0x1c, 0xee, 0x92, 0xeb, 0xab, 0x77, 0xd4, 0x68, 0x61, 0x1a, 0xfb, 0x73, 0x7c, 0x3d,
	0xe0, 0x7f, 0xbb, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x84, 0xce, 0x3e, 0x28, 0x51, 0x01, 0x00,
	0x00,
}
