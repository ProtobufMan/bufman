package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type dockerGroup struct {
	dockerRepoController *controllers.DockerRepoController
}

var DockerGroup = &dockerGroup{
	dockerRepoController: controllers.NewDockerRepoController(),
}

func (group *dockerGroup) CreateDockerRepo(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.CreateDockerRepoRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.dockerRepoController.CreateDockerRepo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *dockerGroup) GetDockerRepo(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.GetDockerRepoRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.dockerRepoController.GetDockerRepo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *dockerGroup) ListDockerRepos(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListDockerReposRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.dockerRepoController.ListDockerRepos(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *dockerGroup) UpdateDockerRepoByID(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.UpdateDockerRepoByIDRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.dockerRepoController.UpdateDockerRepoByID(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
