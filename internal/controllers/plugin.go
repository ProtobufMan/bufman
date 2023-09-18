package controllers

import (
	"context"
	"fmt"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/google/uuid"
)

type PluginController struct {
	pluginService services.PluginService
	userService   services.UserService
	validator     validity.Validator
}

func NewPluginController() *PluginController {
	return &PluginController{
		pluginService: services.NewPluginService(),
		userService:   services.NewUserService(),
		validator:     validity.NewValidator(),
	}
}

func (controller *PluginController) ListCuratedPlugins(ctx context.Context, req *registryv1alpha1.ListCuratedPluginsRequest) (*registryv1alpha1.ListCuratedPluginsResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		fmt.Printf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	plugins, err := controller.pluginService.ListPlugins(ctx, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse(), req.GetIncludeDeprecated())
	if err != nil {
		logger.Errorf("Error list plugins: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(plugins))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListCuratedPluginsResponse{
		Plugins:       plugins.ToProtoPlugins(),
		NextPageToken: nextPageToken,
	}
	return resp, nil
}

func (controller *PluginController) CreateCuratedPlugin(ctx context.Context, req *registryv1alpha1.CreateCuratedPluginRequest) (*registryv1alpha1.CreateCuratedPluginResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 检查插件名称是否合法
	checkErr := controller.validator.CheckPluginName(req.GetName())
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, checkErr
	}

	// 检查版本号是否合法
	checkErr = controller.validator.CheckVersion(req.GetVersion())
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, checkErr
	}

	// 检查reversion
	if req.GetRevision() < 1 {
		checkErr = e.NewInvalidArgumentError("reversion must greater than 0")

		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, checkErr
	}

	// 检查用户名称
	user, checkErr := controller.userService.GetUser(ctx, userID)
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, checkErr
	}
	if user.UserName != req.GetOwner() {
		// 用户id与owner必须对应
		checkErr = e.NewPermissionDeniedError("token and owner mismatch")

		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, checkErr
	}

	// 检查插件数据
	plugin := &model.Plugin{
		ID:          0,
		UserID:      userID,
		UserName:    user.UserName,
		PluginID:    uuid.NewString(),
		PluginName:  req.GetName(),
		Version:     req.GetVersion(),
		Reversion:   req.GetRevision(),
		ImageName:   req.GetImageName(),
		ImageDigest: req.GetImageDigest(),
		Description: req.GetDescription(),
		Visibility:  uint8(req.GetVisibility()),
	}

	plugin, err := controller.pluginService.CreatePlugin(ctx, plugin, req.GetDockerRepoName())
	if err != nil {
		resp := &registryv1alpha1.CreateCuratedPluginResponse{
			Configuration: plugin.ToProtoPlugin(),
		}
		fmt.Printf("Error create plugin: %v\n", err.Error())

		return resp, err
	}

	resp := &registryv1alpha1.CreateCuratedPluginResponse{
		Configuration: plugin.ToProtoPlugin(),
	}
	return resp, nil
}

func (controller *PluginController) GetLatestCuratedPlugin(ctx context.Context, req *registryv1alpha1.GetLatestCuratedPluginRequest) (*registryv1alpha1.GetLatestCuratedPluginResponse, e.ResponseError) {
	version := req.GetVersion()
	reversion := req.GetRevision()

	var plugin *model.Plugin
	var err e.ResponseError
	if version != "" && reversion != 0 {
		// version不为空，reversion不为空
		plugin, err = controller.pluginService.GetLatestPluginWithVersionAndReversion(ctx, req.GetOwner(), req.GetName(), version, reversion)
	} else if version == "" && reversion == 0 {
		// version为空，reversion为空
		plugin, err = controller.pluginService.GetLatestPlugin(ctx, req.GetOwner(), req.GetName())
	} else if version != "" {
		// version不为空，reversion为空
		plugin, err = controller.pluginService.GetLatestPluginWithVersion(ctx, req.GetOwner(), req.GetName(), version)
	} else {
		// version为空，reversion不为空，非法
		err = e.NewInvalidArgumentError("version is empty but reversion is not empty")
	}
	if err != nil {
		fmt.Printf("Error get lastest curated plugin: %v\n", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.GetLatestCuratedPluginResponse{
		Plugin: plugin.ToProtoPlugin(),
	}
	return resp, nil
}
