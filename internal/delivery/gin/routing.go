package gin

import (
	"github.com/DrusGalkin/forum-auth-grpc/pkg/client"
	"github.com/DrusGalkin/forum-client/internal/delivery/gin/handler"
	"github.com/DrusGalkin/forum-client/internal/usecase"
	"github.com/DrusGalkin/forum-client/pkg/wsserver"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/DrusGalkin/forum-client/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(P usecase.PostUseCase, T usecase.ThreadUseCase, authClient *client.AuthClient, hub *wsserver.Hub) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	forumHandler := NewForumHandler(P, T)
	go hub.Run()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v2")
	{
		api.GET("/threads", forumHandler.GetAllThread)
		api.GET("/thread/:id", forumHandler.GetThreadByID)

		authGroup := api.Group("")
		authGroup.Use(handler.AuthMiddleware(authClient))
		{
			authGroup.POST("/threads", forumHandler.CreateThread)
			authGroup.POST("/threads/posts", forumHandler.CreatePost)

			authGroup.GET("/threads/user/:id", forumHandler.GetThreadsByUserID)
			authGroup.GET("/posts/user/:id", forumHandler.GetPostsByUserID)

			authGroup.DELETE("/posts/:id", forumHandler.DeletePostByID)
			authGroup.DELETE("/threads/:id", forumHandler.DeleteTheadByID)

			authGroup.PUT("/threads", forumHandler.EditThread)

			api.GET("/ws/threads/:id", hub.ThreadChat)
		}
	}

	return router
}
