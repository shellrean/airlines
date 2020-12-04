package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"shellrean.com/airlines/domain"
	"shellrean.com/airlines/user/delivery/middleware"
)

type M map[string]interface{}

type responseError struct {
	Message		string	`json:"message"`
}

type planeSeatHandler struct {
	planeSeatUsecase 		domain.PlaneSeatUsecase
}

func NewPlaneSeatHandler(e *echo.Echo, psu domain.PlaneSeatUsecase, mdl *middleware.GoMiddleware) {
	handler := &planeSeatHandler{
		planeSeatUsecase:		psu,
	}

	e.GET("/plane-seats", handler.Index, mdl.Auth)
	e.POST("/plane-seats", handler.Store, mdl.Auth)
	e.PUT("/plane-sets/:id", handler.Update, mdl.Auth)
}

func (psh *planeSeatHandler) Index(c echo.Context) (error) {
	var list []domain.PlaneSeat
	var err error
	
	numS := c.QueryParam("num")
	plane_idS := c.QueryParam("plane")
	num, _ := strconv.Atoi(numS)
	plane_id, _ := strconv.Atoi(plane_idS)

	ctx := c.Request().Context()
	
	if plane_id == 0 {
		list, err = psh.planeSeatUsecase.Fetch(ctx, int64(num))
	} else {
		list, err = psh.planeSeatUsecase.GetByPlaneID(ctx, int64(num), int64(plane_id))
	}
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (psh *planeSeatHandler) Store(c echo.Context) (error) {
	var planeSeat domain.PlaneSeat
	err := c.Bind(&planeSeat)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseError{err.Error()})
	}

	ctx := c.Request().Context()
	err = psh.planeSeatUsecase.Store(ctx, &planeSeat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusCreated, planeSeat)
}

func (psh *planeSeatHandler) Update(c echo.Context) (error) {
	var planeSeat domain.PlaneSeat
	err := c.Bind(&planeSeat)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseError{err.Error()})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, responseError{domain.ErrNotFound.Error()})
	}
	planeSeat.ID = int64(id)
	ctx := c.Request().Context()

	err = psh.PlaneSeatUsecase.Update(planeSeat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusCreated, planeSeat)
}