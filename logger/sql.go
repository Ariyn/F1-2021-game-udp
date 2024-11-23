package logger

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var _ packet.Logger = (*SqlClient)(nil)

type SqlClient struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	Url        string
	client     *sql.DB
	packetChan chan packet.Data
}

func NewSqlClient(url string) (sc *SqlClient, err error) {
	sc = &SqlClient{
		Url:        url,
		packetChan: make(chan packet.Data, 1000),
	}

	sc.client, err = sql.Open("postgres", sc.Url)
	if err != nil {
		return nil, err
	}

	return
}

func (sc *SqlClient) Writer(ctx context.Context, wg *sync.WaitGroup) (c chan<- packet.Data, cancel context.CancelFunc, err error) {
	sc.ctx, cancel = context.WithCancel(ctx)
	sc.wg = wg
	return sc.packetChan, cancel, nil
}

func (sc *SqlClient) Run() {
	defer sc.wg.Done()
	defer func() {
		err := sc.client.Close()
		if err != nil {
			log.Println("failed to close file", err)
		}
	}()

	for packetData := range sc.packetChan {
		b, err := packet.FormatPacket(packetData)
		if err != nil {
			log.Println("failed to format packet", err)
			continue
		}

		header := packetData.GetHeader()
		packetId := packet.Id(header.PacketId)
		sessionUid := header.SessionUid
		frameIdentifier := header.FrameIdentifier
		playerCarIndex := header.PlayerCarIndex
		sessionTime := header.SessionTime

		_, err = sc.client.ExecContext(sc.ctx, "INSERT INTO packets (data, packet_id, packet_type, session_uid, frame_identifier, player_car_index, session_time) VALUES ($1, $2, $3, $4, $5, $6, $7)", b, packetId, packet.NamesById[packetId], fmt.Sprintf("%d", sessionUid), frameIdentifier, playerCarIndex, sessionTime)
		if err != nil {
			log.Println("failed to insert packet data", err)
			continue
		}
	}
}
