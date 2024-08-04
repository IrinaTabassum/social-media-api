package routes

import (
	controller "social-media-api/controllers"
	middleware "social-media-api/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/posts/:id", controller.GetPostByID())
	incomingRoutes.GET("/posts", controller.ListPosts())
	incomingRoutes.POST("/posts" ,middleware.Authentication(), controller.CreatePost())
    incomingRoutes.PUT("/posts/:id" ,middleware.Authentication(), controller.UpdatePost())
	incomingRoutes.DELETE("/posts/:id" ,middleware.Authentication(), controller.DeletePost())
}