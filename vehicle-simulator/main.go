package main

import (
	"github.com/kevinrst/vehicle-simulator/fleet"
)

func main() {
	f := fleet.NewFleet(1000, "localhost:50051")
	f.Start(10)

	select {}
}
