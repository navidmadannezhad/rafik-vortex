package router

import (
	"airun/code-reviewer/controller"

	"github.com/gin-gonic/gin"
)

func attachRoutes(engine *gin.Engine) {
	engine.POST("/api/v1/webhook", controller.RequestPRReview)
	engine.POST("/api/v1/request-ai-instructions", controller.RequestAIInstructions)
}

func RunRouter() {
	engine := gin.Default()
	attachRoutes(engine)
	engine.Run(":8080")
}
