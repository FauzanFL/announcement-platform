package usecase

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
