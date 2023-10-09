package main

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"time"
)

func main() {
	redisLogger := logger.NewRedisClient(context.Background(), "localhost:6379", 0)
	defer redisLogger.CloseListen()

	fmt.Println("start Listening\n\nwaiting inputs")
	for range redisLogger.Listen("LapData") {
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

		now := time.UnixMilli(lapDistance.Timestamp)
		fmt.Print("\u001B[1A\u001B[K\u001B[1A\u001B[K\u001B[1A\u001B[K")

		text := fmt.Sprintf("[%s]-FR[%.0f] | LAP: %.0f - %.2fKm\n", now.Format("2006-01-02 15:04:05"), frameIdentifier.Value, lapNumber.Value, lapDistance.Value)
		text += fmt.Sprintf("TH: %.2f, BR: %.2f, GR: %.0f, ENG: %.0fRPM, DRS:%b\n", throttle.Value, brk.Value, gear.Value, rpm.Value, drs.Value)
		text += fmt.Sprintf("%.2f %.2f %.2f\n", worldPositionX.Value, worldPositionY.Value, worldPositionZ.Value)

		fmt.Print(text)
	}
}
