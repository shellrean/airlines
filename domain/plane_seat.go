package domain

import (
	"context"
	"time"
)

type M map[string]interface{}

type PlaneSeat struct {
	ID 			int64 		`json:"id"`
	Plane 		Plane 		`json:"plane"`
	SeatClass 	M			`json:"seat_class"`
	CreatedAt 	time.Time 	`json:"-"`
	UpdatedAt 	time.Time 	`json:"-"`
}

type PlaneSeatRepository interface {
	Fetch(ctx context.Context, num int64) ([]PlaneSeat, error)
	GetByPlaneId(ctx context.Context, id int64) ([]PlaneSeat, error)
	Store(ctx context.Context, seat *PlaneSeat) (error)
	Update(ctx context.Context, seat *PlaneSeat) (error)
	Delete(ctx context.Context, id int64) (error)
}

type PlaneSeatUsecase interface {
	Fetch(ctx context.Context, num int64) ([]PlaneSeat, error)
}