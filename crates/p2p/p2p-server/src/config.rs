use libp2p::{multiaddr::Protocol, Multiaddr};
use serde::{Deserialize, Serialize};
use std::{str::FromStr, time::Duration};
use tracing::error;

#[derive(Debug, Clone, Deserialize, PartialEq, Eq, Serialize)]
#[serde(default)]
pub struct P2pConfig {
    pub port: u16,
    pub idle_connection_timeout_secs: u64,
    pub boot_nodes: Vec<String>,
    pub peers: Option<Vec<String>>,
    pub topic: String,
    pub max_peers_limit: usize,
}

impl P2pConfig {
    pub fn idle_connection_timeout(&self) -> Duration {
        Duration::from_secs(self.idle_connection_timeout_secs)
    }

    /// Converts string address into proper [`Multiaddr`] struct.
    /// Expected peer address format is - /ip4/{ip}/tcp/{port}
    pub fn peer_addresses(&self) -> eyre::Result<Option<Vec<Multiaddr>>> {
        if let Some(peers) = self.peers.clone() {
            let mut multi_addrs = vec![];

            for peer in peers.iter() {
                let addr = Multiaddr::from_str(peer)?;

                // validate peer address, expected format is /ip4/{ip}/tcp/{port}
                let components = addr.iter().collect::<Vec<_>>();
                if components.len() != 2 {
                    error!(target: "p2p", "Invalid peer address format. Expected - /ip4/(ip)/tcp/(port)");
                    continue;
                }

                match components[0] {
                    Protocol::Ip4(_) => (),
                    _ => {
                        error!(target: "p2p", "Invalid first multiaddr part. Expected to be ip4");
                        continue;
                    }
                };

                match components[1] {
                    Protocol::Tcp(_) => (),
                    _ => {
                        error!(target: "p2p", "Invalid second multiaddr part. Expected to be tcp");
                        continue;
                    }
                };

                // address is in a valid format, store it for future connection attempt
                multi_addrs.push(addr);
            }

            Ok(Some(multi_addrs))
        } else {
            Ok(None)
        }
    }
}

impl Default for P2pConfig {
    fn default() -> Self {
        Self {
            port: 1211,
            idle_connection_timeout_secs: 60,
            boot_nodes: vec![
                "QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN".to_owned(), // TODO: set default values to our node once done
                "QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa".to_owned(),
                "QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb".to_owned(),
                "QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt".to_owned(),
            ],
            peers: None,
            topic: "ramd-topic".to_owned(),
            max_peers_limit: 1,
        }
    }
}
