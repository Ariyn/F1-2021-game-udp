package packet

const MotionDataSize = 1464

type CarMotionData struct {
	WorldPositionX     float32 `json:"m_worldPositionX"`     // World space X position
	WorldPositionY     float32 `json:"m_worldPositionY"`     // World space Y position
	WorldPositionZ     float32 `json:"m_worldPositionZ"`     // World space Z position
	WorldVelocityX     float32 `json:"m_worldVelocityX"`     // Velocity in world space X
	WorldVelocityY     float32 `json:"m_worldVelocityY"`     // Velocity in world space Y
	WorldVelocityZ     float32 `json:"m_worldVelocityZ"`     // Velocity in world space Z
	WorldForwardDirX   uint16  `json:"m_worldForwardDirX"`   // World space forward X direction (normalised)
	WorldForwardDirY   uint16  `json:"m_worldForwardDirY"`   // World space forward Y direction (normalised)
	WorldForwardDirZ   uint16  `json:"m_worldForwardDirZ"`   // World space forward Z direction (normalised)
	WorldRightDirX     uint16  `json:"m_worldRightDirX"`     // World space right X direction (normalised)
	WorldRightDirY     uint16  `json:"m_worldRightDirY"`     // World space right Y direction (normalised)
	WorldRightDirZ     uint16  `json:"m_worldRightDirZ"`     // World space right Z direction (normalised)
	GForceLateral      float32 `json:"m_gForceLateral"`      // Lateral G-Force component
	GForceLongitudinal float32 `json:"m_gForceLongitudinal"` // Longitudinal G-Force component
	GForceVertical     float32 `json:"m_gForceVertical"`     // Vertical G-Force component
	Yaw                float32 `json:"m_yaw"`                // Yaw angle in radians
	Pitch              float32 `json:"m_pitch"`              // Pitch angle in radians
	Roll               float32 `json:"m_roll"`               // Roll angle in radians
}

type MotionData struct {
	Header        Header
	CarMotionData [22]CarMotionData

	SuspensionPosition     [4]float32 `json:"m_suspensionPosition"`     // Note: All wheel arrays have the following order:
	SuspensionVelocity     [4]float32 `json:"m_suspensionVelocity"`     // RL, RR, FL, FR
	SuspensionAcceleration [4]float32 `json:"m_suspensionAcceleration"` // RL, RR, FL, FR
	WheelSpeed             [4]float32 `json:"m_wheelSpeed"`             // Speed of each wheel
	WheelSlip              [4]float32 `json:"m_wheelSlip"`              // Slip ratio for each wheel
	LocalVelocityX         float32    `json:"m_localVelocityX"`         // Velocity in local space
	LocalVelocityY         float32    `json:"m_localVelocityY"`         // Velocity in local space
	LocalVelocityZ         float32    `json:"m_localVelocityZ"`         // Velocity in local space
	AngularVelocityX       float32    `json:"m_angularVelocityX"`       // Angular velocity x-component
	AngularVelocityY       float32    `json:"m_angularVelocityY"`       // Angular velocity y-component
	AngularVelocityZ       float32    `json:"m_angularVelocityZ"`       // Angular velocity z-component
	AngularAccelerationX   float32    `json:"m_angularAccelerationX"`   // Angular velocity x-component
	AngularAccelerationY   float32    `json:"m_angularAccelerationY"`   // Angular velocity y-component
	AngularAccelerationZ   float32    `json:"m_angularAccelerationZ"`   // Angular velocity z-component
	FrontWheelsAngle       float32    `json:"m_frontWheelsAngle"`       // Current front wheels angle in radians
}

func (m MotionData) Player() CarMotionData {
	return m.CarMotionData[m.Header.PlayerCarIndex]
}
