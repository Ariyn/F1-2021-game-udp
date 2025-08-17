package logger

import (
	"context"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ariyn/F1-2021-game-udp/packet"
)

var _ packet.Logger = (*FileClient)(nil)

type FileClient struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	Path       string
	Timestamp  time.Time
	file       *os.File
	packetChan chan packet.Data
}

// RawWriter implements packet.Logger.
func (fc *FileClient) RawWriter(ctx context.Context, wg *sync.WaitGroup) (c chan<- []byte, cancel context.CancelFunc, err error) {
	panic("unimplemented")
}

func NewFileClient(p string, t time.Time) (fc *FileClient, err error) {
	fc = &FileClient{
		Path:       p,
		Timestamp:  t,
		packetChan: make(chan packet.Data, 1000),
	}

	_, err = fc.createFolder(fc.Path)
	if err != nil {
		return
	}

	fc.file, err = os.Create(path.Join(fc.Path, fc.Timestamp.Format("20060102150405")+".data"))
	if err != nil {
		panic(err)
	}

	return
}

func (fc *FileClient) Writer(ctx context.Context, wg *sync.WaitGroup) (c chan<- packet.Data, cancel context.CancelFunc, err error) {
	fc.ctx, cancel = context.WithCancel(ctx)
	fc.wg = wg
	return fc.packetChan, cancel, nil
}

func (fc *FileClient) Run() {
	defer fc.wg.Done()
	defer func() {
		err := fc.file.Close()
		if err != nil {
			log.Println("failed to close file", err)
		}
	}()

	for packetData := range fc.packetChan {
		data, err := packet.FormatPacket(packetData)
		if err != nil {
			log.Println("failed to write packet data", err, packetData.GetHeader().PacketId)
		}

		l, err := fc.file.Write(data)
		if err != nil {
			log.Println("failed to write packet data", err, packetData.GetHeader().PacketId)
		}

		if l != len(data) {
			log.Println("failed to write packet data", err, packetData.GetHeader().PacketId)
		}
	}
}

func (fc *FileClient) createFolder(pathElement ...string) (p string, err error) {
	p = path.Join(pathElement...)
	err = os.MkdirAll(p, 0755)
	return
}
