package usecase

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrForbidden = errors.New("forbidden: admin role required")
	ErrNotFound  = errors.New("announcement not found")
)

const ChannelAnnouncement = "announcement_channel"

type AnnouncementUsecase struct {
	annRepo   repository.AnnouncementRepository
	userRepo  repository.UserRepository
	pubsub    repository.PubSubRepository
}

func NewAnnouncementUsecase(
	annRepo repository.AnnouncementRepository,
	userRepo repository.UserRepository,
	pubsub repository.PubSubRepository,
) *AnnouncementUsecase {
	return &AnnouncementUsecase{annRepo: annRepo, userRepo: userRepo, pubsub: pubsub}
}

func (u *AnnouncementUsecase) requireAdmin(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	if !user.IsAdmin() {
		return nil, ErrForbidden
	}
	return user, nil
}

func (u *AnnouncementUsecase) Create(ctx context.Context, requesterID uuid.UUID, title, content string) (*entity.Announcement, error) {
	if _, err := u.requireAdmin(ctx, requesterID); err != nil {
		return nil, err
	}

	ann := &entity.Announcement{Title: title, Content: content, CreatedBy: requesterID}
	if err := u.annRepo.Create(ctx, ann); err != nil {
		return nil, err
	}


	u.publishEvent(ctx, entity.AnnouncementEvent{Type: "created", Announcement: ann})

	return ann, nil
}

func (u *AnnouncementUsecase) Update(ctx context.Context, requesterID, announcementID uuid.UUID, title, content string) (*entity.Announcement, error) {
	if _, err := u.requireAdmin(ctx, requesterID); err != nil {
		return nil, err
	}

	ann, err := u.annRepo.FindByID(ctx, announcementID)
	if err != nil {
		return nil, ErrNotFound
	}

	ann.Title = title
	ann.Content = content
	if err := u.annRepo.Update(ctx, ann); err != nil {
		return nil, err
	}

	u.publishEvent(ctx, entity.AnnouncementEvent{Type: "updated", Announcement: ann})

	return ann, nil
}

func (u *AnnouncementUsecase) Delete(ctx context.Context, requesterID, announcementID uuid.UUID) error {
	if _, err := u.requireAdmin(ctx, requesterID); err != nil {
		return err
	}

	if err := u.annRepo.Delete(ctx, announcementID); err != nil {
		return err
	}

	u.publishEvent(ctx, entity.AnnouncementEvent{Type: "deleted", ID: announcementID.String()})

	return nil
}

func (u *AnnouncementUsecase) List(ctx context.Context) ([]entity.Announcement, error) {
	return u.annRepo.FindAll(ctx)
}

func (u *AnnouncementUsecase) ListWithReadStatus(ctx context.Context, userID uuid.UUID) ([]repository.AnnouncementWithStatus, error) {
    return u.annRepo.FindAllWithReadStatus(ctx, userID)
}

func (u *AnnouncementUsecase) Get(ctx context.Context, id uuid.UUID) (*entity.Announcement, error) {
	ann, err := u.annRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return ann, nil
}

func (u *AnnouncementUsecase) publishEvent(ctx context.Context, event entity.AnnouncementEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		return
	}
	_ = u.pubsub.Publish(ctx, ChannelAnnouncement, payload)
}

func (u *AnnouncementUsecase) SubscribeToEvents(ctx context.Context) (<-chan string, func(), error) {
	return u.pubsub.Subscribe(ctx, ChannelAnnouncement)
}
