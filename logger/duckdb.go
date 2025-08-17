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
	ctx           context.Context
	wg            *sync.WaitGroup
	Path          string
	client        *sql.DB
	packetChan    chan packet.Data
	rawPacketChan chan []byte
}

func NewDuckDBClient(path string) (dc *DuckDBClient, err error) {
	dc = &DuckDBClient{
		Path:          path,
		packetChan:    make(chan packet.Data, 1000),
		rawPacketChan: make(chan []byte, 1000),
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
			session_time FLOAT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS raw (
			data BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

// RawWriter implements packet.Logger.
func (dc *DuckDBClient) RawWriter(ctx context.Context, wg *sync.WaitGroup) (c chan<- []byte, cancel context.CancelFunc, err error) {
	dc.ctx, cancel = context.WithCancel(ctx)
	dc.wg = wg
	return dc.rawPacketChan, cancel, nil
}

func (dc *DuckDBClient) Run() {
	defer dc.wg.Done()
	defer func() {
		err := dc.client.Close()
		if err != nil {
			log.Println("failed to close duckdb client", err)
		}
	}()

	buffer := make([]([]byte), 10000)
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-dc.ctx.Done():
				tx, err := dc.client.BeginTx(dc.ctx, nil)
				if err != nil {
					mu.Unlock()
					log.Println("failed to begin transaction for raw packet insert", err)
					continue
				}
				stmt, err := tx.PrepareContext(dc.ctx, "INSERT INTO raw (data) VALUES (?)")
				if err != nil {
					tx.Rollback()
					mu.Unlock()
					log.Println("failed to prepare statement for raw packet insert", err)
					continue
				}
				for _, b := range buffer {
					_, err := stmt.Exec(b)
					if err != nil {
						log.Println("failed to insert raw packet data into duckdb", err)
					}
				}
				stmt.Close()
				err = tx.Commit()
				if err != nil {
					mu.Unlock()
					log.Println("failed to commit transaction for raw packet insert", err)
				}
				return
			default:
			}

			mu.Lock()
			if len(buffer) > 3000 {
				tx, err := dc.client.BeginTx(dc.ctx, nil)
				if err != nil {
					mu.Unlock()
					log.Println("failed to begin transaction for raw packet insert", err)
					continue
				}
				stmt, err := tx.PrepareContext(dc.ctx, "INSERT INTO raw (data) VALUES (?)")
				if err != nil {
					tx.Rollback()
					mu.Unlock()
					log.Println("failed to prepare statement for raw packet insert", err)
					continue
				}
				for _, b := range buffer {
					_, err := stmt.Exec(b)
					if err != nil {
						log.Println("failed to insert raw packet data into duckdb", err)
						// continue inserting others
					}
				}
				stmt.Close()
				err = tx.Commit()

				if err != nil {
					mu.Unlock()
					log.Println("failed to commit transaction for raw packet insert", err)
					continue
				}

				buffer = buffer[:0]
			}
			mu.Unlock()
		}
	}()

	for {
		select {
		case <-dc.ctx.Done():
			wg.Wait()
			return
		case packetData := <-dc.packetChan:
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
		case bytes := <-dc.rawPacketChan:
			mu.Lock()
			buffer = append(buffer, bytes)
			mu.Unlock()
		}
	}
}
