use restate_sdk::prelude::*;

use crate::define_restate_service;

define_restate_service! {
    #[restate_sdk::service]
    pub trait CounterService {
        async fn increment() -> Result<i64, HandlerError>;
        async fn decrement() -> Result<i64, HandlerError>;
        async fn get_count() -> Result<i64, HandlerError>;
        async fn reset() -> Result<(), HandlerError>;
        async fn add(value: i64) -> Result<i64, HandlerError>;
    }
}

pub struct CounterServiceServer;

#[allow(unused_variables, unused_mut)]
impl CounterService for CounterServiceServer {
    async fn increment<'a>(&self, mut ctx: Context<'a>) -> Result<i64, HandlerError> {
        todo!()
    }

    async fn decrement<'a>(&self, mut ctx: Context<'a>) -> Result<i64, HandlerError> {
        todo!()
    }

    async fn get_count<'a>(&self, mut ctx: Context<'a>) -> Result<i64, HandlerError> {
        todo!()
    }

    async fn reset<'a>(&self, mut ctx: Context<'a>) -> Result<(), HandlerError> {
        todo!()
    }

    async fn add<'a>(&self, mut ctx: Context<'a>, value: i64) -> Result<i64, HandlerError> {
        todo!()
    }
}
