package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

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
}

func (psh *planeSeatHandler) Index(c echo.Context) (error) {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)

	ctx := c.Request().Context()
	list, err := psh.planeSeatUsecase.Fetch(ctx, int64(num))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseError{err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}
