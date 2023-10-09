package logger

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"time"
)

var _ packet.Logger = (*StdoutClient)(nil)

type lapTime struct {
	number          int
	totalMs         time.Duration
	s1Ms            time.Duration
	s2Ms            time.Duration
	s3Ms            time.Duration
	timesBy10Meters map[int]time.Duration
}

type StdoutClient struct {
	inputChannel     chan packet.PacketData
	ctx              context.Context
	printTicker      *time.Ticker
	isRunning        bool
	printed          bool
	started          bool
	startedAt        time.Time
	currentLapNumber int
	lapTimes         map[int]*lapTime
	lastHeader       packet.Header
	lastTelemetry    packet.CarTelemetry
	lastMotion       packet.CarMotionData
	lastSession      packet.SessionData
	lastLapData      packet.DriverLap
}

func (c *StdoutClient) Writer(ctx context.Context) (channel chan<- packet.PacketData, cancel context.CancelFunc, err error) {
	c.inputChannel = make(chan packet.PacketData, 100)
	c.ctx, cancel = context.WithCancel(ctx)

	newCancel := func() {
		cancel()
		c.started = false
		c.isRunning = false
	}
	return c.inputChannel, newCancel, nil
}

func (c *StdoutClient) printData() {
	for range c.printTicker.C {
		if !c.isRunning {
			c.printTicker.Stop()
			break
		}

		if !c.printed {
			c.printed = true
		} else {
			fmt.Printf("\u001B[1A\u001B[K\033[1A\033[K\u001B[1A\u001B[K\u001B[1A\u001B[K")
			//fmt.Printf("\u001B[1A\u001B[K")
			//fmt.Printf("\u001B[1A\u001B[K")
			//fmt.Printf("\u001B[1A\u001B[K")
		}
		now := c.startedAt.Add(getUnixTime(c.lastHeader.SessionTime))
		ld := c.lastLapData
		hd := c.lastHeader
		td := c.lastTelemetry
		sd := c.lastSession

		currentDistance := int(ld.LapDistance) / 10 * 10

		var lastS1, lastS2, lastS3, lastLapTime string
		var lastDistanceDuration time.Duration
		if lastLt, ok := c.lapTimes[c.currentLapNumber-1]; ok {
			lastS1 = lastLt.s1Ms.String()
			lastS2 = lastLt.s2Ms.String()
			lastS3 = lastLt.s3Ms.String()
			lastLapTime = lastLt.totalMs.String()
			if v, ok := lastLt.timesBy10Meters[currentDistance]; ok {
				lastDistanceDuration = v
			}
		}

		var currS1, currS2, currS3, currLapTime string
		var currDistanceDuration time.Duration
		if currLt, ok := c.lapTimes[c.currentLapNumber]; ok {
			currS1 = currLt.s1Ms.String()
			currS2 = currLt.s2Ms.String()
			currS3 = currLt.s3Ms.String()
			currLapTime = currLt.totalMs.String()
			currDistanceDuration = currLt.totalMs
		}

		delta := " - "
		if currDistanceDuration != 0 && lastDistanceDuration != 0 {
			delta = fmt.Sprintf("%.3f", (currDistanceDuration - lastDistanceDuration).Seconds())
		}

		currDistance := float64(ld.LapDistance) / 1000
		trackLength := float64(sd.TrackLength) / 1000
		text := fmt.Sprintf("[%s]-FR[%d]| LAP: %d\n", now.Format("2006-01-02 15:04:05"), hd.FrameIdentifier, ld.CurrentLapNumber)
		text += fmt.Sprintf("TH: %.2f, BR: %.2f, GR: %d, ENG: %dRPM, DRS:%b\n", td.Throttle, td.Break, td.Gear, td.EngineRPM, td.DRS)
		text += fmt.Sprintf("[%s] - %.2f / %.2f (KM), %s\t(%s)\t[ %s(%s) / %s(%s) / %s(%s) ]\n", delta, currDistance, trackLength, currLapTime, lastLapTime, currS1, lastS1, currS2, lastS2, currS3, lastS3)

		fmt.Print(text)
	}
}

func (c *StdoutClient) Run() {
	c.isRunning = true
	c.printTicker = time.NewTicker(1000 / 60 * time.Millisecond)
	//go c.printData()

	for data := range c.inputChannel {
		switch v := data.(type) {
		case packet.ParticipantData:

		case packet.EventData:
			switch v.Event.(type) {
			case packet.SessionStarted:
				fmt.Println("Session STARTED")
			case packet.SessionEnded:
				fmt.Println("Session Ended")
			}
		case packet.MotionData:
			c.lastMotion = v.Player()
		case packet.SessionData:
			c.lastSession = v
		case packet.CarTelemetryData:
			c.lastTelemetry = v.Player()
		case packet.LapData:
			c.lastLapData = v.Player()

			lapNumber := int(c.lastLapData.CurrentLapNumber)
			if c.currentLapNumber != lapNumber {
				c.currentLapNumber = lapNumber
				c.lapTimes[lapNumber] = &lapTime{
					number:          lapNumber,
					totalMs:         0,
					s1Ms:            0,
					s2Ms:            0,
					s3Ms:            0,
					timesBy10Meters: make(map[int]time.Duration),
				}

				if lt, ok := c.lapTimes[lapNumber-1]; ok {
					lt.totalMs = time.Duration(c.lastLapData.LatestLapTime) * time.Millisecond
					lt.s3Ms = lt.totalMs - lt.s2Ms
				}
			}

			lt := c.lapTimes[lapNumber]
			if lt.s1Ms == 0 && c.lastLapData.Sector1Time != 0 {
				lt.s1Ms = time.Duration(c.lastLapData.Sector1Time) * time.Millisecond
			}
			if lt.s2Ms == 0 && c.lastLapData.Sector2Time != 0 {
				lt.s2Ms = time.Duration(c.lastLapData.Sector2Time) * time.Millisecond
			}

			lt.totalMs = time.Duration(c.lastLapData.CurrentLapTime) * time.Millisecond

			currentLapDistance := int(c.lastLapData.LapDistance)
			if _, ok := lt.timesBy10Meters[currentLapDistance/10*10]; !ok {
				lt.timesBy10Meters[currentLapDistance/10*10] = lt.totalMs
			}
		}
	}
}

func getUnixTime(sessionTime float32) time.Duration {
	return time.Duration(sessionTime*1000) * time.Millisecond
}

func NewStdoutClient() (c *StdoutClient) {
	c = &StdoutClient{
		currentLapNumber: 0,
		lapTimes:         make(map[int]*lapTime),
	}
	return
}
