package main

import (
	"context"
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
var visualizerPath = ""

func main() {
	//visualizerPath = "C:\\Users\\ariyn\\Documents\\go\\src\\github.com\\ariyn\\F1-2021-game-udp\\visualizer"
	network, err := net.ListenPacket("udp", "0.0.0.0:1278")
	if err != nil {
		panic(err)
	}

	storagePath = path.Join(os.TempDir(), "f1")
	err = os.Mkdir(storagePath, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	ctx, c, err := Writer(storagePath)
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

		select {
		case c <- packetData{
			Buf:  buf,
			Size: n,
		}:
		case <-ctx.Done():
			log.Println("new session will be started")

			ctx, c, err = Writer(storagePath)
			if err != nil {
				panic(err)
			}
		}
	}

	// TODO: unreachable code
	close(c)
}

// TODO: use packet models only. let visualizer parse every raw data
func Writer(storagePath string) (ctx context.Context, c chan packetData, err error) {
	l, err := logger.NewLogger(storagePath, time.Now(), 22)
	if err != nil {
		return
	}

	fmt.Printf("%s\n\nSave To %s", FormulaLoggerLogo, l.Path)

	ctx, cancel := context.WithCancel(context.Background())
	c = make(chan packetData, 100)
	go write(cancel, c, l)
	return
}

// TODO: more obvious name
func write(cancel context.CancelFunc, c <-chan packetData, l logger.Logger) {
	// TODO: defer error handling?
	defer l.Close()
	defer cancel()

	neverSavedParticipants := true

	driversLap := make([]int, 22)
	for i := 0; i < 22; i++ {
		driversLap[i] = -1
	}

	drivers := make([]f1.Driver, 22)
	defer func(drivers *[]f1.Driver, driversLap *[]int) {
		data, err := json.MarshalIndent(&drivers, "", "\t")
		if err != nil {
			log.Println(err)
			return
		}

		err = l.WriteText("drivers.json", string(data))
		if err != nil {
			log.Println(err)
		}

		data, err = json.MarshalIndent(&driversLap, "", "\t")
		if err != nil {
			log.Println(err)
			return
		}

		err = l.WriteText("driver-laps.json", string(data))
		if err != nil {
			log.Println(err)
		}
	}(&drivers, &driversLap)

	eventLogger := log.New(os.Stdout, "[EVENT]", log.Ltime)

	for packetData := range c {
		header, err := packet.ParseHeader(packetData.Buf)
		if err != nil {
			panic(err)
		}

		timestamp := int64(time.Duration(header.SessionTime*1000) * time.Millisecond)

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
				// TODO: call visualizer
				if visualizerPath != "" {
					//err = exec.Command("go", "run", "--path", visualizerPath).Run()
					//if err != nil {
					//	panic(err)
					//}
				}
				return
			case "SSTA":
				eventLogger.Println("Session Started!")
			case "CHQF":
				eventLogger.Println("Chequered flag.")
			case "TMPT":
				eventLogger.Println("Teammate is in pits.")
			case "STLG":
				event := packet.StartLights{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("START LIGHT! %d", int(event.NumberLights))
			case "LGOT":
				eventLogger.Printf("%s\nLIGHTS OUT AND A WAY WE GO!", LightsOutLogo)
			case "FTLP":
				event := packet.FastestLap{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("Fastest Lap by %s - %s!", drivers[event.VehicleIndex].Name, time.Duration(event.LapTime*1000)*time.Millisecond)
			case "RTMT":
				event := packet.Retirement{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}
				eventLogger.Printf("Retired %s!", drivers[event.VehicleIndex].Name)
			case "PENA":
				event := packet.Penalty{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}

				gained := ""
				if event.Time != 255 {
					gained = fmt.Sprintf(" gained %ds", event.Time)
				} else if event.PlacesGained != 255 {
					gained = fmt.Sprintf(" gained %d grid positions", event.PlacesGained)
				}

				due := packet.InfringementsById[event.InfringementType]
				if event.OtherVehicleIndex != 255 {
					due += " with " + drivers[event.OtherVehicleIndex].Name
				}
				eventLogger.Printf("Penalty %s to %s, due to %s%s", packet.PenaltiesById[event.PenaltyType], drivers[event.VehicleIndex].Name, due, gained)
			}
			continue
		case packet.ParticipantsId:
			participants := packet.ParticipantData{}
			err = packet.ParsePacket(data, &participants)
			if err != nil {
				panic(err)
			}

			if neverSavedParticipants {
				for carIndex, p := range participants.Participants {
					name := packet.DriverNameById[p.DriverId]
					drivers[carIndex] = f1.Driver{
						Id:         int(p.DriverId),
						Name:       name,
						RaceNumber: int(p.RaceNumber),
						TeamName:   packet.TeamNameById[p.TeamId],
						CarIndex:   carIndex,
						IsAi:       p.IsAiControlled == 1,
					}
				}
				neverSavedParticipants = false
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
					playerMot := f1.GetPlayerMotionData(timestamp, &motion)

					data, err = packet.FormatPacket(playerMot)
					if err != nil {
						panic(err)
					}
				} else {
					mot := f1.GetMotionData(timestamp, m)

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
				smpTlm := f1.SimplifyTelemetry(timestamp, tlm)

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
				smpLap := f1.SimplifyLap(timestamp, lap)
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
		case packet.SessionDataId:
			session := packet.SessionData{}
			err = packet.ParsePacket(data, &session)
			if err != nil {
				panic(err)
			}

			smpSess := f1.SimplifySession(timestamp, session)
			data, err := packet.FormatPacket(smpSess)
			if err != nil {
				panic(err)
			}

			err = l.Write(logger.GeneralData, 0, data)
			if err != nil {
				panic(err)
			}
		}
	}
}
