package middleware

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/labstack/echo"
    jwt "github.com/dgrijalva/jwt-go"

    "shellrean.com/airlines/domain"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
    jwtConfig       domain.JwtConfig
}

// ResponseError represent the response error struct
type ResponseError struct {
    Message string `json:"message"`
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("Access-Control-Allow-Origin", "*")
        return next(c)
    }
}

// Authenticate will handle the Authentication middleware
func (m *GoMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authorizationHeader := c.Request().Header.Get("Authorization")

        unauthorizedMessage := map[string]interface{}{
            "message": "unauthorized",
        }

        if !strings.Contains(authorizationHeader, "Bearer") {
            return c.JSON(http.StatusUnauthorized, unauthorizedMessage)
        }

        signarureKey := []byte(m.jwtConfig.JWT_Key)

        tokenString := strings.Replace(authorizationHeader, "Bearer ","", -1)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Signing method invalid")
            } else if method != JWT_SIGNING_METHOD {
                return nil, fmt.Errorf("Signing method invalid")
            }

            return signarureKey, nil
        })

        if err != nil {
            return c.JSON(http.StatusUnauthorized, ResponseError{err.Error()})
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            return c.JSON(http.StatusUnauthorized, unauthorizedMessage)
        }

        c.Set("userInfo", claims)

        return next(c)
    }
}

// InitMiddleware initialize the middleware
func InitMiddleware(config domain.JwtConfig) *GoMiddleware {
    return &GoMiddleware{
        jwtConfig:      config,
    }
}