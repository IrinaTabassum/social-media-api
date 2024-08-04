package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/likes", middleware.Authentication(), controller.CreateLike())
    incomingRoutes.GET("/likes/:id", middleware.Authentication(), controller.GetLikeByID())
	incomingRoutes.DELETE("/likes/:id", middleware.Authentication(), controller.DeleteLike())
	incomingRoutes.GET("/likes", middleware.Authentication(), controller.GetLikeList())
}