package grpc

import (
	"context"
	"time"

	"github.com/Kevinrestrepoh/control-api/proto/controlpb"
	"github.com/Kevinrestrepoh/control-api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IngestionClient struct {
	client controlpb.ControlServiceClient
}

func NewIngestionClient(addr string) (*IngestionClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &IngestionClient{
		client: controlpb.NewControlServiceClient(conn),
	}, nil
}

func (c *IngestionClient) SendCommand(vehicleID uint32, cmd types.CommandRequest, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := c.client.SendCommand(ctx, &controlpb.SendCommandRequest{
		VehicleId: vehicleID,
		Command: &controlpb.Command{
			Type:  controlpb.Type(controlpb.Type_value[cmd.Type]),
			Value: cmd.Value,
		},
	})

	return err
}
