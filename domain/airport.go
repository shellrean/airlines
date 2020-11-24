package domain

import (
    "context"
    "time"
)

// Airport...
type Airport struct {
    ID          int64       `json:"id"`
    Code        string      `json:"code"`
    Name        string      `json:"name"`
    CreatedAt   time.Time   `json:"-"`
    UpdatedAt   time.Time   `json:"-"`
}

// AirportRepository represent airport's repository
type AirportRepository interface {
    Fetch(ctx context.Context, num int64) ([]Airport, error)
    GetByID(ctx context.Context, id int64) (Airport, error)
    Store(ctx context.Context, airport *Airport) (error)
    Update(ctx context.Context, airport *Airport) (error)
    Delete(ctx context.Context, id int64) (error)
}

// AirportUsecase represent airport's usecase
type AirportUsecase interface {
    Fetch(ctx context.Context, num int64) ([]Airport, error)
    GetByID(ctx context.Context, id int64) (Airport, error)
    Store(ctx context.Context, airport *Airport) (error)
    Update(ctx context.Context, airport *Airport) (error)
    Delete(ctx context.Context, id int64) (error)
}