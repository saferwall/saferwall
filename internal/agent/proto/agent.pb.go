// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: agent.proto

package agent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// PingReply contains the version of the server running and some guest info.
type PingReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The server version.
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// JSON containing the guest system information like OS, hardware, ...
	Sysinfo []byte `protobuf:"bytes,2,opt,name=sysinfo,proto3" json:"sysinfo,omitempty"`
}

func (x *PingReply) Reset() {
	*x = PingReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingReply) ProtoMessage() {}

func (x *PingReply) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingReply.ProtoReflect.Descriptor instead.
func (*PingReply) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{0}
}

func (x *PingReply) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PingReply) GetSysinfo() []byte {
	if x != nil {
		return x.Sysinfo
	}
	return nil
}

// DeployRequest message contains a zip package that includes all
// necessery files.
type DeployRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Zip file containing the sandbox app with its dependencies.
	Package []byte `protobuf:"bytes,1,opt,name=package,proto3" json:"package,omitempty"`
	// Destination path where to deploy the package.
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *DeployRequest) Reset() {
	*x = DeployRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeployRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployRequest) ProtoMessage() {}

func (x *DeployRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployRequest.ProtoReflect.Descriptor instead.
func (*DeployRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{1}
}

func (x *DeployRequest) GetPackage() []byte {
	if x != nil {
		return x.Package
	}
	return nil
}

func (x *DeployRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// DeployReply contains the version of the package that was deployed.
type DeployReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The sandbox version.
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *DeployReply) Reset() {
	*x = DeployReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeployReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployReply) ProtoMessage() {}

func (x *DeployReply) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployReply.ProtoReflect.Descriptor instead.
func (*DeployReply) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{2}
}

func (x *DeployReply) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// The request message containing a sample for analysis.
type AnalyzeFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The binary file content.
	Binary []byte `protobuf:"bytes,1,opt,name=binary,proto3" json:"binary,omitempty"`
	// Configuration used to run the binary.
	// This is basically a JSON serialized byte array that contains
	// the configuration used to run the malware.
	// Example of fields that it contains is:
	//   - Full path to where the binary should be dropped in the guest.
	//   - Arguments used to execute the binary.
	//   - Timeout in seconds for how long to run the binary.
	//   - Country used to tunnel the connections.
	//   - The SHA256 hash of the binary.
	//   - Whether mem dumps should be taken.
	//   - etc ...
	Config []byte `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *AnalyzeFileRequest) Reset() {
	*x = AnalyzeFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFileRequest) ProtoMessage() {}

func (x *AnalyzeFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeFileRequest.ProtoReflect.Descriptor instead.
func (*AnalyzeFileRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{3}
}

func (x *AnalyzeFileRequest) GetBinary() []byte {
	if x != nil {
		return x.Binary
	}
	return nil
}

func (x *AnalyzeFileRequest) GetConfig() []byte {
	if x != nil {
		return x.Config
	}
	return nil
}

// The response message containing the analysis results.
type AnalyzeFileReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// APIs log produced by the sandbox in JSON Lines;
	APITrace    []byte                         `protobuf:"bytes,1,opt,name=api_trace,json=apiTrace,proto3" json:"api_trace,omitempty"`
	Screenshots []*AnalyzeFileReply_Screenshot `protobuf:"bytes,2,rep,name=screenshots,proto3" json:"screenshots,omitempty"`
	Artifacts   []*AnalyzeFileReply_Artifact   `protobuf:"bytes,3,rep,name=artifacts,proto3" json:"artifacts,omitempty"`
	APIBuffers  []*AnalyzeFileReply_APIBuffer  `protobuf:"bytes,4,rep,name=api_buffers,json=apiBuffers,proto3" json:"api_buffers,omitempty"`
	// Agent server log.
	ServerLog []byte `protobuf:"bytes,5,opt,name=server_log,json=serverLog,proto3" json:"server_log,omitempty"`
	// Controller log.
	ControllerLog []byte `protobuf:"bytes,6,opt,name=controller_log,json=controllerLog,proto3" json:"controller_log,omitempty"`
	// Process tree data.
	ProcessTree []byte `protobuf:"bytes,7,opt,name=process_tree,json=processTree,proto3" json:"process_tree,omitempty"`
}

func (x *AnalyzeFileReply) Reset() {
	*x = AnalyzeFileReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeFileReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFileReply) ProtoMessage() {}

func (x *AnalyzeFileReply) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeFileReply.ProtoReflect.Descriptor instead.
func (*AnalyzeFileReply) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4}
}

func (x *AnalyzeFileReply) GetAPITrace() []byte {
	if x != nil {
		return x.APITrace
	}
	return nil
}

func (x *AnalyzeFileReply) GetScreenshots() []*AnalyzeFileReply_Screenshot {
	if x != nil {
		return x.Screenshots
	}
	return nil
}

func (x *AnalyzeFileReply) GetArtifacts() []*AnalyzeFileReply_Artifact {
	if x != nil {
		return x.Artifacts
	}
	return nil
}

func (x *AnalyzeFileReply) GetAPIBuffers() []*AnalyzeFileReply_APIBuffer {
	if x != nil {
		return x.APIBuffers
	}
	return nil
}

func (x *AnalyzeFileReply) GetServerLog() []byte {
	if x != nil {
		return x.ServerLog
	}
	return nil
}

func (x *AnalyzeFileReply) GetControllerLog() []byte {
	if x != nil {
		return x.ControllerLog
	}
	return nil
}

func (x *AnalyzeFileReply) GetProcessTree() []byte {
	if x != nil {
		return x.ProcessTree
	}
	return nil
}

// Screenshots collected during the analysis.
type AnalyzeFileReply_Screenshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id represents an identifier to keep screenshots order.
	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// content contains the image data.
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AnalyzeFileReply_Screenshot) Reset() {
	*x = AnalyzeFileReply_Screenshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeFileReply_Screenshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFileReply_Screenshot) ProtoMessage() {}

func (x *AnalyzeFileReply_Screenshot) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeFileReply_Screenshot.ProtoReflect.Descriptor instead.
func (*AnalyzeFileReply_Screenshot) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4, 0}
}

func (x *AnalyzeFileReply_Screenshot) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AnalyzeFileReply_Screenshot) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

// Artifacts created files or memory dumps.
type AnalyzeFileReply_Artifact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the artifact.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The binary content of the artifact.
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AnalyzeFileReply_Artifact) Reset() {
	*x = AnalyzeFileReply_Artifact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeFileReply_Artifact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFileReply_Artifact) ProtoMessage() {}

func (x *AnalyzeFileReply_Artifact) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeFileReply_Artifact.ProtoReflect.Descriptor instead.
func (*AnalyzeFileReply_Artifact) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4, 1}
}

func (x *AnalyzeFileReply_Artifact) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AnalyzeFileReply_Artifact) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

// API Buffers represents the buffers for parameters of type BYTE*
// that are larger than 4KB.
type AnalyzeFileReply_APIBuffer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the API buffer file.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The binary content of the buffer.
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AnalyzeFileReply_APIBuffer) Reset() {
	*x = AnalyzeFileReply_APIBuffer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeFileReply_APIBuffer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeFileReply_APIBuffer) ProtoMessage() {}

func (x *AnalyzeFileReply_APIBuffer) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeFileReply_APIBuffer.ProtoReflect.Descriptor instead.
func (*AnalyzeFileReply_APIBuffer) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4, 2}
}

func (x *AnalyzeFileReply_APIBuffer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AnalyzeFileReply_APIBuffer) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_agent_proto protoreflect.FileDescriptor

var file_agent_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73,
	0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x09, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79,
	0x73, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x73, 0x79, 0x73,
	0x69, 0x6e, 0x66, 0x6f, 0x22, 0x3d, 0x0a, 0x0d, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x22, 0x27, 0x0a, 0x0b, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x44, 0x0a, 0x12,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x22, 0x95, 0x04, 0x0a, 0x10, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x70, 0x69, 0x5f, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x70, 0x69, 0x54,
	0x72, 0x61, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68,
	0x6f, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x61, 0x6e, 0x64,
	0x62, 0x6f, 0x78, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x2e, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x52,
	0x0b, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x12, 0x40, 0x0a, 0x09,
	0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x66,
	0x61, 0x63, 0x74, 0x52, 0x09, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73, 0x12, 0x44,
	0x0a, 0x0b, 0x61, 0x70, 0x69, 0x5f, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x6e,
	0x61, 0x6c, 0x79, 0x7a, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x41,
	0x50, 0x49, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x42, 0x75, 0x66,
	0x66, 0x65, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x6c,
	0x6f, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x4c, 0x6f, 0x67, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65,
	0x72, 0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x72, 0x65, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x54, 0x72, 0x65, 0x65, 0x1a, 0x36, 0x0a,
	0x0a, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x38, 0x0a, 0x08, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a,
	0x39, 0x0a, 0x09, 0x41, 0x50, 0x49, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0xbc, 0x01, 0x0a, 0x05, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x06, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x12, 0x16, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x73,
	0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x07, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x12,
	0x1b, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73,
	0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x61, 0x66, 0x65, 0x72, 0x77, 0x61, 0x6c,
	0x6c, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_agent_proto_rawDescOnce sync.Once
	file_agent_proto_rawDescData = file_agent_proto_rawDesc
)

func file_agent_proto_rawDescGZIP() []byte {
	file_agent_proto_rawDescOnce.Do(func() {
		file_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_agent_proto_rawDescData)
	})
	return file_agent_proto_rawDescData
}

var file_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_agent_proto_goTypes = []interface{}{
	(*PingReply)(nil),                   // 0: sandbox.PingReply
	(*DeployRequest)(nil),               // 1: sandbox.DeployRequest
	(*DeployReply)(nil),                 // 2: sandbox.DeployReply
	(*AnalyzeFileRequest)(nil),          // 3: sandbox.AnalyzeFileRequest
	(*AnalyzeFileReply)(nil),            // 4: sandbox.AnalyzeFileReply
	(*AnalyzeFileReply_Screenshot)(nil), // 5: sandbox.AnalyzeFileReply.Screenshot
	(*AnalyzeFileReply_Artifact)(nil),   // 6: sandbox.AnalyzeFileReply.Artifact
	(*AnalyzeFileReply_APIBuffer)(nil),  // 7: sandbox.AnalyzeFileReply.APIBuffer
	(*emptypb.Empty)(nil),               // 8: google.protobuf.Empty
}
var file_agent_proto_depIdxs = []int32{
	5, // 0: sandbox.AnalyzeFileReply.screenshots:type_name -> sandbox.AnalyzeFileReply.Screenshot
	6, // 1: sandbox.AnalyzeFileReply.artifacts:type_name -> sandbox.AnalyzeFileReply.Artifact
	7, // 2: sandbox.AnalyzeFileReply.api_buffers:type_name -> sandbox.AnalyzeFileReply.APIBuffer
	8, // 3: sandbox.Agent.Ping:input_type -> google.protobuf.Empty
	1, // 4: sandbox.Agent.Deploy:input_type -> sandbox.DeployRequest
	3, // 5: sandbox.Agent.Analyze:input_type -> sandbox.AnalyzeFileRequest
	0, // 6: sandbox.Agent.Ping:output_type -> sandbox.PingReply
	2, // 7: sandbox.Agent.Deploy:output_type -> sandbox.DeployReply
	4, // 8: sandbox.Agent.Analyze:output_type -> sandbox.AnalyzeFileReply
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_agent_proto_init() }
func file_agent_proto_init() {
	if File_agent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_agent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingReply); i {
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
		file_agent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeployRequest); i {
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
		file_agent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeployReply); i {
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
		file_agent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeFileRequest); i {
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
		file_agent_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeFileReply); i {
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
		file_agent_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeFileReply_Screenshot); i {
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
		file_agent_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeFileReply_Artifact); i {
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
		file_agent_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeFileReply_APIBuffer); i {
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
			RawDescriptor: file_agent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_proto_goTypes,
		DependencyIndexes: file_agent_proto_depIdxs,
		MessageInfos:      file_agent_proto_msgTypes,
	}.Build()
	File_agent_proto = out.File
	file_agent_proto_rawDesc = nil
	file_agent_proto_goTypes = nil
	file_agent_proto_depIdxs = nil
}
