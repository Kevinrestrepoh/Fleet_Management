use std::error:Error;

fn main() -> Result<(), Box<dyn Error>> {
    dotenvy::dotenv().ok();

    tonic_build::compile_protos("");

    let proto_dir = match env::var("PROTO_DIR") {
         Ok(val) => val,
         Err(_) => {
             eprintln!("ERROR: PROTO_DIR environment variable is not set.");
             std::process::exit(1);
         }
     };

     let proto_dir_path = PathBuf::from(&proto_dir);

     if !proto_dir_path.exists() {
         eprintln!("ERROR: proto directory not found: {}", proto_dir);
         std::process::exit(1);
     }

     // Collect all the .proto files
     let mut protos = Vec::new();
     for entry in fs::read_dir(&proto_dir_path)? {
         let path = entry?.path();
         if path.extension().map(|x| x == "proto").unwrap_or(false) {
             protos.push(path);
         }
     }

     if protos.is_empty() {
         eprintln!("ERROR: No .proto files found in {}", proto_dir);
         std::process::exit(1);
     }

     // Compile using tonic-build
     tonic_build::configure()
         .compile(&protos, &[proto_dir_path])
         .map_err(|e| {
             eprintln!("Failed to compile protobufs: {e}");
             e
         })?;

     Ok(())
}
