package logger

import (
	"errors"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

var packetIds = []int{
	int(packet.MotionDataId),
	int(packet.SessionDataId),
	int(packet.LapDataId),
	int(packet.EventId),
	int(packet.ParticipantsId),
	int(packet.CarSetupsId),
	int(packet.CarTelemetryDataId),
	int(packet.CarStatusId),
	int(packet.FinalClassificationId),
	int(packet.LobbyInfoId),
	int(packet.CarDamageId),
	int(packet.SessionHistoryId),
}

const GeneralData = uint8(255)

type Logger struct {
	Path             string
	storage          string
	Timestamp        time.Time
	maximumCarNumber int
	files            [][]*os.File // [driverIndex][packetId]
	generalFile      *os.File
	rawFiles         []*os.File
}

func NewLogger(p string, t time.Time, maxCarNumber int) (l Logger, err error) {
	l = Logger{
		Path:             p,
		Timestamp:        t,
		maximumCarNumber: maxCarNumber,
	}

	l.storage, err = l.createFolder(l.Path, l.Timestamp.Format("2006-01-02"), l.Timestamp.Format("150405"))
	if err != nil {
		return
	}

	rawStorage, err := l.createFolder(l.storage, "raw")
	if err != nil {
		return
	}

	l.files = make([][]*os.File, l.maximumCarNumber)
	for driverIndex := 0; driverIndex < l.maximumCarNumber; driverIndex++ {
		l.files[driverIndex] = make([]*os.File, len(packetIds))
	}

	l.generalFile, err = os.Create(path.Join(l.storage, "general"))
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(packetIds); i++ {
		var f *os.File
		f, err = os.Create(path.Join(rawStorage, strconv.Itoa(packetIds[i])))
		if err != nil {
			return
		}
		l.rawFiles = append(l.rawFiles, f)
	}

	for i := 0; i < l.maximumCarNumber; i++ {
		err = l.NewLap(-1, i)
		if err != nil {
			return
		}
	}

	return
}

func (l Logger) createFolder(pathElement ...string) (p string, err error) {
	p = path.Join(pathElement...)
	err = os.MkdirAll(p, 0755)
	return
}

func (l *Logger) NewLap(lap, driverIndex int) (err error) {
	p, err := l.createFolder(l.storage, strconv.Itoa(lap))
	if err != nil && !os.IsExist(err) {
		return
	}

	for packetIndex, id := range packetIds {
		if l.files[driverIndex][packetIndex] != nil {
			err = l.files[driverIndex][packetIndex].Close()
			if err != nil {
				return
			}
		}

		l.files[driverIndex][packetIndex], err = os.Create(path.Join(p, fmt.Sprintf("%d-%d", driverIndex, id)))
		if err != nil {
			return
		}
	}

	return nil
}

func (l Logger) Write(id uint8, carIndex int, data []byte) (err error) {
	var f *os.File
	if id == GeneralData {
		f = l.generalFile
	} else {
		f = l.files[carIndex][id]
	}
	n, err := f.Write(data)
	if err != nil {
		return
	}

	if n != len(data) {
		err = errors.New("not enough write")
	}
	return
}

func (l Logger) WriteAsync(id uint8, carIndex int, data []byte) {
	go func(id uint8, carIndex int, data []byte) {
		var f *os.File
		if id == GeneralData {
			f = l.generalFile
		} else {
			f = l.files[carIndex][id]
		}
		n, err := f.Write(data)
		if err != nil {
			return
		}

		if n != len(data) {
			err = errors.New("not enough write")
		}
	}(id, carIndex, data)
	return
}

func (l Logger) WriteRaw(id uint8, data []byte) (err error) {
	n, err := l.rawFiles[id].Write(data)
	if err != nil {
		return
	}

	if n != len(data) {
		err = errors.New("not enough write")
	}
	return
}

func (l Logger) WriteRawAsync(id uint8, data []byte) {
	go func(id uint8, data []byte) {
		n, err := l.rawFiles[id].Write(data)
		if err != nil {
			return
		}

		if n != len(data) {
			err = errors.New("not enough write")
		}
	}(id, data)
	return
}

func (l Logger) WriteText(name, value string) (err error) {
	return ioutil.WriteFile(path.Join(l.storage, name), []byte(value), 0755)
}

func (l Logger) WriteTextAsync(name, value string) {
	go ioutil.WriteFile(path.Join(l.storage, name), []byte(value), 0755)
	return
}

func (l Logger) Close() (err error) {
	for _, fs := range l.files {
		for _, f := range fs {
			err = f.Close()
			if err != nil && err != os.ErrClosed {
				return
			}
		}
	}
	for _, f := range l.rawFiles {
		err = f.Close()
		if err != nil && err != os.ErrClosed {
			return
		}
	}
	return
}
