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

		// Parse the full packet body based on the header's packet ID.
		packetData := l.parsePacketBody(header, rawBytes)
		if packetData == nil {
			// This was an unknown or failed packet type, silently ignore.
			continue
		}

		// Send the parsed packet to all logger channels.
		for _, channel := range l.loggerChannels {
			channel <- packetData
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

func (l *Listener) parsePacketBody(header Header, rawBytes []byte) Data {
	var data Data
	var err error

	switch Id(header.PacketId) {
	case MotionDataId:
		var motionData MotionData
		err = ParsePacket(rawBytes, &motionData)
		data = &motionData
	case SessionDataId:
		var sessionData SessionData
		err = ParsePacket(rawBytes, &sessionData)
		data = &sessionData
	case LapDataId:
		var lapData LapData
		err = ParsePacket(rawBytes, &lapData)
		data = &lapData
	case EventId:
		var eventData EventData
		err = ParsePacket(rawBytes, &eventData)
		data = &eventData
	case ParticipantsId:
		var participantsData ParticipantData
		err = ParsePacket(rawBytes, &participantsData)
		data = &participantsData
	case CarSetupsId:
		var carSetupData CarSetupData
		err = ParsePacket(rawBytes, &carSetupData)
		data = &carSetupData
	case CarTelemetryDataId:
		var carTelemetryData CarTelemetryData
		err = ParsePacket(rawBytes, &carTelemetryData)
		data = &carTelemetryData
	case CarStatusId:
		var carStatusData CarStatusData
		err = ParsePacket(rawBytes, &carStatusData)
		data = &carStatusData
	case FinalClassificationId:
		var finalClassificationData FinalClassificationData
		err = ParsePacket(rawBytes, &finalClassificationData)
		data = &finalClassificationData
	case LobbyInfoId:
		var lobbyInfoData LobbyInfoData
		err = ParsePacket(rawBytes, &lobbyInfoData)
		data = &lobbyInfoData
	case CarDamageId:
		var carDamageData CarDamageData
		err = ParsePacket(rawBytes, &carDamageData)
		data = &carDamageData
	case SessionHistoryId:
		var sessionHistoryData SessionHistoryData
		err = ParsePacket(rawBytes, &sessionHistoryData)
		data = &sessionHistoryData
	default:
		log.Printf("unknown packet id: %d", header.PacketId)
		return nil
	}

	if err != nil {
		log.Printf("failed to parse packet body for id %d: %v", header.PacketId, err)
		return nil
	}

	return data
}
