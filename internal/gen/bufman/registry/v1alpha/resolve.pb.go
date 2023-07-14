//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: bufman/registry/v1alpha/resolve.proto

package registryv1alpha

import (
	v1alpha "github.com/ProtobufMan/bufman/internal/gen/bufman/module/v1alpha"
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

type GetModulePinsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModuleReferences []*v1alpha.ModuleReference `protobuf:"bytes,1,rep,name=module_references,json=moduleReferences,proto3" json:"module_references,omitempty"`
	// current_module_pins allows for partial dependency updates by letting clients
	// send a request with the pins for their current module and only the
	// identities of the dependencies they want to update in module_references.
	//
	// When resolving, if a client supplied module pin is:
	// - in the transitive closure of pins resolved from the module_references,
	//   the client supplied module pin will be an extra candidate for tie
	//   breaking.
	// - NOT in the in the transitive closure of pins resolved from the
	//   module_references, it will be returned as is.
	CurrentModulePins []*v1alpha.ModulePin `protobuf:"bytes,2,rep,name=current_module_pins,json=currentModulePins,proto3" json:"current_module_pins,omitempty"`
}

func (x *GetModulePinsRequest) Reset() {
	*x = GetModulePinsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bufman_registry_v1alpha_resolve_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetModulePinsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetModulePinsRequest) ProtoMessage() {}

func (x *GetModulePinsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bufman_registry_v1alpha_resolve_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetModulePinsRequest.ProtoReflect.Descriptor instead.
func (*GetModulePinsRequest) Descriptor() ([]byte, []int) {
	return file_bufman_registry_v1alpha_resolve_proto_rawDescGZIP(), []int{0}
}

func (x *GetModulePinsRequest) GetModuleReferences() []*v1alpha.ModuleReference {
	if x != nil {
		return x.ModuleReferences
	}
	return nil
}

func (x *GetModulePinsRequest) GetCurrentModulePins() []*v1alpha.ModulePin {
	if x != nil {
		return x.CurrentModulePins
	}
	return nil
}

type GetModulePinsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModulePins []*v1alpha.ModulePin `protobuf:"bytes,1,rep,name=module_pins,json=modulePins,proto3" json:"module_pins,omitempty"`
}

func (x *GetModulePinsResponse) Reset() {
	*x = GetModulePinsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bufman_registry_v1alpha_resolve_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetModulePinsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetModulePinsResponse) ProtoMessage() {}

func (x *GetModulePinsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bufman_registry_v1alpha_resolve_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetModulePinsResponse.ProtoReflect.Descriptor instead.
func (*GetModulePinsResponse) Descriptor() ([]byte, []int) {
	return file_bufman_registry_v1alpha_resolve_proto_rawDescGZIP(), []int{1}
}

func (x *GetModulePinsResponse) GetModulePins() []*v1alpha.ModulePin {
	if x != nil {
		return x.ModulePins
	}
	return nil
}

var File_bufman_registry_v1alpha_resolve_proto protoreflect.FileDescriptor

var file_bufman_registry_v1alpha_resolve_proto_rawDesc = []byte{
	0x0a, 0x25, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x1a, 0x22, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x64, 0x75,
	0x6c, 0x65, 0x50, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x53, 0x0a,
	0x11, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61,
	0x6e, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x10, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x73, 0x12, 0x50, 0x0a, 0x13, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x5f, 0x70, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x69,
	0x6e, 0x52, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x50, 0x69, 0x6e, 0x73, 0x22, 0x5a, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x50, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x0b, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x70, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x6d, 0x6f, 0x64, 0x75,
	0x6c, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x50, 0x69, 0x6e, 0x52, 0x0a, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x69, 0x6e, 0x73,
	0x32, 0x85, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x73, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x50, 0x69, 0x6e, 0x73, 0x12, 0x2d, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x47,
	0x65, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x47, 0x65,
	0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x03, 0x90, 0x02, 0x01, 0x42, 0xfd, 0x01, 0x0a, 0x1b, 0x63, 0x6f, 0x6d,
	0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x42, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x4d, 0x61, 0x6e,
	0x2f, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x3b, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xa2, 0x02, 0x03, 0x42,
	0x52, 0x58, 0xaa, 0x02, 0x17, 0x42, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xca, 0x02, 0x17, 0x42,
	0x75, 0x66, 0x6d, 0x61, 0x6e, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xe2, 0x02, 0x23, 0x42, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x5c,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x19, 0x42,
	0x75, 0x66, 0x6d, 0x61, 0x6e, 0x3a, 0x3a, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x3a,
	0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bufman_registry_v1alpha_resolve_proto_rawDescOnce sync.Once
	file_bufman_registry_v1alpha_resolve_proto_rawDescData = file_bufman_registry_v1alpha_resolve_proto_rawDesc
)

func file_bufman_registry_v1alpha_resolve_proto_rawDescGZIP() []byte {
	file_bufman_registry_v1alpha_resolve_proto_rawDescOnce.Do(func() {
		file_bufman_registry_v1alpha_resolve_proto_rawDescData = protoimpl.X.CompressGZIP(file_bufman_registry_v1alpha_resolve_proto_rawDescData)
	})
	return file_bufman_registry_v1alpha_resolve_proto_rawDescData
}

var file_bufman_registry_v1alpha_resolve_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_bufman_registry_v1alpha_resolve_proto_goTypes = []interface{}{
	(*GetModulePinsRequest)(nil),    // 0: bufman.registry.v1alpha.GetModulePinsRequest
	(*GetModulePinsResponse)(nil),   // 1: bufman.registry.v1alpha.GetModulePinsResponse
	(*v1alpha.ModuleReference)(nil), // 2: bufman.module.v1alpha.ModuleReference
	(*v1alpha.ModulePin)(nil),       // 3: bufman.module.v1alpha.ModulePin
}
var file_bufman_registry_v1alpha_resolve_proto_depIdxs = []int32{
	2, // 0: bufman.registry.v1alpha.GetModulePinsRequest.module_references:type_name -> bufman.module.v1alpha.ModuleReference
	3, // 1: bufman.registry.v1alpha.GetModulePinsRequest.current_module_pins:type_name -> bufman.module.v1alpha.ModulePin
	3, // 2: bufman.registry.v1alpha.GetModulePinsResponse.module_pins:type_name -> bufman.module.v1alpha.ModulePin
	0, // 3: bufman.registry.v1alpha.ResolveService.GetModulePins:input_type -> bufman.registry.v1alpha.GetModulePinsRequest
	1, // 4: bufman.registry.v1alpha.ResolveService.GetModulePins:output_type -> bufman.registry.v1alpha.GetModulePinsResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_bufman_registry_v1alpha_resolve_proto_init() }
func file_bufman_registry_v1alpha_resolve_proto_init() {
	if File_bufman_registry_v1alpha_resolve_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bufman_registry_v1alpha_resolve_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetModulePinsRequest); i {
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
		file_bufman_registry_v1alpha_resolve_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetModulePinsResponse); i {
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
			RawDescriptor: file_bufman_registry_v1alpha_resolve_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_bufman_registry_v1alpha_resolve_proto_goTypes,
		DependencyIndexes: file_bufman_registry_v1alpha_resolve_proto_depIdxs,
		MessageInfos:      file_bufman_registry_v1alpha_resolve_proto_msgTypes,
	}.Build()
	File_bufman_registry_v1alpha_resolve_proto = out.File
	file_bufman_registry_v1alpha_resolve_proto_rawDesc = nil
	file_bufman_registry_v1alpha_resolve_proto_goTypes = nil
	file_bufman_registry_v1alpha_resolve_proto_depIdxs = nil
}