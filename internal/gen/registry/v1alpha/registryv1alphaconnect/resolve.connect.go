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

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: registry/v1alpha/resolve.proto

package registryv1alphaconnect

import (
	context "context"
	errors "errors"
	v1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	connect_go "github.com/bufbuild/connect-go"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion1_7_0

const (
	// ResolveServiceName is the fully-qualified name of the ResolveService service.
	ResolveServiceName = "registry.v1alpha.ResolveService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ResolveServiceGetModulePinsProcedure is the fully-qualified name of the ResolveService's
	// GetModulePins RPC.
	ResolveServiceGetModulePinsProcedure = "/registry.v1alpha.ResolveService/GetModulePins"
)

// ResolveServiceClient is a client for the registry.v1alpha.ResolveService service.
type ResolveServiceClient interface {
	// GetModulePins finds all the latest digests and respective dependencies of
	// the provided module references and picks a set of distinct modules pins.
	//
	// Note that module references with commits should still be passed to this function
	// to make sure this function can do dependency resolution.
	//
	// This function also deals with tiebreaking what ModulePin wins for the same repository.
	GetModulePins(context.Context, *connect_go.Request[v1alpha.GetModulePinsRequest]) (*connect_go.Response[v1alpha.GetModulePinsResponse], error)
}

// NewResolveServiceClient constructs a client for the registry.v1alpha.ResolveService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewResolveServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ResolveServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &resolveServiceClient{
		getModulePins: connect_go.NewClient[v1alpha.GetModulePinsRequest, v1alpha.GetModulePinsResponse](
			httpClient,
			baseURL+ResolveServiceGetModulePinsProcedure,
			connect_go.WithIdempotency(connect_go.IdempotencyNoSideEffects),
			connect_go.WithClientOptions(opts...),
		),
	}
}

// resolveServiceClient implements ResolveServiceClient.
type resolveServiceClient struct {
	getModulePins *connect_go.Client[v1alpha.GetModulePinsRequest, v1alpha.GetModulePinsResponse]
}

// GetModulePins calls registry.v1alpha.ResolveService.GetModulePins.
func (c *resolveServiceClient) GetModulePins(ctx context.Context, req *connect_go.Request[v1alpha.GetModulePinsRequest]) (*connect_go.Response[v1alpha.GetModulePinsResponse], error) {
	return c.getModulePins.CallUnary(ctx, req)
}

// ResolveServiceHandler is an implementation of the registry.v1alpha.ResolveService service.
type ResolveServiceHandler interface {
	// GetModulePins finds all the latest digests and respective dependencies of
	// the provided module references and picks a set of distinct modules pins.
	//
	// Note that module references with commits should still be passed to this function
	// to make sure this function can do dependency resolution.
	//
	// This function also deals with tiebreaking what ModulePin wins for the same repository.
	GetModulePins(context.Context, *connect_go.Request[v1alpha.GetModulePinsRequest]) (*connect_go.Response[v1alpha.GetModulePinsResponse], error)
}

// NewResolveServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewResolveServiceHandler(svc ResolveServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	resolveServiceGetModulePinsHandler := connect_go.NewUnaryHandler(
		ResolveServiceGetModulePinsProcedure,
		svc.GetModulePins,
		connect_go.WithIdempotency(connect_go.IdempotencyNoSideEffects),
		connect_go.WithHandlerOptions(opts...),
	)
	return "/registry.v1alpha.ResolveService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ResolveServiceGetModulePinsProcedure:
			resolveServiceGetModulePinsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedResolveServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedResolveServiceHandler struct{}

func (UnimplementedResolveServiceHandler) GetModulePins(context.Context, *connect_go.Request[v1alpha.GetModulePinsRequest]) (*connect_go.Response[v1alpha.GetModulePinsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("registry.v1alpha.ResolveService.GetModulePins is not implemented"))
}
