use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{RwLock, mpsc};
use tonic::Status;

use crate::pb::vehicle::Command;

pub type CommandSender = mpsc::Sender<Result<Command, Status>>;

#[derive(Clone)]
pub struct VehicleRegistry {
    inner: Arc<RwLock<HashMap<u32, CommandSender>>>,
}

impl VehicleRegistry {
    pub fn new() -> Self {
        Self {
            inner: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn register(&self, vehicle_id: u32, sender: CommandSender) {
        self.inner.write().await.insert(vehicle_id, sender);
    }

    pub async fn unregister(&self, vehicle_id: u32) {
        self.inner.write().await.remove(&vehicle_id);
    }

    pub async fn send_command(&self, vehicle_id: u32, cmd: Command) -> Result<(), Status> {
        let map = self.inner.read().await;
        let sender = map
            .get(&vehicle_id)
            .ok_or(Status::not_found("vehicle not connected"))?;
        sender
            .send(Ok(cmd))
            .await
            .map_err(|_| Status::internal("channel closed"))
    }
}
