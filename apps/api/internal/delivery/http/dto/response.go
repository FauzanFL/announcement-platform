package dto

import "github.com/google/uuid"

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type MessageResponse struct {
	Message string `json:"message" example:"operation success"`
}

type AuthResponse struct {
	Message string       `json:"message" example:"operation success"`
	User    UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"       example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string    `json:"username"  example:"john_doe"`
	Role     string    `json:"role"      example:"user"`
}

type RegisterResponse struct {
	ID       uuid.UUID `json:"id"       example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string    `json:"username"  example:"john_doe"`
	Role     string    `json:"role"      example:"user"`
}

type UnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count" example:"5"`
}
