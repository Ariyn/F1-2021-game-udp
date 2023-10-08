package main

import (
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"io/ioutil"
	"path"
	"reflect"
	"strconv"
	"time"
)

type CatTelemetryType int

const (
	StatSteering CatTelemetryType = iota
	StatThrottle
	StatBreak
	StatGear
	StatSpeed
	StatEngineRPM
)

type Ct struct {
	LapDuration time.Duration
	Steer       float32
	Throttle    float32
	Break       float32
	Gear        int
	EngineRpm   int
	Speed       int
}

func loadCarTelemetryData(p string, lap, racingNumber int) (cts []Ct, size int, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), fmt.Sprintf("%d-%d", racingNumber, packet.CarTelemetryDataId)))
	if err != nil {
		return
	}

	//start = time.Unix(0, int64(_t.Timestamp))
	size = packet.Sizeof(reflect.ValueOf(f1.SimplifiedTelemetry{}))

	for i := 0; i < len(b); i += size {
		var t f1.SimplifiedTelemetry
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		//timestamp := time.Unix(0, int64(t.Timestamp))
		//lapTime := timestamp.Sub(start)
		cts = append(cts, Ct{
			LapDuration: time.Duration(t.Timestamp),
			Steer:       t.Steer,
			Throttle:    t.Throttle,
			Break:       t.Break,
			Gear:        int(t.Gear),
			EngineRpm:   int(t.EngineRPM),
			Speed:       int(t.Speed),
		})
	}

	return
}
