package logger

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/ariyn/F1-2021-game-udp/packet"
	_ "github.com/marcboeker/go-duckdb/v2"
)

var _ packet.Logger = (*DuckDBClient)(nil)

type DuckDBClient struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	Path       string
	client     *sql.DB
	packetChan chan packet.Data
}

func NewDuckDBClient(path string) (dc *DuckDBClient, err error) {
	dc = &DuckDBClient{
		Path:       path,
		packetChan: make(chan packet.Data, 1000),
	}

	dc.client, err = sql.Open("duckdb", dc.Path)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = dc.client.Exec(`
		CREATE TABLE IF NOT EXISTS packets (
			data BLOB,
			packet_id UINTEGER,
			packet_type VARCHAR,
			session_uid UBIGINT,
			frame_identifier UINTEGER,
			player_car_index UINTEGER,
			session_time FLOAT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return
}

func (dc *DuckDBClient) Writer(ctx context.Context, wg *sync.WaitGroup) (c chan<- packet.Data, cancel context.CancelFunc, err error) {
	dc.ctx, cancel = context.WithCancel(ctx)
	dc.wg = wg
	return dc.packetChan, cancel, nil
}

func (dc *DuckDBClient) Run() {
	defer dc.wg.Done()
	defer func() {
		err := dc.client.Close()
		if err != nil {
			log.Println("failed to close duckdb client", err)
		}
	}()

	for packetData := range dc.packetChan {
		b, err := packet.FormatPacket(packetData)
		if err != nil {
			continue
		}

		header := packetData.GetHeader()
		packetId := packet.Id(header.PacketId)
		sessionUid := header.SessionUid
		frameIdentifier := header.FrameIdentifier
		playerCarIndex := header.PlayerCarIndex
		sessionTime := header.SessionTime

		_, err = dc.client.ExecContext(dc.ctx, "INSERT INTO packets (data, packet_id, packet_type, session_uid, frame_identifier, player_car_index, session_time) VALUES (?, ?, ?, ?, ?, ?, ?)", b, packetId, packet.NamesById[packetId], sessionUid, frameIdentifier, playerCarIndex, sessionTime)
		if err != nil {
			log.Println("failed to insert packet data into duckdb", err)
			continue
		}
	}
}
