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
// source: bufman/registry/v1alpha/authn.proto

package registryv1alpha

import (
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

type GetCurrentUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetCurrentUserRequest) Reset() {
	*x = GetCurrentUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bufman_registry_v1alpha_authn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCurrentUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCurrentUserRequest) ProtoMessage() {}

func (x *GetCurrentUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bufman_registry_v1alpha_authn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCurrentUserRequest.ProtoReflect.Descriptor instead.
func (*GetCurrentUserRequest) Descriptor() ([]byte, []int) {
	return file_bufman_registry_v1alpha_authn_proto_rawDescGZIP(), []int{0}
}

type GetCurrentUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetCurrentUserResponse) Reset() {
	*x = GetCurrentUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bufman_registry_v1alpha_authn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCurrentUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCurrentUserResponse) ProtoMessage() {}

func (x *GetCurrentUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bufman_registry_v1alpha_authn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCurrentUserResponse.ProtoReflect.Descriptor instead.
func (*GetCurrentUserResponse) Descriptor() ([]byte, []int) {
	return file_bufman_registry_v1alpha_authn_proto_rawDescGZIP(), []int{1}
}

func (x *GetCurrentUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_bufman_registry_v1alpha_authn_proto protoreflect.FileDescriptor

var file_bufman_registry_v1alpha_authn_proto_rawDesc = []byte{
	0x0a, 0x23, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x1a, 0x22,
	0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4b, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x32, 0x86, 0x01, 0x0a, 0x0c, 0x41, 0x75, 0x74,
	0x68, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x76, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2e, 0x2e, 0x62, 0x75,
	0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x62, 0x75,
	0x66, 0x6d, 0x61, 0x6e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x03, 0x90, 0x02,
	0x01, 0x42, 0xfb, 0x01, 0x0a, 0x1b, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x42, 0x0a, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x4d, 0x61, 0x6e, 0x2f, 0x62, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x62, 0x75, 0x66, 0x6d,
	0x61, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x3b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0xa2, 0x02, 0x03, 0x42, 0x52, 0x58, 0xaa, 0x02, 0x17, 0x42, 0x75, 0x66, 0x6d,
	0x61, 0x6e, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0xca, 0x02, 0x17, 0x42, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x5c, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0xe2, 0x02, 0x23,
	0x42, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c,
	0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x19, 0x42, 0x75, 0x66, 0x6d, 0x61, 0x6e, 0x3a, 0x3a, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bufman_registry_v1alpha_authn_proto_rawDescOnce sync.Once
	file_bufman_registry_v1alpha_authn_proto_rawDescData = file_bufman_registry_v1alpha_authn_proto_rawDesc
)

func file_bufman_registry_v1alpha_authn_proto_rawDescGZIP() []byte {
	file_bufman_registry_v1alpha_authn_proto_rawDescOnce.Do(func() {
		file_bufman_registry_v1alpha_authn_proto_rawDescData = protoimpl.X.CompressGZIP(file_bufman_registry_v1alpha_authn_proto_rawDescData)
	})
	return file_bufman_registry_v1alpha_authn_proto_rawDescData
}

var file_bufman_registry_v1alpha_authn_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_bufman_registry_v1alpha_authn_proto_goTypes = []interface{}{
	(*GetCurrentUserRequest)(nil),  // 0: bufman.registry.v1alpha.GetCurrentUserRequest
	(*GetCurrentUserResponse)(nil), // 1: bufman.registry.v1alpha.GetCurrentUserResponse
	(*User)(nil),                   // 2: bufman.registry.v1alpha.User
}
var file_bufman_registry_v1alpha_authn_proto_depIdxs = []int32{
	2, // 0: bufman.registry.v1alpha.GetCurrentUserResponse.user:type_name -> bufman.registry.v1alpha.User
	0, // 1: bufman.registry.v1alpha.AuthnService.GetCurrentUser:input_type -> bufman.registry.v1alpha.GetCurrentUserRequest
	1, // 2: bufman.registry.v1alpha.AuthnService.GetCurrentUser:output_type -> bufman.registry.v1alpha.GetCurrentUserResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_bufman_registry_v1alpha_authn_proto_init() }
func file_bufman_registry_v1alpha_authn_proto_init() {
	if File_bufman_registry_v1alpha_authn_proto != nil {
		return
	}
	file_bufman_registry_v1alpha_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_bufman_registry_v1alpha_authn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCurrentUserRequest); i {
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
		file_bufman_registry_v1alpha_authn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCurrentUserResponse); i {
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
			RawDescriptor: file_bufman_registry_v1alpha_authn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_bufman_registry_v1alpha_authn_proto_goTypes,
		DependencyIndexes: file_bufman_registry_v1alpha_authn_proto_depIdxs,
		MessageInfos:      file_bufman_registry_v1alpha_authn_proto_msgTypes,
	}.Build()
	File_bufman_registry_v1alpha_authn_proto = out.File
	file_bufman_registry_v1alpha_authn_proto_rawDesc = nil
	file_bufman_registry_v1alpha_authn_proto_goTypes = nil
	file_bufman_registry_v1alpha_authn_proto_depIdxs = nil
}