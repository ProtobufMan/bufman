package http_handlers

import (
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authnGroup struct {
	authnController *controllers.AuthnController
}

var AuthnGroup = &authnGroup{
	authnController: controllers.NewAuthnController(),
}

func (group *authnGroup) GetCurrentUser(c *gin.Context) {
	req := &registryv1alpha1.GetCurrentUserRequest{}
	resp, err := group.authnController.GetCurrentUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPResponse(err))
		return
	}

	// 正常返回
	c.JSON(http.StatusOK, NewHTTPResponse(resp))
}
