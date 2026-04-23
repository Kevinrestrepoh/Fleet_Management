use thiserror::Error;
use tonic::Status;

#[derive(Error, Debug)]
pub enum IngestionError {
    #[error("Vehicle {0} not found")]
    VehicleNotFound(u32),

    #[error("Vehicle {0} is already connected")]
    VehicleAlreadyConnected(u32),

    #[error("Command channel closed for vehicle {0}")]
    CommandChannelClosed(u32),
}

pub type Result<T> = std::result::Result<T, IngestionError>;

impl From<IngestionError> for Status {
    fn from(err: IngestionError) -> Self {
        match err {
            IngestionError::VehicleNotFound { .. } => Status::not_found(err.to_string()),
            IngestionError::CommandChannelClosed { .. } => Status::unavailable(err.to_string()),
            _ => Status::internal(err.to_string()),
        }
    }
}
