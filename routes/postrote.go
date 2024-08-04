package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/posts/:id", controller.GetPostByID())
	incomingRoutes.GET("/posts", controller.ListPosts())
    incomingRoutes.Use(middleware.Authentication())
    incomingRoutes.POST("/posts", controller.CreatePost())
    incomingRoutes.PUT("/posts/:id", controller.UpdatePost())
	incomingRoutes.DELETE("/posts/:id", controller.DeletePost())
}