package packet

const SessionHistorySize = 1155

type LapHistory struct {
	LapTimeInMs      uint32 `json:"m_lapTimeInMS"`      // Lap time in milliseconds
	Sector1TimeInMs  uint16 `json:"m_sector1TimeInMS"`  // Sector 1 time in milliseconds
	Sector2TimeInMs  uint16 `json:"m_sector2TimeInMS"`  // Sector 2 time in milliseconds
	Sector3TimeInMs  uint16 `json:"m_sector3TimeInMS"`  // Sector 3 time in milliseconds
	LapValidBitFlags uint8  `json:"m_lapValidBitFlags"` // 0x01 bit set-lap valid, 0x02 bit set-sector 1 valid 0x04 bit set-sector 2 valid, 0x08 bit set-sector 3 valid
}

type TypeStintHistory struct {
	EndLap             uint8 `json:"m_endLap"`             // Lap the tyre usage ends on (255 of current tyre)
	TyreActualCompound uint8 `json:"m_tyreActualCompound"` // Actual tyres used by this driver
	TyreVisualCompound uint8 `json:"m_tyreVisualCompound"` // Visual tyres used by this driver
}

var _ Data = (*SessionHistoryData)(nil)

type SessionHistoryData struct {
	Header        Header `json:"header"`
	CarIndex      uint8  `json:"m_carIdx"`        // Index of the car this lap data relates to
	NumLaps       uint8  `json:"m_numLaps"`       // Num laps in the data (including current partial lap)
	NumTyreStints uint8  `json:"m_numTyreStints"` // Number of tyre stints in the data

	BestLapTimeLapNum uint8 `json:"m_bestLapTimeLapNum"` // Lap the best lap time was achieved on
	BestSector1LapNum uint8 `json:"m_bestSector1LapNum"` // Lap the best Sector 1 time was achieved on
	BestSector2LapNum uint8 `json:"m_bestSector2LapNum"` // Lap the best Sector 2 time was achieved on
	BestSector3LapNum uint8 `json:"m_bestSector3LapNum"` // Lap the best Sector 3 time was achieved on

	LapHistory        [100]LapHistory     `json:"m_lapHistoryData"` // 100 laps of data max
	TyreStintsHistory [8]TypeStintHistory `json:"m_tyreStintsHistoryData"`
}

func (s SessionHistoryData) GetHeader() Header {
	return s.Header
}

func (s SessionHistoryData) Id() Id {
	return SessionHistoryId
}
