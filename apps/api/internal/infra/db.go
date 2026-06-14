package infra

import (
	"announcement-api/internal/config"
	"announcement-api/internal/domain/entity"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)

	if err := db.AutoMigrate(&entity.User{}, &entity.Announcement{}, &entity.Notification{}); err != nil {
		return nil, err
	}

	return db, nil
}
