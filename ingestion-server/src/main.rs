use tonic::transport::Server;

mod grpc;
mod pb;

use grpc::registry::VehicleRegistry;
use grpc::vehicle_service::VehicleTelemetryServiceImpl;
use pb::vehicle::vehicle_telemetry_service_server::VehicleTelemetryServiceServer;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse()?;

    let registry = VehicleRegistry::new();
    let vehicle_service = VehicleTelemetryServiceImpl::new(registry);

    Server::builder()
        .add_service(VehicleTelemetryServiceServer::new(vehicle_service))
        .serve(addr)
        .await?;

    Ok(())
}
