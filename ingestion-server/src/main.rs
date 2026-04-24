use metrics::metrics_aggregator::MetricsAggregator;
use pb::vehicle::vehicle_telemetry_service_server::VehicleTelemetryServiceServer;
use pb::{
    control::control_service_server::ControlServiceServer,
    metrics::metrics_service_server::MetricsServiceServer,
};
use tokio::signal;
use tokio::sync::broadcast;
use tonic::transport::Server;

mod error;
mod grpc;
mod metrics;
mod pb;
mod state;
mod vehicle_registry;

use std::sync::Arc;

use grpc::{
    control_service::ControlServiceImpl, metrics_service::MetricsServiceImpl,
    vehicle_service::VehicleTelemetryServiceImpl,
};
use state::state_store::StateStore;
use vehicle_registry::registry::VehicleRegistry;

use crate::metrics::{prometheus::Metrics, prometheus_server::MetricsServer};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse()?;
    let metrics_addr = "0.0.0.0:9090".parse()?;

    let metrics = Arc::new(Metrics::new());

    let registry = VehicleRegistry::new();
    let state_store = StateStore::new();
    let metrics_aggregator = MetricsAggregator::new(state_store.clone());

    let vehicle_service =
        VehicleTelemetryServiceImpl::new(registry.clone(), state_store, metrics.clone());
    let control_service = ControlServiceImpl::new(registry.clone(), metrics.clone());
    let metrics_service = MetricsServiceImpl::new(metrics_aggregator);

    // Start metrics server
    let metrics_server = MetricsServer::new(metrics.clone());
    let metrics_handle = tokio::spawn(async move {
        if let Err(e) = metrics_server.run(metrics_addr).await {
            println!("Metrics server error: {}", e);
        }
    });

    // Set up graceful shutdown
    let (shutdown_tx, mut shutdown_rx) = broadcast::channel::<()>(1);

    let shutdown_tx_clone = shutdown_tx.clone();
    tokio::spawn(async move {
        if let Err(e) = signal::ctrl_c().await {
            println!("Failed to listen for Ctrl-C: {}", e);
            return;
        }
        if let Err(e) = shutdown_tx_clone.send(()) {
            println!("Failed to send shutdown signal: {}", e);
        }
    });

    println!("Server running on port: {}", addr);

    Server::builder()
        .add_service(VehicleTelemetryServiceServer::new(vehicle_service))
        .add_service(ControlServiceServer::new(control_service))
        .add_service(MetricsServiceServer::new(metrics_service))
        .serve_with_shutdown(addr, async {
            let _ = shutdown_rx.recv().await;
            println!("Shutdown signal received");
        })
        .await?;

    // Cleanup
    registry.shutdown().await;
    metrics_handle.abort();

    Ok(())
}
