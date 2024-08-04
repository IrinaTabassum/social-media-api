package main

import (
    "os"

    // middleware "practice/social-media-api/middleware"
    routes "social-media-api/routes"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    "social-media-api/docs"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        port = "8000"
    }
    
    //router
    router := gin.New()
    router.Use(gin.Logger())
    routes.AuthRoutes(router)
    routes.PostRoutes(router)
    routes.CommentRoutes(router)
    routes.LikeRoutes(router)

    //swagger
    docs.SwaggerInfo.BasePath = "/"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    

    router.Run(":" + port)
}