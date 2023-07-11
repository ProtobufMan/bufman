package handlers

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/util/manifest"
	"github.com/bufbuild/connect-go"
)

type PushServiceHandler struct {
	pushService services.PushService
}

func NewPushServiceHandler() *PushServiceHandler {
	return &PushServiceHandler{
		pushService: services.NewPushService(),
	}
}

func (handler *PushServiceHandler) PushManifestAndBlobs(ctx context.Context, req *connect.Request[registryv1alpha.PushManifestAndBlobsRequest]) (*connect.Response[registryv1alpha.PushManifestAndBlobsResponse], error) {
	if req.Msg.GetDraftName() == constant.DefaultBranch {
		responseError := e.NewInvalidArgumentError("draft (can not be 'main')")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	if req.Msg.GetDraftName() != "" && len(req.Msg.GetTags()) > 0 {
		responseError := e.NewInvalidArgumentError("draft and tags (only choose one)")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	// 读取文件清单
	fileManifest, err := manifest.NewManifestFromProto(ctx, req.Msg.GetManifest())
	if err != nil {
		responseError := e.NewInvalidArgumentError("manifest")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	// 读取文件列表
	fileBlobs, err := manifest.NewBlobSetFromProto(ctx, req.Msg.GetBlobs())
	if err != nil {
		responseError := e.NewInvalidArgumentError("blobs")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	// 检查文件清单和blobs
	err = fileManifest.Range(func(path string, digest manifest.Digest) error {
		_, ok := fileBlobs.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		return nil
	})
	if err != nil {
		responseError := e.NewInvalidArgumentError("blobs and manifest")
		return nil, connect.NewError(responseError.Code(), responseError.Err())
	}

	var commit *model.Commit
	var serviceErr e.ResponseError
	userID := ctx.Value(constant.UserIDKey).(string)
	if req.Msg.DraftName != "" {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithDraft(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, fileBlobs, req.Msg.GetDraftName())
	} else if len(req.Msg.GetTags()) > 0 {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithTags(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, fileBlobs, req.Msg.GetTags())
	} else {
		commit, serviceErr = handler.pushService.PushManifestAndBlobs(userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, fileBlobs)
	}
	if serviceErr != nil {
		return nil, connect.NewError(serviceErr.Code(), serviceErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.PushManifestAndBlobsResponse{
		LocalModulePin: commit.ToProtoLocalModulePin(),
	})
	return resp, nil
}
