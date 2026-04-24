use std::sync::Arc;
use std::time::Instant;

use tokio::sync::mpsc;
use tokio_stream::{StreamExt, wrappers::ReceiverStream};
use tonic::{Request, Response, Status};

use crate::metrics::prometheus::Metrics;
use crate::pb::vehicle::{
    Command, Telemetry, vehicle_telemetry_service_server::VehicleTelemetryService,
};
use crate::state::{state_store::StateStore, vehicle_state::VehicleState};
use crate::vehicle_registry::registry::VehicleRegistry;

pub struct VehicleTelemetryServiceImpl {
    registry: VehicleRegistry,
    state_store: StateStore,
    metrics: Arc<Metrics>,
}

impl VehicleTelemetryServiceImpl {
    pub fn new(registry: VehicleRegistry, state_store: StateStore, metrics: Arc<Metrics>) -> Self {
        Self {
            registry,
            state_store,
            metrics,
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
        let metrics = self.metrics.clone();

        tokio::spawn(async move {
            let mut vehicle_id: Option<u32> = None;

            while let Some(Ok(telemetry)) = inbound.next().await {
                let start_time = Instant::now();
                let vid = telemetry.vehicle_id;

                if vehicle_id.is_none() {
                    if let Err(e) = registry.register(vid, tx.clone()).await {
                        println!("failed to register vehicle {}: {}", vid, e);
                        metrics.connection_errors.inc();
                        continue;
                    }
                    vehicle_id = Some(vid);
                    metrics.total_connections.inc();
                    metrics.active_connections.inc();
                }

                let state = VehicleState {
                    lat: telemetry.lat,
                    lon: telemetry.lon,
                    speed_kmh: telemetry.speed_kmh,
                    battery: telemetry.battery_percent,
                    engine_temp: telemetry.engine_temp_c,
                    last_telemetry_ms: telemetry.timestamp_ms,
                    last_seen: Instant::now(),
                };

                metrics.telemetry_received.inc();
                state_store.upsert(vid, state).await;

                let processing_time = start_time.elapsed().as_secs_f64();
                metrics.telemetry_processing_time.observe(processing_time);
            }

            if let Some(id) = vehicle_id {
                registry.unregister(id).await;
                metrics.active_connections.dec();
            }
        });

        Ok(Response::new(ReceiverStream::new(rx)))
    }
}
