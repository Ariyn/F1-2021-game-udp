package packet

import "context"

type Logger interface {
	Writer(ctx context.Context) (c chan<- PacketData, cancel context.CancelFunc, err error)
	Run()
}
