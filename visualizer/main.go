package main

import (
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"image"
	"image/png"
)

type St struct {
	LapDuration time.Duration
	Steer       float32
	Throttle    float32
	Break       float32
	Gear        int
	EngineRpm   int
	Speed       int
}

type Lt struct {
	SectorDurations  []time.Duration
	TotalLapDuration time.Duration
}

const (
	windowsFontPath = "C:\\Windows\\Fonts"
	LinuxFontPath   = "/Library/Fonts"
)

var (
	sectorColors []color.RGBA
)

var storagePath = path.Join(os.TempDir(), "f1")

func init() {
	s1 := colornames.Indianred
	s1.A = 10

	s2 := colornames.Deepskyblue
	s2.A = 15

	s3 := colornames.Yellow
	s3.A = 15
	sectorColors = append(sectorColors, s1, s2, s3)

	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	storagePath = path.Join(storagePath, "2021-07-17/114711")
	err := os.Mkdir(path.Join(storagePath, "image"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	for lap := 1; lap <= 8; lap++ {
		drawImage(storagePath, lap)
	}
}

func loadCarTelemetryData(p string, lap int) (sts []St, size int, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), strconv.Itoa(int(packet.CarTelemetryDataId))))
	if err != nil {
		return
	}

	var start time.Time
	var _t f1.SimplifiedTelemetry
	err = packet.ParsePacket(b, &_t)
	if err != nil {
		return
	}

	start = time.Unix(0, int64(_t.Timestamp))
	size = packet.Sizeof(reflect.ValueOf(f1.SimplifiedTelemetry{}))

	for i := 0; i < len(b); i += size {
		var t f1.SimplifiedTelemetry
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		timestamp := time.Unix(0, int64(t.Timestamp))
		lapTime := timestamp.Sub(start)
		sts = append(sts, St{
			LapDuration: lapTime,
			Steer:       t.Steer,
			Throttle:    t.Throttle,
			Break:       t.Break,
			Gear:        int(t.Gear),
			EngineRpm:   int(t.EngineRPM),
			Speed:       int(t.Speed),
		})
	}

	return
}

type Mt struct {
	LapDuration         time.Duration
	GForce              f1.Float3d
	WheelSlip           f1.FloatWheels
	WheelSpeed          f1.FloatWheels
	AngularVelocity     f1.Float3d
	AngularAcceleration f1.Float3d
	Position            f1.Float3d
	Heading             f1.Float3d
}

func loadMotionData(p string, lap int) (mts []Mt, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), strconv.Itoa(int(packet.MotionDataId))))
	if err != nil {
		return
	}

	var _t f1.PlayerMotionData
	err = packet.ParsePacket(b, &_t)
	if err != nil {
		return
	}

	size := packet.Sizeof(reflect.ValueOf(f1.PlayerMotionData{}))

	for i := 0; i < len(b); i += size {
		var t f1.PlayerMotionData
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		timestamp := time.Duration(t.Timestamp)
		mts = append(mts, Mt{
			LapDuration:         timestamp,
			GForce:              t.GForce,
			WheelSlip:           t.WheelSlip,
			WheelSpeed:          t.WheelSpeed,
			AngularAcceleration: t.AngularAcceleration,
			AngularVelocity:     t.AngularVelocity,
			Position:            t.Position,
			Heading:             t.Heading,
		})
	}

	return
}

func loadLapTelemetries(p string, lap int) (lt Lt, err error) {
	b, err := ioutil.ReadFile(path.Join(p, strconv.Itoa(lap), strconv.Itoa(int(packet.LapDataId))))
	if err != nil {
		return
	}

	//var start time.Time
	var _t f1.SimplifiedLap
	err = packet.ParsePacket(b, &_t)
	if err != nil {
		return
	}

	//start = time.Unix(0, int64(_t.Timestamp))
	size := packet.Sizeof(reflect.ValueOf(f1.SimplifiedLap{}))

	sector1 := time.Duration(0)
	sector2 := time.Duration(0)
	var lapTime time.Duration
	for i := 0; i < len(b); i += size {
		var t f1.SimplifiedLap
		err = packet.ParsePacket(b[i:i+size], &t)
		if err != nil {
			return
		}

		lapNumber := int(t.CurrentLapNumber)
		if lapNumber != lap {
			log.Printf("previous lap %#v", t)
			continue
		}
		if t.LapDistance < 0 {
			continue
		}

		if t.Sector1Time != 0 && sector1.Milliseconds() == 0 {
			sector1 = time.Duration(t.Sector1Time) * time.Millisecond
		}
		if t.Sector2Time != 0 && sector2.Milliseconds() == 0 {
			sector2 = time.Duration(t.Sector2Time) * time.Millisecond
		}

		if t.CurrentLapTime != 0 {
			lapTime = time.Duration(t.CurrentLapTime) * time.Millisecond
		}

		// TODO: driver status로 pitstop start, end 알아내기
		// t.DriverStatus
	}

	lt = Lt{
		SectorDurations: []time.Duration{
			sector1, sector2, lapTime - sector2 - sector1,
		},
		TotalLapDuration: lapTime,
	}
	return
}

func drawImage(p string, lap int) {
	carTelemetries, stSize, err := loadCarTelemetryData(p, lap)
	if err != nil {
		panic(err)
	}

	lapTelemetries, err := loadLapTelemetries(p, lap)
	if err != nil {
		panic(err)
	}

	motionTelemetries, err := loadMotionData(p, lap)
	if err != nil {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, len(carTelemetries), 1500))

	ctx := gg.NewContextForRGBA(img)
	err = ctx.LoadFontFace(path.Join(windowsFontPath, "Arial.ttf"), 25)
	if err != nil {
		panic(err)
	}

	ctx.SetColor(colornames.White)

	width := img.Rect.Max.X
	height := img.Rect.Max.Y

	totalLapTime := carTelemetries[len(carTelemetries)-1].LapDuration - carTelemetries[0].LapDuration
	log.Printf("lap %d, size:%fkb, lapTime: %s", lap, float32(stSize)/1024, totalLapTime)
	log.Println("     sector1", lapTelemetries.SectorDurations[0], "sector2", lapTelemetries.SectorDurations[1], "sector3", lapTelemetries.SectorDurations[2], "total", lapTelemetries.TotalLapDuration)

	// TODO: drawGrid
	gap := float32(img.Rect.Max.X) / float32(totalLapTime.Seconds())

	previousSectorDuration := time.Duration(0)
	for i, sectorDuration := range lapTelemetries.SectorDurations {
		ctx.SetColor(sectorColors[i])
		ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, sectorDuration.Seconds()*float64(gap), float64(height))
		ctx.Fill()
		previousSectorDuration += sectorDuration
	}

	count := time.Duration(0)
	for x := 0; x < width; x += int(gap * 5) {
		ctx.SetColor(color.Black)
		ctx.DrawString(count.String(), float64(x), 25)

		ctx.SetColor(colornames.Aqua)
		ctx.DrawLine(float64(x), 0, float64(x), float64(height))
		ctx.Stroke()
		count += time.Second * 5
	}

	draw.Draw(img, image.Rect(0, 0, width, 500), drawStat(carTelemetries, StatSteering, 500, NoVerticalGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 0, width, 500), drawStat(carTelemetries, StatSpeed, 500, NoGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 550, width, 650), drawStat(carTelemetries, StatThrottle, 100, NoGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 550, width, 650), drawStat(carTelemetries, StatBreak, 100, NoGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 700, width, 700+20*8), drawStat(carTelemetries, StatGear, 20*8, NoGrid), image.Point{}, draw.Over)

	draw.Draw(img, image.Rect(0, 900, width, height), drawStat(carTelemetries, StatEngineRPM, 150, NoVerticalGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Engine RPM", 1, 920)

	draw.Draw(img, image.Rect(0, 1100, width, 1150), drawMotion(motionTelemetries, MotionWheelSpeedFrontBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	err = ctx.LoadFontFace(windowsFontPath+"/Arial.ttf", 15)
	if err != nil {
		panic(err)
	}
	ctx.SetColor(colornames.Black)
	ctx.DrawString("FL", 1, 1122)
	ctx.DrawString("FR", 1, 1140)

	draw.Draw(img, image.Rect(0, 1150, width, 1200), drawMotion(motionTelemetries, MotionWheelSpeedRearBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("RL", 1, 1172)
	ctx.DrawString("RR", 1, 1190)

	draw.Draw(img, image.Rect(0, 1200, width, 1250), drawMotion(motionTelemetries, MotionWheelSpeedLeftBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("FL", 1, 1222)
	ctx.DrawString("RL", 1, 1240)

	draw.Draw(img, image.Rect(0, 1250, width, 1300), drawMotion(motionTelemetries, MotionWheelSpeedRightBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("FR", 1, 1272)
	ctx.DrawString("RR", 1, 1290)

	f, err := os.Create(path.Join(p, "image", fmt.Sprintf("test-%d.png", lap)))
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

type StatType int
type MotionType int

const (
	StatSteering StatType = iota
	StatThrottle
	StatBreak
	StatGear
	StatSpeed
	StatEngineRPM
)

const (
	MotionWheelSpeedFrontBias MotionType = iota
	MotionWheelSpeedRearBias
	MotionWheelSpeedLeftBias
	MotionWheelSpeedRightBias
)

type Option int

const (
	NoGrid Option = iota
	NoHorizontalGrid
	NoVerticalGrid
	DrawMiddleGrid
)

func drawStat(sts []St, typ StatType, height int, options ...Option) image.Image {
	drawVerticalGrid := true
	drawHorizontalGrid := true
	drawMiddleGrid := false

	for _, o := range options {
		switch o {
		case NoVerticalGrid:
			drawVerticalGrid = false
		case NoHorizontalGrid:
			drawHorizontalGrid = false
		case NoGrid:
			drawHorizontalGrid = false
			drawVerticalGrid = false
		case DrawMiddleGrid:
			drawMiddleGrid = true
		}
	}

	width := len(sts)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	totalLapTime := sts[width-1].LapDuration - sts[0].LapDuration
	ctx := gg.NewContextForRGBA(img)

	if drawMiddleGrid {
		ctx.SetColor(colornames.Lightgray)
		ctx.DrawLine(0, float64(height/2), float64(width), float64(height/2))
		ctx.Stroke()
	}

	if drawVerticalGrid {
		gap := float32(totalLapTime.Seconds()) / float32(img.Rect.Max.X)
		gridGap := int(time.Second * 5 / (time.Millisecond * time.Duration(gap*1000)))
		ctx.SetColor(colornames.Aqua)
		for x := 0; x < width; x += gridGap {
			ctx.DrawLine(float64(x), 0, float64(x), float64(height))
		}
		ctx.Stroke()
	}

	if drawHorizontalGrid {
		ctx.SetColor(colornames.Gray)
		for y := 0; y <= height; y += height / 5 {
			ctx.DrawLine(0, float64(y), float64(width), float64(y))
		}
		ctx.Stroke()
	}

	verticalMargin := 1

	height = height - verticalMargin*2

	hg := height / 8
	hh := height / 2
	hf := float32(height)
	hhf := float32(hh)
	previousY := float64(0)

	for x, t := range sts {
		var c color.RGBA
		var y int
		switch typ {
		case StatSteering:
			y = int(t.Steer*hhf + hhf)
			c = colornames.Black
		case StatSpeed:
			y = -t.Speed + height + verticalMargin
			c = colornames.Green
		case StatBreak:
			y = int(-t.Break*hf) + height + verticalMargin
			c = colornames.Orangered
		case StatThrottle:
			y = int(-t.Throttle*hf) + height + verticalMargin
			c = colornames.Blue
		case StatGear:
			y = -t.Gear*hg + height + verticalMargin
			c = colornames.Pink
		case StatEngineRPM:
			y = int(float32(-t.EngineRpm)/13000*hf + hf)
			c = colornames.Black
		}

		//img.Set(x, y, c)

		currentY := float64(y)
		ctx.SetColor(c)
		if x == 0 {
			ctx.DrawPoint(0, currentY, 1)
		} else {
			ctx.DrawLine(float64(x-1), previousY, float64(x), currentY)
		}
		ctx.Stroke()
		previousY = currentY
	}

	return img
}

func drawMotion(mts []Mt, typ MotionType, height int, options ...Option) image.Image {
	drawVerticalGrid := true
	drawHorizontalGrid := true
	drawMiddleGrid := false
	for _, o := range options {
		switch o {
		case NoVerticalGrid:
			drawVerticalGrid = false
		case NoHorizontalGrid:
			drawHorizontalGrid = false
		case NoGrid:
			drawHorizontalGrid = false
			drawVerticalGrid = false
		case DrawMiddleGrid:
			drawMiddleGrid = true
		}
	}

	width := len(mts)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	totalLapTime := mts[width-1].LapDuration - mts[0].LapDuration
	ctx := gg.NewContextForRGBA(img)

	if drawVerticalGrid {
		gap := float32(totalLapTime.Seconds()) / float32(img.Rect.Max.X)
		gridGap := int(time.Second * 5 / (time.Millisecond * time.Duration(gap*1000)))
		ctx.SetColor(colornames.Aqua)
		for x := 0; x < width; x += gridGap {
			ctx.DrawLine(float64(x), 0, float64(x), float64(height))
		}
		ctx.Stroke()
	}

	if drawHorizontalGrid {
		ctx.SetColor(colornames.Gray)
		for y := 0; y <= height; y += height / 5 {
			ctx.DrawLine(0, float64(y), float64(width), float64(y))
		}
		ctx.Stroke()
	}

	if drawMiddleGrid {
		ctx.SetColor(colornames.Lightgray)
		ctx.DrawLine(0, float64(height/2), float64(width), float64(height/2))
		ctx.Stroke()
	}

	verticalMargin := 1

	height = height - verticalMargin*2

	hh := height / 2
	//hf := float32(height)
	hhf := float32(hh)
	previousY := float64(0)

	for x, m := range mts {
		var c color.RGBA
		var y int
		switch typ {
		case MotionWheelSpeedFrontBias:
			c = colornames.Blueviolet

			frontBias := float32(0)
			if m.WheelSpeed.FR != 0 {
				frontBias = m.WheelSpeed.FL/m.WheelSpeed.FR - 1
			}
			y = int(-frontBias*hhf)*10 + hh
		case MotionWheelSpeedRearBias:
			c = colornames.Blueviolet

			rearBias := float32(0)
			if m.WheelSpeed.RR != 0 {
				rearBias = m.WheelSpeed.RL/m.WheelSpeed.RR - 1
			}
			y = int(-rearBias*hhf)*10 + hh
		case MotionWheelSpeedLeftBias:
			c = colornames.Blueviolet

			leftBias := float32(0)
			if m.WheelSpeed.RL != 0 {
				leftBias = m.WheelSpeed.FL/m.WheelSpeed.RL - 1
			}
			y = int(-leftBias*hhf)*10 + hh
		case MotionWheelSpeedRightBias:
			c = colornames.Blueviolet

			rightBias := float32(0)
			if m.WheelSpeed.RR != 0 {
				rightBias = m.WheelSpeed.FR/m.WheelSpeed.RR - 1
			}
			y = int(-rightBias*hhf)*10 + hh
		}

		currentY := float64(y)
		if x == 0 {
			ctx.DrawPoint(0, currentY, 1)
			ctx.SetColor(c)
		} else {
			ctx.DrawLine(float64(x-1), previousY, float64(x), currentY)
		}
		previousY = currentY
	}
	ctx.Stroke()

	return img
}
