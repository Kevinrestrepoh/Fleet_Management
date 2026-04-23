package fleet

import (
	"log"
	"sync"
	"time"
)

type Fleet struct {
	VehicleCount int
	Addr         string
}

func NewFleet(n int, addr string) *Fleet {
	return &Fleet{VehicleCount: n, Addr: addr}
}

func (f *Fleet) Start(batchSize int) {
	var wg sync.WaitGroup

	for i := 1; i <= f.VehicleCount; i += batchSize {
		end := i + batchSize
		end = min(end, f.VehicleCount+1)

		for j := i; j < end; j++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				StartVehicle(id, f.Addr)
			}(j)
			time.Sleep(time.Millisecond * 20)
		}

		log.Printf("Started batch %d-%d", i, end-1)
		time.Sleep(time.Millisecond * 50)
	}

	log.Printf("Started %d vehicle simulators", f.VehicleCount)
}
