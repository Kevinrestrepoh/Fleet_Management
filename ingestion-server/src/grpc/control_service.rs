use crate::pb::control::{
    SendCommandRequest, SendCommandResponse, control_service_server::ControlService,
};
use crate::pb::vehicle::Command;
use crate::vehicle_registry::registry::VehicleRegistry;
use tonic::{Request, Response, Status};

pub struct ControlServiceImpl {
    registry: VehicleRegistry,
}

impl ControlServiceImpl {
    pub fn new(registry: VehicleRegistry) -> Self {
        Self { registry }
    }
}

#[tonic::async_trait]
impl ControlService for ControlServiceImpl {
    async fn send_command(
        &self,
        request: Request<SendCommandRequest>,
    ) -> Result<Response<SendCommandResponse>, Status> {
        let req = request.into_inner();

        let vehicle_id = req.vehicle_id;
        let command = req
            .command
            .ok_or(Status::invalid_argument("command missing"))?;

        // Forward command to vehicle via registry
        self.registry
            .send_command(
                vehicle_id,
                Command {
                    r#type: command.r#type,
                    value: command.value,
                },
            )
            .await?;

        Ok(Response::new(SendCommandResponse {
            success: true,
            message: "command sent".to_string(),
        }))
    }
}
