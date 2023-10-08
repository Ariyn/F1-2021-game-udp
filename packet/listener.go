package packet

import (
	"context"
	"log"
	"net"
)

const (
	DefaultNetwork = "udp"
	DefaultAddress = "0.0.0.0:1278"
)

type Listener struct {
	ctx           context.Context
	conn          net.PacketConn
	logger        Logger
	loggerChannel chan<- Raw
	loggerCancel  context.CancelFunc
	started       bool
}

func NewListener(ctx context.Context, network, address string, logger Logger) (l Listener, err error) {
	l = Listener{
		ctx:     ctx,
		logger:  logger,
		started: false,
	}
	l.loggerChannel, l.loggerCancel, err = logger.Writer(ctx)
	if err != nil {
		return
	}

	l.conn, err = net.ListenPacket(network, address)
	if err != nil {
		return
	}

	return
}

func (l *Listener) Run() (err error) {
	defer l.conn.Close()
	defer l.loggerCancel()
	defer close(l.loggerChannel)

	go l.logger.Run()

	for {
		buf := make([]byte, 2048) // all telemetry data is under 2048 bytes.
		n, _, err := l.conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		if n == 0 {
			log.Println("buffer size is 0...")
			continue
		}
		if !l.started {
			log.Println("started!")
			l.started = true
		}

		select {
		case l.loggerChannel <- Raw{
			Buf:  buf,
			Size: n,
		}:
		case <-l.ctx.Done():
			log.Println("new session will be started")
			break
		}
	}

	return nil
}
