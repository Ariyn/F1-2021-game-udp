package visualizer

type SimplifiedTelemetry struct {
	TimeStamp    uint64
	Steer        float32
	Throttle     float32
	Break        float32
	Gear         uint8
	EngineRPM    uint16
	Speed        uint16
	Sector       uint8
	DriverStatus uint8
	LapDistance  float32
}

type DriverStatus int

const (
	DriverStatusInGarage DriverStatus = iota
	DriverStatusFlyingLap
	DriverStatusInLap
	DriverStatusOutLap
	DriverStatusOnTrack
)
