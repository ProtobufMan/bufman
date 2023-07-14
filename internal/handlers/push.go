package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/bufman/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/bufbuild/connect-go"
)

type PushServiceHandler struct {
	pushService     services.PushService
	validator       validity.Validator
	downloadService services.DownloadService
}

func NewPushServiceHandler() *PushServiceHandler {
	return &PushServiceHandler{
		pushService:     services.NewPushService(),
		downloadService: services.NewDownloadService(),
		validator:       validity.NewValidator(),
	}
}

func (handler *PushServiceHandler) PushManifestAndBlobs(ctx context.Context, req *connect.Request[registryv1alpha.PushManifestAndBlobsRequest]) (*connect.Response[registryv1alpha.PushManifestAndBlobsResponse], error) {
	// 验证参数

	// 检查tags名称合法性
	var argErr e.ResponseError
	for _, tag := range req.Msg.GetTags() {
		argErr = handler.validator.CheckTagName(tag)
		if argErr != nil {
			return nil, connect.NewError(argErr.Code(), argErr.Err())
		}
	}

	// 检查draft名称合法性
	if req.Msg.GetDraftName() != "" {
		argErr = handler.validator.CheckDraftName(req.Msg.GetDraftName())
		if argErr != nil {
			return nil, connect.NewError(argErr.Code(), argErr.Err())
		}
	}

	// draft和tag只能二选一
	if req.Msg.GetDraftName() != "" && len(req.Msg.GetTags()) > 0 {
		responseError := e.NewInvalidArgumentError("draft and tags (only choose one)")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	// 检查上传文件
	fileManifest, blobSet, bufConfig, checkErr := handler.validator.CheckManifestAndBlobs(ctx, req.Msg.GetManifest(), req.Msg.GetBlobs())
	if checkErr != nil {
		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取依赖文件
	dependentCommits, dependenceErr := handler.pushService.GetDependencies(bufConfig.Build.DependencyModuleReferences)
	if dependenceErr != nil {
		return nil, connect.NewError(dependenceErr.Code(), dependenceErr)
	}

	// 读取依赖文件
	dependentManifests := make([]*manifest.Manifest, 0, len(dependentCommits))
	dependentBlobSets := make([]*manifest.BlobSet, 0, len(dependentCommits))
	for i := 0; i < len(dependentCommits); i++ {
		dependentCommit := dependentCommits[i]
		dependentManifest, dependentBlobSet, downloadErr := handler.downloadService.DownloadManifestAndBlobs(dependentCommit.RepositoryID, dependentCommit.CommitName)
		if downloadErr != nil {
			return nil, connect.NewError(downloadErr.Code(), downloadErr)
		}

		dependentManifests = append(dependentManifests, dependentManifest)
		dependentBlobSets = append(dependentBlobSets, dependentBlobSet)
	}

	// 编译检查
	compileErr := handler.pushService.TryCompile(ctx, fileManifest, blobSet, dependentManifests, dependentBlobSets)
	if compileErr != nil {
		return nil, connect.NewError(compileErr.Code(), compileErr)
	}

	var commit *model.Commit
	var serviceErr e.ResponseError
	userID := ctx.Value(constant.UserIDKey).(string)
	if req.Msg.DraftName != "" {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithDraft(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet, req.Msg.GetDraftName())
	} else if len(req.Msg.GetTags()) > 0 {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithTags(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet, req.Msg.GetTags())
	} else {
		commit, serviceErr = handler.pushService.PushManifestAndBlobs(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet)
	}
	if serviceErr != nil {
		return nil, connect.NewError(serviceErr.Code(), serviceErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.PushManifestAndBlobsResponse{
		LocalModulePin: commit.ToProtoLocalModulePin(),
	})
	return resp, nil
}
