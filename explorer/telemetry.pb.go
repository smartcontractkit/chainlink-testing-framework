package explorer

// original pb https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting/internal/serialization/protobuf/cl_offchainreporting_telemetry.pb.go

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

type TelemetryWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Wrapped:
	//  *TelemetryWrapper_MessageReceived
	//  *TelemetryWrapper_MessageBroadcast
	//  *TelemetryWrapper_MessageSent
	//  *TelemetryWrapper_AssertionViolation
	//  *TelemetryWrapper_RoundStarted
	Wrapped isTelemetryWrapper_Wrapped `protobuf_oneof:"wrapped"`
}

func (x *TelemetryWrapper) Reset() {
	*x = TelemetryWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryWrapper) ProtoMessage() {}

func (x *TelemetryWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryWrapper.ProtoReflect.Descriptor instead.
func (*TelemetryWrapper) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{0}
}

func (m *TelemetryWrapper) GetWrapped() isTelemetryWrapper_Wrapped {
	if m != nil {
		return m.Wrapped
	}
	return nil
}

func (x *TelemetryWrapper) GetMessageReceived() *TelemetryMessageReceived {
	if x, ok := x.GetWrapped().(*TelemetryWrapper_MessageReceived); ok {
		return x.MessageReceived
	}
	return nil
}

func (x *TelemetryWrapper) GetMessageBroadcast() *TelemetryMessageBroadcast {
	if x, ok := x.GetWrapped().(*TelemetryWrapper_MessageBroadcast); ok {
		return x.MessageBroadcast
	}
	return nil
}

func (x *TelemetryWrapper) GetMessageSent() *TelemetryMessageSent {
	if x, ok := x.GetWrapped().(*TelemetryWrapper_MessageSent); ok {
		return x.MessageSent
	}
	return nil
}

func (x *TelemetryWrapper) GetAssertionViolation() *TelemetryAssertionViolation {
	if x, ok := x.GetWrapped().(*TelemetryWrapper_AssertionViolation); ok {
		return x.AssertionViolation
	}
	return nil
}

func (x *TelemetryWrapper) GetRoundStarted() *TelemetryRoundStarted {
	if x, ok := x.GetWrapped().(*TelemetryWrapper_RoundStarted); ok {
		return x.RoundStarted
	}
	return nil
}

type isTelemetryWrapper_Wrapped interface {
	isTelemetryWrapper_Wrapped()
}

type TelemetryWrapper_MessageReceived struct {
	MessageReceived *TelemetryMessageReceived `protobuf:"bytes,1,opt,name=messageReceived,proto3,oneof"`
}

type TelemetryWrapper_MessageBroadcast struct {
	MessageBroadcast *TelemetryMessageBroadcast `protobuf:"bytes,2,opt,name=messageBroadcast,proto3,oneof"`
}

type TelemetryWrapper_MessageSent struct {
	MessageSent *TelemetryMessageSent `protobuf:"bytes,3,opt,name=messageSent,proto3,oneof"`
}

type TelemetryWrapper_AssertionViolation struct {
	AssertionViolation *TelemetryAssertionViolation `protobuf:"bytes,4,opt,name=assertionViolation,proto3,oneof"`
}

type TelemetryWrapper_RoundStarted struct {
	RoundStarted *TelemetryRoundStarted `protobuf:"bytes,5,opt,name=roundStarted,proto3,oneof"`
}

func (*TelemetryWrapper_MessageReceived) isTelemetryWrapper_Wrapped() {}

func (*TelemetryWrapper_MessageBroadcast) isTelemetryWrapper_Wrapped() {}

func (*TelemetryWrapper_MessageSent) isTelemetryWrapper_Wrapped() {}

func (*TelemetryWrapper_AssertionViolation) isTelemetryWrapper_Wrapped() {}

func (*TelemetryWrapper_RoundStarted) isTelemetryWrapper_Wrapped() {}

type TelemetryMessageReceived struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg          *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Sender       uint32          `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *TelemetryMessageReceived) Reset() {
	*x = TelemetryMessageReceived{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryMessageReceived) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryMessageReceived) ProtoMessage() {}

func (x *TelemetryMessageReceived) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryMessageReceived.ProtoReflect.Descriptor instead.
func (*TelemetryMessageReceived) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{1}
}

func (x *TelemetryMessageReceived) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryMessageReceived) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryMessageReceived) GetSender() uint32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

type TelemetryMessageBroadcast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest  []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg           *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	SerializedMsg []byte          `protobuf:"bytes,3,opt,name=serializedMsg,proto3" json:"serializedMsg,omitempty"`
}

func (x *TelemetryMessageBroadcast) Reset() {
	*x = TelemetryMessageBroadcast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryMessageBroadcast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryMessageBroadcast) ProtoMessage() {}

func (x *TelemetryMessageBroadcast) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryMessageBroadcast.ProtoReflect.Descriptor instead.
func (*TelemetryMessageBroadcast) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{2}
}

func (x *TelemetryMessageBroadcast) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryMessageBroadcast) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryMessageBroadcast) GetSerializedMsg() []byte {
	if x != nil {
		return x.SerializedMsg
	}
	return nil
}

type TelemetryMessageSent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest  []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg           *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	SerializedMsg []byte          `protobuf:"bytes,3,opt,name=serializedMsg,proto3" json:"serializedMsg,omitempty"`
	Receiver      uint32          `protobuf:"varint,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
}

func (x *TelemetryMessageSent) Reset() {
	*x = TelemetryMessageSent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryMessageSent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryMessageSent) ProtoMessage() {}

func (x *TelemetryMessageSent) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryMessageSent.ProtoReflect.Descriptor instead.
func (*TelemetryMessageSent) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{3}
}

func (x *TelemetryMessageSent) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryMessageSent) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryMessageSent) GetSerializedMsg() []byte {
	if x != nil {
		return x.SerializedMsg
	}
	return nil
}

func (x *TelemetryMessageSent) GetReceiver() uint32 {
	if x != nil {
		return x.Receiver
	}
	return 0
}

type TelemetryAssertionViolation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Violation:
	//  *TelemetryAssertionViolation_InvalidSignature
	//  *TelemetryAssertionViolation_InvalidSerialization
	Violation isTelemetryAssertionViolation_Violation `protobuf_oneof:"violation"`
}

func (x *TelemetryAssertionViolation) Reset() {
	*x = TelemetryAssertionViolation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryAssertionViolation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryAssertionViolation) ProtoMessage() {}

func (x *TelemetryAssertionViolation) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryAssertionViolation.ProtoReflect.Descriptor instead.
func (*TelemetryAssertionViolation) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{4}
}

func (m *TelemetryAssertionViolation) GetViolation() isTelemetryAssertionViolation_Violation {
	if m != nil {
		return m.Violation
	}
	return nil
}

func (x *TelemetryAssertionViolation) GetInvalidSignature() *TelemetryAssertionViolationInvalidSignature {
	if x, ok := x.GetViolation().(*TelemetryAssertionViolation_InvalidSignature); ok {
		return x.InvalidSignature
	}
	return nil
}

func (x *TelemetryAssertionViolation) GetInvalidSerialization() *TelemetryAssertionViolationInvalidSerialization {
	if x, ok := x.GetViolation().(*TelemetryAssertionViolation_InvalidSerialization); ok {
		return x.InvalidSerialization
	}
	return nil
}

type isTelemetryAssertionViolation_Violation interface {
	isTelemetryAssertionViolation_Violation()
}

type TelemetryAssertionViolation_InvalidSignature struct {
	InvalidSignature *TelemetryAssertionViolationInvalidSignature `protobuf:"bytes,1,opt,name=invalidSignature,proto3,oneof"`
}

type TelemetryAssertionViolation_InvalidSerialization struct {
	InvalidSerialization *TelemetryAssertionViolationInvalidSerialization `protobuf:"bytes,2,opt,name=invalidSerialization,proto3,oneof"`
}

func (*TelemetryAssertionViolation_InvalidSignature) isTelemetryAssertionViolation_Violation() {}

func (*TelemetryAssertionViolation_InvalidSerialization) isTelemetryAssertionViolation_Violation() {}

type TelemetryAssertionViolationInvalidSignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg          *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Sender       uint32          `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *TelemetryAssertionViolationInvalidSignature) Reset() {
	*x = TelemetryAssertionViolationInvalidSignature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryAssertionViolationInvalidSignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryAssertionViolationInvalidSignature) ProtoMessage() {}

func (x *TelemetryAssertionViolationInvalidSignature) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryAssertionViolationInvalidSignature.ProtoReflect.Descriptor instead.
func (*TelemetryAssertionViolationInvalidSignature) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{5}
}

func (x *TelemetryAssertionViolationInvalidSignature) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSignature) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSignature) GetSender() uint32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

type TelemetryAssertionViolationInvalidSerialization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest  []byte `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	SerializedMsg []byte `protobuf:"bytes,2,opt,name=serializedMsg,proto3" json:"serializedMsg,omitempty"`
	Sender        uint32 `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *TelemetryAssertionViolationInvalidSerialization) Reset() {
	*x = TelemetryAssertionViolationInvalidSerialization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryAssertionViolationInvalidSerialization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryAssertionViolationInvalidSerialization) ProtoMessage() {}

func (x *TelemetryAssertionViolationInvalidSerialization) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryAssertionViolationInvalidSerialization.ProtoReflect.Descriptor instead.
func (*TelemetryAssertionViolationInvalidSerialization) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{6}
}

func (x *TelemetryAssertionViolationInvalidSerialization) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSerialization) GetSerializedMsg() []byte {
	if x != nil {
		return x.SerializedMsg
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSerialization) GetSender() uint32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

type TelemetryRoundStarted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Epoch        uint64 `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round        uint64 `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
	Leader       uint64 `protobuf:"varint,4,opt,name=leader,proto3" json:"leader,omitempty"`
	Time         uint64 `protobuf:"varint,5,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *TelemetryRoundStarted) Reset() {
	*x = TelemetryRoundStarted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryRoundStarted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryRoundStarted) ProtoMessage() {}

func (x *TelemetryRoundStarted) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_telemetry_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelemetryRoundStarted.ProtoReflect.Descriptor instead.
func (*TelemetryRoundStarted) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_telemetry_proto_rawDescGZIP(), []int{7}
}

func (x *TelemetryRoundStarted) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryRoundStarted) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *TelemetryRoundStarted) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *TelemetryRoundStarted) GetLeader() uint64 {
	if x != nil {
		return x.Leader
	}
	return 0
}

func (x *TelemetryRoundStarted) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_cl_offchainreporting_telemetry_proto protoreflect.FileDescriptor

var file_cl_offchainreporting_telemetry_proto_rawDesc = []byte{
	0x0a, 0x24, 0x63, 0x6c, 0x5f, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x23, 0x63, 0x6c, 0x5f, 0x6f, 0x66,
	0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd1,
	0x03, 0x0a, 0x10, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x12, 0x57, 0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6f,
	0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x48, 0x00, 0x52, 0x0f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x12, 0x5a, 0x0a, 0x10,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x48, 0x00, 0x52, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x4b, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e,
	0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x60, 0x0a, 0x12, 0x61, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69,
	0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2e, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x41,
	0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x48, 0x00, 0x52, 0x12, 0x61, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69,
	0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x6f, 0x75, 0x6e, 0x64,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x48, 0x00, 0x52, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x42, 0x09, 0x0a, 0x07, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x18, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x12,
	0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67,
	0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x22, 0x9a, 0x01, 0x0a, 0x19, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x22,
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65,
	0x73, 0x74, 0x12, 0x33, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x22, 0xb1, 0x01,
	0x0a, 0x14, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x24, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x73, 0x67,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x65, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x22, 0x92, 0x02, 0x0a, 0x1b, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x41,
	0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x6c, 0x0a, 0x10, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3e, 0x2e, 0x6f, 0x66,
	0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69,
	0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x48, 0x00, 0x52, 0x10, 0x69,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x78, 0x0a, 0x14, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x42, 0x2e,
	0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72,
	0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x48, 0x00, 0x52, 0x14, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x76, 0x69, 0x6f,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x9e, 0x01, 0x0a, 0x2b, 0x54, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x93, 0x01, 0x0a, 0x2f, 0x54, 0x65, 0x6c, 0x65,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69,
	0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12,
	0x24, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x65, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x93, 0x01,
	0x0a, 0x15, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x6f, 0x75, 0x6e, 0x64,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63,
	0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cl_offchainreporting_telemetry_proto_rawDescOnce sync.Once
	file_cl_offchainreporting_telemetry_proto_rawDescData = file_cl_offchainreporting_telemetry_proto_rawDesc
)

func file_cl_offchainreporting_telemetry_proto_rawDescGZIP() []byte {
	file_cl_offchainreporting_telemetry_proto_rawDescOnce.Do(func() {
		file_cl_offchainreporting_telemetry_proto_rawDescData = protoimpl.X.CompressGZIP(file_cl_offchainreporting_telemetry_proto_rawDescData)
	})
	return file_cl_offchainreporting_telemetry_proto_rawDescData
}

var file_cl_offchainreporting_telemetry_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_cl_offchainreporting_telemetry_proto_goTypes = []interface{}{
	(*TelemetryWrapper)(nil),                                // 0: offchainreporting.TelemetryWrapper
	(*TelemetryMessageReceived)(nil),                        // 1: offchainreporting.TelemetryMessageReceived
	(*TelemetryMessageBroadcast)(nil),                       // 2: offchainreporting.TelemetryMessageBroadcast
	(*TelemetryMessageSent)(nil),                            // 3: offchainreporting.TelemetryMessageSent
	(*TelemetryAssertionViolation)(nil),                     // 4: offchainreporting.TelemetryAssertionViolation
	(*TelemetryAssertionViolationInvalidSignature)(nil),     // 5: offchainreporting.TelemetryAssertionViolationInvalidSignature
	(*TelemetryAssertionViolationInvalidSerialization)(nil), // 6: offchainreporting.TelemetryAssertionViolationInvalidSerialization
	(*TelemetryRoundStarted)(nil),                           // 7: offchainreporting.TelemetryRoundStarted
	(*MessageWrapper)(nil),                                  // 8: offchainreporting.MessageWrapper
}
var file_cl_offchainreporting_telemetry_proto_depIdxs = []int32{
	1,  // 0: offchainreporting.TelemetryWrapper.messageReceived:type_name -> offchainreporting.TelemetryMessageReceived
	2,  // 1: offchainreporting.TelemetryWrapper.messageBroadcast:type_name -> offchainreporting.TelemetryMessageBroadcast
	3,  // 2: offchainreporting.TelemetryWrapper.messageSent:type_name -> offchainreporting.TelemetryMessageSent
	4,  // 3: offchainreporting.TelemetryWrapper.assertionViolation:type_name -> offchainreporting.TelemetryAssertionViolation
	7,  // 4: offchainreporting.TelemetryWrapper.roundStarted:type_name -> offchainreporting.TelemetryRoundStarted
	8,  // 5: offchainreporting.TelemetryMessageReceived.msg:type_name -> offchainreporting.MessageWrapper
	8,  // 6: offchainreporting.TelemetryMessageBroadcast.msg:type_name -> offchainreporting.MessageWrapper
	8,  // 7: offchainreporting.TelemetryMessageSent.msg:type_name -> offchainreporting.MessageWrapper
	5,  // 8: offchainreporting.TelemetryAssertionViolation.invalidSignature:type_name -> offchainreporting.TelemetryAssertionViolationInvalidSignature
	6,  // 9: offchainreporting.TelemetryAssertionViolation.invalidSerialization:type_name -> offchainreporting.TelemetryAssertionViolationInvalidSerialization
	8,  // 10: offchainreporting.TelemetryAssertionViolationInvalidSignature.msg:type_name -> offchainreporting.MessageWrapper
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_cl_offchainreporting_telemetry_proto_init() }
func file_cl_offchainreporting_telemetry_proto_init() {
	if File_cl_offchainreporting_telemetry_proto != nil {
		return
	}
	file_cl_offchainreporting_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cl_offchainreporting_telemetry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryWrapper); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryMessageReceived); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryMessageBroadcast); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryMessageSent); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryAssertionViolation); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryAssertionViolationInvalidSignature); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryAssertionViolationInvalidSerialization); i {
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
		file_cl_offchainreporting_telemetry_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryRoundStarted); i {
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
	file_cl_offchainreporting_telemetry_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*TelemetryWrapper_MessageReceived)(nil),
		(*TelemetryWrapper_MessageBroadcast)(nil),
		(*TelemetryWrapper_MessageSent)(nil),
		(*TelemetryWrapper_AssertionViolation)(nil),
		(*TelemetryWrapper_RoundStarted)(nil),
	}
	file_cl_offchainreporting_telemetry_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*TelemetryAssertionViolation_InvalidSignature)(nil),
		(*TelemetryAssertionViolation_InvalidSerialization)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cl_offchainreporting_telemetry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cl_offchainreporting_telemetry_proto_goTypes,
		DependencyIndexes: file_cl_offchainreporting_telemetry_proto_depIdxs,
		MessageInfos:      file_cl_offchainreporting_telemetry_proto_msgTypes,
	}.Build()
	File_cl_offchainreporting_telemetry_proto = out.File
	file_cl_offchainreporting_telemetry_proto_rawDesc = nil
	file_cl_offchainreporting_telemetry_proto_goTypes = nil
	file_cl_offchainreporting_telemetry_proto_depIdxs = nil
}
