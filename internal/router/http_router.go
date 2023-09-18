package router

import (
	"github.com/ProtobufMan/bufman/internal/handlers/http_handlers"
	"github.com/gin-gonic/gin"
)

func InitHTTPRouter() *gin.Engine {
	router := gin.Default()

	authn := router.Group("/authn")
	{
		authn.GET("/current_user", http_handlers.AuthnGroup.GetCurrentUser)
	}

	return router
}
