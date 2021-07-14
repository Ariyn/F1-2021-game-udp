package packet_test

import (
	"github.com/ariyn/f1/packet"
	"github.com/magiconair/properties/assert"
	"reflect"
	"testing"
)

func Test_LapDataSize(t *testing.T) {
	l := packet.LapData{}
	size := packet.Sizeof(reflect.TypeOf(l))
	assert.Equal(t, size, packet.LapDataSize)
}

func Test_MotionDataSize(t *testing.T) {
	l := packet.MotionData{}
	size := packet.Sizeof(reflect.TypeOf(l))
	assert.Equal(t, size, packet.MotionDataSize)
}

func Test_ParsePacketData(t *testing.T) {
	m := packet.MotionData{}
	packet.ParsePacketData(nil, m)
}
