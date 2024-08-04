package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
    incomingRoutes.POST("/likes", controller.CreateLike())
    incomingRoutes.GET("/likes/:id", controller.GetLikeByID())
	incomingRoutes.DELETE("/likes/:id", controller.DeleteLike())
	incomingRoutes.GET("/likes", controller.GetLikeList())
}