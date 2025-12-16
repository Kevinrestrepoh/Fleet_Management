use grpc::control_service::ControlServiceImpl;
use pb::control::control_service_server::ControlServiceServer;
use tonic::transport::Server;

mod grpc;
mod metrics;
mod pb;
mod state;
mod vehicle_registry;

use grpc::vehicle_service::VehicleTelemetryServiceImpl;
use pb::vehicle::vehicle_telemetry_service_server::VehicleTelemetryServiceServer;
use state::state_store::StateStore;
use vehicle_registry::registry::VehicleRegistry;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse()?;

    let registry = VehicleRegistry::new();
    let state_store = StateStore::new();
    let vehicle_service = VehicleTelemetryServiceImpl::new(registry.clone(), state_store.clone());
    let control_service = ControlServiceImpl::new(registry);

    println!("Server running on port: {}", addr);

    Server::builder()
        .add_service(VehicleTelemetryServiceServer::new(vehicle_service))
        .add_service(ControlServiceServer::new(control_service))
        .serve(addr)
        .await?;

    Ok(())
}
