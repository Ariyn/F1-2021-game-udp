package packet

const LapDataSize = 970

type DriverLap struct {
	LatestLapTime                 uint32  `json:"m_lastLapTime"`                 // Last lap time in seconds
	CurrentLapTime                uint32  `json:"m_currentLapTime"`              // Current time around the lap in seconds
	Sector1Time                   uint16  `json:"m_sector1Time"`                 // Sector 1 time in seconds
	Sector2Time                   uint16  `json:"m_sector2Time"`                 // Sector 2 time in seconds
	LapDistance                   float32 `json:"m_lapDistance"`                 // Distance vehicle is around current lap in metres – could be negative if line hasn’t been crossed yet
	TotalDistance                 float32 `json:"m_totalDistance"`               // Total distance travelled in session in metres – could
	SafetyCarDelta                float32 `json:"m_safetyCarDelta"`              // Delta in seconds for safety car
	CarPosition                   uint8   `json:"m_carPosition"`                 // Car race position
	CurrentLapNumber              uint8   `json:"m_currentLapNum"`               // Current lap number
	PitStatus                     uint8   `json:"m_pitStatus"`                   // 0 = none, 1 = pitting, 2 = in pit area
	NumberPitStops                uint8   `json:"m_numPitStops"`                 // Number of pit stops taken in this race
	Sector                        uint8   `json:"m_sector"`                      // 0 = sector1, 1 = sector2, 2 = sector3
	CurrentLapInvalid             uint8   `json:"m_currentLapInvalid"`           // Current lap invalid - 0 = valid, 1 = invalid
	Penalties                     uint8   `json:"m_penalties"`                   // Accumulated time penalties in seconds to be added
	Warnings                      uint8   `json:"m_warnings"`                    // Accumulated number of warnings issued
	UnServedDriveThroughPenalties uint8   `json:"m_numUnservedDriveThroughPens"` // Num drive through pens left to serve
	UnServedStopAndGoPenalties    uint8   `json:"m_numUnservedStopGoPens"`       // Num stop go pens left to serve
	GridPosition                  uint8   `json:"m_gridPosition"`                // Grid position the vehicle started the race in
	DriverStatus                  uint8   `json:"m_driverStatus"`                // Status of driver - 0 = in garage, 1 = flying lap, 2 = in lap, 3 = out lap, 4 = on track
	ResultStatus                  uint8   `json:"m_resultStatus"`                // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = disqualified, 5 = not classified, 6 = retired
	PitLaneTimerActive            uint8   `json:"m_pitLaneTimerActive"`          // Pit lane timing, 0 = inactive, 1 = active
	PitLaneTimeInLane             uint16  `json:"m_pitLaneTimeInLaneInMS"`       // If active, the current time spent in the pit lane in ms
	PitStopTimer                  uint16  `json:"m_pitStopTimerInMS"`            // Time of the actual pit stop in ms
	PitStopShouldServePenalty     uint8   `json:"m_pitStopShouldServePen"`       // Whether the car should serve a penalty at this stop
}

type LapData struct {
	Header     Header
	DriverLaps [22]DriverLap
}

func (l LapData) Player() DriverLap {
	return l.DriverLaps[l.Header.PlayerCarIndex]
}
