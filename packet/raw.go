package packet

type Raw struct {
	H    Header
	Buf  []byte
	Size int
}

func (r *Raw) Id() Id {
	return Id(r.H.PacketId)
}

func (r *Raw) GetHeader() Header {
	return r.H
}
