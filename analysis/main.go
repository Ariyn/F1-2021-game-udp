package main

import (
	"github.com/ariyn/F1-2021-game-udp/packet"
	"io/ioutil"
	"log"
	"os"
)

var path = "D:\\Library\\Download\\2021-08-29\\162631\\raw\\1"

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	index := 0

	for i:=0; i<100; i++ {
		header, err := packet.ParseHeader(data[index:])
		if err != nil {
			panic(err)
		}
		//index += packet.HeaderSize

		log.Printf("%#v", header)
		log.Println(header.PacketFormat, header.MajorGameVersion, header.MinorGameVersion, header.PacketVersion, header.SessionTime, header.FrameIdentifier, packet)

		switch header.PacketId {
		case packet.MotionDataId:
			var md packet.MotionData
			err = packet.ParsePacket(data[index:index+packet.MotionDataSize], md)
			if err != nil {
				panic(err)
			}
			index += packet.MotionDataSize
			os.Exit(0)
		case packet.SessionDataId:
			index += packet.SessionDataSize
		case packet.LapDataId:
			index += packet.LapDataSize
		case packet.EventId:
			index += packet.EventDataSize
		case packet.ParticipantsId:
			index += packet.ParticipantDataSize
		case packet.CarSetupsId:
			index += packet.CarSetupsSize
		case packet.CarTelemetryDataId:
			index += packet.CarTelemetryDataSize
		case packet.CarStatusId:
			index += packet.CarStatusSize
		case packet.FinalClassificationId:
			index += packet.FinalClassificationSize
		case packet.LobbyInfoId:
			index += packet.LobbyInfoSize
		case packet.CarDamageId:
			index += packet.CarDamageSize
		case packet.SessionHistoryId:
			index += packet.SessionHistorySize
		default:
			log.Println(header)
			os.Exit(0)
		}
	}
}

