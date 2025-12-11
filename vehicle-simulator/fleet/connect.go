package fleet

import (
	"context"

	"github.com/Kevinrestrepoh/vehicle-simulator/proto/vehiclepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Stream interface {
	Send(*vehiclepb.Telemetry) error
	SendCommand(vehiclepb.Type, uint32) error
}

type VehicleStream struct {
	stream vehiclepb.VehicleTelemetryService_StreamTelemetryClient
}

func (s *VehicleStream) Send(t *vehiclepb.Telemetry) error {
	return s.stream.Send(t)
}

func (s *VehicleStream) SendCommand(t vehiclepb.Type, val uint32) error {
	cmd := &vehiclepb.Command{
		Type:  t,
		Value: val,
	}
	return s.stream.SendMsg(cmd)
}

func connectVehicleStream(addr string) (Stream, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := vehiclepb.NewVehicleTelemetryServiceClient(conn)
	stream, err := client.StreamTelemetry(context.Background())
	if err != nil {
		return nil, err
	}

	return &VehicleStream{stream}, nil
}
