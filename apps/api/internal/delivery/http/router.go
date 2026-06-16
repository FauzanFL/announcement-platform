package http

import (
	"announcement-api/internal/delivery/http/handler"
	"announcement-api/internal/delivery/http/middleware"
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/usecase"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Dependencies struct {
	JWTSecret string

	AuthUC         *usecase.AuthUsecase
	UserUC         *usecase.UserUsecase
	AnnouncementUC *usecase.AnnouncementUsecase
	NotificationUC *usecase.NotificationUsecase
}

func NewRouter(deps Dependencies, port string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
	}))

	authHandler := handler.NewAuthHandler(deps.AuthUC)
	annHandler := handler.NewAnnouncementHandler(deps.AnnouncementUC)
	notifHandler := handler.NewNotificationHandler(deps.NotificationUC)
	sseHandler := handler.NewSSEHandler(deps.AnnouncementUC, deps.NotificationUC)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to Announcement API"})
	})

	api := r.Group("/api")
	{
		api.GET("/docs/*any", func(ctx *gin.Context) {
			anyParam := ctx.Param("any")

			if anyParam == "" || anyParam == "/" {
				ctx.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
				return
			}

			ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.URL(fmt.Sprintf("http://localhost:%s/api/docs/doc.json", port)),
			)(ctx)
		})

		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		auth := api.Group("/")
		auth.Use(middleware.AuthRequired(deps.JWTSecret))
		auth.Use(middleware.LoadCurrentUser(deps.UserUC))
		{
			auth.GET("/announcements", annHandler.List)
			auth.GET("/announcements/:id", annHandler.Get)

			admin := auth.Group("/")
			admin.Use(middleware.RequireRole(entity.RoleAdmin))
			{
				admin.POST("/announcements", annHandler.Create)
				admin.PUT("/announcements/:id", annHandler.Update)
				admin.DELETE("/announcements/:id", annHandler.Delete)
			}

			auth.GET("/notifications", notifHandler.List)
			auth.GET("/notifications/unread-count", notifHandler.UnreadCount)
			auth.PUT("/notifications/:id/read", notifHandler.MarkRead)
			auth.PUT("/notifications/read-all", notifHandler.MarkAllRead)

			auth.GET("/stream", sseHandler.Stream)
		}
	}

	return r
}
