package repository

import (
	"announcement-api/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindAllByRole(ctx context.Context, role entity.Role) ([]entity.User, error)
}
