package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type commitGroup struct {
	commitController *controllers.CommitController
}

var CommitGroup = &commitGroup{
	commitController: controllers.NewCommitController(),
}

func (group *commitGroup) ListRepositoryCommitsByReference(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListRepositoryCommitsByReferenceRequest{}
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

	resp, err := group.commitController.ListRepositoryCommitsByReference(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *commitGroup) GetRepositoryCommitByReference(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.GetRepositoryCommitByReferenceRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.commitController.GetRepositoryCommitByReference(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *commitGroup) ListRepositoryDraftCommits(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListRepositoryDraftCommitsRequest{}
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

	resp, err := group.commitController.ListRepositoryDraftCommits(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *commitGroup) DeleteRepositoryDraftCommit(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.DeleteRepositoryDraftCommitRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.commitController.DeleteRepositoryDraftCommit(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
