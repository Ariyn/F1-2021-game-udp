package packet

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	//"strconv"
)

// TODO: CarSetupData, CarStatusData, FInalClassificationData, LobbyInfoData, CarDamageData, SessionHistoryData
type Types interface {
	MotionData | SessionData | LapData | ParticipantData | CarTelemetryData |
		EventData | SessionStarted | SessionEnded | FastestLap | Retirement | DRSEnabled | DRSDisabled | TeamMateInPits | Flashback | Buttons
}

type Id uint8

const (
	MotionDataId Id = iota
	SessionDataId
	LapDataId
	EventId
	ParticipantsId
	CarSetupsId
	CarTelemetryDataId
	CarStatusId
	FinalClassificationId
	LobbyInfoId
	CarDamageId
	SessionHistoryId
)

var Ids = []int{
	int(MotionDataId),
	int(SessionDataId),
	int(LapDataId),
	int(EventId),
	int(ParticipantsId),
	int(CarSetupsId),
	int(CarTelemetryDataId),
	int(CarStatusId),
	int(FinalClassificationId),
	int(LobbyInfoId),
	int(CarDamageId),
	int(SessionHistoryId),
}

type PacketData interface {
	Id() Id
	GetHeader() Header
}

// F1-2021
var HeaderSize = 24

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
	err = ParsePacket(b[:HeaderSize], &h)
	return
}

func ParsePacketGeneric[T Types](b []byte) (data T, err error) {
	err = ParsePacket(b, &data)
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
				elemValue.SetUint(parseUint64(b[dataIndex : dataIndex+8]))
			case reflect.Float32:
				elemValue.SetFloat(float64(parseFloat32(b[dataIndex : dataIndex+8])))
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
			fieldValue.SetUint(parseUint64(b[dataIndex : dataIndex+8]))
		case reflect.Float32:
			fieldValue.SetFloat(float64(parseFloat32(b[dataIndex : dataIndex+8])))
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

func FormatPacket(model interface{}) (b []byte, err error) {
	value := reflect.ValueOf(model)
	size := Sizeof(value)
	b = make([]byte, size)

	switch value.Kind() {
	case reflect.Array:
		var elemBytes []byte
		var index int
		for i, n := 0, value.Len(); i < n; i++ {
			elemBytes = nil
			element := value.Index(i)
			elemBytes, err = FormatPacket(element.Interface())
			if err != nil {
				return
			}

			b = copyBytes(b, elemBytes, index)
			index += len(elemBytes)
		}
	//case reflect.Slice:
	case reflect.Struct:
		var fieldBytes []byte
		var index int
		for i, n := 0, value.NumField(); i < n; i++ {
			fieldBytes = nil

			field := value.Field(i)
			fieldBytes, err = FormatPacket(field.Interface())
			if err != nil {
				return
			}

			b = copyBytes(b, fieldBytes, index)
			index += len(fieldBytes)
		}
	case reflect.Int8:
		b[0] = uint8(value.Int())
	case reflect.Uint8:
		b[0] = uint8(value.Uint())
	case reflect.Uint16:
		binary.LittleEndian.PutUint16(b, uint16(value.Uint()))
	case reflect.Uint32:
		binary.LittleEndian.PutUint32(b, uint32(value.Uint()))
	case reflect.Uint64:
		binary.LittleEndian.PutUint64(b, value.Uint())
	case reflect.Float32:
		binary.LittleEndian.PutUint32(b, math.Float32bits(float32(value.Float())))
	}
	return
}

func copyBytes(dest, src []byte, index int) []byte {
	for i := 0; i < len(src); i++ {
		dest[index+i] = src[i]
	}

	return dest
}
