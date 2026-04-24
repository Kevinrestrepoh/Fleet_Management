use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{RwLock, mpsc};
use tokio::task::JoinSet;
use tonic::Status;

use crate::error::{IngestionError, Result};
use crate::pb::vehicle::Command;

pub type CommandSender = mpsc::Sender<std::result::Result<Command, Status>>;

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

    pub async fn register(&self, vehicle_id: u32, sender: CommandSender) -> Result<()> {
        let mut registry = self.inner.write().await;

        if registry.contains_key(&vehicle_id) {
            return Err(IngestionError::VehicleAlreadyConnected(vehicle_id));
        }

        registry.insert(vehicle_id, sender);
        Ok(())
    }

    pub async fn unregister(&self, vehicle_id: u32) {
        self.inner.write().await.remove(&vehicle_id);
    }

    pub async fn send_command(&self, vehicle_id: u32, cmd: Command) -> Result<()> {
        let map = self.inner.read().await;
        let sender = map
            .get(&vehicle_id)
            .ok_or_else(|| IngestionError::VehicleNotFound(vehicle_id))?;

        sender
            .send(Ok(cmd))
            .await
            .map_err(|_| IngestionError::CommandChannelClosed(vehicle_id))
    }

    pub async fn shutdown(&self) {
        println!("Shutting down vehicle registry...");

        let senders: Vec<_> = {
            let mut registry = self.inner.write().await;
            registry.drain().collect()
        };

        let mut set = JoinSet::new();
        for (_, sender) in senders {
            set.spawn(async move {
                let _ = sender
                    .send(Err(Status::unavailable("Server shutting down")))
                    .await;
            });
        }
        set.join_all().await;

        println!("Vehicle registry shutdown complete");
    }
}
