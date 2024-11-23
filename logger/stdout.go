package logger

import (
	"context"
	"fmt"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"sync"
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

type event struct {
	name       string
	occurredAt time.Duration
}

type participant struct {
	isMe        bool
	driverId    int
	name        string
	team        string
	carNumber   int
	grid        int
	lastLaptime time.Duration
	lastLap     packet.Lap
	gapToLeader time.Duration // seconds
	gapToFront  time.Duration // seconds
	gapToMe     time.Duration // seconds
}

type StdoutClient struct {
	inputChannel     chan packet.Data
	ctx              context.Context
	wg               *sync.WaitGroup
	events           []event
	printTicker      *time.Ticker
	printed          bool
	startedAt        time.Time
	currentLapNumber int
	isProcessing     bool
	lapTimes         map[int]*lapTime
	lastHeader       packet.Header
	lastTelemetry    packet.CarTelemetry
	lastMotion       packet.CarMotionData
	lastSession      packet.SessionData
	lastLapData      packet.Lap
	participants     []participant
	grids            []*participant
}

func (c *StdoutClient) Writer(ctx context.Context, wg *sync.WaitGroup) (channel chan<- packet.Data, cancel context.CancelFunc, err error) {
	c.inputChannel = make(chan packet.Data, 100)
	c.ctx, cancel = context.WithCancel(ctx)
	c.wg = wg

	return c.inputChannel, func() {
		cancel()
		c.printTicker.Stop()
	}, nil
}

func (c *StdoutClient) printData() {
	for range c.printTicker.C {
		if !c.printed {
			c.printed = true
		} else {
			fmt.Printf("\033[26A\033[K")
			//fmt.Printf("\u001B[1A\u001B[K\033[1A\033[K\u001B[1A\u001B[K\u001B[1A\u001B[K")

			//for range c.grids {
			//	fmt.Printf("\u001B[1A\u001B[K")
			//}

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

		for i, p := range c.grids {
			if p == nil {
				continue
			}
			text += fmt.Sprintf("P%d. %s(%d) | leader %.03f | interval %.03f | player %.03f\n",
				i+1,
				p.name,
				p.driverId,
				p.gapToLeader.Seconds(),
				p.gapToFront.Seconds(),
				p.gapToMe.Seconds(),
			)
		}

		fmt.Print(text)
	}
}

func (c *StdoutClient) Run() {
	defer c.wg.Done()

	go c.printData()

	c.printTicker = time.NewTicker(1000 / 60 * time.Millisecond)

	for data := range c.inputChannel {
		switch v := data.(type) {
		case packet.EventData:
			switch evt := v.Event.(type) {
			case packet.SessionStarted:
				c.events = append(c.events, event{
					name:       "Session STARTED",
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
				c.startedAt = time.Now()
			case packet.SessionEnded:
				c.events = append(c.events, event{
					name:       "Session ENDED",
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.FastestLap:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("Fastest LAP %s", c.participants[evt.VehicleIndex].name),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.Retirement:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("%s Retired", c.participants[evt.VehicleIndex].name),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.DRSEnabled:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("DRS ENABLED"),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.DRSDisabled:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("DRS DISABLED"),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.Penalty:
				c.events = append(c.events, event{
					name: fmt.Sprintf("Penalty %s to %s due to %s",
						packet.PenaltiesById[evt.PenaltyType],
						c.participants[evt.VehicleIndex].name,
						packet.InfringementsById[evt.InfringementType]),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.StartLights:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("Start Lights - %d", int(evt.NumberLights)),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.LightsOut:
				c.events = append(c.events, event{
					name:       "Lights Out",
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.RaceWinner:
				c.events = append(c.events, event{
					name:       fmt.Sprintf("%s Win the race", c.participants[evt.VehicleIndex].name),
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			case packet.Flashback:
				c.events = append(c.events, event{
					name:       "Flashback",
					occurredAt: getUnixTime(c.lastHeader.SessionTime),
				})
			}
		case packet.ParticipantData:
			for i, p := range v.Participants {
				c.participants[i].isMe = i == int(v.Header.PlayerCarIndex)
				c.participants[i].driverId = int(p.DriverId)
				if p.DriverId == 255 || p.DriverId == 0 {
					c.participants[i].name = p.GetName()
				} else {
					c.participants[i].name = packet.DriverNameById[p.DriverId]
				}
				c.participants[i].team = packet.TeamNameById[p.TeamId]
				c.participants[i].carNumber = int(p.RaceNumber)

				c.grids[i] = &c.participants[i]
			}
		case packet.MotionData:
			c.lastMotion = v.Player()
		case packet.SessionData:
			c.lastSession = v
		case packet.CarTelemetryData:
			c.lastTelemetry = v.Player()
		case packet.LapData:
			c.lastLapData = v.Player()
			leader := v.Leader()

			for index, lap := range v.DriverLaps {
				c.participants[index].lastLaptime = time.Duration(lap.LatestLapTime) * time.Second
				c.participants[index].lastLap = lap
				c.participants[index].grid = int(lap.GridPosition)
				c.grids[index] = &c.participants[index]

				if lap.GridPosition != 1 {
					c.participants[index].gapToLeader = getGapBetween(leader, lap)
				}
				if index != int(v.Header.PlayerCarIndex) {
					c.participants[index].gapToMe = getGapBetween(v.Player(), lap)
				}
			}

			for i := 0; i < len(c.grids)-1; i++ {
				//p := c.grids[i]
				c.grids[i].gapToFront = getGapBetween(c.grids[i].lastLap, c.grids[i+1].lastLap)
				//c.grids[i] = p
			}

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
					lt.totalMs = time.Duration(c.lastLapData.LatestLapTime) * time.Second
					lt.s3Ms = lt.totalMs - lt.s2Ms
				}
			}

			lt := c.lapTimes[lapNumber]
			if lt.s1Ms == 0 && c.lastLapData.Sector1Time != 0 {
				lt.s1Ms = time.Duration(c.lastLapData.Sector1Time) * time.Second
			}
			if lt.s2Ms == 0 && c.lastLapData.Sector2Time != 0 {
				lt.s2Ms = time.Duration(c.lastLapData.Sector2Time) * time.Second
			}

			lt.totalMs = time.Duration(c.lastLapData.CurrentLapTime) * time.Second

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
		printTicker:      time.NewTicker(1000 / 5 * time.Millisecond),
		participants:     make([]participant, 22),
		grids:            make([]*participant, 22),
	}
	return
}

// a is the target. b is the reference
// if a is behind b, the result will be negative
func getGapBetween(a, b packet.Lap) time.Duration {
	if a.CurrentLapNumber == 0 || b.CurrentLapNumber == 0 {
		return 0
	}

	distance := distanceBetween(a, b)
	velocityRel := float64((a.LapDistance / float32(a.CurrentLapTime)) - (b.LapDistance / float32(b.CurrentLapTime)))

	return time.Duration(distance / velocityRel * 1000)
}

func distanceBetween(a, b packet.Lap) float64 {
	if a.CurrentLapNumber == 0 || b.CurrentLapNumber == 0 {
		return 0
	}

	return float64(a.TotalDistance - b.TotalDistance)
}
