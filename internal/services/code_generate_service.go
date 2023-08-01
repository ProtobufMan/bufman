package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufimage"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/plugin"
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/pluginpb"
	"gorm.io/gorm"
)

type CodeGenerateService interface {
	PluginCodeGenerate(ctx context.Context, userID string, req *connect.Request[registryv1alpha1.GenerateCodeRequest]) ([]*pluginpb.CodeGeneratorResponse, e.ResponseError)
}

type CodeGenerateServiceImpl struct {
	pluginMapper     mapper.PluginMapper
	dockerRepoMapper mapper.DockerRepoMapper
}

func NewCodeGenerateService() CodeGenerateService {
	return &CodeGenerateServiceImpl{
		pluginMapper:     &mapper.PluginMapperImpl{},
		dockerRepoMapper: &mapper.DockerRepoMapperImpl{},
	}
}

func (codeGenerateService *CodeGenerateServiceImpl) PluginCodeGenerate(ctx context.Context, userID string, req *connect.Request[registryv1alpha1.GenerateCodeRequest]) ([]*pluginpb.CodeGeneratorResponse, e.ResponseError) {
	image, err := bufimage.NewImageForProto(req.Msg.GetImage())
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	codeGeneratorResponses := make([]*pluginpb.CodeGeneratorResponse, 0, len(req.Msg.GetRequests()))
	pluginGenerationRequests := req.Msg.GetRequests()
	for _, pluginGenerationRequest := range pluginGenerationRequests {
		pluginReference := pluginGenerationRequest.GetPluginReference()

		// 在数据库中查询插件
		owner := pluginReference.GetOwner()
		name := pluginReference.GetName()
		version := pluginReference.GetVersion()
		revision := pluginReference.GetRevision()
		var pluginModel *model.Plugin
		if revision != 0 && version != "" {
			pluginModel, err = codeGenerateService.pluginMapper.FindByNameAndVersionReversion(owner, name, version, revision)
		} else if revision == 0 && version == "" {
			pluginModel, err = codeGenerateService.pluginMapper.FindLastByName(owner, name)
		} else if version != "" {
			pluginModel, err = codeGenerateService.pluginMapper.FindLastByNameAndVersion(owner, name, version)
		} else {
			return nil, e.NewInvalidArgumentError("reversion is not empty but version is empty")
		}

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, e.NewNotFoundError(fmt.Sprintf("plugin %s/%s not found", owner, name))
			}

			return nil, e.NewInternalError("Find Plugin")
		}
		if pluginModel.Visibility == uint8(registryv1alpha1.CuratedPluginVisibility_CURATED_PLUGIN_VISIBILITY_PRIVATE) && pluginModel.UserID != userID {
			// 插件是私有的，且不属于当前用户
			return nil, e.NewPermissionDeniedError(fmt.Sprintf("plugin %s/%s is private", owner, name))
		}

		// get docker repo
		dockerRepo, err := codeGenerateService.dockerRepoMapper.FindByDockerRepoID(pluginModel.DockerRepoID)
		if err != nil {
			return nil, e.NewInternalError(err.Error())
		}

		// options
		options := pluginGenerationRequest.GetOptions()
		if len(options) != 1 {
			return nil, e.NewInvalidArgumentError(fmt.Sprintf("options in %s/%s(length must be 1)", owner, name))
		}
		option := options[0]

		// include imports and include well known types
		includeImports := req.Msg.GetIncludeImports()
		includeWellKnownTypes := req.Msg.GetIncludeWellKnownTypes()
		if pluginGenerationRequest.IncludeImports != nil {
			includeImports = pluginGenerationRequest.GetIncludeImports()
		}
		if pluginGenerationRequest.IncludeWellKnownTypes != nil {
			includeWellKnownTypes = pluginGenerationRequest.GetIncludeWellKnownTypes()
		}

		codeGenerateHelper := plugin.NewCodeGenerateHelper(dockerRepo.Address, dockerRepo.UserName, dockerRepo.Password)
		// get CodeGeneratorRequest
		codeGeneratorRequest := codeGenerateHelper.GetGeneratorRequest(image, option, includeImports, includeWellKnownTypes)
		// generate code
		codeGeneratorResponse, err := codeGenerateHelper.Generate(ctx, pluginModel.PluginName, pluginModel.ImageName, pluginModel.ImageDigest, codeGeneratorRequest)
		if err != nil {
			return nil, e.NewInternalError(err.Error())
		}

		codeGeneratorResponses = append(codeGeneratorResponses, codeGeneratorResponse)
	}

	return codeGeneratorResponses, nil
}
