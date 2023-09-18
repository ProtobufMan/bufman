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

	token := router.Group("/token")
	{
		token.POST("/", http_handlers.TokenGroup.CreateToken)            // 创建token
		token.GET("/:token_id", http_handlers.TokenGroup.GetToken)       // 获取token
		token.POST("/list", http_handlers.TokenGroup.ListTokens)         // 批量查询token
		token.DELETE("/:token_id", http_handlers.TokenGroup.DeleteToken) // 删除tokens
	}

	return router
}
