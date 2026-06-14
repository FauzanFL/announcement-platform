package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Role      Role      `gorm:"type:varchar(10);not null;default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
