package main

import (
	"encoding/json"
	"flag"
	"fmt"
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"image"
	"image/png"
)

const (
	windowsFontPath = "C:\\Windows\\Fonts"
	LinuxFontPath   = "/Library/Fonts"
)

var (
	sectorColors []color.RGBA
	pitLaneColor color.RGBA
)

var storagePath = path.Join(os.TempDir(), "f1")
var argumentPath string

func init() {
	s1 := colornames.Indianred
	s1.A = 10

	s2 := colornames.Deepskyblue
	s2.A = 15

	s3 := colornames.Yellow
	s3.A = 15
	sectorColors = append(sectorColors, s1, s2, s3)

	pitLaneColor = colornames.Gray
	pitLaneColor.A = 120

	log.SetFlags(log.Llongfile | log.LstdFlags)

	flag.StringVar(&argumentPath, "path", "", "/path/to/f1/session/folder")
}

func getOutlinePositions(outline1Lap, outline2Lap int) {
	motionTelemetries, err := loadMotionData(storagePath, outline1Lap, 19)
	if err != nil {
		panic(err)
	}

	start := 0
	if outline1Lap == 1 {
		start = int(float64(len(motionTelemetries)/2) - float64(len(motionTelemetries))*0.2)
	}

	outline1Positions := make([]f1.Float3d, 0)
	for _, mt := range motionTelemetries[start:] {
		outline1Positions = append(outline1Positions, mt.Position)
	}

	motionTelemetries, err = loadMotionData(storagePath, outline2Lap, 19)
	if err != nil {
		panic(err)
	}

	start = 0

	outline2Positions := make([]f1.Float3d, 0)
	for _, mt := range motionTelemetries[start:] {
		outline2Positions = append(outline2Positions, mt.Position)
	}

	f, err := os.Create(path.Join(storagePath, "silverstone-outline.json"))
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal([][]f1.Float3d{outline1Positions, outline2Positions})
	if err != nil {
		panic(err)
	}

	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if argumentPath != "" {
		storagePath = argumentPath
	} else {
		storagePath = path.Join(storagePath, "2021-09-10/185929")
	}

	//getOutlinePositions(1, 3)

	//getPosition(storagePath, 1, 19, 5000)
	//getPosition(storagePath, 2, 19, 5000)
	//getPosition(storagePath, 3, 19, 5000)
	//getPosition(storagePath, 4, 19, 5000)
	//getPosition(storagePath, 5, 19, 6000)
	//getPosition(storagePath, 6, 19, 10000)
	//for i := 1; i < 11; i++ {
	//	getPosition(storagePath, i, 0, 5000)
	//}

	pos := getPositionData(storagePath, 2, 0)
	f, err := os.OpenFile(path.Join(storagePath, "/pos.json"), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(pos)
	if err != nil {
		panic(err)
	}
	f.Write(b)

	os.Exit(4)

	log.Println(storagePath)
	_, laps, err := getDriverData(storagePath)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(path.Join(storagePath, "telemetry-sheet"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	err = os.Mkdir(path.Join(storagePath, "motion"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	//var racingNumbers = []int{66}
	//for _, racingNumber := range racingNumbers {
	for driverIndex := 19; driverIndex <= 19; driverIndex++ {
		//driverIndex := getDriverIndex(drivers, racingNumber)
		//if driverIndex == -1 {
		//	panic("no such racing number driver")
		//}
		totalLap := laps[driverIndex]
		//driverName := packet.DriverNameByRacingNumber[racingNumber]
		driverName := "multiplayer-" + strconv.Itoa(driverIndex)
		err = os.Mkdir(path.Join(storagePath, "telemetry-sheet", driverName), 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}

		err = os.Mkdir(path.Join(storagePath, "motion", driverName), 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}

		for lap := 1; lap <= totalLap; lap++ {
			drawGeneralSheet(storagePath, lap, driverIndex, driverName)
			drawMotionSheet(storagePath, lap, driverIndex, driverName)
		}
	}
}

func drawOutline(name string, ctx *gg.Context) {
	f, err := os.Open(fmt.Sprintf("/home/ariyn/Documents/go/src/github.com/ariyn/F1-2021-game-udp/sample/%s-outline.json", name))
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var silverstone [][]f1.Float3d
	err = json.Unmarshal(b, &silverstone)
	if err != nil {
		panic(err)
	}

	ctx.SetLineWidth(1.0)
	ctx.SetColor(color.White)

	imageSize := float64(ctx.Image().Bounds().Max.X)
	min, _ := f1.Float3d{X: -1000, Z: -1000}, f1.Float3d{X: 1000, Z: 1000}
	width, height := 2000.0, 2000.0

	previousX, previousZ := (float64(silverstone[0][0].X)+math.Abs(float64(min.X)))/width*imageSize, (float64(silverstone[0][0].Z)+math.Abs(float64(min.Z)))/height*imageSize
	for _, p := range silverstone[0][1:] {
		currentX, currentZ := (float64(p.X)+math.Abs(float64(min.X)))/width*imageSize, (float64(p.Z)+math.Abs(float64(min.Z)))/height*imageSize
		ctx.DrawLine(previousX, previousZ, currentX, currentZ)
		previousX, previousZ = currentX, currentZ
	}
	ctx.Stroke()

	previousX, previousZ = (float64(silverstone[1][0].X)+math.Abs(float64(min.X)))/width*imageSize, (float64(silverstone[1][0].Z)+math.Abs(float64(min.Z)))/height*imageSize
	for _, p := range silverstone[1][1:] {
		currentX, currentZ := (float64(p.X)+math.Abs(float64(min.X)))/width*imageSize, (float64(p.Z)+math.Abs(float64(min.Z)))/height*imageSize
		ctx.DrawLine(previousX, previousZ, currentX, currentZ)
		previousX, previousZ = currentX, currentZ
	}
	ctx.Stroke()

}

func getPositionData(storagePath string, lap, driverIndex int) (pos []f1.Float3d) {
	motionTelemetries, err := loadMotionData(storagePath, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	for _, mt := range motionTelemetries {
		pos = append(pos, mt.Position)
	}

	return
}

func getPosition(storagePath string, lap, driverIndex int, imageSize int) (img image.Image) {
	img = image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	ctx := gg.NewContextForRGBA(img.(*image.RGBA))

	maxPosition := float64(imageSize - 1)
	ctx.SetLineWidth(3.0)
	ctx.DrawLine(1, 1, 1, maxPosition)
	ctx.DrawLine(1, maxPosition, maxPosition, maxPosition)
	ctx.DrawLine(maxPosition, maxPosition, maxPosition, 1)
	ctx.DrawLine(maxPosition, 1, 1, 1)
	ctx.DrawLine(1, 1, maxPosition, maxPosition)
	ctx.Stroke()

	//drawOutline("silverstone", ctx)

	motionTelemetries, err := loadMotionData(storagePath, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	carTelemetries, _, err := loadCarTelemetryData(storagePath, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	maxSpeed := float64(36)

	size := float32(1200)

	min, _ := f1.Float3d{X: -size, Z: -size}, f1.Float3d{X: size, Z: size}
	width, height := float64(size*2), float64(size*2)

	newMin := f1.Float3d{}
	posMinX, posMinZ := float32(0.0), float32(0.0)
	ctx.SetLineWidth(1.0)
	//ctx.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	dotSize := 3.0
	previousX, previousZ := 0.0, 0.0
	for index, mt := range motionTelemetries {
		spd := float64(carTelemetries[index].Speed)
		// t 0 -> red, t max -> green

		ctx.SetColor(color.RGBA{R: uint8(255 - (spd / maxSpeed * 25)), G: 255 + uint8(spd/maxSpeed*25), A: 255})

		currentX, currentZ := (float64(mt.Position.X)+math.Abs(float64(min.X)))/width*float64(imageSize), (float64(mt.Position.Z)+math.Abs(float64(min.Z)))/height*float64(imageSize)
		if index == 0 {
			ctx.DrawPoint(currentX, currentZ, dotSize)
		} else {
			ctx.DrawLine(previousX, previousZ, currentX, currentZ)
		}

		if mt.Position.X < posMinX {
			posMinX = mt.Position.X
		}

		if mt.Position.Z < posMinZ {
			posMinZ = mt.Position.Z
		}

		previousX, previousZ = currentX, currentZ
		ctx.Stroke()

		if mt.Position.X < newMin.X {
			newMin = mt.Position
		}
		if mt.Position.Z < newMin.Z {
			newMin = mt.Position
		}
	}

	log.Println(newMin, newMin)
	log.Println(posMinX, posMinZ)

	f, err := os.Create(path.Join(storagePath, fmt.Sprintf("d%d-l%d.png", driverIndex, lap)))
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

	return img
}

func getDriverIndex(ds []f1.Driver, racingNumber int) int {
	for index, d := range ds {
		if d.RaceNumber == racingNumber {
			return index
		}
	}
	return -1
}

func getDriverData(p string) (drivers []f1.Driver, totalLaps []int, err error) {
	b, err := ioutil.ReadFile(path.Join(p, "drivers.json"))
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &drivers)
	if err != nil {
		return
	}

	b, err = ioutil.ReadFile(path.Join(p, "driver-laps.json"))
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &totalLaps)
	if err != nil {
		return
	}

	return
}

func drawGeneralSheet(p string, lap int, driverIndex int, name string) {
	carTelemetries, stSize, err := loadCarTelemetryData(p, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	lapTelemetries, err := loadLapTelemetries(p, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	motionTelemetries, err := loadMotionData(p, lap, driverIndex)
	if err != nil {
		return
	}

	width := len(carTelemetries)
	height := 1500
	//img := image.NewRGBA(image.Rect(0, 0, width, height))

	img := drawSector(lapTelemetries, width, height).(*image.RGBA)
	ctx := gg.NewContextForRGBA(img)
	//err = ctx.LoadFontFace(path.Join(windowsFontPath, "Arial.ttf"), 25)
	//if err != nil {
	//	panic(err)
	//}

	//draw.Draw(img, image.Rect(0, 0, width, height), , image.Point{}, draw.Over)
	totalLapTime := lapTelemetries[len(lapTelemetries)-1].Timestamp - lapTelemetries[0].Timestamp
	log.Printf("lap %d, size:%fkb, lapTime: %s", lap, float32(stSize)/1024, totalLapTime)

	if totalLapTime == 0 {
		return
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
	//err = ctx.LoadFontFace(windowsFontPath+"/Arial.ttf", 15)
	//if err != nil {
	//	panic(err)
	//}
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

	f, err := os.Create(path.Join(p, "telemetry-sheet", name, fmt.Sprintf("%d.png", lap)))
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func drawMotionSheet(p string, lap int, driverIndex int, name string) {
	carTelemetries, stSize, err := loadCarTelemetryData(p, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	lapTelemetries, err := loadLapTelemetries(p, lap, driverIndex)
	if err != nil {
		panic(err)
	}

	motionTelemetries, err := loadMotionData(p, lap, driverIndex)
	if err != nil {
		return
	}

	width := len(carTelemetries)
	height := 1500
	//img := image.NewRGBA(image.Rect(0, 0, width, height))

	img := drawSector(lapTelemetries, width, height).(*image.RGBA)
	ctx := gg.NewContextForRGBA(img)
	//err = ctx.LoadFontFace(path.Join(windowsFontPath, "Arial.ttf"), 25)
	//if err != nil {
	//	panic(err)
	//}

	//draw.Draw(img, image.Rect(0, 0, width, height), , image.Point{}, draw.Over)
	totalLapTime := lapTelemetries[len(lapTelemetries)-1].Timestamp - lapTelemetries[0].Timestamp
	log.Printf("lap %d, size:%fkb, lapTime: %s", lap, float32(stSize)/1024, totalLapTime)

	if totalLapTime == 0 {
		return
	}
	//err = ctx.LoadFontFace(windowsFontPath+"/Arial.ttf", 30)
	//if err != nil {
	//	panic(err)
	//}

	draw.Draw(img, image.Rect(0, 0, width, 500), drawStat(carTelemetries, StatSteering, 500, NoVerticalGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 0, width, 500), drawStat(carTelemetries, StatSpeed, 500, NoGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 550, width, 650), drawStat(carTelemetries, StatThrottle, 100, NoGrid), image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(0, 550, width, 650), drawStat(carTelemetries, StatBreak, 100, NoGrid), image.Point{}, draw.Over)

	draw.Draw(img, image.Rect(0, 700, width, 800), drawMotion(motionTelemetries, MotionGForceLatitude, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("G-force Right", 1, 720)

	draw.Draw(img, image.Rect(0, 800, width, 900), drawMotion(motionTelemetries, MotionGForceLongitude, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("G-force Forward", 1, 820)

	draw.Draw(img, image.Rect(0, 900, width, 1000), drawMotion(motionTelemetries, MotionYaw, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Yaw", 1, 920)

	draw.Draw(img, image.Rect(0, 1000, width, 1100), drawMotion(motionTelemetries, MotionWheelSlipFL, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Slip FL", 1, 1020)

	draw.Draw(img, image.Rect(0, 1100, width, 1200), drawMotion(motionTelemetries, MotionWheelSlipFR, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Slip FR", 1, 1120)

	draw.Draw(img, image.Rect(0, 1200, width, 1300), drawMotion(motionTelemetries, MotionWheelSlipRL, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Slip RL", 1, 1220)

	draw.Draw(img, image.Rect(0, 1300, width, 1400), drawMotion(motionTelemetries, MotionWheelSlipRR, 100, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	ctx.SetColor(colornames.Black)
	ctx.DrawString("Slip RR", 1, 1320)

	//draw.Draw(img, image.Rect(0, 1100, width, 1150), drawMotion(motionTelemetries, MotionWheelSpeedFrontBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	//ctx.SetColor(colornames.Black)
	//ctx.DrawString("FL", 1, 1122)
	//ctx.DrawString("FR", 1, 1140)
	//
	//draw.Draw(img, image.Rect(0, 1150, width, 1200), drawMotion(motionTelemetries, MotionWheelSpeedRearBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	//ctx.SetColor(colornames.Black)
	//ctx.DrawString("RL", 1, 1172)
	//ctx.DrawString("RR", 1, 1190)
	//
	//draw.Draw(img, image.Rect(0, 1200, width, 1250), drawMotion(motionTelemetries, MotionWheelSpeedLeftBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	//ctx.SetColor(colornames.Black)
	//ctx.DrawString("FL", 1, 1222)
	//ctx.DrawString("RL", 1, 1240)
	//
	//draw.Draw(img, image.Rect(0, 1250, width, 1300), drawMotion(motionTelemetries, MotionWheelSpeedRightBias, 50, NoGrid, DrawMiddleGrid), image.Point{}, draw.Over)
	//ctx.SetColor(colornames.Black)
	//ctx.DrawString("FR", 1, 1272)
	//ctx.DrawString("RR", 1, 1290)

	f, err := os.Create(path.Join(p, "motion", name, fmt.Sprintf("%d.png", lap)))
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

type Option int

const (
	NoGrid Option = iota
	NoHorizontalGrid
	NoVerticalGrid
	DrawMiddleGrid
)

func drawSector(lts []Lt, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	ctx := gg.NewContextForRGBA(img)
	//err := ctx.LoadFontFace(path.Join(windowsFontPath, "Arial.ttf"), 25)
	//if err != nil {
	//	panic(err)
	//}

	start := lts[0].Timestamp
	totalLapTime := lts[len(lts)-1].Timestamp - start
	gap := float32(img.Rect.Max.X) / float32(totalLapTime.Seconds())
	//log.Println(totalLapTime, gap)

	previousSectorDuration := time.Duration(0)
	lastSector := -1
	lastStatus := f1.DriverStatusNull
	for _, lt := range lts {
		if lastStatus != lt.DriverStatus {
			// lap paused due to pit lane enter

			if lt.DriverStatus == f1.DriverStatusInGarage {
				//log.Println("driver entered pits! sector3 end!", previousSectorDuration, lt.Timestamp-start)
				ctx.SetColor(pitLaneColor)
				ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-previousSectorDuration).Seconds()*float64(gap), float64(height))
				ctx.Fill()

				ctx.SetColor(colornames.Black)
				ctx.DrawString("PIT", previousSectorDuration.Seconds()*float64(gap)+10, 50)

				lastSector = -1
				previousSectorDuration = lt.Timestamp - start
			} else if lt.DriverStatus == f1.DriverStatusOutLap {
				//log.Println("driver existed pits! sector1 started!", previousSectorDuration, lt.Timestamp-start)
				lastSector = 0
				previousSectorDuration = lt.Timestamp - start
			} else if lt.DriverStatus == f1.DriverStatusFlyingLap && lastStatus == f1.DriverStatusInGarage {
				ctx.SetColor(pitLaneColor)
				ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-start-previousSectorDuration).Seconds()*float64(gap), float64(height))
				ctx.Fill()

				ctx.SetColor(colornames.Black)
				ctx.DrawString("PIT", previousSectorDuration.Seconds()*float64(gap)+10, 50)

				lastSector = 2
				previousSectorDuration = lt.Timestamp - start
			}
			//log.Println("new driver status", lt.DriverStatus, lt.Timestamp-start)

			lastStatus = lt.DriverStatus
			continue
		}

		if lt.Sector != lastSector {
			//log.Println("new sector", lt.SetSector, lt.Timestamp-start)
			if lt.Sector == 0 {
				// lap started
				if lastSector == 2 {
					//log.Println("sector 3 Ended!", previousSectorDuration, lt.Timestamp-start)
					ctx.SetColor(sectorColors[2])
					ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-start-previousSectorDuration).Seconds()*float64(gap), float64(height))
					ctx.Fill()

					ctx.SetColor(colornames.Black)
					ctx.DrawString("S3", previousSectorDuration.Seconds()*float64(gap)+10, 50)
				}
				//log.Println("sector 1 started!", previousSectorDuration, lt.Timestamp-start)

				previousSectorDuration = lt.Timestamp - start
			} else if lt.Sector == 1 {
				// lap sector 1 end

				//log.Println("sector 1 end!", previousSectorDuration, lt.Timestamp-start)
				ctx.SetColor(sectorColors[0])
				ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-start-previousSectorDuration).Seconds()*float64(gap), float64(height))
				ctx.Fill()

				ctx.SetColor(colornames.Black)
				ctx.DrawString("S1", previousSectorDuration.Seconds()*float64(gap)+10, 50)

				previousSectorDuration = lt.Timestamp - start
			} else if lt.Sector == 2 {
				// lap sector 2 end
				//log.Println("sector 2 end!", previousSectorDuration, lt.Timestamp-start)
				ctx.SetColor(sectorColors[1])
				ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-start-previousSectorDuration).Seconds()*float64(gap), float64(height))
				ctx.Fill()

				ctx.SetColor(colornames.Black)
				ctx.DrawString("S2", previousSectorDuration.Seconds()*float64(gap)+10, 50)
			}

			previousSectorDuration = lt.Timestamp - start
			//var sectorTime time.Duration

			//if lt.SetSector == 2 {
			//	sectorTime = lt.Sector1Time
			//} else if lt.SetSector == 3 {
			//	sectorTime = lt.Sector2Time
			//} else {
			//	sectorTime = lt.CurrentLapTime - lt.Sector2Time
			//}
			//
			//log.Println(lastSector, lt.SetSector, lt.CurrentLapTime, sectorTime, lt.Timestamp-startTime)
			//ctx.SetColor(sectorColors[lt.SetSector])
			//ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, sectorTime.Seconds()*float64(gap), float64(height))
			//ctx.Fill()
			lastSector = lt.Sector
			//previousSectorDuration += sectorTime
		}
	}
	if lastSector == 2 {
		// lap sector 3 end (maybe)
		lt := lts[len(lts)-1]

		//log.Println("sector 3 end!", previousSectorDuration, lt.Timestamp-start)
		ctx.SetColor(sectorColors[2])
		ctx.DrawRectangle(previousSectorDuration.Seconds()*float64(gap), 0, (lt.Timestamp-start-previousSectorDuration).Seconds()*float64(gap), float64(height))
		ctx.Fill()

		ctx.SetColor(colornames.Black)
		ctx.DrawString("S3", previousSectorDuration.Seconds()*float64(gap)+10, 50)
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

	return img
}

func drawStat(cts []Ct, typ CatTelemetryType, height int, options ...Option) image.Image {
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

	width := len(cts)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	totalLapTime := cts[width-1].LapDuration - cts[0].LapDuration
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

	ctx.SetLineWidth(3)
	verticalMargin := 1

	height = height - verticalMargin*2

	hg := height / 8
	hh := height / 2
	hf := float32(height)
	hhf := float32(hh)
	previousY := float64(0)

	for x, t := range cts {
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

	ctx.SetLineWidth(3)

	verticalMargin := 1

	height = height - verticalMargin*2

	hh := height / 2
	hf := float32(height)
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
		case MotionGForceLatitude:
			c = colornames.Black
			y = int(-m.GForce.X/4*hhf) + hh
		case MotionGForceLongitude:
			c = colornames.Black
			y = int(m.GForce.Y/4*hhf) + hh
		case MotionYaw:
			c = colornames.Black
			y = int(math.Sin(float64(m.Heading.X))*float64(hhf)) + hh
		case MotionWheelSlipFL:
			c = colornames.Black
			//log.Println(m.WheelSpeed.FL)
			y = int(math.Sin(float64(m.WheelSpeed.FL))*float64(hhf)) + hh
		case MotionWheelSlipFR:
			c = colornames.Black
			//log.Println(m.WheelSpeed.FR)
			y = int(math.Sin(float64(m.WheelSpeed.FR))*float64(hhf)) + hh
		case MotionWheelSlipRL:
			c = colornames.Black
			//log.Println(m.WheelSpeed.RL)
			y = int(math.Sin(float64(m.WheelSpeed.RL))*float64(hhf)) + hh
		case MotionWheelSlipRR:
			c = colornames.Black
			//log.Println(m.WheelSpeed.RR)
			y = int(math.Sin(float64(m.WheelSpeed.RR))*float64(hhf)) + hh
		}

		currentY := float64(y)
		if typ == MotionWheelSpeedLeftBias || typ == MotionWheelSpeedFrontBias || typ == MotionWheelSpeedRearBias || typ == MotionWheelSpeedRightBias {
			currentY = math.Min(float64(hf), currentY)
			currentY = math.Max(-float64(hf), currentY)
		}
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
