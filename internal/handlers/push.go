package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufconfig"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/parser"
	"github.com/ProtobufMan/bufman/internal/core/resolve"
	"github.com/ProtobufMan/bufman/internal/core/storage"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
	"io"
)

type PushServiceHandler struct {
	pushService   services.PushService
	validator     validity.Validator
	resolver      resolve.Resolver
	storageHelper storage.StorageHelper
	protoParser   parser.ProtoParser
}

func NewPushServiceHandler() *PushServiceHandler {
	return &PushServiceHandler{
		pushService:   services.NewPushService(),
		validator:     validity.NewValidator(),
		resolver:      resolve.NewResolver(),
		storageHelper: storage.NewStorageHelper(),
		protoParser:   parser.NewProtoParser(),
	}
}

func (handler *PushServiceHandler) PushManifestAndBlobs(ctx context.Context, req *connect.Request[registryv1alpha1.PushManifestAndBlobsRequest]) (*connect.Response[registryv1alpha1.PushManifestAndBlobsResponse], error) {
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
	fileManifest, blobSet, checkErr := handler.validator.CheckManifestAndBlobs(ctx, req.Msg.GetManifest(), req.Msg.GetBlobs())
	if checkErr != nil {
		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 获取bufConfig
	bufConfigBlob, err := handler.storageHelper.GetBufManConfigFromBlob(ctx, fileManifest, blobSet)
	if err != nil {
		configErr := e.NewInternalError(err.Error())
		return nil, connect.NewError(configErr.Code(), configErr)
	}

	var dependentManifests []*manifest.Manifest
	var dependentBlobSets []*manifest.BlobSet
	if bufConfigBlob != nil {
		// 生成Config
		reader, err := bufConfigBlob.Open(ctx)
		if err != nil {
			return nil, e.NewInternalError(err.Error())
		}
		defer reader.Close()
		configData, err := io.ReadAll(reader)
		if err != nil {
			return nil, e.NewInternalError(err.Error())
		}
		bufConfig, err := bufconfig.GetConfigForData(ctx, configData)
		if err != nil {
			// 无法解析配置文件
			return nil, e.NewInternalError(err.Error())
		}

		// 获取全部依赖commits
		dependentCommits, dependenceErr := handler.resolver.GetAllDependenciesFromBufConfig(ctx, bufConfig)
		if dependenceErr != nil {
			return nil, connect.NewError(dependenceErr.Code(), dependenceErr)
		}

		// 读取依赖文件
		dependentManifests = make([]*manifest.Manifest, 0, len(dependentCommits))
		dependentBlobSets = make([]*manifest.BlobSet, 0, len(dependentCommits))
		for i := 0; i < len(dependentCommits); i++ {
			dependentCommit := dependentCommits[i]
			dependentManifest, dependentBlobSet, getErr := handler.pushService.GetManifestAndBlobSet(ctx, dependentCommit.RepositoryID, dependentCommit.CommitName)
			if getErr != nil {
				return nil, connect.NewError(getErr.Code(), getErr)
			}

			dependentManifests = append(dependentManifests, dependentManifest)
			dependentBlobSets = append(dependentBlobSets, dependentBlobSet)
		}
	}

	// 编译检查
	compileErr := handler.protoParser.TryCompile(ctx, fileManifest, blobSet, dependentManifests, dependentBlobSets)
	if compileErr != nil {
		return nil, connect.NewError(compileErr.Code(), compileErr)
	}

	var commit *model.Commit
	var serviceErr e.ResponseError
	userID := ctx.Value(constant.UserIDKey).(string)
	if req.Msg.DraftName != "" {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithDraft(ctx, userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet, req.Msg.GetDraftName())
	} else if len(req.Msg.GetTags()) > 0 {
		commit, serviceErr = handler.pushService.PushManifestAndBlobsWithTags(ctx, userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet, req.Msg.GetTags())
	} else {
		commit, serviceErr = handler.pushService.PushManifestAndBlobs(ctx, userID, req.Msg.GetOwner(), req.Msg.GetRepository(), fileManifest, blobSet)
	}
	if serviceErr != nil {
		return nil, connect.NewError(serviceErr.Code(), serviceErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.PushManifestAndBlobsResponse{
		LocalModulePin: commit.ToProtoLocalModulePin(),
	})
	return resp, nil
}

func (handler *PushServiceHandler) Push(ctx context.Context, req *connect.Request[registryv1alpha1.PushRequest]) (*connect.Response[registryv1alpha1.PushResponse], error) {
	//TODO implement me
	panic("implement me")
}
