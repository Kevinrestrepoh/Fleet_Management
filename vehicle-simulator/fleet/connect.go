package fleet

import (
	"context"

	"github.com/Kevinrestrepoh/vehicle-simulator/proto/vehiclepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Stream interface {
	Send(*vehiclepb.Telemetry) error
	Recv() (*vehiclepb.Command, error)
	Close() error
}

type VehicleStream struct {
	stream vehiclepb.VehicleTelemetryService_StreamTelemetryClient
}

func (s *VehicleStream) Send(t *vehiclepb.Telemetry) error {
	return s.stream.Send(t)
}

func (s *VehicleStream) Recv() (*vehiclepb.Command, error) {
	return s.stream.Recv()
}

func (s *VehicleStream) Close() error {
	return s.stream.CloseSend()
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
