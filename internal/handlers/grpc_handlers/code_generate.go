package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
	"time"
)

type CodeGenerateServiceHandler struct {
	codeGenerateService services.CodeGenerateService
}

func NewCodeGenerateServiceHandler() *CodeGenerateServiceHandler {
	return &CodeGenerateServiceHandler{
		codeGenerateService: services.NewCodeGenerateService(),
	}
}

func (handler *CodeGenerateServiceHandler) GenerateCode(ctx context.Context, req *connect.Request[registryv1alpha1.GenerateCodeRequest]) (*connect.Response[registryv1alpha1.GenerateCodeResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	codeGeneratorResponses, err := handler.codeGenerateService.PluginCodeGenerate(timeoutCtx, userID, req)
	if err != nil {
		logger.Errorf("Error Plugin Code Generate: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err)
	}

	pluginGenerationResponse := make([]*registryv1alpha1.PluginGenerationResponse, len(codeGeneratorResponses))
	for i := 0; i < len(codeGeneratorResponses); i++ {
		pluginGenerationResponse[i] = connect.NewResponse(&registryv1alpha1.PluginGenerationResponse{
			Response: codeGeneratorResponses[i],
		}).Msg
	}
	resp := connect.NewResponse(&registryv1alpha1.GenerateCodeResponse{
		Responses: pluginGenerationResponse,
	})
	return resp, nil
}
