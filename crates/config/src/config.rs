use std::path::PathBuf;

use ramd_db::config::RocksConfig;
use ramd_jsonrpc_server::config::JsonRpcServerConfig;
use ramd_node::NodeConfig;
use ramd_tracing::config::TracingConfig;
use serde::{Deserialize, Serialize};

/// Directory path for storing all ramd related data
const RAMD_DIR: &str = ".ramd";

/// Directory path for storing ramd config information
const CONFIG_DIR: &str = "config";

const CONFIG_FIILE: &str = "ramd.toml";

/// This struct gathers all config values used across ramd node
#[derive(Debug, Clone, Default, Deserialize, PartialEq, Eq, Serialize)]
pub struct RamdConfig {
    /// Configuration for RAM node
    pub node: NodeConfig,
    /// Configuration for tracing/logging
    pub tracing: TracingConfig,
    /// Configuration for jsonrpc server
    pub json_rpc: JsonRpcServerConfig,
}

impl RamdConfig {
    /// Reads config from default path, returns error if config doesn't exists
    pub fn read() -> eyre::Result<Self> {
        let home_path = std::env::var("HOME")?;
        let ramd_config_path: PathBuf = [home_path.as_str(), RAMD_DIR, CONFIG_DIR, CONFIG_FIILE]
            .iter()
            .collect();

        let config = std::fs::read_to_string(ramd_config_path)
            .map_err(|_| eyre::eyre!("Path doesn't exist"))?;

        let config: RamdConfig = toml::from_str(&config)?;
        Ok(config)
    }

    /// Creates default config if not exists otherwise reads it
    pub fn init_or_read() -> eyre::Result<Self> {
        let config_maybe = RamdConfig::read();
        if let Ok(config) = config_maybe {
            return Ok(config);
        };

        let home_path = std::env::var("HOME")?;

        // create ramd root directory
        let root_dir: PathBuf = [home_path.as_str(), RAMD_DIR].iter().collect();
        std::fs::create_dir_all(&root_dir)?;

        // instantiate ramd config
        let config = RamdConfig {
            node: NodeConfig {
                rocks: RocksConfig::new(root_dir.clone()),
            },
            tracing: TracingConfig::new(root_dir.clone()),
            ..Default::default()
        };

        // create directory to store ramd config
        let config_dir = root_dir.join(CONFIG_DIR);
        std::fs::create_dir(&config_dir)?;

        // create directory for database
        std::fs::create_dir(&config.node.rocks.path)?;

        // create directory for logs
        std::fs::create_dir(&config.tracing.path)?;

        // store initial config values
        let config_path = config_dir.join(CONFIG_FIILE);

        let toml_config = toml::to_string(&config)?;
        std::fs::write(config_path, toml_config)?;

        Ok(config)
    }
}
