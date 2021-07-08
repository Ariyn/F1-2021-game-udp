package main

import (
	"github.com/ariyn/f1/packet"
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
			data := packet.ParsePacket(buf[:n])

			go func(buf []byte) {
				_, err = files[data.Header.PacketId].Write(buf)
				if err != nil {
					panic(err)
				}
			}(buf[:n])
			counters[data.Header.PacketId]++

			if counters[data.Header.PacketId]%100 == 0 {
				log.Println(data.Header.PacketId, counters[data.Header.PacketId])
			}
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
