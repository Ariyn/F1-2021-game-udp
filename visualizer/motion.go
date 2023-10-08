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

type MotionType int

const (
	MotionWheelSpeedFrontBias MotionType = iota
	MotionWheelSpeedRearBias
	MotionWheelSpeedLeftBias
	MotionWheelSpeedRightBias
	MotionGForceLatitude
	MotionGForceLongitude
	MotionYaw
	MotionWheelSlipFL
	MotionWheelSlipFR
	MotionWheelSlipRL
	MotionWheelSlipRR
)

type Mt struct {
	LapDuration         time.Duration
	GForce              f1.Float3d
	WheelSlip           f1.FloatWheels
	WheelSpeed          f1.FloatWheels
	AngularVelocity     f1.Float3d
	AngularAcceleration f1.Float3d
	Position            f1.Float3d
	Heading             f1.Float3d
}

func loadMotionData(p string, lap, racingNumber int) (mts []Mt, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), fmt.Sprintf("%d-%d", racingNumber, packet.MotionDataId)))
	if err != nil {
		return
	}

	var _t f1.PlayerMotionData
	err = packet.ParsePacket(b, &_t)
	if err != nil {
		return
	}

	size := packet.Sizeof(reflect.ValueOf(f1.PlayerMotionData{}))

	for i := 0; i+size < len(b); i += size {
		var t f1.PlayerMotionData
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		timestamp := time.Duration(t.Timestamp)
		mts = append(mts, Mt{
			LapDuration:         timestamp,
			GForce:              t.GForce,
			WheelSlip:           t.WheelSlip,
			WheelSpeed:          t.WheelSpeed,
			AngularAcceleration: t.AngularAcceleration,
			AngularVelocity:     t.AngularVelocity,
			Position:            t.Position,
			Heading:             t.Heading,
		})
	}

	return
}
