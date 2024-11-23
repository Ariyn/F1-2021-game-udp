package packet

const CarDamageSize = 882

// Car Damage Packet
type CarDamage struct {
	TyresWear            [4]uint8 `json:"m_tyresWear"`            // Tyre wear (percentage)
	TyresDamage          [4]uint8 `json:"m_tyresDamage"`          // Tyre damage (percentage)
	BrakesDamage         [4]uint8 `json:"m_brakesDamage"`         // Brakes damage (percentage)
	FrontLeftWingDamage  uint8    `json:"m_frontLeftWingDamage"`  // Front left wing damage (percentage)
	FrontRightWingDamage uint8    `json:"m_frontRightWingDamage"` // Front right wing damage (percentage)
	RearWingDamage       uint8    `json:"m_rearWingDamage"`       // Rear wing damage (percentage)
	FloorDamage          uint8    `json:"m_floorDamage"`          // Floor damage (percentage)
	DiffuserDamage       uint8    `json:"m_diffuserDamage"`       // Diffuser damage (percentage)
	SidepodDamage        uint8    `json:"m_sidepodDamage"`        // Sidepod damage (percentage)
	DrsFault             uint8    `json:"m_drsFault"`             // Indicator for DRS fault, 0 = OK, 1 = fault
	GearBoxDamage        uint8    `json:"m_gearBoxDamage"`        // Gear box damage (percentage)
	EngineDamage         uint8    `json:"m_engineDamage"`         // Engine damage (percentage)
	EngineMGUHWear       uint8    `json:"m_engineMGUHWear"`       // Engine wear MGU-H (percentage)
	EngineESWear         uint8    `json:"m_engineESWear"`         // Engine wear ES (percentage)
	EngineCEWear         uint8    `json:"m_engineCEWear"`         // Engine wear CE (percentage)
	EngineICEWear        uint8    `json:"m_engineICEWear"`        // Engine wear ICE (percentage)
	EngineMGUKWear       uint8    `json:"m_engineMGUKWear"`       // Engine wear MGU-K (percentage)
	EngineTCWear         uint8    `json:"m_engineTCWear"`         // Engine wear TC (percentage)
}

var _ Data = (*CarDamageData)(nil)

type CarDamageData struct {
	Header     Header
	DamageData [22]CarDamage
}

func (c CarDamageData) GetHeader() Header {
	return c.Header
}

func (c CarDamageData) Id() Id {
	return CarDamageId
}
