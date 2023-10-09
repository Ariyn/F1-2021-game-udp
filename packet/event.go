package packet

const EventDataSize = 36

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

var _ PacketData = (*EventData)(nil)

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
