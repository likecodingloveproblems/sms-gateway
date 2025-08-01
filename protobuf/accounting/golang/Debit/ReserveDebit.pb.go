// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: protobuf/accounting/ReserveDebit.proto

package Debit

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReserveDebitRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint64                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        uint64                 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReserveDebitRequest) Reset() {
	*x = ReserveDebitRequest{}
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReserveDebitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReserveDebitRequest) ProtoMessage() {}

func (x *ReserveDebitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReserveDebitRequest.ProtoReflect.Descriptor instead.
func (*ReserveDebitRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_accounting_ReserveDebit_proto_rawDescGZIP(), []int{0}
}

func (x *ReserveDebitRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReserveDebitRequest) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type ReserveDebitResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IsSuccessful  bool                   `protobuf:"varint,1,opt,name=is_successful,json=isSuccessful,proto3" json:"is_successful,omitempty"`
	Reason        string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"` // optional: "user_not_found", "insufficient_balance", or error message
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReserveDebitResponse) Reset() {
	*x = ReserveDebitResponse{}
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReserveDebitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReserveDebitResponse) ProtoMessage() {}

func (x *ReserveDebitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReserveDebitResponse.ProtoReflect.Descriptor instead.
func (*ReserveDebitResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_accounting_ReserveDebit_proto_rawDescGZIP(), []int{1}
}

func (x *ReserveDebitResponse) GetIsSuccessful() bool {
	if x != nil {
		return x.IsSuccessful
	}
	return false
}

func (x *ReserveDebitResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type CancelDebitRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        uint64                 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        uint64                 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CancelDebitRequest) Reset() {
	*x = CancelDebitRequest{}
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelDebitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelDebitRequest) ProtoMessage() {}

func (x *CancelDebitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelDebitRequest.ProtoReflect.Descriptor instead.
func (*CancelDebitRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_accounting_ReserveDebit_proto_rawDescGZIP(), []int{2}
}

func (x *CancelDebitRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CancelDebitRequest) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type CancelDebitResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IsSuccessful  bool                   `protobuf:"varint,1,opt,name=is_successful,json=isSuccessful,proto3" json:"is_successful,omitempty"`
	Reason        string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"` // optional: "user_not_found", or error message
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CancelDebitResponse) Reset() {
	*x = CancelDebitResponse{}
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelDebitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelDebitResponse) ProtoMessage() {}

func (x *CancelDebitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_accounting_ReserveDebit_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelDebitResponse.ProtoReflect.Descriptor instead.
func (*CancelDebitResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_accounting_ReserveDebit_proto_rawDescGZIP(), []int{3}
}

func (x *CancelDebitResponse) GetIsSuccessful() bool {
	if x != nil {
		return x.IsSuccessful
	}
	return false
}

func (x *CancelDebitResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

var File_protobuf_accounting_ReserveDebit_proto protoreflect.FileDescriptor

const file_protobuf_accounting_ReserveDebit_proto_rawDesc = "" +
	"\n" +
	"&protobuf/accounting/ReserveDebit.proto\x12\n" +
	"accounting\"F\n" +
	"\x13ReserveDebitRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x04R\x06userId\x12\x16\n" +
	"\x06amount\x18\x02 \x01(\x04R\x06amount\"S\n" +
	"\x14ReserveDebitResponse\x12#\n" +
	"\ris_successful\x18\x01 \x01(\bR\fisSuccessful\x12\x16\n" +
	"\x06reason\x18\x02 \x01(\tR\x06reason\"E\n" +
	"\x12CancelDebitRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x04R\x06userId\x12\x16\n" +
	"\x06amount\x18\x02 \x01(\x04R\x06amount\"R\n" +
	"\x13CancelDebitResponse\x12#\n" +
	"\ris_successful\x18\x01 \x01(\bR\fisSuccessful\x12\x16\n" +
	"\x06reason\x18\x02 \x01(\tR\x06reason2\xaf\x01\n" +
	"\n" +
	"Accounting\x12Q\n" +
	"\fReserveDebit\x12\x1f.accounting.ReserveDebitRequest\x1a .accounting.ReserveDebitResponse\x12N\n" +
	"\vCancelDebit\x12\x1e.accounting.CancelDebitRequest\x1a\x1f.accounting.CancelDebitResponseB\"Z protobuf/accounting/golang/Debitb\x06proto3"

var (
	file_protobuf_accounting_ReserveDebit_proto_rawDescOnce sync.Once
	file_protobuf_accounting_ReserveDebit_proto_rawDescData []byte
)

func file_protobuf_accounting_ReserveDebit_proto_rawDescGZIP() []byte {
	file_protobuf_accounting_ReserveDebit_proto_rawDescOnce.Do(func() {
		file_protobuf_accounting_ReserveDebit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protobuf_accounting_ReserveDebit_proto_rawDesc), len(file_protobuf_accounting_ReserveDebit_proto_rawDesc)))
	})
	return file_protobuf_accounting_ReserveDebit_proto_rawDescData
}

var file_protobuf_accounting_ReserveDebit_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protobuf_accounting_ReserveDebit_proto_goTypes = []any{
	(*ReserveDebitRequest)(nil),  // 0: accounting.ReserveDebitRequest
	(*ReserveDebitResponse)(nil), // 1: accounting.ReserveDebitResponse
	(*CancelDebitRequest)(nil),   // 2: accounting.CancelDebitRequest
	(*CancelDebitResponse)(nil),  // 3: accounting.CancelDebitResponse
}
var file_protobuf_accounting_ReserveDebit_proto_depIdxs = []int32{
	0, // 0: accounting.Accounting.ReserveDebit:input_type -> accounting.ReserveDebitRequest
	2, // 1: accounting.Accounting.CancelDebit:input_type -> accounting.CancelDebitRequest
	1, // 2: accounting.Accounting.ReserveDebit:output_type -> accounting.ReserveDebitResponse
	3, // 3: accounting.Accounting.CancelDebit:output_type -> accounting.CancelDebitResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protobuf_accounting_ReserveDebit_proto_init() }
func file_protobuf_accounting_ReserveDebit_proto_init() {
	if File_protobuf_accounting_ReserveDebit_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protobuf_accounting_ReserveDebit_proto_rawDesc), len(file_protobuf_accounting_ReserveDebit_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_accounting_ReserveDebit_proto_goTypes,
		DependencyIndexes: file_protobuf_accounting_ReserveDebit_proto_depIdxs,
		MessageInfos:      file_protobuf_accounting_ReserveDebit_proto_msgTypes,
	}.Build()
	File_protobuf_accounting_ReserveDebit_proto = out.File
	file_protobuf_accounting_ReserveDebit_proto_goTypes = nil
	file_protobuf_accounting_ReserveDebit_proto_depIdxs = nil
}
