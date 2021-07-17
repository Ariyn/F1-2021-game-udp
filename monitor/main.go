package main

import (
	"encoding/json"
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"net"
	"os"
	"path"
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

func Writer(storagePath string) (c chan packetData, err error) {
	l, err := logger.NewLogger(storagePath, time.Now(), 22)
	if err != nil {
		return
	}

	fmt.Printf("%s\n\nSave To %s", FormulaLoggerLogo, l.Path)

	c = make(chan packetData, 100)
	go write(c, l)
	return
}

func write(c <-chan packetData, l logger.Logger) {
	//err := l.NewLap(oldLapNumber, 22)
	//if err != nil {
	//	panic(err)
	//}
	defer l.Close()

	neverSavedParticipants := true

	driversLap := make([]int, 22)
	for i := 0; i < 22; i++ {
		driversLap[i] = -1
	}
	driverNames := make([]string, 22)

	eventLogger := log.New(os.Stdout, "[EVENT]", log.Ltime)

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

		// TODO: CarStatusData, FinalClassificationData, Event, CarDamageData, PacketSessionHistoryData
		switch header.PacketId {
		case packet.EventId:
			eventHeader := packet.EventHeaderData{}
			err = packet.ParsePacket(data, &eventHeader)
			if err != nil {
				log.Println("packet.Event", err)
				continue
			}

			// TODO: make these code into const
			switch eventHeader.StringCode() {
			case "SEND":
				eventLogger.Println("Session Ended!")
				panic("session end") // TODO: TEMP CODE
			case "SSTA":
				eventLogger.Println("Session Started!")
			case "CHQF":
				eventLogger.Println("Chequered flag.")
			case "TMPT":
				eventLogger.Println("Teammate is in pits.")
			case "LGOT":
				event := packet.StartLights{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("START LIGHT! %d", int(event.NumberLights))
			case "FTLP":
				event := packet.FastestLap{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("Fastest Lap by %s - %s!", driverNames[event.VehicleIndex], time.Duration(event.LapTime*1000)*time.Millisecond)
			case "RTMT":
				event := packet.Retirement{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("Retired %s!", driverNames[event.VehicleIndex])
			}
			continue
		case packet.ParticipantsId:
			participants := packet.ParticipantData{}
			err = packet.ParsePacket(data, &participants)
			if err != nil {
				panic(err)
			}

			if neverSavedParticipants {
				drivers := make([]Driver, 22)

				for carIndex, p := range participants.Participants {
					name := packet.DriverNameById[p.DriverId]
					drivers[carIndex] = Driver{
						Id:       int(p.DriverId),
						Name:     name,
						TeamName: packet.TeamNameById[p.TeamId],
						CarIndex: carIndex,
						IsAi:     p.IsAiControlled == 1,
					}
					driverNames[carIndex] = name
				}
				neverSavedParticipants = false

				data, err := json.MarshalIndent(drivers, "", "\t")
				if err != nil {
					log.Println(err)
					continue
				}

				err = l.WriteText("drivers.json", string(data))
				if err != nil {
					log.Println(err)
				}
				continue
			}

		case packet.MotionDataId:
			motion := packet.MotionData{}
			err = packet.ParsePacket(data, &motion)
			if err != nil {
				panic(err)
			}

			for carIndex, m := range motion.CarMotionData {
				var data []byte

				if carIndex == int(motion.Header.PlayerCarIndex) {
					playerMot := f1.GetPlayerMotionData(packetData.Timestamp, &motion)

					data, err = packet.FormatPacket(playerMot)
					if err != nil {
						panic(err)
					}
				} else {
					mot := f1.GetMotionData(packetData.Timestamp, m)

					data, err = packet.FormatPacket(mot)
					if err != nil {
						panic(err)
					}
				}

				// TODO: l.Write도 goroutine을 통해 비동기로 싱행되게 수정
				err = l.Write(packet.MotionDataId, carIndex, data)
				if err != nil {
					panic(err)
				}
			}
		case packet.CarTelemetryDataId:
			carTelemetry := packet.CarTelemetryData{}
			err = packet.ParsePacket(data, &carTelemetry)
			if err != nil {
				panic(err)
			}

			for carIndex, tlm := range carTelemetry.CarTelemetries {
				smpTlm := f1.SimplifyTelemetry(packetData.Timestamp, tlm)

				// TODO: 이거 l.Write가 interface로 바로 SmpTlm을 받을 수 있게 수정
				data, err := packet.FormatPacket(smpTlm)
				if err != nil {
					panic(err)
				}

				err = l.Write(packet.CarTelemetryDataId, carIndex, data)
				if err != nil {
					panic(err)
				}
			}
		case packet.LapDataId:
			lap := packet.LapData{}
			err = packet.ParsePacket(data, &lap)
			if err != nil {
				panic(err)
			}

			for carIndex, lap := range lap.DriverLaps {
				smpLap := f1.SimplifyLap(packetData.Timestamp, lap)
				currentLapNumber := int(smpLap.CurrentLapNumber)
				if driversLap[carIndex] != currentLapNumber {
					err = l.NewLap(currentLapNumber, carIndex)
					if err != nil {
						panic(err)
					}

					driversLap[carIndex] = currentLapNumber
				}

				// TODO: 이거 l.Write가 interface로 바로 SmpTlm을 받을 수 있게 수정
				data, err := packet.FormatPacket(smpLap)
				if err != nil {
					panic(err)
				}

				err = l.Write(packet.LapDataId, carIndex, data)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
