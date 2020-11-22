package domain

import (
    "context"
    "time"
)

// Plane...
type Plane struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name" validate:"required"`
    PlaneCode   string      `json:"plane_code"`
    SeatSize    int         `json:"seat_size"`
    CreatedAt   time.Time   `json:"-"`
    UpdatedAt   time.Time   `json:"-"` 
}

// PlaneUsecase represent plane's usecase interface
type PlaneUsecase interface {
    Fetch(ctx context.Context, num int64) ([]Plane, error)
    GetByID(ctx context.Context, id int64) (Plane, error)
    Store(ctx context.Context, plane *Plane) (error)
    Update(ctx context.Context, plane *Plane) (error)
    Delete(ctx context.Context, id int64) (error)
}

// PlaneRepository represent plane's repository interface
type PlaneRepository interface {
    Fetch(ctx context.Context, num int64) ([]Plane, error)
    GetByID(ctx context.Context, id int64) (Plane, error)
    Store(ctx context.Context, plane *Plane) (error)
    Update(ctx context.Context, plane *Plane) (error)
    Delete(ctx context.Context, id int64) (error)
}