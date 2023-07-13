package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmanifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/module/v1alpha"
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
	fileManifest, blobSet, err := handler.downloadService.DownloadManifestAndBlobs(repository.RepositoryID, req.Msg.GetReference())
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	// TODO: 以后改为一致的manifest和blobSet
	protoManifest, protoBlobs, toProtoErr := bufmanifest.ToProtoManifestAndBlobs(ctx, fileManifest, blobSet)
	if toProtoErr != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	blobs := make([]*modulev1alpha.Blob, 0, len(protoBlobs))
	for i := 0; i < len(protoBlobs); i++ {
		blobs = append(blobs, &modulev1alpha.Blob{
			Digest: &modulev1alpha.Digest{
				DigestType: modulev1alpha.DigestType(protoBlobs[i].Digest.DigestType),
				Digest:     protoBlobs[i].Digest.Digest,
			},
			Content: protoBlobs[i].Content,
		})
	}

	resp := connect.NewResponse(&registryv1alpha.DownloadManifestAndBlobsResponse{
		Manifest: &modulev1alpha.Blob{
			Digest: &modulev1alpha.Digest{
				DigestType: modulev1alpha.DigestType(protoManifest.Digest.DigestType),
				Digest:     protoManifest.Digest.Digest,
			},
			Content: protoManifest.Content,
		},
		Blobs: blobs,
	})
	return resp, nil
}
