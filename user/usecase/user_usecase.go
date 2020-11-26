package usecase

import (
    "context"
    "time"
    "fmt"

    // "golang.org/x/sync/errgroup"

    "shellrean.com/airlines/domain"
    "golang.org/x/crypto/bcrypt"
    jwt "github.com/dgrijalva/jwt-go"
)

type userUsecase struct {
    userRepo        domain.UserRepository
    contextTimeout  time.Duration
    jwtConfig       domain.JwtConfig  
}

type MyClaims struct {
    jwt.StandardClaims
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Email       string `json:"email"`
}

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

// NewUserUsecase will create new on articleUsecase object represent of domain.UserUsecase interface
func NewUserUsecase(u domain.UserRepository, timeout time.Duration, config domain.JwtConfig) domain.UserUsecase {
    return &userUsecase{
        userRepo:           u,
        contextTimeout:     timeout,
        jwtConfig:          config,
    }
}

func (u *userUsecase) GetAuthentication(c context.Context, us domain.User) (string, error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.userRepo.GetByEmail(ctx, us.Email)
    if err != nil {
        return "", domain.InvalidUser
    }

    err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(us.Password))
    if err != nil {
        return "", domain.InvalidPassword
    }

    claims := MyClaims{
        StandardClaims: jwt.StandardClaims{
            Issuer: u.jwtConfig.AppName,
            ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour).Unix(),
        },
        ID:         res.ID,
        Name:       res.Name,
        Email:      res.Email,
    }

    token := jwt.NewWithClaims(
        JWT_SIGNING_METHOD,
        claims,
    )

    signedToken, err := token.SignedString([]byte(u.jwtConfig.JWT_Key))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}

func (u *userUsecase) GetAutenticatedUserData(c context.Context, email string) (domain.User, error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return domain.User{}, err
    }

    return res, nil
} 

func (u *userUsecase) RefreshToken(c context.Context, tokenString string) (string, error) {
    _, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    claims := &MyClaims{}
    signarureKey := []byte(u.jwtConfig.JWT_Key)

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Signing method invalid")
        } else if method != JWT_SIGNING_METHOD {
            return nil, fmt.Errorf("Signing method invalid")
        }

        return signarureKey, nil
    })

    expirationTime := time.Now().Add(time.Duration(5) * time.Hour).Unix()
    claims.ExpiresAt = expirationTime
    token = jwt.NewWithClaims(
        JWT_SIGNING_METHOD,
        claims,
    )

    signedToken, err := token.SignedString([]byte(u.jwtConfig.JWT_Key))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}