package main

import (
	"encoding/binary"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"math"
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

	files := make([]*os.File, 8)
	for i := 0; i < 8; i++ {
		f, err := os.Create("/tmp/f1-" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		files[i] = f
	}

	storagePath = os.TempDir()
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
	if _, err = os.Stat(storagePath); os.IsNotExist(err) {
		return
	}

	c = make(chan packetData, 100)
	go write(c, storagePath)
	return
}

func write(c <-chan packetData, storagePath string) {
	oldLapNumber := -1
	f, err := createLapFolder(storagePath, oldLapNumber)
	if err != nil {
		panic(err)
	}

	for packetData := range c {
		header, err := packet.ParseHeader(packetData.Buf)
		if err != nil {
			panic(err)
		}

		data := packetData.Buf[:packetData.Size]
		switch header.PacketId {
		case packet.MotionDataId:
			motion := packet.MotionData{}
			err = packet.ParsePacket(data, &motion)
			if err != nil {
				panic(err)
			}
		case packet.CarTelemetryDataId:
			carTelemetry := packet.CarTelemetryData{}
			err = packet.ParsePacket(data, &carTelemetry)
			if err != nil {
				panic(err)
			}

			playerTlm := carTelemetry.CarTelemetries[int(carTelemetry.Header.PlayerCarIndex)]
			b := make([]byte, 0)
			binary.LittleEndian.PutUint64(b, uint64(packetData.Timestamp))
			binary.LittleEndian.PutUint32(b, math.Float32bits(playerTlm.Steer))
			binary.LittleEndian.PutUint32(b, math.Float32bits(playerTlm.Throttle))
			binary.LittleEndian.PutUint32(b, math.Float32bits(playerTlm.Break))

			n, err := f.Write(b)
			if err != nil {
				panic(err)
			}
			if n != 20 {
				panic("not enough write")
			}
		case packet.LapDataId:
			lap := packet.LapData{}
			err = packet.ParsePacket(data, &lap)
			if err != nil {
				panic(err)
			}

			currentLapNumber := int(lap.DriverLaps[int(lap.Header.PlayerCarIndex)].CurrentLapNumber)
			if currentLapNumber != oldLapNumber {
				f, err = createLapFolder(storagePath, currentLapNumber)
				if err != nil {
					panic(err)
				}

				oldLapNumber = currentLapNumber
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

func createLapFolder(storagePath string, lapNumber int) (f *os.File, err error) {
	newFolder := path.Join(storagePath, strconv.Itoa(lapNumber))
	err = os.Mkdir(newFolder, 0755)
	if err != nil {
		return
	}

	return os.Open(path.Join(newFolder, "data"))
}
