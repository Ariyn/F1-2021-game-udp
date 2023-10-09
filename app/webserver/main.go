package main

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type LapData struct {
	Number    int
	startedAt int64
	EndedAt   int64
}

var lapHistory = make(map[int]LapData)
var currentLapNumber = -1

type TelemetryData struct {
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
}

var data = TelemetryData{}

func read(redisLogger *logger.RedisClient) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		var errorOccurred bool
		redisLogger.GetCurrentSessionUid()

		lapDistance, err := redisLogger.GetTs("lapDistance")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		lapNumber, err := redisLogger.GetTs("lapNumber")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		if currentLapNumber != int(lapNumber.Value) {
			lapHistory[int(lapNumber.Value)] = LapData{
				Number:    int(lapNumber.Value),
				startedAt: lapNumber.Timestamp,
			}

			if v, ok := lapHistory[int(lapNumber.Value)-1]; ok {
				v.EndedAt = lapNumber.Timestamp
				lapHistory[int(lapNumber.Value)-1] = v
			}
		}

		frameIdentifier, err := redisLogger.GetTs("frameIdentifier")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		throttle, err := redisLogger.GetTs("throttle")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		brk, err := redisLogger.GetTs("break")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		gear, err := redisLogger.GetTs("gear")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		rpm, err := redisLogger.GetTs("rpm")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		drs, err := redisLogger.GetTs("drs")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		speed, err := redisLogger.GetTs("speed")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		worldPositionX, err := redisLogger.GetTs("worldPositionX")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}
		worldPositionY, err := redisLogger.GetTs("worldPositionY")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}
		worldPositionZ, err := redisLogger.GetTs("worldPositionZ")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		lapDeltaTime, err := redisLogger.GetTs("lapDeltaTime")
		if err != nil {
			errorOccurred = true
			//log.Println(err)
		}

		if errorOccurred {
			continue
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

func readSpecificLap(redisLogger *logger.RedisClient, lapNumber int) (data []TelemetryData, err error) {
	redisLogger.GetCurrentSessionUid()

	lapData, err := redisLogger.GetLapHistory(lapNumber)
	if err != nil {
		log.Println(err)
		return
	}

	startedAt := int(lapData.StartedAt)
	endedAt := int(lapData.EndedAt)
	log.Println(startedAt, endedAt)

	dataByTimestamp := make(map[int64]TelemetryData)
	lapDistances, _ := redisLogger.GetTsWithRange("lapDistance", startedAt, endedAt)
	for _, ld := range lapDistances {
		dataByTimestamp[ld.Timestamp] = TelemetryData{
			LapDistance: ld.Value,
		}
	}

	lapNumbers, _ := redisLogger.GetTsWithRange("lapNumber", startedAt, endedAt)
	for _, d := range lapNumbers {
		if v, ok := dataByTimestamp[d.Timestamp]; !ok {
			dataByTimestamp[d.Timestamp] = TelemetryData{
				LapNumber: int(int64(d.Value)),
			}
		} else {
			v.LapNumber = int(int64(d.Value))
			dataByTimestamp[d.Timestamp] = v
		}
	}

	frameIdentifiers, _ := redisLogger.GetTsWithRange("frameIdentifier", startedAt, endedAt)
	for _, fi := range frameIdentifiers {
		if v, ok := dataByTimestamp[fi.Timestamp]; !ok {
			dataByTimestamp[fi.Timestamp] = TelemetryData{
				FrameIdentifier: int64(fi.Value),
			}
		} else {
			v.FrameIdentifier = int64(fi.Value)
			dataByTimestamp[fi.Timestamp] = v
		}
	}

	//throttles, _ := redisLogger.GetTsWithRange("throttle", startedAt, endedAt)
	//brks, _ := redisLogger.GetTsWithRange("break", startedAt, endedAt)
	//gears, _ := redisLogger.GetTsWithRange("gear", startedAt, endedAt)
	//rpms, _ := redisLogger.GetTsWithRange("rpm", startedAt, endedAt)
	//drss, _ := redisLogger.GetTsWithRange("drs", startedAt, endedAt)
	//speeds, _ := redisLogger.GetTsWithRange("speed", startedAt, endedAt)

	worldPositionXs, _ := redisLogger.GetTsWithRange("worldPositionX", startedAt, endedAt)
	for _, d := range worldPositionXs {
		if v, ok := dataByTimestamp[d.Timestamp]; !ok {
			dataByTimestamp[d.Timestamp] = TelemetryData{
				WorldX: d.Value,
			}
		} else {
			v.WorldX = d.Value
			dataByTimestamp[d.Timestamp] = v
		}
	}

	worldPositionYs, _ := redisLogger.GetTsWithRange("worldPositionY", startedAt, endedAt)
	for _, d := range worldPositionYs {
		if v, ok := dataByTimestamp[d.Timestamp]; !ok {
			dataByTimestamp[d.Timestamp] = TelemetryData{
				WorldY: d.Value,
			}
		} else {
			v.WorldY = d.Value
			dataByTimestamp[d.Timestamp] = v
		}
	}

	worldPositionZs, _ := redisLogger.GetTsWithRange("worldPositionZ", startedAt, endedAt)
	for _, d := range worldPositionZs {
		if v, ok := dataByTimestamp[d.Timestamp]; !ok {
			dataByTimestamp[d.Timestamp] = TelemetryData{
				WorldZ: d.Value,
			}
		} else {
			v.WorldZ = d.Value
			dataByTimestamp[d.Timestamp] = v
		}
	}

	lapDeltaTimes, _ := redisLogger.GetTsWithRange("lapDeltaTime", startedAt, endedAt)
	for _, d := range lapDeltaTimes {
		if v, ok := dataByTimestamp[d.Timestamp]; !ok {
			dataByTimestamp[d.Timestamp] = TelemetryData{
				LapDeltaTime: d.Value,
			}
		} else {
			v.LapDeltaTime = d.Value
			dataByTimestamp[d.Timestamp] = v
		}
	}

	data = make([]TelemetryData, 0)
	for timestamp, d := range dataByTimestamp {
		d.SessionTime = float64(timestamp)
		data = append(data, d)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].SessionTime < data[j].SessionTime
	})

	return
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
		lapQuery := c.QueryParam("lap")
		lapNumber, err := strconv.Atoi(lapQuery)
		if err != nil && lapQuery != "" {
			log.Println("ERROR", lapQuery)
			return nil
		}

		if lapQuery != "" {
			data, err := readSpecificLap(redisLogger, lapNumber)
			if err != nil {
				panic(err)
			}
			return c.JSON(http.StatusOK, data)
		}

		return c.JSON(http.StatusOK, data)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
