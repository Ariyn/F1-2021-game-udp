package F1_2021_game_udp

import "github.com/ariyn/F1-2021-game-udp/packet"

type Float3d struct {
	X float32
	Y float32
	Z float32
}

type Int3d struct {
	X int
	Y int
	Z int
}

type FloatWheels struct {
	RL float32
	RR float32
	FL float32
	FR float32
}

type IntWheels struct {
	RL int
	RR int
	FL int
	FR int
}

type MotionData struct {
	Position   Float3d
	Velocity   Float3d
	ForwardDir Int3d
	RightDir   Int3d
	GForce     Float3d
	Heading    Float3d
}

type PlayerMotionData struct {
	Timestamp              uint64
	Position               Float3d
	Velocity               Float3d
	ForwardDir             Int3d
	RightDir               Int3d
	GForce                 Float3d
	Heading                Float3d
	SuspensionPosition     FloatWheels
	SuspensionVelocity     FloatWheels
	SuspensionAcceleration FloatWheels
	WheelSpeed             FloatWheels
	WheelSlip              FloatWheels
	LocalVelocity          Float3d
	AngularVelocity        Float3d
	AngularAcceleration    Float3d
	FrontWheelsAngle       float32 // Current front wheels angle in radians
}

func GetPlayerMotionData(timestamp int64, m *packet.MotionData) PlayerMotionData {
	player := m.Player()
	return PlayerMotionData{
		Timestamp: uint64(timestamp),

		Position: Float3d{
			X: player.WorldPositionX,
			Y: player.WorldPositionY,
			Z: player.WorldPositionZ,
		},
		Velocity: Float3d{
			X: player.WorldVelocityX,
			Y: player.WorldVelocityY,
			Z: player.WorldVelocityZ,
		},
		ForwardDir: Int3d{
			X: int(player.WorldForwardDirX),
			Y: int(player.WorldForwardDirY),
			Z: int(player.WorldForwardDirZ),
		},
		RightDir: Int3d{
			X: int(player.WorldRightDirX),
			Y: int(player.WorldRightDirY),
			Z: int(player.WorldRightDirZ),
		},
		GForce: Float3d{
			X: player.GForceLateral,
			Y: player.GForceLongitudinal,
			Z: player.GForceVertical,
		},
		Heading: Float3d{
			X: player.Yaw,
			Y: player.Roll,
			Z: player.Pitch,
		},
		SuspensionPosition: FloatWheels{
			RL: m.SuspensionPosition[0],
			RR: m.SuspensionPosition[1],
			FL: m.SuspensionPosition[2],
			FR: m.SuspensionPosition[3],
		},
		SuspensionVelocity: FloatWheels{
			RL: m.SuspensionVelocity[0],
			RR: m.SuspensionVelocity[1],
			FL: m.SuspensionVelocity[2],
			FR: m.SuspensionVelocity[3],
		},
		SuspensionAcceleration: FloatWheels{
			RL: m.SuspensionAcceleration[0],
			RR: m.SuspensionAcceleration[1],
			FL: m.SuspensionAcceleration[2],
			FR: m.SuspensionAcceleration[3],
		},
		WheelSpeed: FloatWheels{
			RL: m.WheelSpeed[0],
			RR: m.WheelSpeed[1],
			FL: m.WheelSpeed[2],
			FR: m.WheelSpeed[3],
		},
		WheelSlip: FloatWheels{
			RL: m.WheelSlip[0],
			RR: m.WheelSlip[1],
			FL: m.WheelSlip[2],
			FR: m.WheelSlip[3],
		},
		LocalVelocity: Float3d{
			X: m.LocalVelocityX,
			Y: m.LocalVelocityY,
			Z: m.LocalVelocityZ,
		},
		AngularVelocity: Float3d{
			X: m.AngularVelocityX,
			Y: m.AngularVelocityY,
			Z: m.AngularVelocityZ,
		},
		AngularAcceleration: Float3d{
			X: m.AngularAccelerationX,
			Y: m.AngularAccelerationY,
			Z: m.AngularAccelerationZ,
		},
		FrontWheelsAngle: m.FrontWheelsAngle,
	}
}

func GetMotionData(timestamp int64, m packet.CarMotionData) PlayerMotionData {
	return PlayerMotionData{
		Timestamp: uint64(timestamp),

		Position: Float3d{
			X: m.WorldPositionX,
			Y: m.WorldPositionY,
			Z: m.WorldPositionZ,
		},
		Velocity: Float3d{
			X: m.WorldVelocityX,
			Y: m.WorldVelocityY,
			Z: m.WorldVelocityZ,
		},
		ForwardDir: Int3d{
			X: int(m.WorldForwardDirX),
			Y: int(m.WorldForwardDirY),
			Z: int(m.WorldForwardDirZ),
		},
		RightDir: Int3d{
			X: int(m.WorldRightDirX),
			Y: int(m.WorldRightDirY),
			Z: int(m.WorldRightDirZ),
		},
		GForce: Float3d{
			X: m.GForceLateral,
			Y: m.GForceLongitudinal,
			Z: m.GForceVertical,
		},
		Heading: Float3d{
			X: m.Yaw, // TODO: yaw -> Z, roll -> X, pitch -> Y
			Y: m.Roll,
			Z: m.Pitch,
		},
	}
}
