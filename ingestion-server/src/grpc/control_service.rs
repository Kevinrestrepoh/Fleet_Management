use std::sync::Arc;
use std::time::Instant;

use crate::metrics::prometheus::Metrics;
use crate::pb::control::{
    SendCommandRequest, SendCommandResponse, control_service_server::ControlService,
};
use crate::pb::vehicle::Command;
use crate::vehicle_registry::registry::VehicleRegistry;
use tonic::{Request, Response, Status};

pub struct ControlServiceImpl {
    registry: VehicleRegistry,
    metrics: Arc<Metrics>,
}

impl ControlServiceImpl {
    pub fn new(registry: VehicleRegistry, metrics: Arc<Metrics>) -> Self {
        Self { registry, metrics }
    }
}

#[tonic::async_trait]
impl ControlService for ControlServiceImpl {
    async fn send_command(
        &self,
        request: Request<SendCommandRequest>,
    ) -> Result<Response<SendCommandResponse>, Status> {
        let start_time = Instant::now();
        let req = request.into_inner();

        let vehicle_id = req.vehicle_id;

        let command = req
            .command
            .ok_or_else(|| Status::invalid_argument("command missing"))?;

        // Forward command to vehicle via registry
        let result = self
            .registry
            .send_command(
                vehicle_id,
                Command {
                    r#type: command.r#type,
                    value: command.value,
                },
            )
            .await;

        match result {
            Ok(()) => {
                self.metrics.commands_sent.inc();
                let processing_time = start_time.elapsed().as_secs_f64();
                self.metrics
                    .command_processing_time
                    .observe(processing_time);

                Ok(Response::new(SendCommandResponse {
                    success: true,
                    message: "command sent".to_string(),
                }))
            }
            Err(e) => {
                self.metrics.command_errors.inc();
                Err(Status::from(e))
            }
        }
    }
}
