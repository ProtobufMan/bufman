package grpc_handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type DocServiceHandler struct {
	docsService services.DocsService
	validator   validity.Validator
}

func NewDocServiceHandler() *DocServiceHandler {
	return &DocServiceHandler{
		docsService: services.NewDocsService(),
		validator:   validity.NewValidator(),
	}
}

func (handler *DocServiceHandler) GetSourceDirectoryInfo(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceDirectoryInfoRequest]) (*connect.Response[registryv1alpha1.GetSourceDirectoryInfoResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 检查用户权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DocServiceGetSourceDirectoryInfoProcedure)
	if checkErr != nil {
		logger.Errorf("Error Check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取目录结构信息
	directoryInfo, respErr := handler.docsService.GetSourceDirectoryInfo(ctx, repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
		logger.Errorf("Error get source dir info: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetSourceDirectoryInfoResponse{
		Root: directoryInfo.ToProtoFileInfo(),
	})
	return resp, nil
}

func (handler *DocServiceHandler) GetSourceFile(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceFileRequest]) (*connect.Response[registryv1alpha1.GetSourceFileResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 检查用户权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DocServiceGetSourceFileProcedure)
	if checkErr != nil {
		logger.Errorf("Error Check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取源码内容
	content, respErr := handler.docsService.GetSourceFile(ctx, repository.RepositoryID, req.Msg.GetReference(), req.Msg.GetPath())
	if respErr != nil {
		logger.Errorf("Error get source file: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetSourceFileResponse{
		Content: content,
	})
	return resp, nil
}

func (handler *DocServiceHandler) GetModulePackages(ctx context.Context, req *connect.Request[registryv1alpha1.GetModulePackagesRequest]) (*connect.Response[registryv1alpha1.GetModulePackagesResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 检查用户权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DocServiceGetModulePackagesProcedure)
	if checkErr != nil {
		logger.Errorf("Error Check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	modulePackages, respErr := handler.docsService.GetModulePackages(ctx, repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
		logger.Errorf("Error get module packages: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetModulePackagesResponse{
		Name:           req.Msg.GetRepository(),
		ModulePackages: modulePackages,
	})
	return resp, nil
}

func (handler *DocServiceHandler) GetModuleDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetModuleDocumentationRequest]) (*connect.Response[registryv1alpha1.GetModuleDocumentationResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 检查用户权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DocServiceGetModuleDocumentationProcedure)
	if checkErr != nil {
		logger.Errorf("Error Check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	moduleDocumentation, respErr := handler.docsService.GetModuleDocumentation(ctx, repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
		logger.Errorf("Error get module doc: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetModuleDocumentationResponse{
		ModuleDocumentation: moduleDocumentation,
	})
	return resp, nil
}

func (handler *DocServiceHandler) GetPackageDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetPackageDocumentationRequest]) (*connect.Response[registryv1alpha1.GetPackageDocumentationResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 检查用户权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DocServiceGetPackageDocumentationProcedure)
	if checkErr != nil {
		logger.Errorf("Error Check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	packageDocumentation, respErr := handler.docsService.GetPackageDocumentation(ctx, repository.RepositoryID, req.Msg.GetReference(), req.Msg.GetPackageName())
	if respErr != nil {
		logger.Errorf("Error get package doc: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetPackageDocumentationResponse{
		PackageDocumentation: packageDocumentation,
	})
	return resp, nil
}
