use std::time::Instant;

use tokio::sync::mpsc;
use tokio_stream::{StreamExt, wrappers::ReceiverStream};
use tonic::{Request, Response, Status};

use crate::pb::vehicle::{
    Command, Telemetry, vehicle_telemetry_service_server::VehicleTelemetryService,
};
use crate::state::{state_store::StateStore, vehicle_state::VehicleState};
use crate::vehicle_registry::registry::VehicleRegistry;

pub struct VehicleTelemetryServiceImpl {
    registry: VehicleRegistry,
    state_store: StateStore,
}

impl VehicleTelemetryServiceImpl {
    pub fn new(registry: VehicleRegistry, state_store: StateStore) -> Self {
        Self {
            registry,
            state_store,
        }
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
        let state_store = self.state_store.clone();

        tokio::spawn(async move {
            let mut vehicle_id: Option<u32> = None;

            while let Some(Ok(telemetry)) = inbound.next().await {
                let vid = telemetry.vehicle_id;

                if vehicle_id.is_none() {
                    registry.register(vid, tx.clone()).await;
                    vehicle_id = Some(vid);
                }

                let state = VehicleState {
                    speed_kmh: telemetry.speed_kmh,
                    battery: telemetry.battery_percent,
                    engine_temp: telemetry.engine_temp_c,
                    last_seen: Instant::now(),
                };
                println!("telemetry received: {:?}", state);
                state_store.upsert(vid, state).await;
            }

            if let Some(id) = vehicle_id {
                registry.unregister(id).await;
            }
        });

        Ok(Response::new(ReceiverStream::new(rx)))
    }
}
