package packet

const CarTelemetryDataSize = 1347

type CarTelemetry struct {
	Speed                   uint16
	Throttle                float32
	Steer                   float32
	Break                   float32
	Clutch                  uint8
	Gear                    int8
	EngineRPM               uint16
	DRS                     uint8
	RevLightsPercent        uint8
	RevLightsBitValue       uint16
	BreaksTemperature       [4]uint16
	TyresSurfaceTemperature [4]uint8
	TyresInnerTemperature   [4]uint8
	EngineTemperature       uint16
	TyresPressure           [4]float32
	SurfaceType             [4]uint8
}

var _ PacketData = (*CarTelemetryData)(nil)

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
