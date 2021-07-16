package main

import (
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"net"
	"os"
	"strconv"
)

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

	counters := make([]int, 8)
	for {
		buf := make([]byte, 2048) // all telemetry data is under 2048 bytes.
		n, _, err := network.ReadFrom(buf)
		if n == 0 {
			log.Println("buffer size is 0...")
		}
		if n > 0 {
			header, err := packet.ParseHeader(buf[:n])
			if err != nil {
				panic(err)
			}

			var data interface{}
			switch header.PacketId {
			case packet.MotionDataId:
				data = packet.MotionData{}
				err = packet.ParsePacket(buf[:n], &data)
				if err != nil {
					panic(err)
				}
			case packet.CarTelemetryDataId:
				data = packet.CarTelemetryData{}
				err = packet.ParsePacket(buf[:n], &data)
				if err != nil {
					panic(err)
				}
			}

			go func(buf []byte) {
				_, err = files[header.PacketId].Write(buf)
				if err != nil {
					panic(err)
				}
			}(buf[:n])
			counters[header.PacketId]++

			if counters[header.PacketId]%100 == 0 {
				log.Println(header.PacketId, counters[header.PacketId])
			}
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
