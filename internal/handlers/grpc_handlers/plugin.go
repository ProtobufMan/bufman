package grpc_handlers

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
	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
)

type PluginServiceHandler struct {
	pluginService services.PluginService
	userService   services.UserService
	validator     validity.Validator
}

func NewPluginServiceHandler() *PluginServiceHandler {
	return &PluginServiceHandler{
		pluginService: services.NewPluginService(),
		userService:   services.NewUserService(),
		validator:     validity.NewValidator(),
	}
}

func (handler *PluginServiceHandler) ListCuratedPlugins(ctx context.Context, req *connect.Request[registryv1alpha1.ListCuratedPluginsRequest]) (*connect.Response[registryv1alpha1.ListCuratedPluginsResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		fmt.Printf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	plugins, err := handler.pluginService.ListPlugins(ctx, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse(), req.Msg.GetIncludeDeprecated())
	if err != nil {
		logger.Errorf("Error list plugins: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(plugins))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListCuratedPluginsResponse{
		Plugins:       plugins.ToProtoPlugins(),
		NextPageToken: nextPageToken,
	})
	return resp, nil
}

func (handler *PluginServiceHandler) CreateCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.CreateCuratedPluginRequest]) (*connect.Response[registryv1alpha1.CreateCuratedPluginResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 检查插件名称是否合法
	checkErr := handler.validator.CheckPluginName(req.Msg.GetName())
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 检查版本号是否合法
	checkErr = handler.validator.CheckVersion(req.Msg.GetVersion())
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 检查reversion
	if req.Msg.GetRevision() < 1 {
		checkErr = e.NewInvalidArgumentError("reversion must greater than 0")

		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 检查用户名称
	user, checkErr := handler.userService.GetUser(ctx, userID)
	if checkErr != nil {
		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}
	if user.UserName != req.Msg.GetOwner() {
		// 用户id与owner必须对应
		checkErr = e.NewPermissionDeniedError("token and owner mismatch")

		fmt.Printf("Error check: %v\n", checkErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 检查插件数据
	plugin := &model.Plugin{
		ID:          0,
		UserID:      userID,
		UserName:    user.UserName,
		PluginID:    uuid.NewString(),
		PluginName:  req.Msg.GetName(),
		Version:     req.Msg.GetVersion(),
		Reversion:   req.Msg.GetRevision(),
		ImageName:   req.Msg.GetImageName(),
		ImageDigest: req.Msg.GetImageDigest(),
		Description: req.Msg.GetDescription(),
		Visibility:  uint8(req.Msg.GetVisibility()),
	}

	plugin, err := handler.pluginService.CreatePlugin(ctx, plugin, req.Msg.GetDockerRepoName())
	if err != nil {
		resp := connect.NewResponse(&registryv1alpha1.CreateCuratedPluginResponse{
			Configuration: plugin.ToProtoPlugin(),
		})
		fmt.Printf("Error create plugin: %v\n", err.Error())

		return resp, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.CreateCuratedPluginResponse{
		Configuration: plugin.ToProtoPlugin(),
	})
	return resp, nil
}

func (handler *PluginServiceHandler) GetLatestCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[registryv1alpha1.GetLatestCuratedPluginResponse], error) {
	version := req.Msg.GetVersion()
	reversion := req.Msg.GetRevision()

	var plugin *model.Plugin
	var err e.ResponseError
	if version != "" && reversion != 0 {
		// version不为空，reversion不为空
		plugin, err = handler.pluginService.GetLatestPluginWithVersionAndReversion(ctx, req.Msg.GetOwner(), req.Msg.GetName(), version, reversion)
	} else if version == "" && reversion == 0 {
		// version为空，reversion为空
		plugin, err = handler.pluginService.GetLatestPlugin(ctx, req.Msg.GetOwner(), req.Msg.GetName())
	} else if version != "" {
		// version不为空，reversion为空
		plugin, err = handler.pluginService.GetLatestPluginWithVersion(ctx, req.Msg.GetOwner(), req.Msg.GetName(), version)
	} else {
		// version为空，reversion不为空，非法
		err = e.NewInvalidArgumentError("version is empty but reversion is not empty")
	}
	if err != nil {
		fmt.Printf("Error get lastest curated plugin: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetLatestCuratedPluginResponse{
		Plugin: plugin.ToProtoPlugin(),
	})
	return resp, nil
}

func (handler *PluginServiceHandler) DeleteCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[registryv1alpha1.DeleteCuratedPluginResponse], error) {
	//TODO implement me
	panic("implement me")
}
