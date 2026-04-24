use axum::{Router, extract::State, http::StatusCode, response::IntoResponse, routing::get};
use std::sync::Arc;
use std::time::Instant;

use crate::metrics::prometheus::Metrics;

pub struct MetricsServer {
    metrics: Arc<Metrics>,
    start_time: Instant,
}

impl MetricsServer {
    pub fn new(metrics: Arc<Metrics>) -> Self {
        Self {
            metrics,
            start_time: Instant::now(),
        }
    }

    pub fn router(&self) -> Router {
        Router::new()
            .route("/metrics", get(metrics_handler))
            .with_state(self.metrics.clone())
            .route("/health", get(health_handler))
            .with_state(self.start_time)
    }

    pub async fn run(self, addr: String) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        println!("Metrics running on port: {}", addr);

        let listener = tokio::net::TcpListener::bind(addr).await?;
        axum::serve(listener, self.router()).await?;

        Ok(())
    }
}

async fn metrics_handler(
    State(metrics): State<Arc<Metrics>>,
) -> Result<impl IntoResponse, StatusCode> {
    match metrics.encode() {
        Ok(output) => Ok(output),
        Err(_) => Err(StatusCode::INTERNAL_SERVER_ERROR),
    }
}

async fn health_handler(State(start_time): State<Instant>) -> impl IntoResponse {
    let uptime = start_time.elapsed().as_secs();
    format!("Uptime: {} seconds", uptime)
}
