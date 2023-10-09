package main

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

var data = struct {
	SessionTime     float64
	FrameIdentifier int64
	LapDistance     float64
	LapNumber       int
	Throttle        float64
	Brk             float64
	Gear            int
	Rpm             int
	Drs             int
	Speed           int
	Delta           float64
	WorldX          float64
	WorldY          float64
	WorldZ          float64
	LapDeltaTime    float64
}{}

func read(redisLogger *logger.RedisClient) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		lapDistance, err := redisLogger.GetTs("lapDistance")
		if err != nil {
			panic(err)
		}

		lapNumber, err := redisLogger.GetTs("lapNumber")
		if err != nil {
			panic(err)
		}

		frameIdentifier, err := redisLogger.GetTs("frameIdentifier")
		if err != nil {
			panic(err)
		}

		throttle, err := redisLogger.GetTs("throttle")
		if err != nil {
			panic(err)
		}

		brk, err := redisLogger.GetTs("break")
		if err != nil {
			panic(err)
		}

		gear, err := redisLogger.GetTs("gear")
		if err != nil {
			panic(err)
		}

		rpm, err := redisLogger.GetTs("rpm")
		if err != nil {
			panic(err)
		}

		drs, err := redisLogger.GetTs("drs")
		if err != nil {
			panic(err)
		}

		speed, err := redisLogger.GetTs("speed")
		if err != nil {
			panic(err)
		}

		worldPositionX, err := redisLogger.GetTs("worldPositionX")
		if err != nil {
			panic(err)
		}
		worldPositionY, err := redisLogger.GetTs("worldPositionY")
		if err != nil {
			panic(err)
		}
		worldPositionZ, err := redisLogger.GetTs("worldPositionZ")
		if err != nil {
			panic(err)
		}

		lapDeltaTime, err := redisLogger.GetTs("lapDeltaTime")
		if err != nil {
			panic(err)
		}

		data.FrameIdentifier = int64(frameIdentifier.Value)
		data.LapNumber = int(int64(lapNumber.Value))
		data.LapDistance = lapDistance.Value
		data.Throttle = throttle.Value
		data.Brk = brk.Value
		data.Gear = int(int64(gear.Value))
		data.Rpm = int(int64(rpm.Value))
		data.Drs = int(int64(drs.Value))
		data.Speed = int(int64(speed.Value))
		data.WorldX = worldPositionX.Value
		data.WorldY = worldPositionY.Value
		data.WorldZ = worldPositionZ.Value
		data.LapDeltaTime = lapDeltaTime.Value
	}
}

func main() {
	redisLogger := logger.NewRedisClient(context.Background(), "localhost:6379", 0)
	defer redisLogger.CloseListen()
	go read(redisLogger)
	fmt.Println("start Listening\n\nwaiting inputs")

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, data)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
