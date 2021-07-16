package packet

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	//"strconv"
)

const (
	MotionDataId uint8 = iota
	SessionDataId
	LapDataId
	EventId
	ParticipantsId
	CarSetupsId
	CarTelemetryId
	CarStatusId
)

// F1-2021
type Header struct {
	PacketFormat            uint16  `json:"m_packetFormat"`
	MajorGameVersion        uint8   `json:"m_gameMajorVersion"`
	MinorGameVersion        uint8   `json:"m_gameMinorVersion"`
	PacketVersion           uint8   `json:"m_packetVersion"`           // Version of this packet type, all start from 1
	PacketId                uint8   `json:"m_packetId"`                // Identifier for the packet type, see below
	SessionUid              uint64  `json:"m_sessionUID"`              // Unique identifier for the session
	SessionTime             float32 `json:"m_sessionTime"`             // Session timestamp
	FrameIdentifier         uint32  `json:"m_frameIdentifier"`         // Identifier for the frame the data was retrieved on
	PlayerCarIndex          uint8   `json:"m_playerCarIndex"`          // Index of player's car in the array
	SecondaryPlayerCarIndex uint8   `json:"m_secondaryPlayerCarIndex"` // Index of secondary player's car in the array (splitscreen) 255 if no second player
}

var (
	errNotEnoughData = errors.New("not enough data")
	errNotPointer    = errors.New("not pointer")
)

func ParseHeader(b []byte) (h Header, err error) {
	h = Header{}
	err = ParsePacket(b, &h)
	return
}

func ParsePacket(b []byte, model interface{}) (err error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return errNotPointer
	}

	if len(b) < Sizeof(value) {
		return errNotEnoughData
	}

	subValue := value.Elem()

	dataIndex := 0
	size := 0
	if subValue.Kind() == reflect.Array {
		firstElem := subValue.Index(0)
		kind := firstElem.Kind()
		size = Sizeof(firstElem)
		for i, n := 0, subValue.Len(); i < n; i, dataIndex = i+1, dataIndex+size {
			elemValue := subValue.Index(i)
			if !elemValue.IsValid() || !elemValue.CanSet() {
				continue
			}

			if kind != reflect.Struct && kind != reflect.Array {
				size = Sizeof(elemValue)
			}

			switch elemValue.Kind() {
			case reflect.Int8:
				elemValue.SetInt(int64(parseInt8(b[dataIndex])))
			case reflect.Uint8:
				elemValue.SetUint(uint64(parseUint8(b[dataIndex])))
			case reflect.Uint16:
				elemValue.SetUint(uint64(parseUint16(b[dataIndex : dataIndex+2])))
			case reflect.Uint32:
				elemValue.SetUint(uint64(parseUint32(b[dataIndex : dataIndex+4])))
			case reflect.Uint64:
				elemValue.SetUint(uint64(parseUint32(b[dataIndex : dataIndex+8])))
			case reflect.Float32:
				elemValue.SetUint(uint64(parseFloat32(b[dataIndex : dataIndex+8])))
			case reflect.Struct:
				fallthrough
			case reflect.Array:
				err = ParsePacket(b[dataIndex:], elemValue.Addr().Interface())
				if err != nil {
					return
				}
			}
		}
		return
	}

	dataIndex = 0
	size = 0

	for i, n := 0, subValue.NumField(); i < n; i, dataIndex, size = i+1, dataIndex+size, 0 {
		fieldValue := subValue.Field(i)
		kind := fieldValue.Kind()
		if kind != reflect.Struct && kind != reflect.Array {
			size = Sizeof(fieldValue)
		}

		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}
		switch fieldValue.Kind() {
		case reflect.Int8:
			fieldValue.SetInt(int64(parseInt8(b[dataIndex])))
		case reflect.Uint8:
			fieldValue.SetUint(uint64(parseUint8(b[dataIndex])))
		case reflect.Uint16:
			fieldValue.SetUint(uint64(parseUint16(b[dataIndex : dataIndex+2])))
		case reflect.Uint32:
			fieldValue.SetUint(uint64(parseUint32(b[dataIndex : dataIndex+4])))
		case reflect.Uint64:
			fieldValue.SetUint(uint64(parseUint32(b[dataIndex : dataIndex+8])))
		case reflect.Float32:
			fieldValue.SetUint(uint64(parseFloat32(b[dataIndex : dataIndex+8])))
		case reflect.Struct:
			fallthrough
		case reflect.Array:
			err = ParsePacket(b[dataIndex:], fieldValue.Addr().Interface())
			if err != nil {
				return
			}

			dataIndex += Sizeof(fieldValue)
		}
	}
	return
}

//func ParseHeader(b []byte) (h Header) {
//	packetFormat := binary.LittleEndian.Uint16(b[:2])
//	if packetFormat != 2019 {
//		panic("not f1 2019")
//	}
//	h.PacketFormat = int(packetFormat)
//
//	majorVersion := parseUint8(b[2])
//	minorVersion := parseUint8(b[3])
//	h.GameVersion = strconv.Itoa(int(majorVersion)) + "." + strconv.Itoa(int(minorVersion))
//
//	h.PacketVersion = parseUint8(b[4])
//	h.PacketId = Id(parseUint8(b[5]))
//	h.SessionUid = binary.LittleEndian.Uint64(b[6:14])
//	h.SessionTime = parseFloat32(b[14:18])
//	h.FrameIdentifier = binary.LittleEndian.Uint32(b[18:22])
//	h.PlayerCarIndex = parseUint8(b[22])
//	return
//}

func parseInt8(b byte) (r int8) {
	return int8(binary.LittleEndian.Uint16(append([]byte{}, b, 0)))
}
func parseUint8(b byte) (r uint8) {
	return uint8(binary.LittleEndian.Uint16(append([]byte{}, b, 0)))
}

func parseUint16(b []byte) (r uint16) {
	return binary.LittleEndian.Uint16(b)
}

func parseUint32(b []byte) (r uint32) {
	return binary.LittleEndian.Uint32(b)
}

func parseUint64(b []byte) (r uint64) {
	return binary.LittleEndian.Uint64(b)
}

func parseFloat32(b []byte) (r float32) {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

type Data interface{}

type Packet struct {
	Header Header
	Data   Data
}
