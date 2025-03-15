package infrastructure

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	routes := router.Group("webhook")
	{
		routes.POST("/process-pull-request", HandlePullRequestEvent)
		routes.POST("/process-deploy", HandleDeployEvent) 
	}
}