package packet

import (
	"context"
	"log"
	"net"
	"sync"
)

const (
	DefaultNetwork = "udp"
	DefaultAddress = "0.0.0.0:1946"
)

type Listener struct {
	ctx            context.Context
	conn           net.PacketConn
	loggers        []Logger
	loggerChannels []chan<- Data
	loggerCancels  []context.CancelFunc
	waitGroup      *sync.WaitGroup
	started        bool
}

func NewListener(ctx context.Context, network, address string, loggers ...Logger) (l *Listener, err error) {
	l = &Listener{
		ctx:       ctx,
		loggers:   loggers,
		started:   false,
		waitGroup: &sync.WaitGroup{},
	}

	l.waitGroup.Add(len(loggers))
	for _, logger := range loggers {
		channel, cancel, err := logger.Writer(ctx, l.waitGroup)
		if err != nil {
			return l, err
		}

		l.loggerChannels = append(l.loggerChannels, channel)
		l.loggerCancels = append(l.loggerCancels, cancel)
	}

	l.conn, err = net.ListenPacket(network, address)
	if err != nil {
		return
	}

	return
}

func (l *Listener) Run() (err error) {
	defer l.conn.Close()
	defer func() {
		for _, cancel := range l.loggerCancels {
			cancel()
		}
	}()
	defer func() {
		for _, channel := range l.loggerChannels {
			close(channel)
		}
	}()

	for _, logger := range l.loggers {
		go logger.Run()
	}

	for {
		buf := make([]byte, 2048) // all telemetry data is under 2048 bytes.
		n, _, err := l.conn.ReadFrom(buf)
		if err != nil {
			// If the context is cancelled, the connection will be closed and ReadFrom will return an error.
			// This is expected, so we just exit gracefully.
			if l.ctx.Err() != nil {
				break
			}
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

		// It's crucial to slice the buffer to the actual number of bytes read.
		rawBytes := buf[:n]

		header, err := ParseHeader(rawBytes)
		if err != nil {
			log.Printf("failed to parse packet header: %v", err)
			continue
		}

		// Create a raw packet struct containing the full buffer and header.
		// This is much faster than parsing the entire packet body.
		data := &Raw{
			H:    header,
			Buf:  rawBytes,
			Size: n,
		}

		// Send the raw packet to all logger channels.
		for _, channel := range l.loggerChannels {
			channel <- data
		}

		select {
		case <-l.ctx.Done():
			log.Println("context cancelled, shutting down listener")
			break
		default:
		}
	}

	l.waitGroup.Wait()
	log.Println("listener shut down gracefully")
	return nil
}
