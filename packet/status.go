package packet

const CarStatusSize = 1058

// Car Status Packet
type CarStatus struct {
	TractionControl         uint8   `json:"m_tractionControl"`         // 0 (off), 1 (medium), 2 (high)
	AntiLockBrakes          uint8   `json:"m_antiLockBrakes"`          // 0 (off), 1 (on)
	FuelMix                 uint8   `json:"m_fuelMix"`                 // Fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
	FrontBrakeBias          uint8   `json:"m_frontBrakeBias"`          // Front brake bias (percentage)
	PitLimiterStatus        uint8   `json:"m_pitLimiterStatus"`        // Pit limiter status - 0 = off, 1 = on
	FuelInTank              float32 `json:"m_fuelInTank"`              // Current fuel mass
	FuelCapacity            float32 `json:"m_fuelCapacity"`            // Fuel capacity
	FuelRemainingLaps       float32 `json:"m_fuelRemainingLaps"`       // Fuel remaining in terms of laps (value on MFD)
	MaxRPM                  uint16  `json:"m_maxRPM"`                  // Cars max RPM, point of rev limiter
	IdleRPM                 uint16  `json:"m_idleRPM"`                 // Cars idle RPM
	MaxGears                uint8   `json:"m_maxGears"`                // Maximum number of gears
	DRSAllowed              uint8   `json:"m_drsAllowed"`              // 0 = not allowed, 1 = allowed
	DRSActivationDistance   uint16  `json:"m_drsActivationDistance"`   // 0 = DRS not available, non-zero - DRS will be available in [X] metres
	ActualTyreCompound      uint8   `json:"m_actualTyreCompound"`      // F1 Modern - 16 = C5, 17 = C4, 18 = C3, 19 = C2, 20 = C1, 7 = inter, 8 = wet
	VisualTyreCompound      uint8   `json:"m_visualTyreCompound"`      // F1 visual (can be different from actual compound) - 16 = soft, 17 = medium, 18 = hard, 7 = inter, 8 = wet
	TyresAgeLaps            uint8   `json:"m_tyresAgeLaps"`            // Age in laps of the current set of tyres
	VehicleFiaFlags         int8    `json:"m_vehicleFiaFlags"`         // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
	ERSStoreEnergy          float32 `json:"m_ersStoreEnergy"`          // ERS energy store in Joules
	ERSDeployMode           uint8   `json:"m_ersDeployMode"`           // ERS deployment mode, 0 = none, 1 = medium, 2 = hotlap, 3 = overtake
	ERSHarvestedThisLapMGUK float32 `json:"m_ersHarvestedThisLapMGUK"` // ERS energy harvested this lap by MGU-K
	ERSHarvestedThisLapMGUH float32 `json:"m_ersHarvestedThisLapMGUH"` // ERS energy harvested this lap by MGU-H
	ERSDeployedThisLap      float32 `json:"m_ersDeployedThisLap"`      // ERS energy deployed this lap
	NetworkPaused           uint8   `json:"m_networkPaused"`           // Whether the car is paused in a network game
}

var _ Data = (*CarStatusData)(nil)

type CarStatusData struct {
	Header      Header
	CarStatuses [22]CarStatus
}

func (c CarStatusData) GetHeader() Header {
	return c.Header
}

func (c CarStatusData) Id() Id {
	return CarStatusId
}
