// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: discovery/message.proto

package discovery

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

// Message represents a message broadcasted by a peer.
type Message struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        []byte                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                                 // Peer ID.
	SessionId int64                  `protobuf:"varint,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"` // Session ID.
	Seqnum    int64                  `protobuf:"varint,3,opt,name=seqnum,proto3" json:"seqnum,omitempty"`                        // Sequence number (indicates that the message is the nth message in the session).
	Timestamp int64                  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                  // Timestamp of the message (unix timestamp in milliseconds).
	// Types that are valid to be assigned to Payload:
	//
	//	*Message_Announcement
	//	*Message_HelloRequest
	//	*Message_HelloResult
	Payload       isMessage_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_discovery_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_discovery_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_discovery_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Message) GetSessionId() int64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *Message) GetSeqnum() int64 {
	if x != nil {
		return x.Seqnum
	}
	return 0
}

func (x *Message) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Message) GetPayload() isMessage_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Message) GetAnnouncement() *Announcement {
	if x != nil {
		if x, ok := x.Payload.(*Message_Announcement); ok {
			return x.Announcement
		}
	}
	return nil
}

func (x *Message) GetHelloRequest() *HelloRequest {
	if x != nil {
		if x, ok := x.Payload.(*Message_HelloRequest); ok {
			return x.HelloRequest
		}
	}
	return nil
}

func (x *Message) GetHelloResult() *HelloResult {
	if x != nil {
		if x, ok := x.Payload.(*Message_HelloResult); ok {
			return x.HelloResult
		}
	}
	return nil
}

type isMessage_Payload interface {
	isMessage_Payload()
}

type Message_Announcement struct {
	Announcement *Announcement `protobuf:"bytes,10,opt,name=announcement,proto3,oneof"`
}

type Message_HelloRequest struct {
	HelloRequest *HelloRequest `protobuf:"bytes,11,opt,name=hello_request,json=helloRequest,proto3,oneof"`
}

type Message_HelloResult struct {
	HelloResult *HelloResult `protobuf:"bytes,12,opt,name=hello_result,json=helloResult,proto3,oneof"`
}

func (*Message_Announcement) isMessage_Payload() {}

func (*Message_HelloRequest) isMessage_Payload() {}

func (*Message_HelloResult) isMessage_Payload() {}

// Announces a peer.
type Announcement struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Peer          *Peer                  `protobuf:"bytes,1,opt,name=peer,proto3" json:"peer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Announcement) Reset() {
	*x = Announcement{}
	mi := &file_discovery_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Announcement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Announcement) ProtoMessage() {}

func (x *Announcement) ProtoReflect() protoreflect.Message {
	mi := &file_discovery_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Announcement.ProtoReflect.Descriptor instead.
func (*Announcement) Descriptor() ([]byte, []int) {
	return file_discovery_message_proto_rawDescGZIP(), []int{1}
}

func (x *Announcement) GetPeer() *Peer {
	if x != nil {
		return x.Peer
	}
	return nil
}

// HelloRequest is a request to say hello to a peer.
type HelloRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Peer          *Peer                  `protobuf:"bytes,1,opt,name=peer,proto3" json:"peer,omitempty"` // Information of the sender.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	mi := &file_discovery_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_discovery_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_discovery_message_proto_rawDescGZIP(), []int{2}
}

func (x *HelloRequest) GetPeer() *Peer {
	if x != nil {
		return x.Peer
	}
	return nil
}

// HelloResult is a result of saying hello to a peer.
type HelloResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Peer          *Peer                  `protobuf:"bytes,1,opt,name=peer,proto3" json:"peer,omitempty"` // Information of the receiver.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HelloResult) Reset() {
	*x = HelloResult{}
	mi := &file_discovery_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HelloResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResult) ProtoMessage() {}

func (x *HelloResult) ProtoReflect() protoreflect.Message {
	mi := &file_discovery_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResult.ProtoReflect.Descriptor instead.
func (*HelloResult) Descriptor() ([]byte, []int) {
	return file_discovery_message_proto_rawDescGZIP(), []int{3}
}

func (x *HelloResult) GetPeer() *Peer {
	if x != nil {
		return x.Peer
	}
	return nil
}

var File_discovery_message_proto protoreflect.FileDescriptor

var file_discovery_message_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x69, 0x6f, 0x2e, 0x69, 0x6e,
	0x68, 0x69, 0x62, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x1a, 0x14, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf1, 0x02, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x71, 0x6e, 0x75,
	0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x65, 0x71, 0x6e, 0x75, 0x6d, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x4f, 0x0a,
	0x0c, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x6f, 0x2e, 0x69, 0x6e, 0x68, 0x69, 0x62, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x0c, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x50,
	0x0a, 0x0d, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x6f, 0x2e, 0x69, 0x6e, 0x68, 0x69, 0x62,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x48, 0x00, 0x52, 0x0c, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x4d, 0x0a, 0x0c, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6f, 0x2e, 0x69, 0x6e, 0x68, 0x69,
	0x62, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x48, 0x00, 0x52, 0x0b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42,
	0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x4a, 0x04, 0x08, 0x05, 0x10, 0x0a,
	0x22, 0x45, 0x0a, 0x0c, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x35, 0x0a, 0x04, 0x70, 0x65, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x69, 0x6f, 0x2e, 0x69, 0x6e, 0x68, 0x69, 0x62, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f,
	0x72, 0x75, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x50, 0x65, 0x65,
	0x72, 0x52, 0x04, 0x70, 0x65, 0x65, 0x72, 0x22, 0x45, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x04, 0x70, 0x65, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6f, 0x2e, 0x69, 0x6e, 0x68, 0x69, 0x62,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x04, 0x70, 0x65, 0x65, 0x72, 0x22, 0x44,
	0x0a, 0x0b, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x35, 0x0a,
	0x04, 0x70, 0x65, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6f,
	0x2e, 0x69, 0x6e, 0x68, 0x69, 0x62, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x6d, 0x6f, 0x72, 0x75, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x04,
	0x70, 0x65, 0x65, 0x72, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x68, 0x69, 0x62, 0x69, 0x74, 0x6f, 0x72, 0x31, 0x32, 0x31, 0x37,
	0x2f, 0x6d, 0x6f, 0x72, 0x75, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_discovery_message_proto_rawDescOnce sync.Once
	file_discovery_message_proto_rawDescData []byte
)

func file_discovery_message_proto_rawDescGZIP() []byte {
	file_discovery_message_proto_rawDescOnce.Do(func() {
		file_discovery_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_discovery_message_proto_rawDesc), len(file_discovery_message_proto_rawDesc)))
	})
	return file_discovery_message_proto_rawDescData
}

var file_discovery_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_discovery_message_proto_goTypes = []any{
	(*Message)(nil),      // 0: io.inhibitor.moru.discovery.Message
	(*Announcement)(nil), // 1: io.inhibitor.moru.discovery.Announcement
	(*HelloRequest)(nil), // 2: io.inhibitor.moru.discovery.HelloRequest
	(*HelloResult)(nil),  // 3: io.inhibitor.moru.discovery.HelloResult
	(*Peer)(nil),         // 4: io.inhibitor.moru.discovery.Peer
}
var file_discovery_message_proto_depIdxs = []int32{
	1, // 0: io.inhibitor.moru.discovery.Message.announcement:type_name -> io.inhibitor.moru.discovery.Announcement
	2, // 1: io.inhibitor.moru.discovery.Message.hello_request:type_name -> io.inhibitor.moru.discovery.HelloRequest
	3, // 2: io.inhibitor.moru.discovery.Message.hello_result:type_name -> io.inhibitor.moru.discovery.HelloResult
	4, // 3: io.inhibitor.moru.discovery.Announcement.peer:type_name -> io.inhibitor.moru.discovery.Peer
	4, // 4: io.inhibitor.moru.discovery.HelloRequest.peer:type_name -> io.inhibitor.moru.discovery.Peer
	4, // 5: io.inhibitor.moru.discovery.HelloResult.peer:type_name -> io.inhibitor.moru.discovery.Peer
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_discovery_message_proto_init() }
func file_discovery_message_proto_init() {
	if File_discovery_message_proto != nil {
		return
	}
	file_discovery_peer_proto_init()
	file_discovery_message_proto_msgTypes[0].OneofWrappers = []any{
		(*Message_Announcement)(nil),
		(*Message_HelloRequest)(nil),
		(*Message_HelloResult)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_discovery_message_proto_rawDesc), len(file_discovery_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_discovery_message_proto_goTypes,
		DependencyIndexes: file_discovery_message_proto_depIdxs,
		MessageInfos:      file_discovery_message_proto_msgTypes,
	}.Build()
	File_discovery_message_proto = out.File
	file_discovery_message_proto_goTypes = nil
	file_discovery_message_proto_depIdxs = nil
}
