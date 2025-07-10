use axum::{routing::get, Router};
use gradewise_backend::{api::counter, durable::service::counter::HttpCounterServiceClient};

const RESTATE_SERVER_BASE_URL: &str = "http://restate-server:8080";

#[tokio::main]
async fn main() {
    // initialize tracing
    tracing_subscriber::fmt::init();

    // build our application with a route
    let app = Router::new()
        .route("/", get(hello_world))
        .route("/health", get(|| async { "Healthy!" }))
        .nest("/counter", counter::routes())
        .nest("/service", Router::new()
            .route("/add", get(call_service))
            .route("/count", get(get_count_service))
        );

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn hello_world() -> &'static str {
    "Hello, World!"
}

async fn call_service() -> String {
    let client = HttpCounterServiceClient::new(RESTATE_SERVER_BASE_URL.to_string());

    println!("Calling the counter service...");
    let response = client
        .increment().await;
    println!("Received response from the counter service...");
    match response {
        Ok(res) => format!("Counter incremented successfully: {res:?}"),
        Err(e) => format!("Failed to increment counter: {e}"),
    }
}

async fn get_count_service() -> String {
    let client = HttpCounterServiceClient::new(RESTATE_SERVER_BASE_URL.to_string());
    println!("Calling the counter service get_count...");
    let response = client.get_count().await;
    println!("Received response from the counter service get_count...");
    match response {
        Ok(res) => format!("Current counter value: {res:?}"),
        Err(e) => format!("Failed to get counter: {e}"),
    }
}
