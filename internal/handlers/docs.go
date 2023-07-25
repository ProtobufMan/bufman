package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/bufbuild/connect-go"
)

type DocsServiceHandler struct {
}

func (handler *DocsServiceHandler) GetSourceDirectoryInfo(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceDirectoryInfoRequest]) (*connect.Response[registryv1alpha1.GetSourceDirectoryInfoResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *DocsServiceHandler) GetSourceFile(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceFileRequest]) (*connect.Response[registryv1alpha1.GetSourceFileResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *DocsServiceHandler) GetModulePackages(ctx context.Context, req *connect.Request[registryv1alpha1.GetModulePackagesRequest]) (*connect.Response[registryv1alpha1.GetModulePackagesResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *DocsServiceHandler) GetModuleDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetModuleDocumentationRequest]) (*connect.Response[registryv1alpha1.GetModuleDocumentationResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *DocsServiceHandler) GetPackageDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetPackageDocumentationRequest]) (*connect.Response[registryv1alpha1.GetPackageDocumentationResponse], error) {
	//TODO implement me
	panic("implement me")
}
