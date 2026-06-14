package postgres

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) CreateOne(ctx context.Context, notif *entity.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *notificationRepo) CreateBatch(ctx context.Context, notifs []entity.Notification) error {
	if len(notifs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&notifs).Error
}

func (r *notificationRepo) FindByUser(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {
	var notifs []entity.Notification
	err := r.db.WithContext(ctx).
		Preload("Announcement").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&notifs).Error

	return notifs, err
}

func (r *notificationRepo) ExistsForUserAndAnnouncement(ctx context.Context, userID, announcementID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND announcement_id = ?", userID, announcementID).Count(&count).Error

	return count > 0, err
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
}

func (r *notificationRepo) MarkRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ? AND user_id = ?", id, userID).Updates(map[string]interface{}{"is_read": true, "read_at": now})

	return result.RowsAffected, result.Error
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ?", userID).Updates(map[string]interface{}{"is_read": true}).Error
}

func (r *notificationRepo) DeleteByAnnouncementID(ctx context.Context, announcementID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("announcement_id = ?", announcementID).Delete(&entity.Notification{}).Error
}
