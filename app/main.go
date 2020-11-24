package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/url"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/labstack/echo"
    "github.com/spf13/viper"

    "shellrean.com/airlines/domain"
    _userHttpDelivery "shellrean.com/airlines/user/delivery/http"
    _userDeliveryMiddleware "shellrean.com/airlines/user/delivery/middleware"
    _userRepo "shellrean.com/airlines/user/repository/mysql"
    _userUsecase "shellrean.com/airlines/user/usecase"

    _planeHttpDelivery "shellrean.com/airlines/plane/delivery/http"
    _planeRepo "shellrean.com/airlines/plane/repository/mysql"
    _planeUsecase "shellrean.com/airlines/plane/usecase"

    _airportHttpDelivery "shellrean.com/airlines/airport/delivery/http"
    _airportRepo "shellrean.com/airlines/airport/repository/mysql"
    _airportUsecase "shellrean.com/airlines/airport/usecase"
)

func init() {
    viper.SetConfigFile(`config.json`)
    err := viper.ReadInConfig()
    if err != nil {
        panic(err)
    }

    if viper.GetBool(`debug`) {
        log.Println("Service RUN on DEBUG mode")
    }
}

func main() {
    dbHost := viper.GetString(`database.host`)
    dbPort := viper.GetString(`database.port`)
    dbUser := viper.GetString(`database.user`)
    dbPass := viper.GetString(`database.pass`)
    dbName := viper.GetString(`database.name`)

    connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

    val := url.Values{}
    val.Add("parseTime", "1")
    val.Add("loc","Asia/Jakarta")
    dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
    dbConn, err := sql.Open(`mysql`, dsn)

    if err != nil {
        log.Fatal(err)
    }
    
    err = dbConn.Ping()
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        err := dbConn.Close()
        if err != nil {
            log.Fatal(err)
        }
    }()
    
    jwtConfig := domain.JwtConfig{
        AppName : viper.GetString(`application_name`),
        JWT_Key: viper.GetString(`jwt_signature_key`),
    }

    // Initialize handler
    e := echo.New()

    // Initialize middleware
    middl := _userDeliveryMiddleware.InitMiddleware(jwtConfig)
    e.Use(middl.CORS)

    // Initialize repository
    userRepo := _userRepo.NewMysqlUserRepository(dbConn)
    planeRepo := _planeRepo.NewMysqlPlaneRepository(dbConn)
    airportRepo := _airportRepo.NewMysqlAirportRepository(dbConn)

    timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

    // Initialize usecase
    userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext, jwtConfig)
    _userHttpDelivery.NewUserHandler(e, userUsecase, middl)

    planeUsecase := _planeUsecase.NewPlaneUsecase(planeRepo, timeoutContext)
    _planeHttpDelivery.NewPlaneHandler(e, planeUsecase, middl)

    airportUsecase := _airportUsecase.NewAirportUsecase(airportRepo, timeoutContext)
    _airportHttpDelivery.NewAirportHandler(e, airportUsecase, middl)

    // Serve Server
    log.Fatal(e.Start(viper.GetString("server.address")))
}