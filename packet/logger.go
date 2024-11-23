package packet

import (
	"context"
	"sync"
)

type Logger interface {
	Writer(ctx context.Context, wg *sync.WaitGroup) (c chan<- Data, cancel context.CancelFunc, err error)
	Run()
}
