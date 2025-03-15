package infrastructure

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {


	routes := router.Group("pull_request")

	{
		routes.POST("/process", HandlePullRequestEvent)
	}

}
