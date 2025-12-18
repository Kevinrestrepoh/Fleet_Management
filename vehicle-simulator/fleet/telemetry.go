package fleet

import (
	"math/rand"
	"time"

	"github.com/Kevinrestrepoh/vehicle-simulator/proto/vehiclepb"
)

func (v *Vehicle) telemetryLoop() {
	for v.running {
		v.move()
		t := v.generateTelemetry()

		_ = v.stream.Send(t)

		time.Sleep(v.UpdateInterval)
	}
}

func (v *Vehicle) move() {
	// Random small movement
	v.Lat += (rand.Float64() - 0.5) * 0.0005
	v.Lon += (rand.Float64() - 0.5) * 0.0005

	v.Speed += (rand.Float64() - 0.5) * 5
	if v.Speed < 0 {
		v.Speed = 0
	}

	// Natural cooling
	v.EngineTemp += (v.Speed / 120.0) * 0.3
	v.EngineTemp += (70 - v.EngineTemp) * 0.02
	v.EngineTemp += (rand.Float64() - 0.5) * 0.2

	drainChance := 0.05 + (v.Speed / 120.0 * 0.15)
	if rand.Float64() < drainChance {
		v.Battery--
	}

	if v.Battery < 0 {
		v.Battery = 0
		v.running = false
		v.stream.Close()
	}
}

func (v *Vehicle) generateTelemetry() *vehiclepb.Telemetry {
	return &vehiclepb.Telemetry{
		VehicleId:      v.ID,
		TimestampMs:    time.Now().UnixMilli(),
		Lat:            v.Lat,
		Lon:            v.Lon,
		SpeedKmh:       float32(v.Speed),
		BatteryPercent: uint32(v.Battery),
		EngineTempC:    float32(v.EngineTemp),
	}
}
