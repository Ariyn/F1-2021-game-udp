package packet

import (
	"log"
	"reflect"
)

var (
	HeaderKind = reflect.TypeOf(Header{}).Kind()
)

func ParsePacket(b []byte) (p Packet) {
	p.Header = ParseHeader(b[:23])
	switch p.Header.PacketId {
	case MotionDataId:
		p.Data = ParsePacketData(b[23:], MotionData{})
	}

	return
}

func ParsePacketData(data []byte, model interface{}) (d Data) {
	typ := reflect.TypeOf(model)

	for i := 0; i < typ.NumField(); i++ {
		log.Println(typ.Field(i))
		field := typ.Field(i)
		switch field.Type.Kind() {
		case HeaderKind:
			log.Println("header")
		case reflect.Array:
			log.Println("array", field.Type.Elem())
		}
	}

	return
}
