package packet

import (
	"encoding/binary"
	"math"
	"strconv"
)

type Id uint8

const (
	MotionDataId Id = iota
	SessionDataId
	LapDataId
	EventId
	ParticipantsId
	CarSetupsId
	CarTelemetryId
	CarStatusId
)

//https://stackoverflow.com/questions/16330490/in-go-how-can-i-convert-a-struct-to-a-byte-array/48134333
type Header struct {
	PacketFormat int
	GameVersion  string

	PacketVersion   uint8   `json:"m_packetVersion"`   // Version of this packet type, all start from 1
	PacketId        Id      `json:"m_packetId"`        // Identifier for the packet type, see below
	SessionUid      uint64  `json:"m_sessionUID"`      // Unique identifier for the session
	SessionTime     float32 `json:"m_sessionTime"`     // Session timestamp
	FrameIdentifier uint32  `json:"m_frameIdentifier"` // Identifier for the frame the data was retrieved on
	PlayerCarIndex  uint8   `json:"m_playerCarIndex"`  // Index of player's car in the array
}

func ParseHeader(b []byte) (h Header) {
	packetFormat := binary.LittleEndian.Uint16(b[:2])
	if packetFormat != 2019 {
		panic("not f1 2019")
	}
	h.PacketFormat = int(packetFormat)

	majorVersion := parseUint8(b[2])
	minorVersion := parseUint8(b[3])
	h.GameVersion = strconv.Itoa(int(majorVersion)) + "." + strconv.Itoa(int(minorVersion))

	h.PacketVersion = parseUint8(b[4])
	h.PacketId = Id(parseUint8(b[5]))
	h.SessionUid = binary.LittleEndian.Uint64(b[6:14])
	h.SessionTime = parseFloat32(b[14:18])
	h.FrameIdentifier = binary.LittleEndian.Uint32(b[18:22])
	h.PlayerCarIndex = parseUint8(b[22])
	return
}

func parseUint8(b byte) (r uint8) {
	return uint8(binary.LittleEndian.Uint16(append([]byte{}, b, 0)))
}

func parseFloat32(b []byte) (r float32) {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

type Data interface{}

type Packet struct {
	Header Header
	Data   Data
}
