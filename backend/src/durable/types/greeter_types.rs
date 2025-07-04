use serde::Deserialize;

#[derive(Deserialize)]
pub struct GreetRequest {
    pub name: String,
}
