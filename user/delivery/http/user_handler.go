package http

import (
    "net/http"
    // "strconv"
    "fmt"
    "strings"

    "github.com/labstack/echo/v4"
    validator "gopkg.in/go-playground/validator.v9"
    jwt "github.com/dgrijalva/jwt-go"

    "shellrean.com/airlines/domain"
    "shellrean.com/airlines/user/delivery/middleware"
)

// ResponseError represent the response error struct
type ResponseError struct {
    Message string `json:"message"`
}

// UserHandler respresent the httphandler for user
type UserHandler struct {
    UserUsecase     domain.UserUsecase
}

// NewUserHandler will initialize the articles resources endpoint
func NewUserHandler(e *echo.Echo, us domain.UserUsecase, middl *middleware.GoMiddleware) {
    handler := &UserHandler{
        UserUsecase: us,
    }
    
    e.GET("/user-data", handler.Authenticated, middl.Auth)
    e.POST("/login", handler.Signin)
    e.GET("/refresh-token", handler.Refresh)
}

func (u *UserHandler) Signin(c echo.Context) error {
    var user domain.User
    err := c.Bind(&user)
    if err != nil {
        return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
    }
    user.Name = "ok"

    var ok bool
    if ok, err = isRequestValid(&user); !ok {
        return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
    }

    var token string

    ctx := c.Request().Context()
    token, err = u.UserUsecase.GetAuthentication(ctx, user)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
    }

    res := map[string]interface{}{
        "token": token,
    }

    return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) Authenticated(c echo.Context) error {
    ctx := c.Request().Context()

    userInfo := c.Get("userInfo").(jwt.MapClaims)
    email := userInfo["email"].(string)

    fmt.Println(userInfo)

    user, err := u.UserUsecase.GetAutenticatedUserData(ctx, email)
    if err != nil {
        return c.JSON(getStatusCode(err), ResponseError{err.Error()})
    }

    res := map[string]interface{}{
        "email": user.Email,
        "name": user.Name,
    }

    return c.JSON(http.StatusOK, res)
}

func (u *UserHandler) Refresh(c echo.Context) error {
    ctx := c.Request().Context()

    authorizationHeader := c.Request().Header.Get("Authorization")
    unauthorizedMessage := map[string]interface{}{
        "message": "unauthorized",
    }

    if !strings.Contains(authorizationHeader, "Bearer") {
        return c.JSON(http.StatusUnauthorized, unauthorizedMessage)
    }
    
    tokenString := strings.Replace(authorizationHeader, "Bearer ","", -1)

    token, err := u.UserUsecase.RefreshToken(ctx, tokenString)
    
    if err != nil {
        return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
    }

    res := map[string]interface{}{
        "refresh_token": token,
    }

    return c.JSON(http.StatusOK, res)
}

func isRequestValid(u *domain.User) (bool, error) {
    validate := validator.New()
    err := validate.Struct(u)
    if err != nil {
        return false, err
    }
    return true, nil
}

func getStatusCode(err error) int {
    if err == nil {
        return http.StatusOK
    }

    switch err {
    case domain.ErrInternalServerError:
        return http.StatusInternalServerError
    case domain.ErrNotFound:
        return http.StatusNotFound
    case domain.ErrConflict:
        return http.StatusConflict
    default:
        return http.StatusInternalServerError
    }
}