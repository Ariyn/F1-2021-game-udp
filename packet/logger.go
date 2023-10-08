package packet

import "context"

type Logger interface {
	Writer(ctx context.Context) (c chan<- Raw, cancel context.CancelFunc, err error)
	Run()
}
