use gradewise_backend::durable::service::counter::{CounterService, CounterServiceServer};
use restate_sdk::prelude::*;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    // Start an HTTP server that serves the Counter Service
    HttpServer::new(
        Endpoint::builder()
            .bind(CounterServiceServer.serve())
            .build(),
    )
    .listen_and_serve("0.0.0.0:8080".parse().unwrap())
    .await;

    println!("Server ready and listening...")
}
