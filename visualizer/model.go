package visualizer

import f1 "github.com/ariyn/F1-2021-game-udp"

type SimplifiedTelemetry struct {
	TimeStamp               uint64
	Steer                   float32
	Throttle                float32
	Break                   float32
	Gear                    int8
	EngineRPM               uint16
	Speed                   uint16
	DRS                     uint8
	BreaksTemperature       f1.Int16Wheel
	TyresSurfaceTemperature f1.Int8Wheel
	TyresInnerTemperature   f1.Int8Wheel
	EngineTemperature       uint16
	TyresPressure           f1.FloatWheels
	SurfaceType             f1.Int8Wheel
}

type SimplifiedLap struct {
	LatestLapTime    uint32
	CurrentLapTime   uint32
	Sector1Time      uint16
	Sector2Time      uint16
	LapDistance      float32
	CurrentLapNumber uint8
	Sector           uint8
	DriverStatus     uint8
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
