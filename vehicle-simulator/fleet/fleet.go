package fleet

import (
	"log"
	"sync"
)

type Fleet struct {
	VehicleCount int
	Addr         string
}

func NewFleet(n int, addr string) *Fleet {
	return &Fleet{VehicleCount: n, Addr: addr}
}

func (f *Fleet) Start() {
	var wg sync.WaitGroup

	for i := 1; i <= f.VehicleCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			StartVehicle(id, f.Addr)
		}(i)
	}

	log.Printf("Started %d vehicle simulators", f.VehicleCount)
}
