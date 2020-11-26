package domain

import (
	"context"
	"time"
)

type SeatClass struct {
	Business	uint8 		`json:"business"`
	Exclusive	uint8 		`json:"exclusive"`
	Economic 	uint8 		`json:"economic"`
}

type PlaneSeat struct {
	ID 			int64 		`json:"id"`
	Plane 		Plane 		`json:"plane"`
	SeatClass 	SeatClass	`json:"seat_class"`
	CreatedAt 	time.Time 	`json:"-"`
	UpdatedAt 	time.Time 	`json:"-"`
}

type PlaneSeatRepository interface {
	Fetch(ctx context.Context, num int64) ([]PlaneSeat, error)
	GetByPlaneID(ctx context.Context, num int64, id int64) ([]PlaneSeat, error)
	Store(ctx context.Context, seat *PlaneSeat) (error)
	Update(ctx context.Context, seat *PlaneSeat) (error)
	Delete(ctx context.Context, id int64) (error)
}

type PlaneSeatUsecase interface {
	Fetch(ctx context.Context, num int64) ([]PlaneSeat, error)
	Store(ctx context.Context, seat *PlaneSeat) (error)
	GetByPlaneID(ctx context.Context, num int64, id int64) ([]PlaneSeat, error)
}