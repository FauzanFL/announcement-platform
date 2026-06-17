// @title           Announcement API
// @version         1.0
// @description     API for announcement app. Using Server-Sent Events (SSE), RBAC, and real-time notification.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    API Support
// @contact.email   support@announcement-app.dev
//
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
//
// @host            localhost:8080
// @BasePath        /api
//
// @schemes         http https
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Input your Bearer token: **Bearer &lt;token&gt;**

package main

import (
	"announcement-api/internal/config"
	deliveryHttp "announcement-api/internal/delivery/http"
	"announcement-api/internal/infra"
	"announcement-api/internal/infra/seeder"
	"announcement-api/internal/repository/postgres"
	redisrepo "announcement-api/internal/repository/redis"
	"announcement-api/internal/usecase"
	"log"

	_ "announcement-api/docs"
)

func main() {
	cfg := config.Load()

	db, err := infra.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	seeder.Run(db)

	redisClient := infra.NewRedisClient(cfg)

	userRepo := postgres.NewUserRepository(db)
	annRepo := postgres.NewAnnouncementRepository(db)
	notifRepo := postgres.NewNotificationRepository(db)
	pubsubRepo := redisrepo.NewPubSubRepository(redisClient)

	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	userUC := usecase.NewUserUsecase(userRepo)
	annUC := usecase.NewAnnouncementUsecase(annRepo, userRepo, pubsubRepo)
	notifUC := usecase.NewNotificationUsecase(notifRepo, annRepo)

	r := deliveryHttp.NewRouter(deliveryHttp.Dependencies{
		JWTSecret:      cfg.JWTSecret,
		AuthUC:         authUC,
		UserUC:         userUC,
		AnnouncementUC: annUC,
		NotificationUC: notifUC,
	}, cfg.APIPort)

	log.Printf("server running on :%s", cfg.APIPort)

	if err := r.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
