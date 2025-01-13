# Eth peer manager

Eth peer manager is a sidecar service that uses the `admin` namespace api calls to check peer availability.
If peer is unreachable it will be removed. There also is an additional feature that provides peer injection functionality.

## Configuration
Configuration management is done by Viper and ENV vars with `EPM_` prefix. 

Allowed ENV vars:
- EPM_GETH_API - URL to connect to the Geth/Erigon API: `http://127.0.0.1:8545`
- EPM_LOG_LEVEL - the default level is `info`
- EPM_PROBE_TIMEOUT - the amount of time in seconds to wait for the peer TCP response 
- EPM_RUN_INTERVAL - the amount of time to wait between runs
- EPM_ADDITIONAL_PEERS - string or list as a string of peers separated by `,` to add on each run

## Configuration
How to run locally

It is very important to add these flags to your execution client.
```sh
--http.api web3,eth,txpool,net,engine,admin
--http.port=8545
```
You need to add the admin api so we can remove bad peers.

```sh
docker run --name eth-peer-manager \
  --restart unless-stopped \
  --stop-timeout 300 \
  --network "{{ docker_network_name }}" \
  -p 127.0.0.1:9093:9090 \
  -e EPM_GETH_API=http://{{ berachain_execution_client_name }}:{{ berachain_execution_rpc_port }} \
  "{{ eth_peer_manager_image }}"
```

This can run on a docker network or on local host.
```sh
docker run --name eth-peer-manager \
  --restart unless-stopped \
  --stop-timeout 300 \
  --network host \
  -e EPM_GETH_API=http://127.0.0.1:8545 \
  "{{ eth_peer_manager_image }}"
```