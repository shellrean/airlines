package http

import (
	"net/http"
	"strconv"
	
	"github.com/labstack/echo/v4"

	"shellrean.com/airlines/domain"
	"shellrean.com/airlines/user/delivery/middleware"
)

type responseError struct {
	Message			string
}

type airportHandler struct {
	airportUsecase 		domain.AirportUsecase
}

func NewAirportHandler(e *echo.Echo, airUcase domain.AirportUsecase, middl *middleware.GoMiddleware) {
	hander := &airportHandler {
		airportUsecase:			airUcase,
	}

	e.GET("/airports", hander.Index, middl.Auth)
	e.GET("/airports/:id", hander.Show, middl.Auth)
	e.POST("/airports", hander.Store, middl.Auth)
	e.PUT("/airports/:id", hander.Update, middl.Auth)
	e.DELETE("/airports/:id", hander.Destroy, middl.Auth)
}

func (a *airportHandler) Index(c echo.Context) (error) {
	num := c.QueryParam("num")
	limit, _ := strconv.Atoi(num)
	
	ctx := c.Request().Context()
	list, err := a.airportUsecase.Fetch(ctx, int64(limit))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (a *airportHandler) Show(c echo.Context) (error) {
	idS := c.Param("id")
	
	id, err := strconv.Atoi(idS)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseError{domain.ErrBadParamInput.Error()})
	}

	ctx := c.Request().Context()
	res, err := a.airportUsecase.GetByID(ctx, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (a *airportHandler) Store(c echo.Context) (error) {
	var airport domain.Airport
	err := c.Bind(&airport)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseError{err.Error()})
	}

	ctx := c.Request().Context()
	err = a.airportUsecase.Store(ctx, &airport)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusCreated, airport)
}

func (a *airportHandler) Update(c echo.Context) (error) {
	idS := c.Param("id")
	
	id, err := strconv.Atoi(idS)
	if err != nil {
		return c.JSON(http.StatusNotFound, responseError{domain.ErrNotFound.Error()})
	}

	var airport domain.Airport
	err = c.Bind(&airport)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, responseError{err.Error()})
	}
	airport.ID = int64(id)

	ctx := c.Request().Context()
	err = a.airportUsecase.Update(ctx, &airport)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusOK, airport)
}

func (a *airportHandler) Destroy(c echo.Context) (error) {
	idS := c.Param("id")
	
	id, err := strconv.Atoi(idS)
	if err != nil {
		return c.JSON(http.StatusNotFound, responseError{err.Error()})
	}

	ctx := c.Request().Context()
	err = a.airportUsecase.Delete(ctx, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}