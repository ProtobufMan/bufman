package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type repositoryGroup struct {
	repositoryController *controllers.RepositoryController
}

var RepositoryGroup = &repositoryGroup{
	repositoryController: controllers.NewRepositoryController(),
}

func (group *repositoryGroup) CreateRepositoryByFullName(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.CreateRepositoryByFullNameRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.CreateRepositoryByFullName(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) GetRepository(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.GetRepositoryRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.GetRepository(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) ListRepositories(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListRepositoriesRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.ListRepositories(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) DeleteRepository(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.DeleteRepositoryRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.DeleteRepository(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) ListUserRepositories(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListUserRepositoriesRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}
	bindErr = c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.ListUserRepositories(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) ListRepositoriesUserCanAccess(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListRepositoriesUserCanAccessRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.ListRepositoriesUserCanAccess(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) DeprecateRepositoryByName(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.DeprecateRepositoryByNameRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.DeprecateRepositoryByName(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) UndeprecateRepositoryByName(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.UndeprecateRepositoryByNameRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.UndeprecateRepositoryByName(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *repositoryGroup) UpdateRepositorySettingsByName(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.UpdateRepositorySettingsByNameRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.repositoryController.UpdateRepositorySettingsByName(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
