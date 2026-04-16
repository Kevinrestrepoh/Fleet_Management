use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;
use tonic::{Request, Response, Status};

use crate::metrics::metrics_aggregator::MetricsAggregator;
use crate::pb;
use crate::pb::metrics::{
    FleetMetrics, FleetVehicles, metrics_service_server::MetricsService,
};

pub struct MetricsServiceImpl {
    aggregator: MetricsAggregator,
}

impl MetricsServiceImpl {
    pub fn new(aggregator: MetricsAggregator) -> Self {
        Self { aggregator }
    }
}

#[tonic::async_trait]
impl MetricsService for MetricsServiceImpl {
    type StreamFleetMetricsStream = ReceiverStream<Result<FleetMetrics, Status>>;
    type StreamFleetVehiclesStream = ReceiverStream<Result<FleetVehicles, Status>>;

    async fn stream_fleet_metrics(
        &self,
        _request: Request<pb::metrics::Empty>,
    ) -> Result<Response<Self::StreamFleetMetricsStream>, Status> {
        let (tx, rx) = mpsc::channel(16);
        let aggregator = self.aggregator.clone();

        tokio::spawn(async move {
            loop {
                let metrics = aggregator.compute().await;

                if tx.send(Ok(metrics)).await.is_err() {
                    break; // client disconnected
                }

                tokio::time::sleep(std::time::Duration::from_secs(1)).await;
            }
        });

        Ok(Response::new(ReceiverStream::new(rx)))
    }

    async fn stream_fleet_vehicles(
        &self,
        _request: Request<pb::metrics::Empty>,
    ) -> Result<Response<Self::StreamFleetVehiclesStream>, Status> {
        let (tx, rx) = mpsc::channel(32);
        let aggregator = self.aggregator.clone();

        tokio::spawn(async move {
            loop {
                let fleet = aggregator.fleet_vehicles_snapshot().await;

                if tx.send(Ok(fleet)).await.is_err() {
                    break;
                }

                tokio::time::sleep(std::time::Duration::from_millis(300)).await;
            }
        });

        Ok(Response::new(ReceiverStream::new(rx)))
    }
}
