pub type HandlerResult<T> = Result<T, anyhow::Error>;

// Global error type for all Restate clients
#[derive(Debug, thiserror::Error)]
pub enum DurableClientError {
    #[error("HTTP error: {0}")]
    HttpError(reqwest::StatusCode),
    #[error("Request error: {0}")]
    RequestError(#[from] reqwest::Error),
    #[error("Serialization error: {0}")]
    SerializationError(#[from] serde_json::Error),
    #[error("Service error: {0}")]
    ServiceError(String),
}

// Helper macro for methods with no parameters (no request body)
#[macro_export]
macro_rules! generate_method_no_params {
    ($service_trait:ident, $method_name:ident, $return_type:ty) => {
        paste::paste! {
            pub async fn $method_name(&self) -> Result<$return_type, $crate::durable::service::rules::DurableClientError> {
                let service_name = stringify!($service_trait);
                // For objects, we need to include an object key. Using "default" as the key.
                let url = format!("{}/{}/default/{}", self.base_url, service_name, stringify!($method_name));

                let response = self.client
                    .post(&url)
                    .send()
                    .await
                    .map_err($crate::durable::service::rules::DurableClientError::RequestError)?;

                let status = response.status();
                let text = response.text().await.map_err($crate::durable::service::rules::DurableClientError::RequestError)?;
                if status.is_success() {
                    serde_json::from_str(&text).map_err($crate::durable::service::rules::DurableClientError::SerializationError)
                } else {
                    if !text.is_empty() {
                        Err($crate::durable::service::rules::DurableClientError::ServiceError(text))
                    } else {
                        Err($crate::durable::service::rules::DurableClientError::HttpError(status))
                    }
                }
            }
        }
    };
}

// Helper macro for methods with parameters (JSON request body)
#[macro_export]
macro_rules! generate_method_with_params {
    ($service_trait:ident, $method_name:ident, $return_type:ty, $($param_name:ident: $param_type:ty),*) => {
        paste::paste! {
            pub async fn $method_name(&self, $($param_name: $param_type),*) -> Result<$return_type, $crate::durable::service::rules::DurableClientError> {
                let service_name = stringify!($service_trait);
                // For objects, we need to include an object key. Using "default" as the key.
                let url = format!("{}/{}/default/{}", self.base_url, service_name, stringify!($method_name));

                // Build request payload with parameters
                let payload = {
                    #[allow(unused_mut)]
                    let mut map = serde_json::Map::new();
                    $(
                        map.insert(stringify!($param_name).to_string(), serde_json::to_value(&$param_name).map_err($crate::durable::service::rules::DurableClientError::SerializationError)?);
                    )*
                    serde_json::Value::Object(map)
                };

                let response = self.client
                    .post(&url)
                    .header("Content-Type", "application/json")
                    .json(&payload)
                    .send()
                    .await
                    .map_err($crate::durable::service::rules::DurableClientError::RequestError)?;

                let status = response.status();
                let text = response.text().await.map_err($crate::durable::service::rules::DurableClientError::RequestError)?;
                if status.is_success() {
                    serde_json::from_str(&text).map_err($crate::durable::service::rules::DurableClientError::SerializationError)
                } else {
                    if !text.is_empty() {
                        Err($crate::durable::service::rules::DurableClientError::ServiceError(text))
                    } else {
                        Err($crate::durable::service::rules::DurableClientError::HttpError(status))
                    }
                }
            }
        }
    };
}

// Macro that defines a Restate service trait or object and generates an HTTP client implementation for it.
//
// Usage:
//
// ```
// define_restate_service! {
//     #[restate_sdk::service] // or #[restate_sdk::object]
//     pub trait UserService {
//         async fn create_user(name: String, email: String) -> Result<User, HandlerError>;
//         async fn get_user(id: i64) -> Result<User, HandlerError>;
//         async fn update_user(id: i64, name: String) -> Result<User, HandlerError>;
//         async fn delete_user(id: i64) -> Result<(), HandlerError>;
//     }
// }
// ```
//
// This will:
// - Define the trait `UserService` (with all provided attributes and visibility)
// - Generate an HTTP client struct `HttpUserServiceClient` with methods matching the trait, returning `Result<_, DurableClientError>`
//
// Example generated client usage:
// ```
// let client = HttpUserServiceClient::new("http://localhost:8080".to_string());
// let user = client.create_user("Alice".to_string(), "alice@example.com".to_string()).await?;
// ```
#[macro_export]
macro_rules! define_restate_service {
    (
        $(#[$meta:meta])* // Allow doc comments and attributes (including #[restate_sdk::service] or #[restate_sdk::object])
        $vis:vis trait $service_trait:ident {
            $(
                $(#[$method_meta:meta])* // Allow method-level attributes like #[shared]
                async fn $method_name:ident($($param_name:ident: $param_type:ty),*) -> Result<$return_type:ty, HandlerError>;
            )*
        }
    ) => {
        paste::paste! {
            // Paste the trait definition as-is
            $(#[$meta])* // This will pass through #[restate_sdk::service] or #[restate_sdk::object]
            $vis trait $service_trait {
                $(
                    $(#[$method_meta])* // Pass through method-level attributes
                    async fn $method_name($($param_name: $param_type),*) -> Result<$return_type, HandlerError>;
                )*
            }

            // Generate HTTP client struct
            pub struct [<Http $service_trait Client>] {
                client: reqwest::Client,
                base_url: String,
            }

            impl [<Http $service_trait Client>] {
                pub fn new(base_url: String) -> Self {
                    Self {
                        client: reqwest::Client::new(),
                        base_url,
                    }
                }

                pub fn with_client(client: reqwest::Client, base_url: String) -> Self {
                    Self {
                        client,
                        base_url,
                    }
                }

                // Generate convenience methods that match the service trait
                $(
                    $crate::generate_method_impl!($service_trait, $method_name, $return_type, $($param_name: $param_type),*);
                )*
            }
        }
    };
}

// Helper macro to generate method implementations
#[macro_export]
macro_rules! generate_method_impl {
    // No parameters
    ($service_trait:ident, $method_name:ident, $return_type:ty, ) => {
        paste::paste! {
            pub async fn $method_name(&self) -> Result<$return_type, $crate::durable::service::rules::DurableClientError> {
                let service_name = stringify!($service_trait);
                // For objects, we need to include an object key. Using "default" as the key.
                let url = format!("{}/{}/default/{}", self.base_url, service_name, stringify!($method_name));

                let response = self.client
                    .post(&url)
                    .send()
                    .await
                    .map_err($crate::durable::service::rules::DurableClientError::RequestError)?;

                let status = response.status();
                let text = response.text().await.map_err($crate::durable::service::rules::DurableClientError::RequestError)?;
                if status.is_success() {
                    serde_json::from_str(&text).map_err($crate::durable::service::rules::DurableClientError::SerializationError)
                } else {
                    if !text.is_empty() {
                        Err($crate::durable::service::rules::DurableClientError::ServiceError(text))
                    } else {
                        Err($crate::durable::service::rules::DurableClientError::HttpError(status))
                    }
                }
            }
        }
    };

    // With parameters
    ($service_trait:ident, $method_name:ident, $return_type:ty, $($param_name:ident: $param_type:ty),+) => {
        paste::paste! {
            pub async fn $method_name(&self, $($param_name: $param_type),+) -> Result<$return_type, $crate::durable::service::rules::DurableClientError> {
                let service_name = stringify!($service_trait);
                // For objects, we need to include an object key. Using "default" as the key.
                let url = format!("{}/{}/default/{}", self.base_url, service_name, stringify!($method_name));

                // Build request payload with parameters
                let payload = {
                    #[allow(unused_mut)]
                    let mut map = serde_json::Map::new();
                    $(
                        map.insert(stringify!($param_name).to_string(), serde_json::to_value(&$param_name).map_err($crate::durable::service::rules::DurableClientError::SerializationError)?);
                    )*
                    serde_json::Value::Object(map)
                };

                let response = self.client
                    .post(&url)
                    .header("Content-Type", "application/json")
                    .json(&payload)
                    .send()
                    .await
                    .map_err($crate::durable::service::rules::DurableClientError::RequestError)?;

                let status = response.status();
                let text = response.text().await.map_err($crate::durable::service::rules::DurableClientError::RequestError)?;
                if status.is_success() {
                    serde_json::from_str(&text).map_err($crate::durable::service::rules::DurableClientError::SerializationError)
                } else {
                    if !text.is_empty() {
                        Err($crate::durable::service::rules::DurableClientError::ServiceError(text))
                    } else {
                        Err($crate::durable::service::rules::DurableClientError::HttpError(status))
                    }
                }
            }
        }
    };
}
