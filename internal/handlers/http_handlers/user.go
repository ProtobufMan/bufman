package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userGroup struct {
	userController *controllers.UserController
}

var UserGroup = &userGroup{
	userController: controllers.NewUserController(),
}

func (group *userGroup) CreateUser(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.CreateUserRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	// 查询用户
	resp, err := group.userController.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
		return
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *userGroup) GetUser(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.GetUserRequest{}
	bindErr := c.ShouldBindUri(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.userController.GetUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}

func (group *userGroup) ListUsers(c *gin.Context) {
	// 绑定参数
	req := &registryv1alpha1.ListUsersRequest{}
	bindErr := c.ShouldBindJSON(req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, NewHTTPResponse(bindErr))
		return
	}

	resp, err := group.userController.ListUsers(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
