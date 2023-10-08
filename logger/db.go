package logger

import (
	"context"
	"log"
	"strconv"
	//f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Header struct {
	gorm.Model
	PacketId        uint8   `gorm:"uniqueIndex:header,sort:desc"`
	SessionUid      string  `gorm:"uniqueIndex:header,sort:desc"`
	SessionTime     float32 `gorm:"uniqueIndex:header,sort:desc"`
	FrameIdentifier uint32  `gorm:"uniqueIndex:header,sort:desc"`
	PlayerCarIndex  uint8
}

type CarTelemetry struct {
	gorm.Model
	Header                    Header `gorm:"foreignKey:ID"`
	IsPlayer                  bool
	Speed                     uint16
	Throttle                  float32
	Steer                     float32
	Break                     float32
	Clutch                    uint8
	Gear                      int8
	EngineRPM                 uint16
	DRS                       bool
	RevLightsPercent          uint8
	RevLightsBitValue         uint16
	BreaksTemperatureFL       uint16
	BreaksTemperatureFR       uint16
	BreaksTemperatureRL       uint16
	BreaksTemperatureRR       uint16
	TyresSurfaceTemperatureFL uint8
	TyresSurfaceTemperatureFR uint8
	TyresSurfaceTemperatureRL uint8
	TyresSurfaceTemperatureRR uint8
	TyresInnerTemperatureFL   uint8
	TyresInnerTemperatureFR   uint8
	TyresInnerTemperatureRL   uint8
	TyresInnerTemperatureRR   uint8
	EngineTemperature         uint16
	TyresPressureFL           float32
	TyresPressureFR           float32
	TyresPressureRL           float32
	TyresPressureRR           float32
	SurfaceTypeFL             uint8
	SurfaceTypeFR             uint8
	SurfaceTypeRL             uint8
	SurfaceTypeRR             uint8
}

var _ packet.Logger = (*DBLogger)(nil)

type DBLogger struct {
	db           *gorm.DB
	inputChannel chan packet.Raw
	ctx          context.Context
}

func NewDBLogger() (*DBLogger, error) {
	ormDB, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logger := &DBLogger{
		db: ormDB,
	}

	err = ormDB.AutoMigrate(&Header{}, &CarTelemetry{})
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *DBLogger) Writer(ctx context.Context) (c chan<- packet.Raw, cancel context.CancelFunc, err error) {
	l.inputChannel = make(chan packet.Raw, 0)
	l.ctx, cancel = context.WithCancel(ctx)
	return l.inputChannel, cancel, nil
}

func (l *DBLogger) Run() {
	ct := CarTelemetry{
		Header: Header{
			PacketId:        1,
			SessionUid:      "12345",
			SessionTime:     0.1,
			FrameIdentifier: 1,
		},
	}
	err := l.db.Create(&ct.Header).Create(ct).Commit().Error
	if err != nil {
		log.Fatal(err)
	}

	ct2 := CarTelemetry{
		Header: Header{
			PacketId:        1,
			SessionUid:      "12345",
			SessionTime:     0.1,
			FrameIdentifier: 1,
		},
	}
	err = l.db.Create(&ct2.Header).Create(ct2).Commit().Error
	if err != nil {
		log.Fatal(err)
	}

	// TODO: listen from inputChannel
	for data := range l.inputChannel {
		header, err := packet.ParseHeader(data.Buf)
		if err != nil {
			panic(err)
		}

		//timestamp := int64(time.Duration(header.SessionTime*1000) * time.Millisecond)
		log.Println(header.SessionTime)
		//now := time.Unix(timestamp, 0)

		switch header.PacketId {
		case packet.CarTelemetryDataId:
			//log.Println(now.Format("2006-01-02 15:04:05"))
			carTelemetry := packet.CarTelemetryData{}
			err = packet.ParsePacket(data.Buf, &carTelemetry)
			if err != nil {
				panic(err)
			}

			for carIndex, tlm := range carTelemetry.CarTelemetries {
				ct := CarTelemetry{}
				ct.Speed = tlm.Speed
				ct.Throttle = tlm.Throttle
				ct.Steer = tlm.Steer
				ct.Break = tlm.Break
				ct.Clutch = tlm.Clutch
				ct.Gear = tlm.Gear
				ct.EngineRPM = tlm.EngineRPM
				ct.DRS = tlm.DRS == 1
				ct.RevLightsPercent = tlm.RevLightsPercent
				ct.RevLightsBitValue = tlm.RevLightsBitValue
				ct.BreaksTemperatureRL = tlm.BreaksTemperature[0]
				ct.BreaksTemperatureRR = tlm.BreaksTemperature[1]
				ct.BreaksTemperatureFL = tlm.BreaksTemperature[2]
				ct.BreaksTemperatureFR = tlm.BreaksTemperature[3]
				ct.TyresSurfaceTemperatureRL = tlm.TyresSurfaceTemperature[0]
				ct.TyresSurfaceTemperatureRR = tlm.TyresSurfaceTemperature[1]
				ct.TyresSurfaceTemperatureFL = tlm.TyresSurfaceTemperature[2]
				ct.TyresSurfaceTemperatureFR = tlm.TyresSurfaceTemperature[3]
				ct.TyresInnerTemperatureRL = tlm.TyresInnerTemperature[0]
				ct.TyresInnerTemperatureRR = tlm.TyresInnerTemperature[1]
				ct.TyresInnerTemperatureFL = tlm.TyresInnerTemperature[2]
				ct.TyresInnerTemperatureFR = tlm.TyresInnerTemperature[3]
				ct.EngineTemperature = tlm.EngineTemperature
				ct.TyresPressureRL = tlm.TyresPressure[0]
				ct.TyresPressureRR = tlm.TyresPressure[1]
				ct.TyresPressureFL = tlm.TyresPressure[2]
				ct.TyresPressureFR = tlm.TyresPressure[3]
				ct.SurfaceTypeRL = tlm.SurfaceType[0]
				ct.SurfaceTypeRR = tlm.SurfaceType[1]
				ct.SurfaceTypeFL = tlm.SurfaceType[2]
				ct.SurfaceTypeFR = tlm.SurfaceType[3]

				ct.Header = Header{
					PacketId:        header.PacketId,
					SessionUid:      strconv.FormatUint(header.SessionUid, 10),
					SessionTime:     header.SessionTime,
					FrameIdentifier: header.FrameIdentifier,
					PlayerCarIndex:  uint8(carIndex),
				}
				ct.IsPlayer = carIndex == int(header.PlayerCarIndex)

				l.db.Create(&ct)
			}
		}
	}
}

func parseEvent(data packet.Raw) (err error) {
	eventHeader := packet.EventHeaderData{}
	err = packet.ParsePacket(data.Buf, &eventHeader)
	if err != nil {
		return
	}

	/*
		SSTA : SESSION STARTED
		SEND : SESSION ENDED
		FTLP : FASTEST LAP
		RTMT : RETIREMENT
		DRSE : DRS ENABLED
		DRSD : DRS DISABLES
		TMPT : TEAMMATES IN PITS
		CHQF : CHEUERED FLAGS
		RCWN : RACE WINNER ANNOUNCED
		PENA : PENALTY ISSUED
		SPTP : SPEED TRAP TRIGGERED BY FASTEST SPEED
		STLG : START LIGHTS
		LGOT : LIGHTS OUT
		DTSV : DRIVE THROUGH PENALTY SERVED
		SGSV : STOP AND GO PENALTY SERVED
		FLBK : FLASHBACK ACTIVATED
		BUTN : BUTTON STATUS CHANGED
	*/
	switch eventHeader.StringCode() {
	case "SSTA":

	}
	return
}
