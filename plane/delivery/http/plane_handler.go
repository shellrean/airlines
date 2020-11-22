package http

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo"

    "shellrean.com/airlines/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
    Message string `json:"message"`
}

// PlaneHandler represent the httphandler for plane
type PlaneHandler struct {
    PlaneUsecase        domain.PlaneUsecase
}

// NewPlaneHandler will initialize the plane resources endpoint
func NewPlaneHandler(e *echo.Echo, pu domain.PlaneUsecase) {
    handler := &PlaneHandler{
        PlaneUsecase: pu,
    }

    e.GET("/planes", handler.FetchPlanes)
    e.POST("/planes", handler.CreatePlane)
    e.GET("/planes/:id", handler.GetPlaneByID)
    e.PUT("/planes/:id", handler.UpdatePlane)
    e.DELETE("/planes/:id", handler.DeletePlane)
}

func (p *PlaneHandler) FetchPlanes(c echo.Context) error {
    numS := c.QueryParam("num")
    num, _ := strconv.Atoi(numS)

    ctx := c.Request().Context()
    listPlane, err := p.PlaneUsecase.Fetch(ctx, int64(num))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
    }

    return c.JSON(http.StatusOK, listPlane)
}

func (p *PlaneHandler) CreatePlane(c echo.Context) (err error) {
    var plane domain.Plane
    err = c.Bind(&plane)
    if err != nil {
        return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
    }

    ctx := c.Request().Context()

    err = p.PlaneUsecase.Store(ctx, &plane)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
    }

    return c.JSON(http.StatusCreated, plane)
}

func (p *PlaneHandler) GetPlaneByID(c echo.Context) (err error) {
    idP, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
    }

    id := int64(idP)

    ctx := c.Request().Context()
    res, err := p.PlaneUsecase.GetByID(ctx, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
    }

    return c.JSON(http.StatusOK, res)
}

func (p *PlaneHandler) UpdatePlane(c echo.Context) (err error) {
    idP, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusNotFound, ResponseError{domain.ErrNotFound.Error()})
    }

    id := int64(idP)

    var plane domain.Plane
    err = c.Bind(&plane)
    if err != nil {
        return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
    }

    plane.ID = id

    ctx := c.Request().Context()

    err = p.PlaneUsecase.Update(ctx, &plane)
    if err != nil {
        return c.JSON(http.StatusInternalServerError,ResponseError{err.Error()})
    }

    return c.JSON(http.StatusOK, plane)
}

func (p *PlaneHandler) DeletePlane(c echo.Context) (err error) {
    idP, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusNotFound, ResponseError{domain.ErrNotFound.Error()})
    }

    id := int64(idP)

    ctx := c.Request().Context()

    err = p.PlaneUsecase.Delete(ctx, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
    }

    return c.NoContent(http.StatusNoContent)
}