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
		user.POST("/create", http_handlers.UserGroup.CreateUser) // 创建用户
		user.GET("/:id", http_handlers.UserGroup.GetUser)        // 查询用户
		user.POST("/list", http_handlers.UserGroup.ListUsers)    // 批量查询用户
	}

	token := router.Group("/token")
	{
		token.POST("/create", http_handlers.TokenGroup.CreateToken)      // 创建token
		token.GET("/:token_id", http_handlers.TokenGroup.GetToken)       // 获取token
		token.POST("/list", http_handlers.TokenGroup.ListTokens)         // 批量查询token
		token.DELETE("/:token_id", http_handlers.TokenGroup.DeleteToken) // 删除tokens
	}

	repository := router.Group("/repository")
	{
		repository.POST("/create", http_handlers.RepositoryGroup.CreateRepositoryByFullName)             // 创建repository
		repository.GET("/:id", http_handlers.RepositoryGroup.GetRepository)                              // 获取repository
		repository.POST("/list", http_handlers.RepositoryGroup.ListRepositories)                         // 批量查询所有repository
		repository.DELETE("/:id", http_handlers.RepositoryGroup.DeleteRepository)                        // 删除repository
		repository.POST("/list/:user_id", http_handlers.RepositoryGroup.ListUserRepositories)            // 批量查询用户的repository
		repository.POST("/list_accessible", http_handlers.RepositoryGroup.ListRepositoriesUserCanAccess) // 批量查询当前用户可访问的repository
		repository.PUT("/deprecate", http_handlers.RepositoryGroup.DeprecateRepositoryByName)            // 弃用repository
		repository.PUT("/undeprecate", http_handlers.RepositoryGroup.UndeprecateRepositoryByName)        // 解除弃用
		repository.PUT("/update", http_handlers.RepositoryGroup.UpdateRepositorySettingsByName)          // 更新repository
	}

	return router
}
