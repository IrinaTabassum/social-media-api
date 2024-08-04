package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/comments", middleware.Authentication(), controller.CreateComment())
	incomingRoutes.GET("/comments/:id", middleware.Authentication(), controller.GetCommentByID())
	incomingRoutes.GET("/comments", middleware.Authentication(), controller.GetCommentList())
    incomingRoutes.PUT("/comments/:id", middleware.Authentication(), controller.UpdateComment())
	incomingRoutes.DELETE("/comments/:id", middleware.Authentication(), controller.DeleteComment())
}