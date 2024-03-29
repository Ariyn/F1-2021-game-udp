package packet_test

import (
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_LapDataSize(t *testing.T) {
	l := packet.LapData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	assert.Equal(t, size, packet.LapDataSize)
}

func Test_MotionDataSize(t *testing.T) {
	l := packet.MotionData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	assert.Equal(t, size, packet.MotionDataSize)
}

func Test_SessionDataSize(t *testing.T) {
	l := packet.SessionData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	assert.Equal(t, size, packet.SessionDataSize)
}

func Test_ParticipantDataSize(t *testing.T) {
	l := packet.ParticipantData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	assert.Equal(t, size, packet.ParticipantDataSize)
}

func Test_EventHeaderDataSize(t *testing.T) {
	l := packet.EventData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	log.Println(size)
}

func Test_CarTelemetryDataSize(t *testing.T) {
	l := packet.CarTelemetryData{}
	size := packet.Sizeof(reflect.ValueOf(l))
	assert.Equal(t, size, packet.CarTelemetryDataSize)
}

func TestParseModel(t *testing.T) {
	type _testModel2 struct {
		C [4]uint8
		D uint32
	}
	type _testModel3 struct {
		E int8
	}
	type _testModel struct {
		A uint32
		B uint8
		C _testModel2
		E [4]_testModel3
	}

	testModel := _testModel{}
	err := packet.ParsePacket([]byte{255, 0, 0, 0, 2, 3, 3, 3, 3, 6, 0, 0, 0, 8, 255, 8, 255}, &testModel)
	assert.NoError(t, err)
	assert.Equal(t, _testModel{
		A: 255,
		B: 2,
		C: _testModel2{
			C: [4]uint8{3, 3, 3, 3},
			D: 6,
		},
		E: [4]_testModel3{
			{8}, {-1}, {8}, {-1},
		},
	}, testModel)
}

func TestFormatPacket(t *testing.T) {
	type _testModel2 struct {
		C [4]uint8
		D uint32
	}
	type _testModel3 struct {
		E int8
	}
	type _testModel struct {
		A uint32
		B uint8
		C _testModel2
		E [4]_testModel3
	}

	testModel := _testModel{
		A: 255,
		B: 2,
		C: _testModel2{
			C: [4]uint8{3, 3, 3, 3},
			D: 6,
		},
		E: [4]_testModel3{
			{8}, {-1}, {8}, {-1},
		},
	}

	formattedBytes, err := packet.FormatPacket(testModel)
	assert.NoError(t, err)
	assert.Equal(t, []byte{255, 0, 0, 0, 2, 3, 3, 3, 3, 6, 0, 0, 0, 8, 255, 8, 255}, formattedBytes)

	formattedBytes, err = packet.FormatPacket(uint16(8))
	assert.NoError(t, err)
	assert.Equal(t, []byte{8, 0}, formattedBytes)

	formattedBytes, err = packet.FormatPacket(_testModel3{8})
	assert.NoError(t, err)
	assert.Equal(t, []byte{8}, formattedBytes)
}
