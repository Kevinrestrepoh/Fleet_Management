use prometheus::{Histogram, IntCounter, IntGauge, Registry, TextEncoder};

#[derive(Clone)]
pub struct Metrics {
    // Connection metrics
    pub active_connections: IntGauge,
    pub total_connections: IntCounter,
    pub connection_errors: IntCounter,

    // Telemetry metrics
    pub telemetry_received: IntCounter,
    pub telemetry_processing_time: Histogram,

    // Command metrics
    pub commands_sent: IntCounter,
    pub command_errors: IntCounter,
    pub command_processing_time: Histogram,

    registry: Registry,
}

impl Metrics {
    pub fn new() -> Self {
        let registry = Registry::new();

        let active_connections = IntGauge::new(
            "fleet_active_connections",
            "Number of currently active vehicle connections",
        )
        .unwrap();
        registry
            .register(Box::new(active_connections.clone()))
            .unwrap();

        let total_connections = IntCounter::new(
            "fleet_total_connections",
            "Total number of connections since server start",
        )
        .unwrap();
        registry
            .register(Box::new(total_connections.clone()))
            .unwrap();

        let connection_errors = IntCounter::new(
            "fleet_connection_errors",
            "Total number of connection errors",
        )
        .unwrap();
        registry
            .register(Box::new(connection_errors.clone()))
            .unwrap();

        let telemetry_received = IntCounter::new(
            "fleet_telemetry_received_total",
            "Total number of telemetry messages received",
        )
        .unwrap();
        registry
            .register(Box::new(telemetry_received.clone()))
            .unwrap();

        let telemetry_processing_time = Histogram::with_opts(
            prometheus::HistogramOpts::new(
                "fleet_telemetry_processing_seconds",
                "Time spent processing telemetry messages",
            )
            .buckets(vec![
                0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
            ]),
        )
        .unwrap();
        registry
            .register(Box::new(telemetry_processing_time.clone()))
            .unwrap();

        let commands_sent = IntCounter::new(
            "fleet_commands_sent_total",
            "Total number of commands sent to vehicles",
        )
        .unwrap();
        registry.register(Box::new(commands_sent.clone())).unwrap();

        let command_errors = IntCounter::new(
            "fleet_command_errors_total",
            "Total number of command sending errors",
        )
        .unwrap();
        registry.register(Box::new(command_errors.clone())).unwrap();

        let command_processing_time = Histogram::with_opts(
            prometheus::HistogramOpts::new(
                "fleet_command_processing_seconds",
                "Time spent processing commands",
            )
            .buckets(vec![
                0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
            ]),
        )
        .unwrap();
        registry
            .register(Box::new(command_processing_time.clone()))
            .unwrap();

        let registry_size =
            IntGauge::new("fleet_registry_size", "Current size of vehicle registry").unwrap();
        registry.register(Box::new(registry_size.clone())).unwrap();

        let server_uptime =
            IntGauge::new("fleet_server_uptime_seconds", "Server uptime in seconds").unwrap();
        registry.register(Box::new(server_uptime.clone())).unwrap();

        Self {
            active_connections,
            total_connections,
            connection_errors,
            telemetry_received,
            telemetry_processing_time,
            commands_sent,
            command_errors,
            command_processing_time,
            registry,
        }
    }

    pub fn encode(&self) -> Result<String, prometheus::Error> {
        let encoder = TextEncoder::new();
        let metric_families = self.registry.gather();
        encoder.encode_to_string(&metric_families)
    }
}
