package fleet

import (
	"log"
	"time"

	"github.com/Kevinrestrepoh/vehicle-simulator/proto/vehiclepb"
)

func (v *Vehicle) commandLoop() {
	for {
		cmd, err := v.stream.Recv()
		if err != nil {
			log.Printf("[veh %d] command stream closed: %v", v.ID, err)
			v.running = false
			v.stream.Close()
			return
		}

		switch cmd.Type {
		case vehiclepb.Type_SHUTDOWN:
			log.Printf("[veh %d] received SHUTDOWN", v.ID)
			v.running = false
			v.stream.Close()
			return

		case vehiclepb.Type_UPDATE_RATE:
			v.UpdateInterval = time.Millisecond * time.Duration(cmd.Value)
			log.Printf("[veh %d] update interval -> %dms", v.ID, cmd.Value)

		case vehiclepb.Type_PING:
			log.Printf("[veh %d] received PING", v.ID)

		default:
			log.Printf("[veh %d] unknown command", v.ID)
		}
	}
}
