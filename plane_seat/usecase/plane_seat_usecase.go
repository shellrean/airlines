package usecase

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
	
	"shellrean.com/airlines/domain"
)

type planeSeatUsecase struct {
	planeRepo 			domain.PlaneRepository
	planeSeatRepo 		domain.PlaneSeatRepository
	contextTimeout 		time.Duration
}

func NewPlaneSeatUsecase(pr domain.PlaneRepository, psr domain.PlaneSeatRepository, timeout time.Duration) domain.PlaneSeatUsecase{
	return &planeSeatUsecase{
		planeRepo: 		pr,
		planeSeatRepo:	psr,
		contextTimeout:	timeout,
	}
}

func (psu *planeSeatUsecase) fillPlaneDetails(c context.Context, data []domain.PlaneSeat) ([]domain.PlaneSeat, error) {
	g, ctx := errgroup.WithContext(c)

	// Get the plane's id
	mapPlanes := map[int64]domain.Plane{}

	for _, planeSeat := range data {
		mapPlanes[planeSeat.Plane.ID] = domain.Plane{}
	}

	// Using goroutine to fetch the plane's detail
	chanPlane := make(chan domain.Plane)
	for planeID := range mapPlanes {
		planeID := planeID
		g.Go(func() error {
			res, err := psu.planeRepo.GetByID(ctx, planeID)
			if err != nil {
				return err
			}
			chanPlane <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			return
		}
		close(chanPlane)
	}()

	for plane := range chanPlane {
		if plane != (domain.Plane{}) {
			mapPlanes[plane.ID] = plane
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Merge the plane's data
	for index, item := range data {
		if a, ok := mapPlanes[item.Plane.ID]; ok {
			data[index].Plane = a
		}
	}

	return data, nil
}

func (psu *planeSeatUsecase) Fetch(c context.Context, num int64) ([]domain.PlaneSeat, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, psu.contextTimeout)
	defer cancel()

	res, err := psu.planeSeatRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	res, err = psu.fillPlaneDetails(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (psu *planeSeatUsecase) GetByPlaneID(c context.Context, num int64, id int64) ([]domain.PlaneSeat, error) {
	if num == 0 {
		num = 10
	}
	
	ctx, cancel := context.WithTimeout(c, psu.contextTimeout)
	defer cancel()

	res, err := psu.planeSeatRepo.GetByPlaneID(ctx, num, id)
	if err != nil {
		return nil, err
	}

	res, err = psu.fillPlaneDetails(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (psu *planeSeatUsecase) Store(c context.Context, planeseat *domain.PlaneSeat) (error) {
	ctx, cancel := context.WithTimeout(c, psu.contextTimeout)
	defer cancel()

	plane, err := psu.planeRepo.GetByID(ctx, planeseat.Plane.ID)
	if err != nil {
		return err
	}

	planeseat.Plane = plane
	planeseat.CreatedAt = time.Now()
	planeseat.UpdatedAt = time.Now()

	err = psu.planeSeatRepo.Store(ctx, planeseat)
	if err != nil {
		return err
	}

	return nil
}