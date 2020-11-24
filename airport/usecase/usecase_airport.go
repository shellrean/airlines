package usecase

import (
	"context"
	"time"
	
	"shellrean.com/airlines/domain"
)

type airportUsecase struct {
	airportRepo 		domain.AirportRepository
	contextTimeout 		time.Duration
}

// NewAirportUsecase represent domain.AirportUsecase interface
func NewAirportUsecase(airRepo domain.AirportRepository, timeout time.Duration) domain.AirportUsecase {
	return &airportUsecase{
		airportRepo:		airRepo,
		contextTimeout:		timeout,
	}
}

func (a *airportUsecase) Fetch(ctx context.Context, num int64) ([]domain.Airport, error) {
	if num == 0 {
		num = 10
	}
	
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	list, err := a.airportRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (a *airportUsecase) GetByID(ctx context.Context, id int64) (domain.Airport, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res, err := a.airportRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Airport{}, nil
	}

	return res, nil
}

func (a *airportUsecase) Store(ctx context.Context, airport *domain.Airport) (error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	airport.CreatedAt = time.Now()
	airport.UpdatedAt = time.Now()

	return a.airportRepo.Store(ctx, airport)
}

func (a *airportUsecase) Update(ctx context.Context, airport *domain.Airport) (error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	airport.UpdatedAt = time.Now()

	return a.airportRepo.Update(ctx, airport)
}

func (a *airportUsecase) Delete(ctx context.Context, id int64) (error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res, err := a.airportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if res == (domain.Airport{}) {
		return domain.ErrNotFound
	}

	return a.airportRepo.Delete(ctx, id)
}