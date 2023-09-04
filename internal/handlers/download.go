package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmanifest"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type DownloadServiceHandler struct {
	downloadService services.DownloadService
	validator       validity.Validator
}

func NewDownloadServiceHandler() *DownloadServiceHandler {
	return &DownloadServiceHandler{
		downloadService: services.NewDownloadService(),
		validator:       validity.NewValidator(),
	}
}

func (handler *DownloadServiceHandler) DownloadManifestAndBlobs(ctx context.Context, req *connect.Request[registryv1alpha1.DownloadManifestAndBlobsRequest]) (*connect.Response[registryv1alpha1.DownloadManifestAndBlobsResponse], error) {
	// 检查用户权限
	userID, _ := ctx.Value(constant.UserIDKey).(string)
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	if checkErr != nil {
		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取对应文件内容、文件清单
	fileManifest, blobSet, err := handler.downloadService.DownloadManifestAndBlobs(ctx, repository.RepositoryID, req.Msg.GetReference())
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	// 转为响应格式
	protoManifest, protoBlobs, toProtoErr := bufmanifest.ToProtoManifestAndBlobs(ctx, fileManifest, blobSet)
	if toProtoErr != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	resp := connect.NewResponse(&registryv1alpha1.DownloadManifestAndBlobsResponse{
		Manifest: protoManifest,
		Blobs:    protoBlobs,
	})
	return resp, nil
}

func (handler *DownloadServiceHandler) Download(ctx context.Context, req *connect.Request[registryv1alpha1.DownloadRequest]) (*connect.Response[registryv1alpha1.DownloadResponse], error) {
	//TODO implement me
	panic("implement me")
}
