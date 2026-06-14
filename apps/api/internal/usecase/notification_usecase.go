package usecase

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationUsecase struct {
	notifRepo repository.NotificationRepository
}

func NewNotificationUsecase(notifRepo repository.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{notifRepo: notifRepo}
}

func (u *NotificationUsecase) List(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {
	return u.notifRepo.FindByUser(ctx, userID)
}

func (u *NotificationUsecase) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return u.notifRepo.CountUnread(ctx, userID)
}

func (u *NotificationUsecase) MarkRead(ctx context.Context, notifID uuid.UUID, userID uuid.UUID) error {
	affected, err := u.notifRepo.MarkRead(ctx, notifID, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

func (u *NotificationUsecase) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return u.notifRepo.MarkAllRead(ctx, userID)
}
