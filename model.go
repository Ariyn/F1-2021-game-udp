package F1_2021_game_udp

type Float3d struct {
	X float32
	Y float32
	Z float32
}

type Int163d struct {
	X int16
	Y int16
	Z int16
}

type FloatWheels struct {
	RL float32
	RR float32
	FL float32
	FR float32
}

type MotionData struct {
	Position   Float3d
	Velocity   Float3d
	ForwardDir Int163d
	RightDir   Int163d
	GForce     Float3d
	Heading    Float3d
}

type PlayerMotionData struct {
	MotionData
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
