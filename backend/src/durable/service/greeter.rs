use crate::durable::types::greeter_types::GreetRequest;
use restate_sdk::prelude::*;

#[restate_sdk::service]
pub trait GreeterService {
    async fn greet(request: Json<GreetRequest>) -> Result<String, HandlerError>;
}

pub struct GreeterServiceImpl;

impl GreeterService for GreeterServiceImpl {
    async fn greet(
        &self,
        _ctx: Context<'_>,
        Json(request): Json<GreetRequest>, // extract from Json<T> to T
    ) -> Result<String, HandlerError> {
        Ok(format!("Hello, {}!", request.name))
    }
}
