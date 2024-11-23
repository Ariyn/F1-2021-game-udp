package packet

const CarTelemetryDataSize = 1347

// Car Telemetry Packet
type CarTelemetry struct {
	Speed                   uint16     // Speed of car in kilometres per hour
	Throttle                float32    // Amount of throttle applied (0.0 to 1.0)
	Steer                   float32    // Steering (-1.0 (full lock left) to 1.0 (full lock right))
	Break                   float32    // Amount of brake applied (0.0 to 1.0)
	Clutch                  uint8      // Amount of clutch applied (0 to 100)
	Gear                    int8       // Gear selected (1-8, N=0, R=-1)
	EngineRPM               uint16     // Engine RPM
	DRS                     uint8      // 0 = off, 1 = on
	RevLightsPercent        uint8      // Rev lights indicator (percentage)
	RevLightsBitValue       uint16     // Rev lights (bit 0 = leftmost LED, bit 14 = rightmost LED)
	BreaksTemperature       [4]uint16  // Brakes temperature (celsius)
	TyresSurfaceTemperature [4]uint8   // Tyres surface temperature (celsius)
	TyresInnerTemperature   [4]uint8   // Tyres inner temperature (celsius)
	EngineTemperature       uint16     // Engine temperature (celsius)
	TyresPressure           [4]float32 // Tyres pressure (PSI)
	SurfaceType             [4]uint8   // Driving surface, see appendices
}

var _ Data = (*CarTelemetryData)(nil)

type CarTelemetryData struct {
	Header                       Header
	CarTelemetries               [22]CarTelemetry
	MFDPanelIndex                uint8 `json:"m_mfdPanelIndex"`                // Index of MFD panel open - 255 = MFD closed, Single player, race – 0 = Car setup, 1 = Pits, 2 = Damage, 3 =  Engine, 4 = Temperatures, May vary depending on game mode
	MFDPanelIndexSecondaryPlayer uint8 `json:"m_mfdPanelIndexSecondaryPlayer"` // See above
	SuggestGear                  int8  `json:"m_suggestedGear"`                // Suggested gear for the player (1-8) 0 if no gear suggested
}

func (c CarTelemetryData) GetHeader() Header {
	return c.Header
}

func (c CarTelemetryData) Id() Id {
	return CarTelemetryDataId
}

func (c CarTelemetryData) Player() CarTelemetry {
	return c.CarTelemetries[c.Header.PlayerCarIndex]
}
