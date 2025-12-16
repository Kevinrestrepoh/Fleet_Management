use std::time::Instant;

#[derive(Clone, Debug)]
pub struct VehicleState {
    pub speed_kmh: f32,
    pub battery: u32,
    pub engine_temp: f32,
    pub last_seen: Instant,
}
