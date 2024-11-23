package packet

const EventDataSize = 36

// Event String Codes
type EventCode string

const (
	SSTA EventCode = "SSTA" // Sent when the session starts
	SEND EventCode = "SEND" // Sent when the session ends
	FTLP EventCode = "FTLP" // When a driver achieves the fastest lap
	RTMT EventCode = "RTMT" // When a driver retires
	DRSE EventCode = "DRSE" // Race control have enabled DRS
	DRSD EventCode = "DRSD" // Race control have disabled DRS
	TMPT EventCode = "TMPT" // Your teammate has entered the pits
	CHQF EventCode = "CHQF" // The chequered flag has been waved
	RCWN EventCode = "RCWN" // The race winner is announced
	PENA EventCode = "PENA" // A penalty has been issued – details in event
	SPTP EventCode = "SPTP" // Speed trap has been triggered by fastest speed
	STLG EventCode = "STLG" // Start lights – number shown
	LGTO EventCode = "LGTO" // Lights out
	DTSV EventCode = "DTSV" // Drive through penalty served
	SGSV EventCode = "SGSV" // Stop and go penalty served
	FLBK EventCode = "FLBK" // Flashback activated
	BUTN EventCode = "BUTN" // Button status changed

)

type EventSpecificData interface {
	StringCode() EventCode
}

var _ EventSpecificData = (*SessionStarted)(nil)

type SessionStarted struct{}

func (SessionStarted) StringCode() EventCode {
	return SSTA
}

var _ EventSpecificData = (*SessionEnded)(nil)

type SessionEnded struct{}

func (SessionEnded) StringCode() EventCode {
	return SEND
}

var _ EventSpecificData = (*FastestLap)(nil)

type FastestLap struct {
	VehicleIndex uint8   `json:"vehicleIdx"` // Vehicle index of car achieving fastest lap
	LapTime      float32 `json:"lapTime"`    // SetLap time is in seconds
}

func (FastestLap) StringCode() EventCode {
	return FTLP
}

var _ EventSpecificData = (*Retirement)(nil)

type Retirement struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of car retiring
}

func (Retirement) StringCode() EventCode {
	return RTMT
}

var _ EventSpecificData = (*DRSEnabled)(nil)

type DRSEnabled struct{}

func (DRSEnabled) StringCode() EventCode {
	return DRSE
}

var _ EventSpecificData = (*DRSDisabled)(nil)

type DRSDisabled struct{}

func (DRSDisabled) StringCode() EventCode {
	return DRSD
}

var _ EventSpecificData = (*TeamMateInPits)(nil)

type TeamMateInPits struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of teammate
}

func (TeamMateInPits) StringCode() EventCode {
	return TMPT
}

var _ EventSpecificData = (*ChequeredFlag)(nil)

type ChequeredFlag struct{}

func (ChequeredFlag) StringCode() EventCode {
	return CHQF
}

var _ EventSpecificData = (*RaceWinner)(nil)

type RaceWinner struct {
	VehicleIndex uint8 `json:"vehicleIdx"` // Vehicle index of the race winner
}

func (RaceWinner) StringCode() EventCode {
	return RCWN
}

var _ EventSpecificData = (*Penalty)(nil)

type Penalty struct {
	PenaltyType       uint8 // Penalty type – see Appendices
	InfringementType  uint8 // Infringement type – see Appendices
	VehicleIndex      uint8 // Vehicle index of the car the penalty is applied to
	OtherVehicleIndex uint8 // Vehicle index of the other car involved
	Time              uint8 // Time gained, or time spent doing action in seconds
	LapNumber         uint8 // Lap the penalty occurred on
	PlacesGained      uint8 // Number of places gained by this
}

func (Penalty) StringCode() EventCode {
	return PENA
}

var _ EventSpecificData = (*SpeedTrap)(nil)

type SpeedTrap struct {
	VehicleIndex            uint8   // Vehicle index of the vehicle triggering speed trap
	Speed                   float32 // Top speed achieved in kilometres per hour
	OverallFastestInSession uint8   // Overall fastest speed in session = 1, otherwise 0
	DriverFastestInSession  uint8   // Fastest speed for driver in session = 1, otherwise 0
}

func (SpeedTrap) StringCode() EventCode {
	return SPTP
}

var _ EventSpecificData = (*StartLights)(nil)

type StartLights struct {
	NumberLights uint8 // Number of lights showing
}

func (StartLights) StringCode() EventCode {
	return STLG
}

var _ EventSpecificData = (*LightsOut)(nil)

type LightsOut struct{}

func (LightsOut) StringCode() EventCode {
	return LGTO
}

var _ EventSpecificData = (*DriveThroughPenaltyServed)(nil)

type DriveThroughPenaltyServed struct {
	VehicleIndex uint8 // Vehicle index of the vehicle serving drive through
}

func (DriveThroughPenaltyServed) StringCode() EventCode {
	return DTSV
}

var _ EventSpecificData = (*StopGoPenaltyServed)(nil)

type StopGoPenaltyServed struct {
	VehicleIndex uint8 // Vehicle index of the vehicle serving stop go
}

func (StopGoPenaltyServed) StringCode() EventCode {
	return SGSV
}

var _ EventSpecificData = (*Flashback)(nil)

type Flashback struct {
	FlashbackFrameIdentifier uint32  // Frame identifier flashed back to
	FlashbackSessionTime     float32 // Session time flashed back to
}

func (f Flashback) StringCode() EventCode {
	return FLBK
}

var _ EventSpecificData = (*Buttons)(nil)

type Buttons struct {
	ButtonStatus uint32 // Bit flags specifying which buttons are being pressed currently - see appendices
}

func (b Buttons) StringCode() EventCode {
	return BUTN
}

var _ Data = (*EventData)(nil)

// Event Packet
// This packet gives details of events that happen during the course of a session.
type EventData struct {
	Header          Header
	EventStringCode [4]uint8
	Event           EventSpecificData
}

func (d EventData) Id() Id {
	return EventId
}

func (d EventData) GetHeader() Header {
	return d.Header
}

func (d EventData) StringCode() EventCode {
	return EventCode([]byte{d.EventStringCode[0], d.EventStringCode[1], d.EventStringCode[2], d.EventStringCode[3]})
}
