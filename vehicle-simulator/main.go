package main

import (
	"github.com/Kevinrestrepoh/vehicle-simulator/fleet"
)

func main() {
	f := fleet.NewFleet(100, "localhost:50051")
	f.Start()

	select {}
}
