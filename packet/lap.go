package packet

type Lap struct {
	LatestLapTime     float32 `json:"m_lastLapTime"`       // Last lap time in seconds
	CurrentLapTime    float32 `json:"m_currentLapTime"`    // Current time around the lap in seconds
	BestLapTime       float32 `json:"m_bestLapTime"`       // Best lap time of the session in seconds
	Sector1Time       float32 `json:"m_sector1Time"`       // Sector 1 time in seconds
	Sector2Time       float32 `json:"m_sector2Time"`       // Sector 2 time in seconds
	LapDistance       float32 `json:"m_lapDistance"`       // Distance vehicle is around current lap in metres – could be negative if line hasn’t been crossed yet
	TotalDistance     float32 `json:"m_totalDistance"`     // Total distance travelled in session in metres – could
	SafetyCarDelta    float32 `json:"m_safetyCarDelta"`    // Delta in seconds for safety car
	CarPosition       uint8   `json:"m_carPosition"`       // Car race position
	CurrentLapNumber  uint8   `json:"m_currentLapNum"`     // Current lap number
	PitStatus         uint8   `json:"m_pitStatus"`         // 0 = none, 1 = pitting, 2 = in pit area
	Sector            uint8   `json:"m_sector"`            // 0 = sector1, 1 = sector2, 2 = sector3
	CurrentLapInvalid uint8   `json:"m_currentLapInvalid"` // Current lap invalid - 0 = valid, 1 = invalid
	Penalties         uint8   `json:"m_penalties"`         // Accumulated time penalties in seconds to be added
	GridPosition      uint8   `json:"m_gridPosition"`      // Grid position the vehicle started the race in
	DriverStatus      uint8   `json:"m_driverStatus"`      // Status of driver - 0 = in garage, 1 = flying lap, 2 = in lap, 3 = out lap, 4 = on track
	ResultStatus      uint8   `json:"m_resultStatus"`      // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = disqualified, 5 = not classified, 6 = retired
}
type LapData struct {
}
