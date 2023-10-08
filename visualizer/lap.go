package main

import (
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"io/ioutil"
	"log"
	"path"
	"reflect"
	"strconv"
	"time"
)

type Lt struct {
	Timestamp        time.Duration
	CurrentLapTime   time.Duration
	Sector1Time      time.Duration
	Sector2Time      time.Duration
	LapDistance      float32
	CurrentLapNumber int
	Sector           int
	DriverStatus     f1.DriverStatus
}

func loadLapTelemetries(p string, lap, racingNumber int) (lts []Lt, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), fmt.Sprintf("%d-%d", racingNumber, packet.LapDataId)))
	if err != nil {
		return
	}

	//var start time.Time
	var _t f1.SimplifiedLap
	err = packet.ParsePacket(b, &_t)
	if err != nil {
		return
	}

	//start = time.Unix(0, int64(_t.Timestamp))
	size := packet.Sizeof(reflect.ValueOf(f1.SimplifiedLap{}))

	for i := 0; i < len(b); i += size {
		var t f1.SimplifiedLap
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		lapNumber := int(t.CurrentLapNumber)
		if lapNumber != lap {
			log.Printf("previous lap %#v", t)
			//continue
		}
		if t.LapDistance < 0 {
			continue
		}

		lts = append(lts, Lt{
			Timestamp:        time.Duration(t.Timestamp),
			CurrentLapTime:   time.Duration(t.CurrentLapTime) * time.Millisecond,
			Sector1Time:      time.Duration(t.Sector1Time) * time.Millisecond,
			Sector2Time:      time.Duration(t.Sector2Time) * time.Millisecond,
			LapDistance:      t.LapDistance,
			CurrentLapNumber: int(t.CurrentLapNumber),
			Sector:           int(t.Sector),
			DriverStatus:     f1.DriverStatus(t.DriverStatus),
		})
	}

	return
}
