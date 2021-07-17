package main

import (
	f1 "github.com/ariyn/F1-2021-game-udp"
	logger "github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

var storagePath = "/tmp"

// need packet decoder, encoder

// folder structure
// session_uid/p1/packet_name
func main() {
	network, err := net.ListenPacket("udp", "0.0.0.0:1278")
	if err != nil {
		panic(err)
	}

	storagePath = path.Join(os.TempDir(), "f1")
	err = os.Mkdir(storagePath, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	log.Println(storagePath)
	c, err := Writer(storagePath)
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 2048) // all telemetry data is under 2048 bytes.
		n, _, err := network.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		if n == 0 {
			log.Println("buffer size is 0...")
			continue
		}

		c <- packetData{
			Buf:       buf,
			Size:      n,
			Timestamp: time.Now().UnixNano(),
		}
	}

	close(c)
}

type packetData struct {
	Buf       []byte
	Size      int
	Timestamp int64
}

func Writer(storagePath string) (c chan packetData, err error) {
	l, err := logger.NewLogger(storagePath, time.Now())
	if err != nil {
		return
	}

	c = make(chan packetData, 100)
	go write(c, l)
	return
}

func write(c <-chan packetData, l logger.Logger) {
	oldLapNumber := -1
	err := l.NewLap(oldLapNumber)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for packetData := range c {
		header, err := packet.ParseHeader(packetData.Buf)
		if err != nil {
			panic(err)
		}

		data := packetData.Buf[:packetData.Size]
		err = l.WriteRaw(header.PacketId, data)
		if err != nil {
			panic(err)
		}

		switch header.PacketId {
		case packet.MotionDataId:
			motion := packet.MotionData{}
			err = packet.ParsePacket(data, &motion)
			if err != nil {
				panic(err)
			}

			playerMot := f1.GetPlayerMotionData(packetData.Timestamp, &motion)

			data, err := packet.FormatPacket(playerMot)
			if err != nil {
				panic(err)
			}

			err = l.Write(packet.MotionDataId, data)
			if err != nil {
				panic(err)
			}
		case packet.CarTelemetryDataId:
			carTelemetry := packet.CarTelemetryData{}
			err = packet.ParsePacket(data, &carTelemetry)
			if err != nil {
				panic(err)
			}

			playerTlm := carTelemetry.Player()
			smpTlm := f1.SimplifyTelemetry(packetData.Timestamp, playerTlm)

			// TODO: 이거 l.Write가 interface로 바로 SmpTlm을 받을 수 있게 수정
			data, err := packet.FormatPacket(smpTlm)
			if err != nil {
				panic(err)
			}

			err = l.Write(packet.CarTelemetryDataId, data)
			if err != nil {
				panic(err)
			}
		case packet.LapDataId:
			lap := packet.LapData{}
			err = packet.ParsePacket(data, &lap)
			if err != nil {
				panic(err)
			}

			playerLap := lap.Player()

			currentLapNumber := int(playerLap.CurrentLapNumber)
			if currentLapNumber != oldLapNumber {
				err = l.NewLap(currentLapNumber)
				if err != nil {
					panic(err)
				}
				oldLapNumber = currentLapNumber
			}

			smpLap := f1.SimplifyLap(packetData.Timestamp, playerLap)

			data, err := packet.FormatPacket(smpLap)
			if err != nil {
				panic(err)
			}

			err = l.Write(packet.LapDataId, data)
			if err != nil {
				panic(err)
			}
		}

		//go func(buf []byte) {
		//	_, err = files[header.PacketId].Write(buf)
		//	if err != nil {
		//		panic(err)
		//	}
		//}(buf[:n])
		//counters[header.PacketId]++
		//
		//if counters[header.PacketId]%100 == 0 {
		//	log.Println(header.PacketId, counters[header.PacketId])
		//}
	}
}

func createLapFolder(storagePath string, lapNumber int, dataType uint8) (f *os.File, err error) {
	newFolder := path.Join(storagePath, strconv.Itoa(lapNumber))
	err = os.Mkdir(newFolder, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	return os.Create(path.Join(newFolder, strconv.Itoa(int(dataType))))
}
