package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	router := gin.Default()
	InitGRPCRouter(router)
	InitHTTPRouter(router)

	return router
}
