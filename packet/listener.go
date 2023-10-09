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
	ctx            context.Context
	conn           net.PacketConn
	loggers        []Logger
	loggerChannels []chan<- PacketData
	loggerCancels  []context.CancelFunc
	started        bool
}

func NewListener(ctx context.Context, network, address string, loggers ...Logger) (l Listener, err error) {
	l = Listener{
		ctx:     ctx,
		loggers: loggers,
		started: false,
	}

	for _, logger := range loggers {
		channel, cancel, err := logger.Writer(ctx)
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

		buf = buf[:n]
		header, err := ParseHeader(buf)
		if err != nil {
			panic(err)
		}

		var data PacketData
		switch Id(header.PacketId) {
		case MotionDataId:
			data, err = ParsePacketGeneric[MotionData](buf)
			if err != nil {
				panic(err)
			}
		case SessionDataId:
			data, err = ParsePacketGeneric[SessionData](buf)
			if err != nil {
				panic(err)
			}
		case CarTelemetryDataId:
			data, err = ParsePacketGeneric[CarTelemetryData](buf)
			if err != nil {
				panic(err)
			}
		case LapDataId:
			data, err = ParsePacketGeneric[LapData](buf)
			if err != nil {
				panic(err)
			}
		case ParticipantsId:
			data, err = ParsePacketGeneric[ParticipantData](buf)
			if err != nil {
				panic(err)
			}
		case EventId:
			data, err = ParsePacketGeneric[EventData](buf)
			if err != nil {
				panic(err)
			}

			v := data.(EventData)
			switch v.StringCode() {
			case SSTA:
				v.Event, err = ParsePacketGeneric[SessionStarted](buf[HeaderSize+4:])
			case SEND:
				v.Event, err = ParsePacketGeneric[SessionEnded](buf[HeaderSize+4:])
			case FTLP:
				v.Event, err = ParsePacketGeneric[FastestLap](buf[HeaderSize+4:])
			case FLBK:
				v.Event, err = ParsePacketGeneric[Flashback](buf[HeaderSize+4:])
			case BUTN:
				v.Event, err = ParsePacketGeneric[Buttons](buf[HeaderSize+4:])
			}
			if err != nil {
				panic(err)
			}

			data = v
		}

		if data != nil {
			for _, channel := range l.loggerChannels {
				channel <- data
			}
		}

		select {
		case <-l.ctx.Done():
			log.Println("new session will be started")
			break
		default:
		}
	}

	return nil
}
