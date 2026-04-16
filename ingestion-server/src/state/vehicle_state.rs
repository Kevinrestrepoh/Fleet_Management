use std::time::Instant;

#[derive(Clone, Debug)]
pub struct VehicleState {
    pub lat: f64,
    pub lon: f64,
    pub speed_kmh: f32,
    pub battery: u32,
    pub engine_temp: f32,
    pub last_telemetry_ms: i64,
    pub last_seen: Instant,
}
