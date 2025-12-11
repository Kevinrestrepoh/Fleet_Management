package fleet

import (
	"log"
	"math/rand/v2"
	"time"
)

type Vehicle struct {
	ID             uint32
	Lat, Lon       float64
	Speed          float64
	Battery        int
	EngineTemp     float64
	UpdateInterval time.Duration

	stream Stream
}

func StartVehicle(id int, addr string) {
	v := &Vehicle{
		ID:             uint32(id),
		Lat:            37.7749 + rand.Float64()*0.1,
		Lon:            -122.4194 + rand.Float64()*0.1,
		Speed:          30 + rand.Float64()*40,
		Battery:        100,
		EngineTemp:     60 + rand.Float64()*10,
		UpdateInterval: 500 * time.Millisecond,
	}

	stream, err := connectVehicleStream(addr)
	if err != nil {
		log.Printf("veh %d failed to connect: %v", id, err)
		return
	}
	v.stream = stream

	go v.autonomousBehaviorLoop()

	v.telemetryLoop()
}
