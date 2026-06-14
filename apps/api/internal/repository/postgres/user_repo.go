package postgres

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAllByRole(ctx context.Context, role entity.Role) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
