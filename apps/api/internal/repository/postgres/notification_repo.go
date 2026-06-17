package postgres

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
  "gorm.io/gorm/clause"
)

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) MarkRead(ctx context.Context, announcementID uuid.UUID, userID uuid.UUID) error {
    receipt := entity.Notification{
        UserID:         userID,
        AnnouncementID: announcementID,
        ReadAt:         time.Now(),
    }
    return r.db.WithContext(ctx).
        Clauses(clause.OnConflict{DoNothing: true}).
        Create(&receipt).Error
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
    now := time.Now()
    return r.db.WithContext(ctx).Exec(`
        INSERT INTO notifications (id, user_id, announcement_id, read_at)
        SELECT gen_random_uuid(), ?, a.id, ?
        FROM announcements a
        WHERE NOT EXISTS (
            SELECT 1 FROM notifications n
            WHERE n.user_id = ? AND n.announcement_id = a.id
        )
    `, userID, now, userID).Error
}

func (r *notificationRepo) IsRead(ctx context.Context, announcementID uuid.UUID, userID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entity.Notification{}).
        Where("user_id = ? AND announcement_id = ?", userID, announcementID).
        Count(&count).Error
    return count > 0, err
}

func (r *notificationRepo) ListReadIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]entity.Notification, error) {
    var receipts []entity.Notification
    if err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Find(&receipts).Error; err != nil {
        return nil, err
    }

    result := make(map[uuid.UUID]entity.Notification, len(receipts))
    for _, r := range receipts {
        result[r.AnnouncementID] = r
    }
    return result, nil
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Raw(`
        SELECT COUNT(*)
        FROM announcements a
        WHERE NOT EXISTS (
            SELECT 1 FROM notifications n
            WHERE n.user_id = ? AND n.announcement_id = a.id
        )
    `, userID).Scan(&count).Error
    return count, err
}