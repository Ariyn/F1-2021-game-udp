package packet

const FinalClassificationSize = 839

// Final Classification Packet
type FinalClassification struct {
	Position         uint8    `json:"m_position"`         // Finishing position
	NumLaps          uint8    `json:"m_numLaps"`          // Number of laps completed
	GridPosition     uint8    `json:"m_gridPosition"`     // Grid position of the car
	Points           uint8    `json:"m_points"`           // Number of points scored
	NumPitStops      uint8    `json:"m_numPitStops"`      // Number of pit stops made
	ResultStatus     uint8    `json:"m_resultStatus"`     // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = didnotfinish, 5 = disqualified, 6 = not classified, 7 = retired
	BestLapTimeInMS  uint32   `json:"m_bestLapTimeInMS"`  // Best lap time of the session in milliseconds
	TotalRaceTime    float64  `json:"m_totalRaceTime"`    // Total race time in seconds without penalties
	PenaltiesTime    uint8    `json:"m_penaltiesTime"`    // Total penalties accumulated in seconds
	NumPenalties     uint8    `json:"m_numPenalties"`     // Number of penalties applied to this driver
	NumTyreStints    uint8    `json:"m_numTyreStints"`    // Number of tyres stints up to maximum
	TyreStintsActual [8]uint8 `json:"m_tyreStintsActual"` // Actual tyres used by this driver
	TyreStintsVisual [8]uint8 `json:"m_tyreStintsVisual"` // Visual tyres used by this driver
}

var _ Data = (*FinalClassificationData)(nil)

type FinalClassificationData struct {
	Header          Header
	Classifications [22]FinalClassification
}

func (f FinalClassificationData) GetHeader() Header {
	return f.Header
}

func (f FinalClassificationData) Id() Id {
	return FinalClassificationId
}
