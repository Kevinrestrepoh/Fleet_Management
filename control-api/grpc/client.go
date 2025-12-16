package grpc

import (
	"context"
	"time"

	"github.com/Kevinrestrepoh/control-api/proto/controlpb"
	"github.com/Kevinrestrepoh/control-api/proto/metricspb"
	"github.com/Kevinrestrepoh/control-api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	control controlpb.ControlServiceClient
	metrics metricspb.MetricsServiceClient
}

func NewIngestionClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		control: controlpb.NewControlServiceClient(conn),
		metrics: metricspb.NewMetricsServiceClient(conn),
	}, nil
}

func (c *Client) SendCommand(vehicleID uint32, cmd types.CommandRequest, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := c.control.SendCommand(ctx, &controlpb.SendCommandRequest{
		VehicleId: vehicleID,
		Command: &controlpb.Command{
			Type:  controlpb.Type(controlpb.Type_value[cmd.Type]),
			Value: cmd.Value,
		},
	})

	return err
}

func (c *Client) Stream(ctx context.Context) (metricspb.MetricsService_StreamFleetMetricsClient, error) {
	return c.metrics.StreamFleetMetrics(ctx, &metricspb.Empty{})
}
