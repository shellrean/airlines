package usecase

import (
    "context"
    "time"

    "shellrean.com/airlines/domain"
)

type planeUsecase struct {
    planeRepo       domain.PlaneRepository
    contextTimeout  time.Duration
}

// NewPlaneUsecase will create an planeUsecase object represent of domain.PlaneUsecase interface
func NewPlaneUsecase(p domain.PlaneRepository, timeout time.Duration) domain.PlaneUsecase {
    return &planeUsecase{
        planeRepo:      p,
        contextTimeout: timeout,
    }
}

func (p *planeUsecase) Fetch(c context.Context, num int64) (res []domain.Plane, err error) {
    if num == 0 {
        num = 10
    }

    ctx, cancel := context.WithTimeout(c, p.contextTimeout)
    defer cancel()

    res, err = p.planeRepo.Fetch(ctx, num)
    if err != nil {
        return nil, err
    }

    return
}

func (p *planeUsecase) GetByID(c context.Context, id int64) (res domain.Plane, err error) {
    ctx, cancel := context.WithTimeout(c, p.contextTimeout)
    defer cancel()

    res, err = p.planeRepo.GetByID(ctx, id)
    if err != nil {
        return domain.Plane{}, err
    }

    return
}

func (p *planeUsecase) Store(c context.Context, dp *domain.Plane) (err error) {
    ctx, cancel := context.WithTimeout(c, p.contextTimeout)
    defer cancel()

    dp.CreatedAt = time.Now()
    dp.UpdatedAt = time.Now()

    err = p.planeRepo.Store(ctx, dp)
    return
}

func (p *planeUsecase) Update(c context.Context, dp *domain.Plane) (err error) {
    ctx, cancel := context.WithTimeout(c, p.contextTimeout)
    defer cancel()

    dp.UpdatedAt = time.Now()

    return p.planeRepo.Update(ctx, dp)
}

func (p *planeUsecase) Delete(c context.Context, id int64) (err error) {
    ctx, cancel := context.WithTimeout(c, p.contextTimeout)
    defer cancel()

    existedPlane, err := p.planeRepo.GetByID(ctx, id)
    if err != nil {
        return
    }

    if existedPlane == (domain.Plane{}) {
        return domain.ErrNotFound
    }

    return p.planeRepo.Delete(ctx, id)
}