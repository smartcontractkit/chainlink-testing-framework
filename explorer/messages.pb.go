package explorer

// original pb https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting/internal/serialization/protobuf/cl_offchainreporting_messages.pb.go

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

type MessageNewEpoch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch uint64 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *MessageNewEpoch) Reset() {
	*x = MessageNewEpoch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageNewEpoch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageNewEpoch) ProtoMessage() {}

func (x *MessageNewEpoch) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageNewEpoch.ProtoReflect.Descriptor instead.
func (*MessageNewEpoch) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{0}
}

func (x *MessageNewEpoch) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type MessageObserveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Epoch uint64 `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *MessageObserveReq) Reset() {
	*x = MessageObserveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageObserveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageObserveReq) ProtoMessage() {}

func (x *MessageObserveReq) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageObserveReq.ProtoReflect.Descriptor instead.
func (*MessageObserveReq) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{1}
}

func (x *MessageObserveReq) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageObserveReq) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type Observation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Observation) Reset() {
	*x = Observation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Observation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Observation) ProtoMessage() {}

func (x *Observation) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Observation.ProtoReflect.Descriptor instead.
func (*Observation) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{2}
}

func (x *Observation) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type SignedObservation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Observation *Observation `protobuf:"bytes,1,opt,name=observation,proto3" json:"observation,omitempty"`
	Signature   []byte       `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *SignedObservation) Reset() {
	*x = SignedObservation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedObservation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedObservation) ProtoMessage() {}

func (x *SignedObservation) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedObservation.ProtoReflect.Descriptor instead.
func (*SignedObservation) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{3}
}

func (x *SignedObservation) GetObservation() *Observation {
	if x != nil {
		return x.Observation
	}
	return nil
}

func (x *SignedObservation) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type AttributedSignedObservation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SignedObservation *SignedObservation `protobuf:"bytes,1,opt,name=signedObservation,proto3" json:"signedObservation,omitempty"`
	Observer          uint32             `protobuf:"varint,2,opt,name=observer,proto3" json:"observer,omitempty"`
}

func (x *AttributedSignedObservation) Reset() {
	*x = AttributedSignedObservation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttributedSignedObservation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributedSignedObservation) ProtoMessage() {}

func (x *AttributedSignedObservation) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributedSignedObservation.ProtoReflect.Descriptor instead.
func (*AttributedSignedObservation) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{4}
}

func (x *AttributedSignedObservation) GetSignedObservation() *SignedObservation {
	if x != nil {
		return x.SignedObservation
	}
	return nil
}

func (x *AttributedSignedObservation) GetObserver() uint32 {
	if x != nil {
		return x.Observer
	}
	return 0
}

type MessageObserve struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch             uint64             `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round             uint64             `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	SignedObservation *SignedObservation `protobuf:"bytes,3,opt,name=signedObservation,proto3" json:"signedObservation,omitempty"`
}

func (x *MessageObserve) Reset() {
	*x = MessageObserve{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageObserve) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageObserve) ProtoMessage() {}

func (x *MessageObserve) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageObserve.ProtoReflect.Descriptor instead.
func (*MessageObserve) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{5}
}

func (x *MessageObserve) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageObserve) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageObserve) GetSignedObservation() *SignedObservation {
	if x != nil {
		return x.SignedObservation
	}
	return nil
}

type MessageReportReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch                        uint64                         `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round                        uint64                         `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	AttributedSignedObservations []*AttributedSignedObservation `protobuf:"bytes,3,rep,name=attributedSignedObservations,proto3" json:"attributedSignedObservations,omitempty"`
}

func (x *MessageReportReq) Reset() {
	*x = MessageReportReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReportReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReportReq) ProtoMessage() {}

func (x *MessageReportReq) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReportReq.ProtoReflect.Descriptor instead.
func (*MessageReportReq) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{6}
}

func (x *MessageReportReq) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageReportReq) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageReportReq) GetAttributedSignedObservations() []*AttributedSignedObservation {
	if x != nil {
		return x.AttributedSignedObservations
	}
	return nil
}

type AttributedObservation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Observation *Observation `protobuf:"bytes,1,opt,name=observation,proto3" json:"observation,omitempty"`
	Observer    uint32       `protobuf:"varint,2,opt,name=observer,proto3" json:"observer,omitempty"`
}

func (x *AttributedObservation) Reset() {
	*x = AttributedObservation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttributedObservation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributedObservation) ProtoMessage() {}

func (x *AttributedObservation) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributedObservation.ProtoReflect.Descriptor instead.
func (*AttributedObservation) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{7}
}

func (x *AttributedObservation) GetObservation() *Observation {
	if x != nil {
		return x.Observation
	}
	return nil
}

func (x *AttributedObservation) GetObserver() uint32 {
	if x != nil {
		return x.Observer
	}
	return 0
}

type AttestedReportOne struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AttributedObservations []*AttributedObservation `protobuf:"bytes,1,rep,name=attributedObservations,proto3" json:"attributedObservations,omitempty"`
	Signature              []byte                   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *AttestedReportOne) Reset() {
	*x = AttestedReportOne{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttestedReportOne) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttestedReportOne) ProtoMessage() {}

func (x *AttestedReportOne) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttestedReportOne.ProtoReflect.Descriptor instead.
func (*AttestedReportOne) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{8}
}

func (x *AttestedReportOne) GetAttributedObservations() []*AttributedObservation {
	if x != nil {
		return x.AttributedObservations
	}
	return nil
}

func (x *AttestedReportOne) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type MessageReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch  uint64             `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round  uint64             `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Report *AttestedReportOne `protobuf:"bytes,3,opt,name=report,proto3" json:"report,omitempty"`
}

func (x *MessageReport) Reset() {
	*x = MessageReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReport) ProtoMessage() {}

func (x *MessageReport) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReport.ProtoReflect.Descriptor instead.
func (*MessageReport) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{9}
}

func (x *MessageReport) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageReport) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageReport) GetReport() *AttestedReportOne {
	if x != nil {
		return x.Report
	}
	return nil
}

type AttestedReportMany struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AttributedObservations []*AttributedObservation `protobuf:"bytes,1,rep,name=attributedObservations,proto3" json:"attributedObservations,omitempty"`
	Signatures             [][]byte                 `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *AttestedReportMany) Reset() {
	*x = AttestedReportMany{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttestedReportMany) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttestedReportMany) ProtoMessage() {}

func (x *AttestedReportMany) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttestedReportMany.ProtoReflect.Descriptor instead.
func (*AttestedReportMany) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{10}
}

func (x *AttestedReportMany) GetAttributedObservations() []*AttributedObservation {
	if x != nil {
		return x.AttributedObservations
	}
	return nil
}

func (x *AttestedReportMany) GetSignatures() [][]byte {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type MessageFinal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch  uint64              `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round  uint64              `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Report *AttestedReportMany `protobuf:"bytes,3,opt,name=report,proto3" json:"report,omitempty"`
}

func (x *MessageFinal) Reset() {
	*x = MessageFinal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageFinal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageFinal) ProtoMessage() {}

func (x *MessageFinal) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageFinal.ProtoReflect.Descriptor instead.
func (*MessageFinal) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{11}
}

func (x *MessageFinal) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageFinal) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageFinal) GetReport() *AttestedReportMany {
	if x != nil {
		return x.Report
	}
	return nil
}

// MessageFinalEcho has exactly the same fields as MessageFinal
type MessageFinalEcho struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Final *MessageFinal `protobuf:"bytes,1,opt,name=final,proto3" json:"final,omitempty"`
}

func (x *MessageFinalEcho) Reset() {
	*x = MessageFinalEcho{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageFinalEcho) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageFinalEcho) ProtoMessage() {}

func (x *MessageFinalEcho) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageFinalEcho.ProtoReflect.Descriptor instead.
func (*MessageFinalEcho) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{12}
}

func (x *MessageFinalEcho) GetFinal() *MessageFinal {
	if x != nil {
		return x.Final
	}
	return nil
}

type MessageWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//  *MessageWrapper_MessageNewEpoch
	//  *MessageWrapper_MessageObserveReq
	//  *MessageWrapper_MessageObserve
	//  *MessageWrapper_MessageReportReq
	//  *MessageWrapper_MessageReport
	//  *MessageWrapper_MessageFinal
	//  *MessageWrapper_MessageFinalEcho
	Msg isMessageWrapper_Msg `protobuf_oneof:"msg"`
}

func (x *MessageWrapper) Reset() {
	*x = MessageWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cl_offchainreporting_messages_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageWrapper) ProtoMessage() {}

func (x *MessageWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_cl_offchainreporting_messages_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageWrapper.ProtoReflect.Descriptor instead.
func (*MessageWrapper) Descriptor() ([]byte, []int) {
	return file_cl_offchainreporting_messages_proto_rawDescGZIP(), []int{13}
}

func (m *MessageWrapper) GetMsg() isMessageWrapper_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *MessageWrapper) GetMessageNewEpoch() *MessageNewEpoch {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageNewEpoch); ok {
		return x.MessageNewEpoch
	}
	return nil
}

func (x *MessageWrapper) GetMessageObserveReq() *MessageObserveReq {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageObserveReq); ok {
		return x.MessageObserveReq
	}
	return nil
}

func (x *MessageWrapper) GetMessageObserve() *MessageObserve {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageObserve); ok {
		return x.MessageObserve
	}
	return nil
}

func (x *MessageWrapper) GetMessageReportReq() *MessageReportReq {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageReportReq); ok {
		return x.MessageReportReq
	}
	return nil
}

func (x *MessageWrapper) GetMessageReport() *MessageReport {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageReport); ok {
		return x.MessageReport
	}
	return nil
}

func (x *MessageWrapper) GetMessageFinal() *MessageFinal {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageFinal); ok {
		return x.MessageFinal
	}
	return nil
}

func (x *MessageWrapper) GetMessageFinalEcho() *MessageFinalEcho {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageFinalEcho); ok {
		return x.MessageFinalEcho
	}
	return nil
}

type isMessageWrapper_Msg interface {
	isMessageWrapper_Msg()
}

type MessageWrapper_MessageNewEpoch struct {
	MessageNewEpoch *MessageNewEpoch `protobuf:"bytes,2,opt,name=messageNewEpoch,proto3,oneof"`
}

type MessageWrapper_MessageObserveReq struct {
	MessageObserveReq *MessageObserveReq `protobuf:"bytes,3,opt,name=messageObserveReq,proto3,oneof"`
}

type MessageWrapper_MessageObserve struct {
	MessageObserve *MessageObserve `protobuf:"bytes,4,opt,name=messageObserve,proto3,oneof"`
}

type MessageWrapper_MessageReportReq struct {
	MessageReportReq *MessageReportReq `protobuf:"bytes,5,opt,name=messageReportReq,proto3,oneof"`
}

type MessageWrapper_MessageReport struct {
	MessageReport *MessageReport `protobuf:"bytes,6,opt,name=messageReport,proto3,oneof"`
}

type MessageWrapper_MessageFinal struct {
	MessageFinal *MessageFinal `protobuf:"bytes,7,opt,name=messageFinal,proto3,oneof"`
}

type MessageWrapper_MessageFinalEcho struct {
	MessageFinalEcho *MessageFinalEcho `protobuf:"bytes,8,opt,name=messageFinalEcho,proto3,oneof"`
}

func (*MessageWrapper_MessageNewEpoch) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageObserveReq) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageObserve) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageReportReq) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageReport) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageFinal) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageFinalEcho) isMessageWrapper_Msg() {}

var File_cl_offchainreporting_messages_proto protoreflect.FileDescriptor

var file_cl_offchainreporting_messages_proto_rawDesc = []byte{
	0x0a, 0x23, 0x63, 0x6c, 0x5f, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x27, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63,
	0x68, 0x22, 0x3f, 0x0a, 0x11, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f,
	0x63, 0x68, 0x22, 0x23, 0x0a, 0x0b, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x73, 0x0a, 0x11, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x0b,
	0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0b, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x8d, 0x01, 0x0a,
	0x1b, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x52, 0x0a, 0x11,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x73,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x90, 0x01, 0x0a,
	0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x52, 0x0a, 0x11, 0x73,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0xb2, 0x01, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x72, 0x0a, 0x1c, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x1c, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x22, 0x75, 0x0a, 0x15, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a,
	0x0b, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x93, 0x01, 0x0a, 0x11,
	0x41, 0x74, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x4f, 0x6e,
	0x65, 0x12, 0x60, 0x0a, 0x16, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x4f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64,
	0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x16, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x22, 0x79, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x3c,
	0x0a, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69,
	0x6e, 0x67, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x4f, 0x6e, 0x65, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x96, 0x01, 0x0a,
	0x12, 0x41, 0x74, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x4d,
	0x61, 0x6e, 0x79, 0x12, 0x60, 0x0a, 0x16, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x16, 0x61,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x79, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x12, 0x3d, 0x0a, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x22, 0x49, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x45, 0x63, 0x68, 0x6f, 0x12, 0x35, 0x0a, 0x05, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46,
	0x69, 0x6e, 0x61, 0x6c, 0x52, 0x05, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x22, 0xc7, 0x04, 0x0a, 0x0e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x4e,
	0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x48, 0x00, 0x52, 0x0f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x54,
	0x0a, 0x11, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x52, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x66, 0x66, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x52, 0x65, 0x71, 0x48,
	0x00, 0x52, 0x11, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x12, 0x4b, 0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6f,
	0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x48,
	0x00, 0x52, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x12, 0x51, 0x0a, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6f, 0x66,
	0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71,
	0x48, 0x00, 0x52, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x48, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6f, 0x66,
	0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x48, 0x00, 0x52,
	0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x45,
	0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x48, 0x00, 0x52, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x51, 0x0a, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x45, 0x63, 0x68, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x23, 0x2e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x45, 0x63, 0x68, 0x6f, 0x48, 0x00, 0x52, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46,
	0x69, 0x6e, 0x61, 0x6c, 0x45, 0x63, 0x68, 0x6f, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x4a,
	0x04, 0x08, 0x01, 0x10, 0x02, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cl_offchainreporting_messages_proto_rawDescOnce sync.Once
	file_cl_offchainreporting_messages_proto_rawDescData = file_cl_offchainreporting_messages_proto_rawDesc
)

func file_cl_offchainreporting_messages_proto_rawDescGZIP() []byte {
	file_cl_offchainreporting_messages_proto_rawDescOnce.Do(func() {
		file_cl_offchainreporting_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_cl_offchainreporting_messages_proto_rawDescData)
	})
	return file_cl_offchainreporting_messages_proto_rawDescData
}

var file_cl_offchainreporting_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_cl_offchainreporting_messages_proto_goTypes = []interface{}{
	(*MessageNewEpoch)(nil),             // 0: offchainreporting.MessageNewEpoch
	(*MessageObserveReq)(nil),           // 1: offchainreporting.MessageObserveReq
	(*Observation)(nil),                 // 2: offchainreporting.Observation
	(*SignedObservation)(nil),           // 3: offchainreporting.SignedObservation
	(*AttributedSignedObservation)(nil), // 4: offchainreporting.AttributedSignedObservation
	(*MessageObserve)(nil),              // 5: offchainreporting.MessageObserve
	(*MessageReportReq)(nil),            // 6: offchainreporting.MessageReportReq
	(*AttributedObservation)(nil),       // 7: offchainreporting.AttributedObservation
	(*AttestedReportOne)(nil),           // 8: offchainreporting.AttestedReportOne
	(*MessageReport)(nil),               // 9: offchainreporting.MessageReport
	(*AttestedReportMany)(nil),          // 10: offchainreporting.AttestedReportMany
	(*MessageFinal)(nil),                // 11: offchainreporting.MessageFinal
	(*MessageFinalEcho)(nil),            // 12: offchainreporting.MessageFinalEcho
	(*MessageWrapper)(nil),              // 13: offchainreporting.MessageWrapper
}
var file_cl_offchainreporting_messages_proto_depIdxs = []int32{
	2,  // 0: offchainreporting.SignedObservation.observation:type_name -> offchainreporting.Observation
	3,  // 1: offchainreporting.AttributedSignedObservation.signedObservation:type_name -> offchainreporting.SignedObservation
	3,  // 2: offchainreporting.MessageObserve.signedObservation:type_name -> offchainreporting.SignedObservation
	4,  // 3: offchainreporting.MessageReportReq.attributedSignedObservations:type_name -> offchainreporting.AttributedSignedObservation
	2,  // 4: offchainreporting.AttributedObservation.observation:type_name -> offchainreporting.Observation
	7,  // 5: offchainreporting.AttestedReportOne.attributedObservations:type_name -> offchainreporting.AttributedObservation
	8,  // 6: offchainreporting.MessageReport.report:type_name -> offchainreporting.AttestedReportOne
	7,  // 7: offchainreporting.AttestedReportMany.attributedObservations:type_name -> offchainreporting.AttributedObservation
	10, // 8: offchainreporting.MessageFinal.report:type_name -> offchainreporting.AttestedReportMany
	11, // 9: offchainreporting.MessageFinalEcho.final:type_name -> offchainreporting.MessageFinal
	0,  // 10: offchainreporting.MessageWrapper.messageNewEpoch:type_name -> offchainreporting.MessageNewEpoch
	1,  // 11: offchainreporting.MessageWrapper.messageObserveReq:type_name -> offchainreporting.MessageObserveReq
	5,  // 12: offchainreporting.MessageWrapper.messageObserve:type_name -> offchainreporting.MessageObserve
	6,  // 13: offchainreporting.MessageWrapper.messageReportReq:type_name -> offchainreporting.MessageReportReq
	9,  // 14: offchainreporting.MessageWrapper.messageReport:type_name -> offchainreporting.MessageReport
	11, // 15: offchainreporting.MessageWrapper.messageFinal:type_name -> offchainreporting.MessageFinal
	12, // 16: offchainreporting.MessageWrapper.messageFinalEcho:type_name -> offchainreporting.MessageFinalEcho
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_cl_offchainreporting_messages_proto_init() }
func file_cl_offchainreporting_messages_proto_init() {
	if File_cl_offchainreporting_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cl_offchainreporting_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageNewEpoch); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageObserveReq); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Observation); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedObservation); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttributedSignedObservation); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageObserve); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReportReq); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttributedObservation); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttestedReportOne); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReport); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttestedReportMany); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageFinal); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageFinalEcho); i {
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
		file_cl_offchainreporting_messages_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageWrapper); i {
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
	file_cl_offchainreporting_messages_proto_msgTypes[13].OneofWrappers = []interface{}{
		(*MessageWrapper_MessageNewEpoch)(nil),
		(*MessageWrapper_MessageObserveReq)(nil),
		(*MessageWrapper_MessageObserve)(nil),
		(*MessageWrapper_MessageReportReq)(nil),
		(*MessageWrapper_MessageReport)(nil),
		(*MessageWrapper_MessageFinal)(nil),
		(*MessageWrapper_MessageFinalEcho)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cl_offchainreporting_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cl_offchainreporting_messages_proto_goTypes,
		DependencyIndexes: file_cl_offchainreporting_messages_proto_depIdxs,
		MessageInfos:      file_cl_offchainreporting_messages_proto_msgTypes,
	}.Build()
	File_cl_offchainreporting_messages_proto = out.File
	file_cl_offchainreporting_messages_proto_rawDesc = nil
	file_cl_offchainreporting_messages_proto_goTypes = nil
	file_cl_offchainreporting_messages_proto_depIdxs = nil
}
