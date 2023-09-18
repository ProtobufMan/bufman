package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pluginGroup struct {
	pluginController *controllers.PluginController
}

var PluginGroup = &pluginGroup{
	pluginController: controllers.NewPluginController(),
}

func (group *pluginGroup) ListCuratedPlugins(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListCuratedPluginsRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.pluginController.ListCuratedPlugins(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *pluginGroup) CreateCuratedPlugin(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.CreateCuratedPluginRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.pluginController.CreateCuratedPlugin(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *pluginGroup) GetLatestCuratedPlugin(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.GetLatestCuratedPluginRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.pluginController.GetLatestCuratedPlugin(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
