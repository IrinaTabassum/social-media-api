package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
    incomingRoutes.POST("/comments", controller.CreateComment())
	incomingRoutes.GET("/comments/:id", controller.GetCommentByID())
	incomingRoutes.GET("/comments", controller.GetCommentList())
    incomingRoutes.PUT("/comments/:id", controller.UpdateComment())
	incomingRoutes.DELETE("/comments/:id", controller.DeleteComment())
    
}