// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: node/v1/node.proto

package nodev1

import (
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// GuardianSet represents a new guardian set to be submitted to and signed by the node.
// During the genesis procedure, this data structure will be assembled using off-chain collaborative tooling
// like GitHub using a human-readable encoding, so readability is a concern.
type GuardianSetUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Index of the current guardian set to be replaced.
	CurrentSetIndex uint32 `protobuf:"varint,1,opt,name=current_set_index,json=currentSetIndex,proto3" json:"current_set_index,omitempty"`
	// UNIX timestamp (s) of the VAA to be created. The timestamp is informational and will be part
	// of the VAA submitted to the chain. It's part of the VAA digest and has to be identical across nodes.
	//
	// For lockups, the timestamp identifies the block that the lockup belongs to. For guardian set updates,
	// we create the VAA manually. Best practice is to pick a timestamp which roughly matches the expected
	// genesis ceremony data.
	//
	// The actual on-chain guardian set creation timestamp will be set when the VAA is accepted on each chain.
	//
	// This is a uint32 to match the on-chain timestamp representation. This becomes a problem in 2106 (sorry).
	Timestamp uint32                        `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Guardians []*GuardianSetUpdate_Guardian `protobuf:"bytes,3,rep,name=guardians,proto3" json:"guardians,omitempty"`
}

func (x *GuardianSetUpdate) Reset() {
	*x = GuardianSetUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_v1_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuardianSetUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuardianSetUpdate) ProtoMessage() {}

func (x *GuardianSetUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_node_v1_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuardianSetUpdate.ProtoReflect.Descriptor instead.
func (*GuardianSetUpdate) Descriptor() ([]byte, []int) {
	return file_node_v1_node_proto_rawDescGZIP(), []int{0}
}

func (x *GuardianSetUpdate) GetCurrentSetIndex() uint32 {
	if x != nil {
		return x.CurrentSetIndex
	}
	return 0
}

func (x *GuardianSetUpdate) GetTimestamp() uint32 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *GuardianSetUpdate) GetGuardians() []*GuardianSetUpdate_Guardian {
	if x != nil {
		return x.Guardians
	}
	return nil
}

type SubmitGuardianSetVAARequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GuardianSet *GuardianSetUpdate `protobuf:"bytes,1,opt,name=guardian_set,json=guardianSet,proto3" json:"guardian_set,omitempty"`
}

func (x *SubmitGuardianSetVAARequest) Reset() {
	*x = SubmitGuardianSetVAARequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_v1_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitGuardianSetVAARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitGuardianSetVAARequest) ProtoMessage() {}

func (x *SubmitGuardianSetVAARequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_v1_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitGuardianSetVAARequest.ProtoReflect.Descriptor instead.
func (*SubmitGuardianSetVAARequest) Descriptor() ([]byte, []int) {
	return file_node_v1_node_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitGuardianSetVAARequest) GetGuardianSet() *GuardianSetUpdate {
	if x != nil {
		return x.GuardianSet
	}
	return nil
}

type SubmitGuardianSetVAAResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Canonical digest of the submitted VAA.
	Digest []byte `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *SubmitGuardianSetVAAResponse) Reset() {
	*x = SubmitGuardianSetVAAResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_v1_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitGuardianSetVAAResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitGuardianSetVAAResponse) ProtoMessage() {}

func (x *SubmitGuardianSetVAAResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_v1_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitGuardianSetVAAResponse.ProtoReflect.Descriptor instead.
func (*SubmitGuardianSetVAAResponse) Descriptor() ([]byte, []int) {
	return file_node_v1_node_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitGuardianSetVAAResponse) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

// GuardianKey specifies the on-disk format for a node's guardian key.
type GuardianKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// description is an optional, free-form description text set by the operator.
	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	// data is the binary representation of the secp256k1 private key.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// pubkey is a human-readable representation of the key, included for operator convenience.
	Pubkey string `protobuf:"bytes,3,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}

func (x *GuardianKey) Reset() {
	*x = GuardianKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_v1_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuardianKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuardianKey) ProtoMessage() {}

func (x *GuardianKey) ProtoReflect() protoreflect.Message {
	mi := &file_node_v1_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuardianKey.ProtoReflect.Descriptor instead.
func (*GuardianKey) Descriptor() ([]byte, []int) {
	return file_node_v1_node_proto_rawDescGZIP(), []int{3}
}

func (x *GuardianKey) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GuardianKey) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GuardianKey) GetPubkey() string {
	if x != nil {
		return x.Pubkey
	}
	return ""
}

// List of guardian set members.
type GuardianSetUpdate_Guardian struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Guardian key pubkey. Stored as hex string with 0x prefix for human readability -
	// this is the canonical Ethereum representation.
	Pubkey string `protobuf:"bytes,1,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	// Optional descriptive name. Not stored on any chain, purely informational.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GuardianSetUpdate_Guardian) Reset() {
	*x = GuardianSetUpdate_Guardian{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_v1_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuardianSetUpdate_Guardian) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuardianSetUpdate_Guardian) ProtoMessage() {}

func (x *GuardianSetUpdate_Guardian) ProtoReflect() protoreflect.Message {
	mi := &file_node_v1_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuardianSetUpdate_Guardian.ProtoReflect.Descriptor instead.
func (*GuardianSetUpdate_Guardian) Descriptor() ([]byte, []int) {
	return file_node_v1_node_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GuardianSetUpdate_Guardian) GetPubkey() string {
	if x != nil {
		return x.Pubkey
	}
	return ""
}

func (x *GuardianSetUpdate_Guardian) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_node_v1_node_proto protoreflect.FileDescriptor

var file_node_v1_node_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd8, 0x01, 0x0a, 0x11,
	0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x74,
	0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x41, 0x0a, 0x09, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61,
	0x6e, 0x53, 0x65, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x47, 0x75, 0x61, 0x72, 0x64,
	0x69, 0x61, 0x6e, 0x52, 0x09, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x73, 0x1a, 0x36,
	0x0a, 0x08, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x75,
	0x62, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x75, 0x62, 0x6b,
	0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5c, 0x0a, 0x1b, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65, 0x74, 0x56, 0x41, 0x41, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x0c, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61,
	0x6e, 0x5f, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65,
	0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x67, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61,
	0x6e, 0x53, 0x65, 0x74, 0x22, 0x36, 0x0a, 0x1c, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x47, 0x75,
	0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65, 0x74, 0x56, 0x41, 0x41, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x5b, 0x0a, 0x0b,
	0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x32, 0x75, 0x0a, 0x0e, 0x4e, 0x6f, 0x64,
	0x65, 0x50, 0x72, 0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x64, 0x12, 0x63, 0x0a, 0x14, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65, 0x74,
	0x56, 0x41, 0x41, 0x12, 0x24, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69, 0x61, 0x6e, 0x53, 0x65, 0x74, 0x56,
	0x41, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x47, 0x75, 0x61, 0x72, 0x64, 0x69,
	0x61, 0x6e, 0x53, 0x65, 0x74, 0x56, 0x41, 0x41, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x16, 0x5a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x76,
	0x31, 0x3b, 0x6e, 0x6f, 0x64, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_v1_node_proto_rawDescOnce sync.Once
	file_node_v1_node_proto_rawDescData = file_node_v1_node_proto_rawDesc
)

func file_node_v1_node_proto_rawDescGZIP() []byte {
	file_node_v1_node_proto_rawDescOnce.Do(func() {
		file_node_v1_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_v1_node_proto_rawDescData)
	})
	return file_node_v1_node_proto_rawDescData
}

var file_node_v1_node_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_node_v1_node_proto_goTypes = []interface{}{
	(*GuardianSetUpdate)(nil),            // 0: node.v1.GuardianSetUpdate
	(*SubmitGuardianSetVAARequest)(nil),  // 1: node.v1.SubmitGuardianSetVAARequest
	(*SubmitGuardianSetVAAResponse)(nil), // 2: node.v1.SubmitGuardianSetVAAResponse
	(*GuardianKey)(nil),                  // 3: node.v1.GuardianKey
	(*GuardianSetUpdate_Guardian)(nil),   // 4: node.v1.GuardianSetUpdate.Guardian
}
var file_node_v1_node_proto_depIdxs = []int32{
	4, // 0: node.v1.GuardianSetUpdate.guardians:type_name -> node.v1.GuardianSetUpdate.Guardian
	0, // 1: node.v1.SubmitGuardianSetVAARequest.guardian_set:type_name -> node.v1.GuardianSetUpdate
	1, // 2: node.v1.NodePrivileged.SubmitGuardianSetVAA:input_type -> node.v1.SubmitGuardianSetVAARequest
	2, // 3: node.v1.NodePrivileged.SubmitGuardianSetVAA:output_type -> node.v1.SubmitGuardianSetVAAResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_node_v1_node_proto_init() }
func file_node_v1_node_proto_init() {
	if File_node_v1_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_v1_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuardianSetUpdate); i {
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
		file_node_v1_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitGuardianSetVAARequest); i {
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
		file_node_v1_node_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitGuardianSetVAAResponse); i {
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
		file_node_v1_node_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuardianKey); i {
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
		file_node_v1_node_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuardianSetUpdate_Guardian); i {
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
			RawDescriptor: file_node_v1_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_node_v1_node_proto_goTypes,
		DependencyIndexes: file_node_v1_node_proto_depIdxs,
		MessageInfos:      file_node_v1_node_proto_msgTypes,
	}.Build()
	File_node_v1_node_proto = out.File
	file_node_v1_node_proto_rawDesc = nil
	file_node_v1_node_proto_goTypes = nil
	file_node_v1_node_proto_depIdxs = nil
}
