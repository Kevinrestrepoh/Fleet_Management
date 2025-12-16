use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

use super::vehicle_state::VehicleState;

#[derive(Clone)]
pub struct StateStore {
    inner: Arc<RwLock<HashMap<u32, VehicleState>>>,
}

impl StateStore {
    pub fn new() -> Self {
        Self {
            inner: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn upsert(&self, vehicle_id: u32, state: VehicleState) {
        self.inner.write().await.insert(vehicle_id, state);
    }

    pub async fn snapshot(&self) -> HashMap<u32, VehicleState> {
        self.inner.read().await.clone()
    }
}
