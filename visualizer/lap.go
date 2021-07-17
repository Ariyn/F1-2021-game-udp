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
	SectorDurations  []time.Duration
	TotalLapDuration time.Duration
}

func loadLapTelemetries(p string, lap, racingNumber int) (lt Lt, err error) {
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

	sector1 := time.Duration(0)
	sector2 := time.Duration(0)
	var lapTime time.Duration
	for i := 0; i < len(b); i += size {
		var t f1.SimplifiedLap
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		lapNumber := int(t.CurrentLapNumber)
		if lapNumber != lap {
			log.Printf("previous lap %#v", t)
			continue
		}
		if t.LapDistance < 0 {
			continue
		}

		if t.Sector1Time != 0 && sector1.Milliseconds() == 0 {
			sector1 = time.Duration(t.Sector1Time) * time.Millisecond
		}
		if t.Sector2Time != 0 && sector2.Milliseconds() == 0 {
			sector2 = time.Duration(t.Sector2Time) * time.Millisecond
		}

		if t.CurrentLapTime != 0 {
			lapTime = time.Duration(t.CurrentLapTime) * time.Millisecond
		}

		// TODO: driver status로 pitstop start, end 알아내기
		// t.DriverStatus
	}

	lt = Lt{
		SectorDurations: []time.Duration{
			sector1, sector2, lapTime - sector2 - sector1,
		},
		TotalLapDuration: lapTime,
	}
	return
}
