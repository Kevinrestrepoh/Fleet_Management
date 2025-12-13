package fleet

import (
	"log"
	"math/rand"
	"time"

	"github.com/Kevinrestrepoh/vehicle-simulator/proto/vehiclepb"
)

func (v *Vehicle) autonomousBehaviorLoop() {
	for {
		delay := time.Duration(10+rand.Intn(20)) * time.Second
		time.Sleep(delay)

		// If battery low, shutdown
		if v.Battery < 5 {
			log.Printf("[veh %d] auto-shutdown: battery low", v.ID)
			v.stream.SendCommand(vehiclepb.Type_SHUTDOWN, 0)
			return
		}

		// Maybe update rate
		if rand.Intn(3) == 0 {
			newRate := uint32(100 + rand.Intn(900))
			log.Printf("[veh %d] auto CHANGE RATE -> %dms", v.ID, newRate)
			v.stream.SendCommand(vehiclepb.Type_UPDATE_RATE, newRate)
			v.UpdateInterval = time.Millisecond * time.Duration(newRate)
		}

		// Occasionally PING
		if rand.Intn(4) == 0 {
			log.Printf("[veh %d] auto PING", v.ID)
			v.stream.SendCommand(vehiclepb.Type_PING, 0)
		}
	}
}
