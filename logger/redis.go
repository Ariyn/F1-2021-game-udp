package logger

import (
	"context"
	"encoding/json"
	"github.com/ariyn/F1-2021-game-udp/packet"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type LapHistory struct {
	Number    int
	StartedAt int64
	EndedAt   int64
}

var _ packet.Logger = (*RedisClient)(nil)

type RedisClient struct {
	ctx                      context.Context
	client                   *redis.Client
	pubsub                   *redis.PubSub
	started                  bool
	startedAt                time.Time
	inputChannel             chan packet.PacketData
	currentLapNumber         int
	distance                 int
	currentLapTimeByDistance map[int]float64
	currentLapDelta          map[int]float64
	lastLapDelta             map[int]float64
	currentSessionUid        string
}

func (r *RedisClient) updateCurrentSessionUid(sessionUid uint64) error {
	uid := strconv.FormatUint(sessionUid, 10)
	if uid != r.currentSessionUid {
		r.currentSessionUid = uid
		r.initRedis()
		return r.setRaw("currentSessionId", uid)
	}

	return nil
}

func (r *RedisClient) GetCurrentSessionUid() error {
	uid := r.client.Get(r.ctx, "currentSessionId").Val()
	if uid != r.currentSessionUid {
		r.currentSessionUid = uid
	}

	return nil
}

func (r *RedisClient) initRedis() {
	_ = r.createTs("frameIdentifier")

	_ = r.createTs("worldPositionX")
	_ = r.createTs("worldPositionY")
	_ = r.createTs("worldPositionZ")

	_ = r.createTs("throttle")
	_ = r.createTs("break")
	_ = r.createTs("steer")
	_ = r.createTs("gear")
	_ = r.createTs("rpm")
	_ = r.createTs("drs")
	_ = r.createTs("speed")

	_ = r.createTs("lapNumber")
	_ = r.createTs("totalDistance")
	_ = r.createTs("lapDistance")
	_ = r.createTs("lapSector")
	_ = r.createTs("lapTime")
	_ = r.createTs("lapDeltaTime")
}

func (r *RedisClient) Writer(ctx context.Context) (c chan<- packet.PacketData, cancel context.CancelFunc, err error) {
	r.inputChannel = make(chan packet.PacketData, 100)
	r.ctx, cancel = context.WithCancel(ctx)

	newCancel := func() {
		cancel()
	}

	return r.inputChannel, newCancel, nil
}

func (r *RedisClient) Run() {
	for data := range r.inputChannel {
		header := data.GetHeader()
		r.updateCurrentSessionUid(header.SessionUid)

		sessionDuration := getUnixTime(header.SessionTime)
		if !r.started {
			r.started = true
			r.startedAt = time.Now().Add(-sessionDuration)
		}
		now := time.Now().UnixMilli()

		switch v := data.(type) {
		case packet.EventData:
			switch v.Event.(type) {
			case packet.SessionStarted:
				r.started = true
				r.startedAt = time.Now()
			}
		case packet.MotionData:
			md := v.Player()
			r.addTs("worldPositionX", now, float64(md.WorldPositionX))
			r.addTs("worldPositionY", now, float64(md.WorldPositionY))
			r.addTs("worldPositionZ", now, float64(md.WorldPositionZ))
		//case packet.SessionData:
		//	sd := v
		case packet.CarTelemetryData:
			lt := v.Player()
			r.addTs("throttle", now, float64(lt.Throttle))
			r.addTs("break", now, float64(lt.Break))
			r.addTs("steer", now, float64(lt.Steer))
			r.addTs("gear", now, float64(lt.Gear))
			r.addTs("rpm", now, float64(lt.EngineRPM))
			r.addTs("drs", now, float64(lt.DRS))
			r.addTs("speed", now, float64(lt.Speed))
			_ = r.publish("CarTelemetryData", "")
		case packet.LapData:
			ld := v.Player()
			if int(ld.CurrentLapNumber) != r.currentLapNumber {
				r.currentLapNumber = int(ld.CurrentLapNumber)
				r.lastLapDelta = r.currentLapDelta
				r.currentLapDelta = make(map[int]float64)
				r.currentLapTimeByDistance = make(map[int]float64)

				if lh, err := r.GetLapHistory(int(ld.CurrentLapNumber) - 1); err == nil {
					lh.EndedAt = now
					r.setLapHistory(int(ld.CurrentLapNumber)-1, lh)
				}

				lh := LapHistory{
					Number:    int(ld.CurrentLapNumber),
					StartedAt: now,
				}
				r.setLapHistory(int(ld.CurrentLapNumber), lh)
			}

			r.addTs("lapNumber", now, float64(ld.CurrentLapNumber))
			r.addTs("totalDistance", now, float64(ld.TotalDistance))
			r.addTs("lapDistance", now, float64(ld.LapDistance))
			r.addTs("lapSector", now, float64(ld.Sector))
			r.addTs("lapTime", now, float64(ld.CurrentLapTime))

			lapDistanceBy10Meters := int(int64(ld.LapDistance)) / r.distance * r.distance
			if _, ok := r.currentLapDelta[lapDistanceBy10Meters]; !ok {
				r.currentLapTimeByDistance[lapDistanceBy10Meters] = float64(ld.CurrentLapTime) / 1000
				r.currentLapDelta[lapDistanceBy10Meters] = r.currentLapTimeByDistance[lapDistanceBy10Meters] - r.currentLapTimeByDistance[lapDistanceBy10Meters-r.distance]

				if _, ok2 := r.lastLapDelta[lapDistanceBy10Meters]; !ok2 {
					r.addTs("lapDeltaTime", now, 0)
				} else {
					r.addTs("lapDeltaTime", now, r.currentLapDelta[lapDistanceBy10Meters]-r.lastLapDelta[lapDistanceBy10Meters])
				}
			}
			_ = r.publish("LapData", "")
		}

		r.addTs("frameIdentifier", now, float64(header.FrameIdentifier))
		err := r.publish("now", strconv.FormatInt(now, 10))
		if err != nil {
			panic(err)
		}
	}
}

func (r *RedisClient) GetLapHistory(number int) (lh LapHistory, err error) {
	data, err := r.client.HGet(r.ctx, "lapHistory", strconv.Itoa(number)).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), &lh)
	return
}

func (r *RedisClient) setLapHistory(number int, lh LapHistory) (err error) {
	b, err := json.Marshal(lh)
	if err != nil {
		return
	}
	return r.client.HSet(r.ctx, "lapHistory", strconv.Itoa(number), string(b)).Err()
}

func (r *RedisClient) hSet(key string, field, value string) error {
	return r.client.HSet(r.ctx, r.getSessionKey(key), field, value).Err()
}

func (r *RedisClient) setRaw(key string, val any) error {
	return r.client.Set(r.ctx, key, val, 0).Err()
}

func (r *RedisClient) set(key string, val any) error {
	return r.client.Set(r.ctx, r.getSessionKey(key), val, 0).Err()
}

func (r *RedisClient) createTs(key string) error {
	return r.client.TSCreate(r.ctx, r.getSessionKey(key)).Err()
}

func (r *RedisClient) addTs(key string, timestampMs int64, value float64) error {
	return r.client.TSAdd(r.ctx, r.getSessionKey(key), timestampMs, value).Err()
}

func (r *RedisClient) GetTs(key string) (val redis.TSTimestampValue, err error) {
	return r.client.TSGet(r.ctx, r.getSessionKey(key)).Result()
}

func (r *RedisClient) GetTsWithRange(key string, from, to int) (val []redis.TSTimestampValue, err error) {
	return r.client.TSRange(r.ctx, r.getSessionKey(key), from, to).Result()
}

func (r *RedisClient) publish(key, val string) error {
	return r.client.Publish(r.ctx, r.getSessionKey(key), val).Err()
}

func (r *RedisClient) Listen(key string) <-chan *redis.Message {
	r.pubsub = r.client.Subscribe(r.ctx, r.getSessionKey(key))
	return r.pubsub.Channel()
}

func (r *RedisClient) CloseListen() error {
	return r.pubsub.Close()
}

func (r *RedisClient) getSessionKey(key string) string {
	return r.currentSessionUid + "-" + key
}

func NewRedisClient(ctx context.Context, url string, db int) (c *RedisClient) {
	c = &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: url,
			DB:   db,
		}),
		ctx:      ctx,
		distance: 50,
	}

	return
}
