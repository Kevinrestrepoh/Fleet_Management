use tokio::sync::mpsc;
use tokio_stream::{StreamExt, wrappers::ReceiverStream};
use tonic::{Request, Response, Status};

use crate::pb::vehicle::{
    Command, Telemetry, vehicle_telemetry_service_server::VehicleTelemetryService,
};

use super::registry::VehicleRegistry;

pub struct VehicleTelemetryServiceImpl {
    registry: VehicleRegistry,
}

impl VehicleTelemetryServiceImpl {
    pub fn new(registry: VehicleRegistry) -> Self {
        Self { registry }
    }
}

#[tonic::async_trait]
impl VehicleTelemetryService for VehicleTelemetryServiceImpl {
    type StreamTelemetryStream = ReceiverStream<Result<Command, Status>>;

    async fn stream_telemetry(
        &self,
        request: Request<tonic::Streaming<Telemetry>>,
    ) -> Result<Response<Self::StreamTelemetryStream>, Status> {
        let mut inbound = request.into_inner();

        let (tx, rx) = mpsc::channel::<Result<Command, Status>>(32);

        let registry = self.registry.clone();

        tokio::spawn(async move {
            let mut vehicle_id: Option<u32> = None;

            while let Some(Ok(telemetry)) = inbound.next().await {
                let vid = telemetry.vehicle_id;

                if vehicle_id.is_none() {
                    registry.register(vid, tx.clone()).await;
                    vehicle_id = Some(vid);
                }

                // Process telemetry (for now just log)
                println!(
                    "telemetry received: vehicle={} lat={} lon={} battery={}",
                    vid, telemetry.lat, telemetry.lon, telemetry.battery_percent
                );
            }

            if let Some(id) = vehicle_id {
                registry.unregister(id).await;
            }
        });

        Ok(Response::new(ReceiverStream::new(rx)))
    }
}
