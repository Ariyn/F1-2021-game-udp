package packet

import (
	"encoding/binary"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type IntegerVector3 struct {
	X int
	Y int
	Z int
}

type CarMotionData struct {
	WorldPosition   Vector3 `json:"worldPosition"`
	WorldVelocity   Vector3 `json:"worldVelocity"`
	WorldForwardDir IntegerVector3
	WorldRightDir   IntegerVector3
	GForce          Vector3
	Rotation        Vector3
}

type MotionData struct {
	CarMotionData [20]CarMotionData
}

func ParseMotionData(b []byte) (m MotionData) {
	for i := 0; i < 20; i++ {
		m.CarMotionData[i] = parseCarMotionData(b[i*60 : (i+1)*60])
	}

	return
}

func parseCarMotionData(b []byte) (c CarMotionData) {
	c.WorldPosition.X = float64(parseFloat32(b[:4]))
	c.WorldPosition.Y = float64(parseFloat32(b[4:8]))
	c.WorldPosition.Z = float64(parseFloat32(b[8:12]))

	c.WorldVelocity.X = float64(parseFloat32(b[12:16]))
	c.WorldVelocity.Y = float64(parseFloat32(b[16:20]))
	c.WorldVelocity.Z = float64(parseFloat32(b[20:24]))

	c.WorldForwardDir.X = int(binary.LittleEndian.Uint16(b[24:26]))
	c.WorldForwardDir.Y = int(binary.LittleEndian.Uint16(b[26:28]))
	c.WorldForwardDir.Z = int(binary.LittleEndian.Uint16(b[28:30]))

	c.WorldRightDir.X = int(binary.LittleEndian.Uint16(b[30:32]))
	c.WorldRightDir.Y = int(binary.LittleEndian.Uint16(b[32:34]))
	c.WorldRightDir.Z = int(binary.LittleEndian.Uint16(b[34:36]))

	c.GForce.X = float64(parseFloat32(b[36:40]))
	c.GForce.Y = float64(parseFloat32(b[40:44]))
	c.GForce.Z = float64(parseFloat32(b[44:48]))

	c.Rotation.X = float64(parseFloat32(b[48:52]))
	c.Rotation.Y = float64(parseFloat32(b[52:56]))
	c.Rotation.Z = float64(parseFloat32(b[56:60]))
	return c
}
