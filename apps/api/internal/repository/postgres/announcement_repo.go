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
	if err := r.db.WithContext(ctx).Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}
