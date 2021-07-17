package logger

import (
	"errors"
	"github.com/ariyn/F1-2021-game-udp/packet"
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

type Logger struct {
	Path      string
	storage   string
	Timestamp time.Time
	files     []*os.File
	rawFiles  []*os.File
}

func NewLogger(p string, t time.Time) (l Logger, err error) {
	l = Logger{
		Path:      p,
		Timestamp: t,
	}

	for i := 0; i < len(packetIds); i++ {
		l.files = append(l.files, nil)

		var f *os.File
		f, err = os.Create(path.Join(l.storage, "raw", strconv.Itoa(packetIds[i])))
		if err != nil {
			return
		}
		l.rawFiles = append(l.rawFiles, f)
	}

	err = l.init()

	return
}

func (l *Logger) init() (err error) {
	l.storage, err = l.createFolder(l.Path, l.Timestamp.Format("2006-01-02 15:04"))
	if err != nil {
		return
	}

	err = l.NewLap(-1)
	return
}

func (l Logger) createFolder(pathElement ...string) (p string, err error) {
	p = path.Join(pathElement...)
	err = os.Mkdir(p, 0755)
	return
}

func (l *Logger) NewLap(lap int) (err error) {
	p, err := l.createFolder(l.storage, strconv.Itoa(lap))
	if err != nil {
		return
	}

	for i, id := range packetIds {
		err = l.files[i].Close()
		if err != nil {
			return
		}
		l.files[i], err = os.Create(path.Join(p, strconv.Itoa(id)))
		if err != nil {
			return
		}
	}

	return nil
}

func (l Logger) Write(id uint8, data []byte) (err error) {
	n, err := l.files[id].Write(data)
	if err != nil {
		return
	}

	if n != len(data) {
		err = errors.New("not enough write")
	}
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

func (l Logger) Close() (err error) {
	for _, f := range l.files {
		err = f.Close()
		if err != nil && err != os.ErrClosed {
			return
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
