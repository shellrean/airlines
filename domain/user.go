package domain

import (
    "context"
    "time"
)

// User...
type User struct {
    ID          int64       `json:"id"`
    Email       string      `json:"email" validate:"required,email"`
    Name        string      `json:"name" validate:"required"`
    Password    string      `json:"password" validate:"required"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// UserUsecase respresent the user's usecase
type UserUsecase interface {
    GetAuthentication(ctx context.Context, us User) (string, error)
    GetAutenticatedUserData(ctx context.Context, email string) (User, error)
    RefreshToken(ctx context.Context, token string) (string, error)
}

// UserRepository represent the user's repository 
type UserRepository interface {
    GetByEmail(ctx context.Context, email string) (User, error)
}