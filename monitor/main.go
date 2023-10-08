package main

import (
	"context"
	"encoding/json"
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/logger"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"os"
	"strconv"
	"time"
)

var visualizerPath = ""

func main() {
	//db, err := sql.Open("sqlite3", "/tmp/f1")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//lgr, err := logger.NewDBLogger(db)
	lgr, err := logger.NewDBLogger()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := packet.NewListener(context.Background(), packet.DefaultNetwork, packet.DefaultAddress, lgr)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	fmt.Println("monitor start")
	defer fmt.Println("monitor ended")

	if err := listener.Run(); err != nil {
		log.Fatal(err)
	}
}

// TODO: use packet models only. let visualizer parse every raw data
func Writer(storagePath string) (ctx context.Context, c chan packetData, err error) {
	l, err := logger.NewLogger(storagePath, time.Now(), 22)
	if err != nil {
		return
	}

	fmt.Printf("%s\n\nSave To %s\n", FormulaLoggerLogo, l.Path)

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
		l.WriteRawAsync(header.PacketId, data)

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
				eventLogger.Println("SetSession Ended!")
				// TODO: call visualizer
				if visualizerPath != "" {
					//err = exec.Command("go", "run", "--path", visualizerPath).Run()
					//if err != nil {
					//	panic(err)
					//}
				}
				return
			case "SSTA":
				eventLogger.Println("SetSession Started!")
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
				eventLogger.Printf("Fastest SetLap by %s - %s!", drivers[event.VehicleIndex].Name, time.Duration(event.LapTime*1000)*time.Millisecond)
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
				} else if event.PlacesGained != 0 {
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
					id := p.DriverId
					if p.IsAiControlled == 0 {
						log.Println(carIndex, p.GetName())
						name = "player " + strconv.Itoa(carIndex+1)
						id = p.NetworkId
					}
					drivers[carIndex] = f1.Driver{
						Id:         int(id),
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

				// TODO: goroutine error handling
				l.WriteAsync(packet.MotionDataId, carIndex, data)
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

				l.WriteAsync(packet.CarTelemetryDataId, carIndex, data)
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

				l.WriteAsync(packet.LapDataId, carIndex, data)
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

			l.WriteAsync(logger.GeneralData, 0, data)
		}
	}
}
