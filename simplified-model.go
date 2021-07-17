package F1_2021_game_udp

import (
	"github.com/ariyn/F1-2021-game-udp/packet"
)

type SimplifiedTelemetry struct {
	Timestamp               uint64
	Steer                   float32
	Throttle                float32
	Break                   float32
	Gear                    int8
	EngineRPM               uint16
	Speed                   uint16
	DRS                     uint8
	BreaksTemperature       Int16Wheel
	TyresSurfaceTemperature Int8Wheel
	TyresInnerTemperature   Int8Wheel
	EngineTemperature       uint16
	TyresPressure           FloatWheels
	SurfaceType             Int8Wheel
}

func SimplifyTelemetry(timestamp int64, telemetry packet.CarTelemetry) SimplifiedTelemetry {
	return SimplifiedTelemetry{
		Timestamp: uint64(timestamp),
		Steer:     telemetry.Steer,
		Throttle:  telemetry.Throttle,
		Break:     telemetry.Break,
		Gear:      telemetry.Gear,
		EngineRPM: telemetry.EngineRPM,
		Speed:     telemetry.Speed,
		DRS:       telemetry.DRS,
		BreaksTemperature: Int16Wheel{
			RL: int16(telemetry.BreaksTemperature[0]),
			RR: int16(telemetry.BreaksTemperature[1]),
			FL: int16(telemetry.BreaksTemperature[2]),
			FR: int16(telemetry.BreaksTemperature[3]),
		},
		TyresSurfaceTemperature: Int8Wheel{
			RL: int8(telemetry.TyresSurfaceTemperature[0]),
			RR: int8(telemetry.TyresSurfaceTemperature[1]),
			FL: int8(telemetry.TyresSurfaceTemperature[2]),
			FR: int8(telemetry.TyresSurfaceTemperature[3]),
		},
		TyresInnerTemperature: Int8Wheel{
			RL: int8(telemetry.TyresInnerTemperature[0]),
			RR: int8(telemetry.TyresInnerTemperature[1]),
			FL: int8(telemetry.TyresInnerTemperature[2]),
			FR: int8(telemetry.TyresInnerTemperature[3]),
		},
		EngineTemperature: 0,
		TyresPressure: FloatWheels{
			RL: telemetry.TyresPressure[0],
			RR: telemetry.TyresPressure[1],
			FL: telemetry.TyresPressure[2],
			FR: telemetry.TyresPressure[3],
		},
		SurfaceType: Int8Wheel{
			RL: int8(telemetry.SurfaceType[0]),
			RR: int8(telemetry.SurfaceType[1]),
			FL: int8(telemetry.SurfaceType[2]),
			FR: int8(telemetry.SurfaceType[3]),
		},
	}
}

type SimplifiedLap struct {
	Timestamp        uint64
	LatestLapTime    uint32
	CurrentLapTime   uint32
	Sector1Time      uint16
	Sector2Time      uint16
	LapDistance      float32
	CurrentLapNumber uint8
	Sector           uint8
	DriverStatus     uint8
}

func SimplifyLap(timestamp int64, lap packet.DriverLap) SimplifiedLap {
	return SimplifiedLap{
		Timestamp:        uint64(timestamp),
		LatestLapTime:    lap.LatestLapTime,
		CurrentLapTime:   lap.CurrentLapTime,
		Sector1Time:      lap.Sector1Time,
		Sector2Time:      lap.Sector2Time,
		LapDistance:      lap.LapDistance,
		CurrentLapNumber: lap.CurrentLapNumber,
		Sector:           lap.Sector,
		DriverStatus:     lap.DriverStatus,
	}
}

//type SimplifiedCarStatus struct {
//	TractionControl uint8
//}

type DriverStatus int

const (
	DriverStatusInGarage DriverStatus = iota
	DriverStatusFlyingLap
	DriverStatusInLap
	DriverStatusOutLap
	DriverStatusOnTrack
)
