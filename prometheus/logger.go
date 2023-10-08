package prometheus

import (
	f1 "github.com/ariyn/F1-2021-game-udp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
)

type MetricName string

const (
	TotalDistance  MetricName = "TotalDistance"
	Position       MetricName = "Position"
	Throttle       MetricName = "Throttle"
	Break          MetricName = "Break"
	Speed          MetricName = "Speed"
	EngineRPM      MetricName = "EngineRPM"
	Gear           MetricName = "Gear"
	Steer          MetricName = "Steer"
	WorldVelocityX MetricName = "WorldVelocityX"
	WorldVelocityY MetricName = "WorldVelocityY"
	WorldVelocityZ MetricName = "WorldVelocityZ"
	Lap            MetricName = "Lap"
	Sector         MetricName = "Sector"
	PitStatus      MetricName = "PitStatus"
	Distance       MetricName = "Distance"
	Delta          MetricName = "Delta"
)

type CarNumber int

type Logger struct {
	maximumCarNumber int
	session          string
	LapCounter       map[CarNumber]prometheus.Counter
	gauges           map[MetricName]*prometheus.GaugeVec
}

//  TODO: logger는 기본만 제공하고, 실제로 사용하는 녀석이 label과 값을 마음껏 지정할 수 있도록 해보자.
func NewLogger(maxCarNumber int) (l Logger, err error) {
	l = Logger{
		maximumCarNumber: maxCarNumber,
		LapCounter:       make(map[CarNumber]prometheus.Counter),
		gauges:           make(map[MetricName]*prometheus.GaugeVec),
	}

	gaugeNames := []MetricName{
		TotalDistance,
		Position,
		Throttle,
		Break,
		Speed,
		EngineRPM,
		Gear,
		Steer,
		WorldVelocityX,
		WorldVelocityY,
		WorldVelocityZ,
		Lap,
		Sector,
		PitStatus,
		Distance,
		Delta,
	}

	for _, gn := range gaugeNames {
		l.gauges[gn] = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "F1_2021",
			Subsystem: "SINGLE_PLAY",
			Name:      string(gn),
		}, []string{
			"Team",
			"Name",
			"Lap",
			"Session",
		})
	}

	return
}

func (l Logger) Run() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func (l Logger) Write(dataId uint8, carIndex int, data []byte) {

}

// TODO: lap, session, sector로 label 붙일 것
func (l Logger) Driver(driver f1.Driver) NewTypeLogger {
	return NewTypeLogger{
		driver: driver,
		gauges: l.gauges,
	}
}

func (l NewTypeLogger) SetSession(session string) NewTypeLogger {
	l.session = session

	return l
}

func (l NewTypeLogger) SetLapNumber(lap uint8) NewTypeLogger {
	l.lapNumber = int(lap)

	return l
}

func (l NewTypeLogger) Init() NewTypeLogger {
	driver := l.driver
	lapNumber := strconv.Itoa(l.lapNumber)
	session := l.session

	l.totalDistance = l.gauges[TotalDistance].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.position = l.gauges[Position].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.worldVelocityX = l.gauges[WorldVelocityX].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.worldVelocityY = l.gauges[WorldVelocityY].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.worldVelocityZ = l.gauges[WorldVelocityZ].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.throttle = l.gauges[Throttle].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.brk = l.gauges[Break].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.speed = l.gauges[Speed].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.engineRPM = l.gauges[EngineRPM].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gear = l.gauges[Gear].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.steer = l.gauges[Steer].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.lap = l.gauges[Lap].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.sector = l.gauges[Sector].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.pitStatus = l.gauges[PitStatus].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.distance = l.gauges[Distance].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.delta = l.gauges[Delta].WithLabelValues(driver.TeamName, driver.Name, lapNumber, session)

	return l
}

func (l NewTypeLogger) Finish() {
	driver := l.driver
	lapNumber := strconv.Itoa(l.lapNumber)
	session := l.session

	l.gauges[TotalDistance].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Position].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[WorldVelocityX].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[WorldVelocityY].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[WorldVelocityZ].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Throttle].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Break].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Speed].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[EngineRPM].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Gear].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Steer].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Lap].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Sector].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[PitStatus].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Distance].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
	l.gauges[Delta].DeleteLabelValues(driver.TeamName, driver.Name, lapNumber, session)
}

type NewTypeLogger struct {
	driver         f1.Driver
	session        string
	gauges         map[MetricName]*prometheus.GaugeVec
	category       uint8
	lapNumber      int
	delete         func()
	totalDistance  prometheus.Gauge
	position       prometheus.Gauge
	worldVelocityX prometheus.Gauge
	worldVelocityY prometheus.Gauge
	worldVelocityZ prometheus.Gauge
	throttle       prometheus.Gauge
	brk            prometheus.Gauge
	speed          prometheus.Gauge
	engineRPM      prometheus.Gauge
	gear           prometheus.Gauge
	steer          prometheus.Gauge
	lap            prometheus.Gauge
	sector         prometheus.Gauge
	pitStatus      prometheus.Gauge
	distance       prometheus.Gauge
	delta          prometheus.Gauge
}

func (d NewTypeLogger) Category(category uint8) NewTypeLogger {
	d.category = category
	return d
}

func (d NewTypeLogger) TotalDistance(distance float32) {
	d.totalDistance.Set(float64(distance))
}

func (d NewTypeLogger) Position(position uint8) {
	d.position.Set(float64(position))
}

func (d NewTypeLogger) WorldVelocityX(velocity float32) {
	d.worldVelocityX.Set(float64(velocity))
}

func (d NewTypeLogger) WorldVelocityY(velocity float32) {
	d.worldVelocityY.Set(float64(velocity))
}

func (d NewTypeLogger) WorldVelocityZ(velocity float32) {
	d.worldVelocityZ.Set(float64(velocity))
}

func (d NewTypeLogger) Throttle(velocity float32) {
	d.throttle.Set(float64(velocity))
}

func (d NewTypeLogger) Break(velocity float32) {
	d.brk.Set(float64(velocity))
}

func (d NewTypeLogger) Speed(velocity uint16) {
	d.speed.Set(float64(velocity))
}

func (d NewTypeLogger) EngineRPM(velocity uint16) {
	d.engineRPM.Set(float64(velocity))
}

func (d NewTypeLogger) Gear(velocity int8) {
	d.gear.Set(float64(velocity))
}

func (d NewTypeLogger) Steer(velocity float32) {
	d.steer.Set(float64(velocity))
}

func (d *NewTypeLogger) Lap(lap uint8) {
	d.lap.Set(float64(lap))
}

func (d NewTypeLogger) Sector(sector uint8) {
	d.sector.Set(float64(sector))
}

func (d NewTypeLogger) PitStatus(status uint8) {
	d.pitStatus.Set(float64(status))
}

func (d NewTypeLogger) Distance(distance float64) {
	d.distance.Set(distance)
}

func (d NewTypeLogger) Delta(delta float64) {
	d.delta.Set(delta)
}

func (d NewTypeLogger) GetLapNumber() int {
	return d.lapNumber
}

func (d NewTypeLogger) NeedInit() bool {
	return d.lap == nil
}
