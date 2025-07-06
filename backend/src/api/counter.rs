use std::sync::{Arc, RwLock};

use axum::{extract::State, routing::get, Json, Router};

#[derive(Clone)]
struct AppState {
    number: i64,
}

impl AppState {
    pub fn new() -> Self {
        AppState { number: 0 }
    }
}

pub fn routes<S>() -> Router<S> {
    // create a `Router` that holds our state
    Router::new()
        .route("/", get(get_handler))
        .route("/add", get(increment_handler))
        // provide the state so the router can access it
        .with_state(Arc::new(RwLock::new(AppState::new())))
}

/* TODO: check if this is thread-safe (ref. to RwLock doc) */

async fn get_handler(
    // access the state via the `State` extractor
    // extracting a state of the wrong type results in a compile error
    State(state): State<Arc<RwLock<AppState>>>,
) -> Json<i64> {
    // use `state`...
    let guard = state.read().unwrap();
    Json(guard.number)
}

async fn increment_handler(
    // access the state via the `State` extractor
    // extracting a state of the wrong type results in a compile error
    State(state): State<Arc<RwLock<AppState>>>,
) -> Json<i64> {
    let mut guard = state.write().unwrap();
    guard.number += 1;
    Json(guard.number)
}
