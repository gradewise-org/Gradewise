use std::{
    thread,
    time::{Duration, SystemTime},
};

use axum::{extract::Query, response::IntoResponse};
use axum::{http::StatusCode, routing::get, Json, Router};
use rand::Rng;
use serde::{Deserialize, Serialize};

#[tokio::main]
async fn main() {
    // initialize tracing
    tracing_subscriber::fmt::init();

    // build our application with a route
    let app = Router::new()
        // `GET /` goes to `root`
        .route("/", get(root))
        // `POST /users` goes to `create_user`
        .route("/email", get(email_handler))
        // catch-all route that prints the path
        .route("/{*wildcard}", get(catch_all));

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

use axum::extract::Path;

async fn catch_all(Path(path): Path<String>) -> String {
    format!("Requested path: /{}", path)
}

// basic handler that responds with a static string
async fn root() -> &'static str {
    "Hello, World!"
}

#[derive(Deserialize)]
struct EmailRequest {
    to: String,
    body: String,
}

#[derive(Serialize)]
#[serde(tag = "status")]
pub enum EmailResponse {
    #[serde(rename = "success")]
    Success { timestamp: SystemTime },
    #[serde(rename = "failure")]
    Failure { error: String },
}

impl IntoResponse for EmailResponse {
    fn into_response(self) -> axum::response::Response {
        match self {
            EmailResponse::Success { timestamp } => {
                (StatusCode::OK, Json(EmailResponse::Success { timestamp }))
            }
            EmailResponse::Failure { error } => (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(EmailResponse::Failure { error }),
            ),
        }
        .into_response()
    }
}

async fn email_handler(Query(params): Query<EmailRequest>) -> EmailResponse {
    println!("Sending an email to '{}'...", params.to);

    let mut rng = rand::rng();
    if rng.random::<f32>() < 0.3 {
        return EmailResponse::Failure {
            error: format!("Failed to send email to {}", params.to),
        };
    }

    thread::sleep(Duration::from_secs(3));
    println!(
        "SENT an email to '{}'! It reads '{}'.",
        params.to, params.body
    );

    EmailResponse::Success {
        timestamp: SystemTime::now(),
    }
}
