use std::time::{SystemTime, UNIX_EPOCH};

use crate::pb::metrics::FleetMetrics;
use crate::state::state_store::StateStore;

#[derive(Clone)]
pub struct MetricsAggregator {
    store: StateStore,
}

impl MetricsAggregator {
    pub fn new(store: StateStore) -> Self {
        Self { store }
    }

    pub async fn compute(&self) -> FleetMetrics {
        let snapshot = self.store.snapshot().await;

        let total = snapshot.len() as u32;
        let mut active = 0;
        let mut low_battery = 0;

        let mut speed_sum = 0.0;
        let mut temp_sum = 0.0;

        for v in snapshot.values() {
            active += 1;

            if v.battery < 10 {
                low_battery += 1;
            }

            speed_sum += v.speed_kmh;
            temp_sum += v.engine_temp;
        }

        let count = snapshot.len().max(1) as f32;

        FleetMetrics {
            total_vehicles: total,
            active_vehicles: active,
            low_battery_vehicles: low_battery,
            avg_speed_kmh: speed_sum / count,
            avg_engine_temp_c: temp_sum / count,
            timestamp_ms: SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_millis() as i64,
        }
    }
}
