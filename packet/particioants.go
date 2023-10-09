package packet

const ParticipantDataSize = 1257

type Participant struct {
	IsAiControlled        uint8     `json:"m_aiControlled"`  // Whether the vehicle is AI (1) or Human (0) controlled
	DriverId              uint8     `json:"m_driverId"`      // Driver id - see appendix, 255 if network human
	NetworkId             uint8     `json:"m_networkId"`     // Network id – unique identifier for network players
	TeamId                uint8     `json:"m_teamId"`        // Team id - see appendix
	IsMyTeam              uint8     `json:"m_myTeam"`        // My team flag – 1 = My Team, 0 = otherwise
	RaceNumber            uint8     `json:"m_raceNumber"`    // Race number of the car
	Nationality           uint8     `json:"m_nationality"`   // Nationality of the driver
	Name                  [48]uint8 `json:"m_name"`          // Name of participant in UTF-8 format – null terminated Will be truncated with … (U+2026) if too long
	IsTelemetryRestricted uint8     `json:"m_yourTelemetry"` // The player's UDP setting, 0 = restricted, 1 = public
}

func (p Participant) GetName() string {
	b := make([]byte, 48)
	for i, c := range p.Name {
		b[i] = c
	}

	return string(b)
}

var _ PacketData = (*ParticipantData)(nil)

type ParticipantData struct {
	Header             Header
	NumberOfActiveCars uint8
	Participants       [22]Participant
}

func (p ParticipantData) GetHeader() Header {
	return p.Header
}

func (p ParticipantData) Id() Id {
	return ParticipantsId
}
