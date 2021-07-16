package packet

const SessionDataSize = 625

type MarshalZone struct {
	ZoneStart float32 `json:"m_zoneStart"` // Fraction (0..1) of way through the lap the marshal zone starts
	ZoneFlag  int8    `json:"m_zoneFlag"`  // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
}

type WeatherForecastSample struct {
	SessionType            uint8 `json:"m_sessionType"`            // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = Time Trial
	TimeOffset             uint8 `json:"m_timeOffset"`             // Time in minutes the forecast is for
	Weather                uint8 `json:"m_weather"`                // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain, 4 = heavy rain, 5 = storm
	TrackTemperature       uint8 `json:"m_trackTemperature"`       // Track temp. in degrees Celsius
	TrackTemperatureChange uint8 `json:"m_trackTemperatureChange"` // Track temp. change – 0 = up, 1 = down, 2 = no change
	AirTemperature         uint8 `json:"m_airTemperature"`         // Air temp. in degrees celsius
	AirTemperatureChange   uint8 `json:"m_airTemperatureChange"`   // Air temp. change – 0 = up, 1 = down, 2 = no change
	RainPercentage         uint8 `json:"m_rainPercentage"`         // Rain percentage (0-100)

}

type SessionData struct {
	Header                       Header
	Weather                      uint8                     `json:"m_weather"`                   // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain, 4 = heavy rain, 5 = storm
	TrackTemperature             uint8                     `json:"m_trackTemperature"`          // Track temp. in degrees celsius
	AirTemperature               uint8                     `json:"m_airTemperature"`            // Air temp. in degrees celsius
	TotalLaps                    uint8                     `json:"m_totalLaps"`                 // Total number of laps in this race
	TrackLength                  uint16                    `json:"m_trackLength"`               // Track length in metres
	SessionType                  uint8                     `json:"m_sessionType"`               // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = R3, 13 = Time Trial
	TrackId                      uint8                     `json:"m_trackId"`                   // -1 for unknown, 0-21 for tracks, see appendix
	FormulaId                    uint8                     `json:"m_formula"`                   // Formula, 0 = F1 Modern, 1 = F1 Classic, 2 = F2, 3 = F1 Generic
	SessionTimeLeft              uint16                    `json:"m_sessionTimeLeft"`           // Time left in session in seconds
	SessionDuration              uint16                    `json:"m_sessionDuration"`           // Session duration in seconds
	PitSpeedLimit                uint8                     `json:"m_pitSpeedLimit"`             // Pit speed limit in kilometres per hour
	GamePaused                   uint8                     `json:"m_gamePaused"`                // Whether the game is paused
	IsSpectating                 uint8                     `json:"m_isSpectating"`              // Whether the player is spectating
	SpectatorCarIndex            uint8                     `json:"m_spectatorCarIndex"`         // Index of the car being spectated
	SliProNativeSupport          uint8                     `json:"m_sliProNativeSupport"`       // SLI Pro support, 0 = inactive, 1 = active
	NumberMarshalZones           uint8                     `json:"m_numMarshalZones"`           // Number of marshal zones to follow
	MarshalZones                 [21]MarshalZone           `json:"m_marshalZones"`              // List of marshal zones – max 21
	SafetyCarStatus              uint8                     `json:"m_safetyCarStatus"`           // 0 = no safety car, 1 = full, 2 = virtual, 3 = formation lap
	IsNetworkGame                uint8                     `json:"m_networkGame"`               // 0 = offline, 1 = online
	NumberWeatherForecastSamples uint8                     `json:"m_numWeatherForecastSamples"` // Number of weather samples to follow
	WeatherForecastSamples       [56]WeatherForecastSample `json:"m_weatherForecastSamples"`    // Array of weather forecast samples
	ForecastAccuracy             uint8                     `json:"m_forecastAccuracy"`          // 0 = Perfect, 1 = Approximate
	AiDifficulty                 uint8                     `json:"m_aiDifficulty"`              // AI Difficulty rating – 0-110
	SeasonLinkIdentifier         uint32                    `json:"m_seasonLinkIdentifier"`      // Identifier for season - persists across saves
	WeekendLinkIdentifier        uint32                    `json:"m_weekendLinkIdentifier"`     // Identifier for weekend - persists across saves
	SessionLinkIdentifier        uint32                    `json:"m_sessionLinkIdentifier"`     // Identifier for session - persists across saves
	PitStopWindowIdealLap        uint8                     `json:"m_pitStopWindowIdealLap"`     // Ideal lap to pit on for current strategy (player)
	PitStopWindowLatestLap       uint8                     `json:"m_pitStopWindowLatestLap"`    // Latest lap to pit on for current strategy (player)
	PitStopRejoinPosition        uint8                     `json:"m_pitStopRejoinPosition"`     // Predicted position to rejoin at (player)
	SteeringAssist               uint8                     `json:"m_steeringAssist"`            // 0 = off, 1 = on
	BreakingAssist               uint8                     `json:"m_brakingAssist"`             // 0 = off, 1 = low, 2 = medium, 3 = high
	GearboxAssist                uint8                     `json:"m_gearboxAssist"`             // 1 = manual, 2 = manual & suggested gear, 3 = auto
	PitAssist                    uint8                     `json:"m_pitAssist"`                 // 0 = off, 1 = on
	PitReleaseAssist             uint8                     `json:"m_pitReleaseAssist"`          // 0 = off, 1 = on
	ERSAssist                    uint8                     `json:"m_ERSAssist"`                 // 0 = off, 1 = on
	DRSAssist                    uint8                     `json:"m_DRSAssist"`                 // 0 = off, 1 = on
	DynamicRacingLine            uint8                     `json:"m_dynamicRacingLine"`         // 0 = off, 1 = corners only, 2 = full
	DynamicRacingLineType        uint8                     `json:"m_dynamicRacingLineType"`     // 0 = 2D, 1 = 3D
}
