package main

import (
	"bytes" // 바이트 슬라이스 처리를 위해 추가
	"database/sql"
	"encoding/binary" // binary.Read를 위해 추가
	"fmt"
	// "github.com/ariyn/F1-2021-game-udp/logger" // Assuming this is no longer needed for Redis
	"github.com/ariyn/F1-2021-game-udp/packet" // 메인 패킷 패키지
	_ "github.com/marcboeker/go-duckdb"         // DuckDB driver
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// type LapData struct {
// 	Number    int
// 	startedAt int64
// 	EndedAt   int64
// }

// var lapHistory = make(map[int]LapData)
// var currentLapNumber = -1

type TelemetryData struct {
	SessionTime     float64
	FrameIdentifier int64
	LapDistance     float64
	LapNumber       int
	Throttle        float64
	Brk             float64
	Gear            int
	Rpm             int
	Drs             int
	Speed           int
	Delta           float64
	WorldX          float64
	WorldY          float64
	WorldZ          float64
	LapDeltaTime    float64
}

var data = TelemetryData{} // This will be populated from DuckDB later

// func read(redisLogger *logger.RedisClient) {
// 	ticker := time.NewTicker(time.Millisecond * 100)
// 	for range ticker.C {
// 		var errorOccurred bool
// 		redisLogger.GetCurrentSessionUid()

// 		lapDistance, err := redisLogger.GetTs("lapDistance")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		lapNumber, err := redisLogger.GetTs("lapNumber")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		// if currentLapNumber != int(lapNumber.Value) {
// 		// 	lapHistory[int(lapNumber.Value)] = LapData{
// 		// 		Number:    int(lapNumber.Value),
// 		// 		startedAt: lapNumber.Timestamp,
// 		// 	}

// 		// 	if v, ok := lapHistory[int(lapNumber.Value)-1]; ok {
// 		// 		v.EndedAt = lapNumber.Timestamp
// 		// 		lapHistory[int(lapNumber.Value)-1] = v
// 		// 	}
// 		// }

// 		frameIdentifier, err := redisLogger.GetTs("frameIdentifier")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		throttle, err := redisLogger.GetTs("throttle")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		brk, err := redisLogger.GetTs("break")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		gear, err := redisLogger.GetTs("gear")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		rpm, err := redisLogger.GetTs("rpm")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		drs, err := redisLogger.GetTs("drs")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		speed, err := redisLogger.GetTs("speed")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		worldPositionX, err := redisLogger.GetTs("worldPositionX")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}
// 		worldPositionY, err := redisLogger.GetTs("worldPositionY")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}
// 		worldPositionZ, err := redisLogger.GetTs("worldPositionZ")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		lapDeltaTime, err := redisLogger.GetTs("lapDeltaTime")
// 		if err != nil {
// 			errorOccurred = true
// 			//log.Println(err)
// 		}

// 		if errorOccurred {
// 			continue
// 		}

// 		data.FrameIdentifier = int64(frameIdentifier.Value)
// 		data.LapNumber = int(int64(lapNumber.Value))
// 		data.LapDistance = lapDistance.Value
// 		data.Throttle = throttle.Value
// 		data.Brk = brk.Value
// 		data.Gear = int(int64(gear.Value))
// 		data.Rpm = int(int64(rpm.Value))
// 		data.Drs = int(int64(drs.Value))
// 		data.Speed = int(int64(speed.Value))
// 		data.WorldX = worldPositionX.Value
// 		data.WorldY = worldPositionY.Value
// 		data.WorldZ = worldPositionZ.Value
// 		data.LapDeltaTime = lapDeltaTime.Value
// 	}
// }

// func readSpecificLap(redisLogger *logger.RedisClient, lapNumber int) (data []TelemetryData, err error) {
// 	// redisLogger.GetCurrentSessionUid()

// 	// lapData, err := redisLogger.GetLapHistory(lapNumber)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return
// 	// }

// 	// startedAt := int(lapData.StartedAt)
// 	// endedAt := int(lapData.EndedAt)
// 	// log.Println(startedAt, endedAt)

// 	// dataByTimestamp := make(map[int64]TelemetryData)
// 	// lapDistances, _ := redisLogger.GetTsWithRange("lapDistance", startedAt, endedAt)
// 	// for _, ld := range lapDistances {
// 	// 	dataByTimestamp[ld.Timestamp] = TelemetryData{
// 	// 		LapDistance: ld.Value,
// 	// 	}
// 	// }

// 	// lapNumbers, _ := redisLogger.GetTsWithRange("lapNumber", startedAt, endedAt)
// 	// for _, d := range lapNumbers {
// 	// 	if v, ok := dataByTimestamp[d.Timestamp]; !ok {
// 	// 		dataByTimestamp[d.Timestamp] = TelemetryData{
// 	// 			LapNumber: int(int64(d.Value)),
// 	// 		}
// 	// 	} else {
// 	// 		v.LapNumber = int(int64(d.Value))
// 	// 		dataByTimestamp[d.Timestamp] = v
// 	// 	}
// 	// }

// 	// frameIdentifiers, _ := redisLogger.GetTsWithRange("frameIdentifier", startedAt, endedAt)
// 	// for _, fi := range frameIdentifiers {
// 	// 	if v, ok := dataByTimestamp[fi.Timestamp]; !ok {
// 	// 		dataByTimestamp[fi.Timestamp] = TelemetryData{
// 	// 			FrameIdentifier: int64(fi.Value),
// 	// 		}
// 	// 	} else {
// 	// 		v.FrameIdentifier = int64(fi.Value)
// 	// 		dataByTimestamp[fi.Timestamp] = v
// 	// 	}
// 	// }

// 	// //throttles, _ := redisLogger.GetTsWithRange("throttle", startedAt, endedAt)
// 	// //brks, _ := redisLogger.GetTsWithRange("break", startedAt, endedAt)
// 	// //gears, _ := redisLogger.GetTsWithRange("gear", startedAt, endedAt)
// 	// //rpms, _ := redisLogger.GetTsWithRange("rpm", startedAt, endedAt)
// 	// //drss, _ := redisLogger.GetTsWithRange("drs", startedAt, endedAt)
// 	// //speeds, _ := redisLogger.GetTsWithRange("speed", startedAt, endedAt)

// 	// worldPositionXs, _ := redisLogger.GetTsWithRange("worldPositionX", startedAt, endedAt)
// 	// for _, d := range worldPositionXs {
// 	// 	if v, ok := dataByTimestamp[d.Timestamp]; !ok {
// 	// 		dataByTimestamp[d.Timestamp] = TelemetryData{
// 	// 			WorldX: d.Value,
// 	// 		}
// 	// 	} else {
// 	// 		v.WorldX = d.Value
// 	// 		dataByTimestamp[d.Timestamp] = v
// 	// 	}
// 	// }

// 	// worldPositionYs, _ := redisLogger.GetTsWithRange("worldPositionY", startedAt, endedAt)
// 	// for _, d := range worldPositionYs {
// 	// 	if v, ok := dataByTimestamp[d.Timestamp]; !ok {
// 	// 		dataByTimestamp[d.Timestamp] = TelemetryData{
// 	// 			WorldY: d.Value,
// 	// 		}
// 	// 	} else {
// 	// 		v.WorldY = d.Value
// 	// 		dataByTimestamp[d.Timestamp] = v
// 	// 	}
// 	// }

// 	// worldPositionZs, _ := redisLogger.GetTsWithRange("worldPositionZ", startedAt, endedAt)
// 	// for _, d := range worldPositionZs {
// 	// 	if v, ok := dataByTimestamp[d.Timestamp]; !ok {
// 	// 		dataByTimestamp[d.Timestamp] = TelemetryData{
// 	// 			WorldZ: d.Value,
// 	// 		}
// 	// 	} else {
// 	// 		v.WorldZ = d.Value
// 	// 		dataByTimestamp[d.Timestamp] = v
// 	// 	}
// 	// }

// 	// lapDeltaTimes, _ := redisLogger.GetTsWithRange("lapDeltaTime", startedAt, endedAt)
// 	// for _, d := range lapDeltaTimes {
// 	// 	if v, ok := dataByTimestamp[d.Timestamp]; !ok {
// 	// 		dataByTimestamp[d.Timestamp] = TelemetryData{
// 	// 			LapDeltaTime: d.Value,
// 	// 		}
// 	// 	} else {
// 	// 		v.LapDeltaTime = d.Value
// 	// 		dataByTimestamp[d.Timestamp] = v
// 	// 	}
// 	// }

// 	// data = make([]TelemetryData, 0)
// 	// for timestamp, d := range dataByTimestamp {
// 	// 	d.SessionTime = float64(timestamp)
// 	// 	data = append(data, d)
// 	// }

// 	// sort.Slice(data, func(i, j int) bool {
// 	// 	return data[i].SessionTime < data[j].SessionTime
// 	// })
// 	log.Println("readSpecificLap will be reimplemented for DuckDB")
// 	return // Placeholder
// }

// getLapMotionTimeline fetches motion data for a specific lap number
func getLapMotionTimeline(lapNumber int) ([]TelemetryData, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// Query to get the player_car_index for the session.
	// This is a simplification. In a real scenario, you might get this from session data or a specific event.
	// For now, let's assume the primary player is often index 0, or we can try to infer it.
	// However, packet_header.player_car_index in the 'packets' table is more reliable per packet.

	rows, err := db.Query("SELECT data FROM packets WHERE packet_id = 0 ORDER BY frame_identifier ASC")
	if err != nil {
		return nil, fmt.Errorf("querying motion packets: %w", err)
	}
	defer rows.Close()

	var results []TelemetryData
	for rows.Next() {
		var rawData []byte
		if err := rows.Scan(&rawData); err != nil {
			log.Printf("Error scanning raw motion data: %v", err)
			continue
		}

		parsedPacket, err := packet.UnmarshalBinary(rawData)
		if err != nil {
			// Log and continue if a packet is corrupted or not unmarshalable
			// log.Printf("Error unmarshalling motion packet: %v", err)
			continue
		}

		if motionPacket, ok := parsedPacket.(*packet.PacketMotionData); ok {
			header := motionPacket.Header
			// Assuming we're interested in the data for the car that generated this packet.
			// If packets table has a specific player_car_index column for filtering, use that in SQL.
			// Otherwise, use the index from the packet's header.
			playerIdx := header.PlayerCarIndex
			if playerIdx >= uint8(len(motionPacket.CarMotionData)) {
				// log.Printf("Player index %d out of bounds for CarMotionData length %d", playerIdx, len(motionPacket.CarMotionData))
				continue
			}

			carData := motionPacket.CarMotionData[playerIdx]

			if int(carData.CurrentLapNum) == lapNumber {
				results = append(results, TelemetryData{
					SessionTime:     header.SessionTime,
					FrameIdentifier: int64(header.FrameIdentifier),
					LapDistance:     float64(carData.LapDistance),
					LapNumber:       int(carData.CurrentLapNum),
					WorldX:          float64(carData.WorldPositionX),
					WorldY:          float64(carData.WorldPositionY),
					WorldZ:          float64(carData.WorldPositionZ),
					// Throttle, Brk, Gear, Rpm, Drs, Speed are not in PacketMotionData
				})
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error over motion packet rows: %w", err)
	}
	return results, nil
}

// apiGetLapMotionTimelineHandler handles requests for lap timeline data
func apiGetLapMotionTimelineHandler(c echo.Context) error {
	lapNumberStr := c.Param("lapNumber")
	lapNumber, err := strconv.Atoi(lapNumberStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lap number"})
	}

	timelineData, err := getLapMotionTimeline(lapNumber)
	if err != nil {
		log.Printf("Error getting lap motion timeline for lap %d: %v", lapNumber, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve lap timeline data"})
	}
	if len(timelineData) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("No motion data found for lap %d", lapNumber)})
	}
	return c.JSON(http.StatusOK, timelineData)
}

var db *sql.DB // Global DB connection object

func main() {
	var err error
	// Connect to DuckDB in read-only mode
	db, err = sql.Open("duckdb", "f1_telemetry.db?access_mode=READ_ONLY")
	if err != nil {
		log.Fatalf("Failed to connect to DuckDB: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to DuckDB (read-only)")
	// redisLogger := logger.NewRedisClient(context.Background(), "localhost:6379", 0)
	// defer redisLogger.CloseListen()
	// go read(redisLogger) // This function is now commented out
	// fmt.Println("start Listening\n\nwaiting inputs") // This message might be misleading now

	e := echo.New()

	// Serve static files from the "frontend" directory
	e.Static("/frontend", "frontend")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// New API endpoint for specific lap motion timeline
	e.GET("/api/lap/:lapNumber/timeline", apiGetLapMotionTimelineHandler)

	e.GET("/", func(c echo.Context) error {
		lapQuery := c.QueryParam("lap")

		if lapQuery != "" {
			lapNumber, err := strconv.Atoi(lapQuery)
			if err != nil {
				log.Println("ERROR: lap query parameter is not a valid number:", lapQuery)
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lap number specified in query."})
			}
			// Fetch and return data for the specified lap
			timelineData, err := getLapMotionTimeline(lapNumber)
			if err != nil {
				log.Printf("Error getting lap motion timeline for lap %d from root: %v", lapNumber, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve lap timeline data."})
			}
			if len(timelineData) == 0 {
				return c.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("No motion data found for lap %d.", lapNumber)})
			}
			return c.JSON(http.StatusOK, timelineData)
		} else {
			// No lap specified, try to get the most recent lap's data.
			// This requires finding the max lap number first.
			var maxLapNum int
			// Query for the highest lap number in PacketMotionData
			// Note: This assumes CurrentLapNum is accurately recorded and sensible.
			// It might be better to get this from PacketLapData if available and more reliable.
			row := db.QueryRow("SELECT MAX(CAST(json_extract_path_text(packet.PacketMotionData.m_carMotionData[header.m_playerCarIndex].m_currentLapNum, '$') AS INTEGER)) FROM packets packet WHERE packet.packet_id = 0")
			// The above query is an example of how one might try to extract nested JSON data if it were stored as JSON.
			// However, our 'data' column is a blob. We need to parse it.
			// A simpler approach for now: query all lap numbers and find max in Go.
			// This is inefficient but will work for a demo.
			// For a more robust solution, PacketLapData (ID 2) should be the source of truth for lap counts.

			// Simplification: Default to lap 1 or a message if determining max lap is too complex here.
			// For this iteration, let's try to fetch data for a default lap (e.g., lap 1)
			// or return a message if that's not found.
			defaultLapToShow := 1 // Or fetch a list of available laps and pick the latest.
			timelineData, err := getLapMotionTimeline(defaultLapToShow)
			if err != nil || len(timelineData) == 0 {
                // Attempt to find the overall latest single piece of motion data as a fallback
                var latestMotion TelemetryData
                var latestFrameId uint32 = 0 // Store the frame identifier of the latest packet

                rows, err_q := db.Query("SELECT data FROM packets WHERE packet_id = 0 ORDER BY frame_identifier DESC LIMIT 1")
                if err_q == nil {
                    defer rows.Close()
                    if rows.Next() {
                        var rawData []byte
                        if err_scan := rows.Scan(&rawData); err_scan == nil {
                            parsedPacket, err_parse := packet.UnmarshalBinary(rawData)
                            if err_parse == nil {
                                if motionPacket, ok := parsedPacket.(*packet.PacketMotionData); ok {
                                    header := motionPacket.Header
                                    playerIdx := header.PlayerCarIndex
                                    if playerIdx < uint8(len(motionPacket.CarMotionData)) {
                                        carData := motionPacket.CarMotionData[playerIdx]
                                        // Check if this frame is later than what we've seen
                                        if header.FrameIdentifier > latestFrameId {
                                            latestFrameId = header.FrameIdentifier
                                            latestMotion = TelemetryData{
                                                SessionTime:     header.SessionTime,
                                                FrameIdentifier: int64(header.FrameIdentifier),
                                                LapDistance:     float64(carData.LapDistance),
                                                LapNumber:       int(carData.CurrentLapNum),
                                                WorldX:          float64(carData.WorldPositionX),
                                                WorldY:          float64(carData.WorldPositionY),
                                                WorldZ:          float64(carData.WorldPositionZ),
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                     if latestFrameId > 0 {
                        return c.JSON(http.StatusOK, []TelemetryData{latestMotion}) // Return as a slice
                    }
                }
                // If still no data, return a generic message or empty data for lap 1
				log.Printf("No motion data found for default lap %d and no latest packet available: %v", defaultLapToShow, err)
				return c.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("No motion data found for default lap %d. Try /api/lap/LAP_NUMBER/timeline", defaultLapToShow)})
			}
			return c.JSON(http.StatusOK, timelineData)
		}
	})
	e.Logger.Fatal(e.Start(":1323"))
}
