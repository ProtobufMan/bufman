package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/validity"
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

func (handler *DownloadServiceHandler) DownloadManifestAndBlobs(ctx context.Context, req *connect.Request[registryv1alpha.DownloadManifestAndBlobsRequest]) (*connect.Response[registryv1alpha.DownloadManifestAndBlobsResponse], error) {
	// 检查用户权限
	userID, _ := ctx.Value(constant.UserIDKey).(string)
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	if checkErr != nil {
		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取对应文件内容、文件清单
	fileManifest, fileBlobs, err := handler.downloadService.DownloadManifestAndBlobs(repository.RepositoryID, req.Msg.GetReference())
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha.DownloadManifestAndBlobsResponse{
		Manifest: fileManifest,
		Blobs:    fileBlobs,
	})
	return resp, nil
}
