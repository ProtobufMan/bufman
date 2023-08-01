package plugin

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufimage"
	"google.golang.org/protobuf/types/pluginpb"
)

type CodeGenerateHelper interface {
	GetGeneratorRequest(image bufimage.Image, option string, includeImports, includeWellKnownTypes bool) *pluginpb.CodeGeneratorRequest
	Generate(ctx context.Context, pluginName, image string, codeGeneratorRequest *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

func NewCodeGenerateHelper(address, username, password string) CodeGenerateHelper {
	return &CodeGenerateHelperImpl{
		address:  address,
		username: username,
		password: password,
	}
}

type CodeGenerateHelperImpl struct {
	address  string
	username string
	password string
}

func (helper *CodeGenerateHelperImpl) GetGeneratorRequest(image bufimage.Image, option string, includeImports, includeWellKnownTypes bool) *pluginpb.CodeGeneratorRequest {
	return bufimage.ImageToCodeGeneratorRequest(image, option, nil, includeImports, includeWellKnownTypes)
}

func (helper *CodeGenerateHelperImpl) Generate(ctx context.Context, pluginName, image string, codeGeneratorRequest *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {

	// 连接docker repo
	d, err := NewDocker(helper.address, helper.username, helper.password)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	codeGeneratorResponse, err := d.GenerateCode(ctx, pluginName, image, codeGeneratorRequest)
	if err != nil {
		return nil, err
	}

	return codeGeneratorResponse, nil
}
