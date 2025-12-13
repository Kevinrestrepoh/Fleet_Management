use std::{error::Error, fs, path::PathBuf};

fn main() -> Result<(), Box<dyn Error>> {
    let proto_dir = PathBuf::from("../proto/");

    if !proto_dir.exists() {
        panic!("proto directory not found: {:?}", proto_dir);
    }

    // Collect all the .proto files
    let mut protos = Vec::new();
    for entry in fs::read_dir(&proto_dir)? {
        let path = entry?.path();
        if path.extension().is_some_and(|ext| ext == "proto") {
            protos.push(path);
        }
    }

    if protos.is_empty() {
        panic!("proto directory not found: {:?}", proto_dir);
    }

    // Compile using tonic-build
    tonic_prost_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(&protos, &[proto_dir])
        .map_err(|e| {
            eprintln!("Failed to compile protobufs: {e}");
            e
        })?;

    Ok(())
}
