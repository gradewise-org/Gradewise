mod durable;

use crate::durable::service::{greeter::GreeterService, GreeterServiceImpl};
use restate_sdk::prelude::*;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    HttpServer::new(Endpoint::builder().bind(GreeterServiceImpl.serve()).build())
        .listen_and_serve("0.0.0.0:9080".parse().unwrap())
        .await;
}
