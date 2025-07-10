use restate_sdk::prelude::*;

use crate::define_restate_service;

define_restate_service! {
    #[restate_sdk::object]
    pub trait CounterService {
        async fn increment() -> Result<i64, HandlerError>;
        async fn decrement() -> Result<i64, HandlerError>;
        #[shared]
        async fn get_count() -> Result<i64, HandlerError>;
        async fn reset() -> Result<(), HandlerError>;
        async fn add(value: i64) -> Result<i64, HandlerError>;
    }
}

pub struct CounterServiceServer;

impl CounterService for CounterServiceServer {
    async fn increment<'a>(&self, ctx: ObjectContext<'a>) -> Result<i64, HandlerError> {
        let current_count = ctx.get::<i64>("count").await?.unwrap_or(0);
        println!("Incrementing!");
        let new_count = current_count + 1;
        ctx.set("count", new_count);
        Ok(new_count)
    }

    async fn decrement<'a>(&self, ctx: ObjectContext<'a>) -> Result<i64, HandlerError> {
        let current_count = ctx.get::<i64>("count").await?.unwrap_or(0);
        let new_count = current_count - 1;
        ctx.set("count", new_count);
        Ok(new_count)
    }

    async fn get_count<'a>(&self, ctx: SharedObjectContext<'a>) -> Result<i64, HandlerError> {
        let count = ctx.get::<i64>("count").await?.unwrap_or(0);
        Ok(count)
    }

    async fn reset<'a>(&self, ctx: ObjectContext<'a>) -> Result<(), HandlerError> {
        ctx.set("count", 0);
        Ok(())
    }

    async fn add<'a>(&self, ctx: ObjectContext<'a>, value: i64) -> Result<i64, HandlerError> {
        let current_count = ctx.get::<i64>("count").await?.unwrap_or(0);
        let new_count = current_count + value;
        ctx.set("count", new_count);
        Ok(new_count)
    }
}
