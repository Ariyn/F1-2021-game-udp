package packet

import (
	"context"
	"sync"
)

type Logger interface {
	RawWriter(ctx context.Context, wg *sync.WaitGroup) (c chan<- []byte, cancel context.CancelFunc, err error)
	Writer(ctx context.Context, wg *sync.WaitGroup) (c chan<- Data, cancel context.CancelFunc, err error)
	Run()
}
