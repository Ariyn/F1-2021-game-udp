package packet

const EventDataSize = 36

type EventType interface {
	StringCode() string
}

type FastestLap struct {
	VehicleIndex uint8   `json:"vehicleIdx"` // Vehicle index of car achieving fastest lap
	LapTime      float32 `json:"lapTime"`    // SetLap time is in seconds
}

func (FastestLap) StringCode() string {
	return "FTLP"
}

type Retirement struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of car retiring
}

func (Retirement) StringCode() string {
	return "RTMT"
}

type TeamMateInPits struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of team mate
}

type RaceWinner struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of the race winner
}

type Penalty struct {
	PenaltyType       uint8
	InfringementType  uint8
	VehicleIndex      uint8
	OtherVehicleIndex uint8 // if not 255
	Time              uint8 // if not 255
	LapNumber         uint8
	PlacesGained      uint8 // if not 255
}

type SpeedTrap struct {
	VehicleIndex            uint8
	Speed                   float32
	OverallFastestInSession uint8
	DriverFastestInSession  uint8
}

type StartLights struct {
	NumberLights uint8
}

type DriveThroughPenaltyServed struct {
	VehicleIndex uint8
}

type StopGoPenaltyServed struct {
	VehicleIndex uint8
}

type Flashback struct {
	FlashbackFrameIdentifier uint32
	FlashbackSessionTime     float32
}

type Buttons struct {
	ButtonStatus uint32
}

type SessionStarted struct{}

func (SessionStarted) StringCode() string {
	return "SSTA"
}

type SessionEnded struct{}

func (SessionEnded) StringCode() string {
	return "SEND"
}

type EventHeaderData struct {
	Header          Header
	EventStringCode [4]uint8
}

func (d EventHeaderData) StringCode() string {
	return string([]byte{d.EventStringCode[0], d.EventStringCode[1], d.EventStringCode[2], d.EventStringCode[3]})
}
