package router

import (
	"github.com/ProtobufMan/bufman/internal/handlers/http_handlers"
	"github.com/gin-gonic/gin"
)

func InitHTTPRouter() *gin.Engine {
	router := gin.Default()

	authn := router.Group("/authn")
	{
		authn.GET("/current_user", http_handlers.AuthnGroup.GetCurrentUser) // 根据token获取当前用户信息
	}

	user := router.Group("/user")
	{
		user.POST("/", http_handlers.UserGroup.CreateUser)    // 创建用户
		user.GET("/:id", http_handlers.UserGroup.GetUser)     // 查询用户
		user.POST("/list", http_handlers.UserGroup.ListUsers) // 批量查询用户
	}

	return router
}
