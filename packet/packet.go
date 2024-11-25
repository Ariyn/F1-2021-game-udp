package packet

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	//"strconv"
)

type Types interface {
	Header |
		MotionData |
		SessionData |
		LapData |
		ParticipantData |
		CarTelemetryData |
		EventData |
		SessionStarted |
		SessionEnded |
		FastestLap |
		Retirement |
		DRSEnabled |
		DRSDisabled |
		TeamMateInPits |
		Flashback |
		Buttons |
		CarSetupData |
		CarStatusData |
		FinalClassificationData |
		LobbyInfoData |
		CarDamageData |
		SessionHistoryData |
		RaceWinner |
		SpeedTrap |
		Penalty |
		StartLights |
		DriveThroughPenaltyServed |
		StopGoPenaltyServed |
		ChequeredFlag |
		LightsOut
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

var NamesById = map[Id]string{
	MotionDataId:          "MotionData",
	SessionDataId:         "SessionData",
	LapDataId:             "LapData",
	EventId:               "Event",
	ParticipantsId:        "Participants",
	CarSetupsId:           "CarSetups",
	CarTelemetryDataId:    "CarTelemetryData",
	CarStatusId:           "CarStatus",
	FinalClassificationId: "FinalClassification",
	LobbyInfoId:           "LobbyInfo",
	CarDamageId:           "CarDamage",
	SessionHistoryId:      "SessionHistory",
}

type Data interface {
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
	errInvalid       = errors.New("invalid value")
)

func ParseHeader(b []byte) (h Header, err error) {
	err = ParsePacket(b[:HeaderSize], &h)
	return
}

func ParsePacketGeneric[T Types](b []byte) (data T, err error) {
	err = ParsePacket(b, &data)
	return
}

func ParsePacket[T Types](b []byte, model *T) (err error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return errNotPointer
	}

	if len(b) < Sizeof(value) {
		return errNotEnoughData
	}

	_, err = Parse(value.Elem(), b)
	return
}

func Parse(elemValue reflect.Value, b []byte) (parsed reflect.Value, err error) {
	if !elemValue.IsValid() || !elemValue.CanSet() {
		return elemValue, errInvalid
	}

	switch elemValue.Kind() {
	case reflect.Int8:
		elemValue.SetInt(int64(parseInt8(b[0])))
	case reflect.Uint8:
		elemValue.SetUint(uint64(parseUint8(b[0])))
	case reflect.Int16:
		elemValue.SetInt(int64(parseInt16(b[:2])))
	case reflect.Uint16:
		elemValue.SetUint(uint64(parseUint16(b[:2])))
	case reflect.Int32:
		elemValue.SetInt(int64(parseInt32(b[:4])))
	case reflect.Uint32:
		elemValue.SetUint(uint64(parseUint32(b[:4])))
	case reflect.Int64:
		elemValue.SetInt(parseInt64(b[:8]))
	case reflect.Uint64:
		elemValue.SetUint(parseUint64(b[:8]))
	case reflect.Float32:
		elemValue.SetFloat(float64(parseFloat32(b[:8])))
	case reflect.Struct:
		_struct, err := parseStruct(b, elemValue)
		if err != nil {
			return elemValue, err
		}

		elemValue.Set(_struct)
	case reflect.Array:
		elem, err := parseArray(b, elemValue)
		if err != nil {
			return elemValue, err
		}
		elemValue.Set(elem)
	}

	return elemValue, nil
}

func parseInt8(b byte) (r int8) {
	return int8(binary.LittleEndian.Uint16(append([]byte{}, b, 0)))
}
func parseUint8(b byte) (r uint8) {
	return uint8(binary.LittleEndian.Uint16(append([]byte{}, b, 0)))
}

func parseInt16(b []byte) (r int16) {
	return int16(binary.LittleEndian.Uint16(b))
}
func parseUint16(b []byte) (r uint16) {
	return binary.LittleEndian.Uint16(b)
}

func parseInt32(b []byte) (r int32) {
	return int32(binary.LittleEndian.Uint32(b))
}
func parseUint32(b []byte) (r uint32) {
	return binary.LittleEndian.Uint32(b)
}

func parseInt64(b []byte) (r int64) {
	return int64(binary.LittleEndian.Uint64(b))
}
func parseUint64(b []byte) (r uint64) {
	return binary.LittleEndian.Uint64(b)
}

func parseFloat32(b []byte) (r float32) {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

func parseArray(b []byte, array reflect.Value) (v reflect.Value, err error) {
	size := 0
	dataIndex := 0

	for i, n := 0, array.Len(); i < n; i, dataIndex = i+1, dataIndex+size {
		v, err = Parse(array.Index(i), b[dataIndex:])
		if err != nil {
			return array, err
		}

		array.Index(i).Set(v)
		size = Sizeof(array.Index(i))
	}

	return array, nil
}

func parseStruct(b []byte, _struct reflect.Value) (v reflect.Value, err error) {
	dataIndex := 0
	size := 0

	for i, n := 0, _struct.NumField(); i < n; i, dataIndex, size = i+1, dataIndex+size, 0 {
		_v, err := Parse(_struct.Field(i), b[dataIndex:])
		if err != nil {
			return _struct, err
		}

		_struct.Field(i).Set(_v)
		size = Sizeof(_struct.Field(i))
	}

	return _struct, nil
}

func FormatPacket(model interface{}) (b []byte, err error) {
	value := reflect.ValueOf(model)
	size := Sizeof(value)

	if size == -1 {
		return nil, errors.New("invalid size")
	}
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
