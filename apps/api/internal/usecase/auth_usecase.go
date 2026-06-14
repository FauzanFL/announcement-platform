package usecase

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already exists")
)

type AuthUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (u *AuthUsecase) Register(ctx context.Context, username, password string, role entity.Role) (*entity.User, error) {
	if role != entity.RoleAdmin && role != entity.RoleUser {
		role = entity.RoleUser
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: username,
		Password: string(hashed),
		Role:     role,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, ErrUsernameTaken
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, username, password string) (string, *entity.User, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}
