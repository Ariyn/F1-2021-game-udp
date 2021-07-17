package main

import (
	"encoding/binary"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/ariyn/F1-2021-game-udp/visualizer"
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
	defer f.Close()

	for packetData := range c {
		header, err := packet.ParseHeader(packetData.Buf)
		if err != nil {
			panic(err)
		}

		log.Printf("%#v", header)

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

			log.Println(uint64(packetData.Timestamp))
			playerTlm := carTelemetry.CarTelemetries[int(carTelemetry.Header.PlayerCarIndex)]
			data := visualizer.SimplifiedTelemetry{
				TimeStamp: uint64(packetData.Timestamp),
				Steer:     playerTlm.Steer,
				Throttle:  playerTlm.Throttle,
				Break:     playerTlm.Break,
				Gear:      playerTlm.Gear,
				EngineRPM: playerTlm.EngineRPM,
				Speed:     playerTlm.Speed,
				DRS:       playerTlm.DRS,
				BreaksTemperature: f1.Int16Wheel{
					RL: int16(playerTlm.BreaksTemperature[0]),
					RR: int16(playerTlm.BreaksTemperature[1]),
					FL: int16(playerTlm.BreaksTemperature[2]),
					FR: int16(playerTlm.BreaksTemperature[3]),
				},
				TyresSurfaceTemperature: f1.Int8Wheel{
					RL: int8(playerTlm.TyresSurfaceTemperature[0]),
					RR: int8(playerTlm.TyresSurfaceTemperature[1]),
					FL: int8(playerTlm.TyresSurfaceTemperature[2]),
					FR: int8(playerTlm.TyresSurfaceTemperature[3]),
				},
				TyresInnerTemperature: f1.Int8Wheel{
					RL: int8(playerTlm.TyresInnerTemperature[0]),
					RR: int8(playerTlm.TyresInnerTemperature[1]),
					FL: int8(playerTlm.TyresInnerTemperature[2]),
					FR: int8(playerTlm.TyresInnerTemperature[3]),
				},
				EngineTemperature: 0,
				TyresPressure: f1.FloatWheels{
					RL: playerTlm.TyresPressure[0],
					RR: playerTlm.TyresPressure[1],
					FL: playerTlm.TyresPressure[2],
					FR: playerTlm.TyresPressure[3],
				},
				SurfaceType: f1.Int8Wheel{
					RL: int8(playerTlm.SurfaceType[0]),
					RR: int8(playerTlm.SurfaceType[1]),
					FL: int8(playerTlm.SurfaceType[2]),
					FR: int8(playerTlm.SurfaceType[3]),
				},
			}

			n, err := f.Write(b)
			if err != nil {
				panic(err)
			}
			if n != 25 {
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
				f, err = createLapFolder(storagePath, currentLapNumber, packet.LapDataId)
				if err != nil {
					panic(err)
				}
				defer f.Close()

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

func createLapFolder(storagePath string, lapNumber int, dataType uint8) (f *os.File, err error) {
	newFolder := path.Join(storagePath, strconv.Itoa(lapNumber))
	err = os.Mkdir(newFolder, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	return os.Create(path.Join(newFolder, strconv.Itoa(int(dataType))))
}
