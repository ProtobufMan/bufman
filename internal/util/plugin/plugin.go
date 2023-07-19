package plugin

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/buf/bufcli"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufimage"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufpluginexec"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufwasm"
	"github.com/ProtobufMan/bufman-cli/private/pkg/app"
	"github.com/ProtobufMan/bufman-cli/private/pkg/command"
	"github.com/ProtobufMan/bufman/internal/constant"
	"google.golang.org/protobuf/types/pluginpb"
	"path/filepath"
)

type CodeGenerateHelper interface {
	GetGeneratorRequest(image bufimage.Image, option string, includeImports, includeWellKnownTypes bool) *pluginpb.CodeGeneratorRequest
	Generate(ctx context.Context, container app.Container, pluginName string, codeGeneratorRequest *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

func NewCodeGenerateHelper() CodeGenerateHelper {
	return &CodeGenerateHelperImpl{}
}

type CodeGenerateHelperImpl struct {
}

func (helper *CodeGenerateHelperImpl) GetGeneratorRequest(image bufimage.Image, option string, includeImports, includeWellKnownTypes bool) *pluginpb.CodeGeneratorRequest {
	return bufimage.ImageToCodeGeneratorRequest(image, option, nil, includeImports, includeWellKnownTypes)
}

func (helper *CodeGenerateHelperImpl) Generate(ctx context.Context, container app.Container, pluginName string, codeGeneratorRequest *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	storageosProvider := bufcli.NewStorageosProvider(false)
	runner := command.NewRunner()
	wasmPluginExecutor, err := bufwasm.NewPluginExecutor("")
	if err != nil {
		return nil, err
	}

	// generate
	pluginPath := filepath.Join(constant.PluginSaveDir, pluginName)
	var generateOptions []bufpluginexec.GenerateOption
	generateOptions = append(generateOptions, bufpluginexec.GenerateWithPluginPath(pluginPath))
	generateOptions = append(generateOptions, bufpluginexec.GenerateWithWASMEnabled())
	codeGeneratorResponse, err := bufpluginexec.NewGenerator(nil, storageosProvider, runner, wasmPluginExecutor).Generate(ctx, container, pluginName, []*pluginpb.CodeGeneratorRequest{codeGeneratorRequest}, generateOptions...)
	if err != nil {
		return nil, err
	}

	return codeGeneratorResponse, nil
}
