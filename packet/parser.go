package packet

func ParsePacket(b []byte) (p Packet) {
	p.Header = ParseHeader(b[:23])
	switch p.Header.PacketId {
	case MotionDataId:
		p.Data = ParseMotionData(b[23:])
	}

	return
}
