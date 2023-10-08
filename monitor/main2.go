package main

import (
	"context"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/ariyn/F1-2021-game-udp/prometheus"
	"log"
	"math"
	"strconv"
)

func init() {

}

func Writer2() (ctx context.Context, c chan packetData, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	c = make(chan packetData, 100)
	go write2(c, cancel)
	return
}

func write2(c chan packetData, cancel context.CancelFunc) {
	defer cancel()

	var logger map[int]prometheus.NewTypeLogger
	logger = make(map[int]prometheus.NewTypeLogger)

	neverSavedParticipants := true
	drivers := make(map[int]f1.Driver, 0)

	l, err := prometheus.NewLogger(22)
	if err != nil {
		panic(err)
	}

	go l.Run()

	playerSpeed := float64(0)
	retired := make(map[int]bool)

	paused := false
	for packetData := range c {
		header, err := packet.ParseHeader(packetData.Buf)
		if err != nil {
			panic(err)
		}

		data := packetData.Buf[:packetData.Size]

		var previousDriverLogger, nextDriverLogger prometheus.NewTypeLogger

		switch header.PacketId {
		case packet.SessionDataId:
			session := packet.SessionData{}
			err = packet.ParsePacket(data, &session)
			if err != nil {
				panic(err)
			}

			paused = session.GamePaused == 1
		case packet.ParticipantsId:
			participants := packet.ParticipantData{}
			err = packet.ParsePacket(data, &participants)
			if err != nil {
				panic(err)
			}

			if neverSavedParticipants {
				for carIndex, p := range participants.Participants[:20] {
					name := packet.DriverNameById[p.DriverId]
					id := p.DriverId
					if p.IsAiControlled == 0 {
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

				previousDriverLogger = l.Driver(f1.Driver{Name: "previous"}).
					SetSession(strconv.FormatUint(participants.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()

				nextDriverLogger = l.Driver(f1.Driver{Name: "next"}).
					SetSession(strconv.FormatUint(participants.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()

				neverSavedParticipants = false
				continue
			}

		case packet.MotionDataId:
			motion := packet.MotionData{}

			if !paused {
				err = packet.ParsePacket(data, &motion)
				if err != nil {
					panic(err)
				}
			}

			ld, ok := logger[int(motion.Header.PlayerCarIndex)]
			if !ok {
				ld = l.Driver(drivers[int(motion.Header.PlayerCarIndex)]).
					SetSession(strconv.FormatUint(motion.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()

				logger[int(motion.Header.PlayerCarIndex)] = ld
			}

			ld.WorldVelocityX(motion.Player().WorldVelocityX)
			ld.WorldVelocityY(motion.Player().WorldVelocityY)
			ld.WorldVelocityZ(motion.Player().WorldVelocityZ)

			//speed := playerSpeed * 1000 / 3600
			//for index, m := range motion.CarMotionData[:len(drivers)] {
			//	if m == motion.Player() {
			//		continue
			//	}
			//	d := distance(m, motion.Player())
			//
			//	log.Println(drivers[index].Name, d, d/speed, speed)
			//	ld.Distance(d)
			//	ld.Delta(d / speed)
			//}

		case packet.CarTelemetryDataId:
			carTelemetry := packet.CarTelemetryData{}
			if !paused {
				err = packet.ParsePacket(data, &carTelemetry)
				if err != nil {
					panic(err)
				}
			}

			ld, ok := logger[int(carTelemetry.Header.PlayerCarIndex)]
			if !ok {
				ld = l.Driver(drivers[int(carTelemetry.Header.PlayerCarIndex)]).
					SetSession(strconv.FormatUint(carTelemetry.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()

				logger[int(carTelemetry.Header.PlayerCarIndex)] = ld
			}

			ld.Throttle(carTelemetry.Player().Throttle)
			ld.Break(carTelemetry.Player().Break)
			ld.Speed(carTelemetry.Player().Speed)
			ld.EngineRPM(carTelemetry.Player().EngineRPM)
			ld.Gear(carTelemetry.Player().Gear)
			ld.Steer(carTelemetry.Player().Steer)

			playerSpeed = float64(carTelemetry.Player().Speed)

		case packet.LapDataId:
			lap := packet.LapData{}

			if !paused {
				err = packet.ParsePacket(data, &lap)
				if err != nil {
					panic(err)
				}
			}

			if previousDriverLogger.NeedInit() {
				previousDriverLogger = l.Driver(f1.Driver{Name: "previous"}).
					SetSession(strconv.FormatUint(lap.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()
			}

			if nextDriverLogger.NeedInit() {
				nextDriverLogger = l.Driver(f1.Driver{Name: "next"}).
					SetSession(strconv.FormatUint(lap.Header.SessionUid, 10)).
					SetLapNumber(1).
					Init()
			}

			playerPosition := int(lap.Player().CarPosition)
			speed := playerSpeed * 1000 / 3600
			for carIndex, driverLap := range lap.DriverLaps[:len(drivers)] {
				if _, ok := retired[carIndex]; ok {
					continue
				}

				ld, ok := logger[carIndex]
				if !ok {
					ld = l.Driver(drivers[carIndex]).
						SetSession(strconv.FormatUint(lap.Header.SessionUid, 10)).
						SetLapNumber(driverLap.CurrentLapNumber).
						Init()

					logger[carIndex] = ld
				}

				var previousDriver, nextDriver bool
				if int(driverLap.CarPosition)-playerPosition == 1 {
					previousDriver = true
				}

				if int(driverLap.CarPosition)-playerPosition == -1 {
					nextDriver = true
				}

				ld.TotalDistance(driverLap.TotalDistance)
				ld.Position(driverLap.CarPosition)

				ld.Lap(driverLap.CurrentLapNumber)
				ld.Sector(driverLap.Sector)
				ld.PitStatus(driverLap.PitStatus)

				if driverLap != lap.Player() {
					d := float64(driverLap.TotalDistance - lap.Player().TotalDistance)

					ld.Distance(d)

					delta := d / speed
					ld.Delta(delta)

					if previousDriver {
						previousDriverLogger.Delta(delta)
					}
					if nextDriver {
						nextDriverLogger.Delta(delta)
					}
				}
			}

		case packet.EventId:
			eventHeader := packet.EventHeaderData{}
			err = packet.ParsePacket(data, &eventHeader)
			if err != nil {
				log.Println("packet.Event", err)
				continue
			}

			// TODO: make these code into const
			switch eventHeader.StringCode() {
			case "RTMT":
				event := packet.Retirement{}
				err = packet.ParsePacket(data[packet.HeaderSize+4:], &event)
				if err != nil {
					log.Println(err)
					continue
				}

				index := int(event.VehicleIndex)
				if ld, ok := logger[int(event.VehicleIndex)]; ok {
					ld.Finish()
				}

				retired[index] = true

				log.Println(drivers[index].Name, "retired")
			}
		}
	}
}

func distance(a, b packet.CarMotionData) (distance float64) {
	x := a.WorldPositionX - b.WorldPositionX
	y := a.WorldPositionY - b.WorldPositionY
	z := a.WorldPositionZ - b.WorldPositionZ

	return math.Sqrt(float64(x*x + y*y + z*z))
}
