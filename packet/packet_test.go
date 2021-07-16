package packet_test

import (
	"github.com/ariyn/f1/packet"
	"github.com/stretchr/testify/assert"
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
