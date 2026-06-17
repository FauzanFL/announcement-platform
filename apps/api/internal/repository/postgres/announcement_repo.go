package postgres

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type announcementRepo struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) repository.AnnouncementRepository {
	return &announcementRepo{db: db}
}

func (r *announcementRepo) Create(ctx context.Context, announcement *entity.Announcement) error {
	return r.db.WithContext(ctx).Create(announcement).Error
}

func (r *announcementRepo) Update(ctx context.Context, announcement *entity.Announcement) error {
	return r.db.WithContext(ctx).Save(announcement).Error
}

func (r *announcementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Announcement{}).Error
}

func (r *announcementRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Announcement, error) {
	var announcement entity.Announcement
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&announcement).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *announcementRepo) FindAll(ctx context.Context) ([]entity.Announcement, error) {
	var announcements []entity.Announcement
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}

func (r *announcementRepo) FindAllWithReadStatus(ctx context.Context, userID uuid.UUID) ([]repository.AnnouncementWithStatus, error) {
    type row struct {
        entity.Announcement
        IsRead bool
    }

    var rows []row
    err := r.db.WithContext(ctx).Raw(`
        SELECT a.*, 
               CASE WHEN n.user_id IS NOT NULL THEN true ELSE false END AS is_read
        FROM announcements a
        LEFT JOIN notifications n
               ON n.announcement_id = a.id AND n.user_id = ?
        ORDER BY a.created_at DESC
    `, userID).Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    result := make([]repository.AnnouncementWithStatus, len(rows))
    for i, r := range rows {
        result[i] = repository.AnnouncementWithStatus{
            Announcement: r.Announcement,
            IsRead:       r.IsRead,
        }
    }
    return result, nil
}

func (r *announcementRepo) CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
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