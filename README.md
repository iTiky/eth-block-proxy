# Ethereum block caching proxy

Application is a caching proxy for `eth_getBlockByNumber` and `eth_blockNumber` Eth JSON RPC calls.
Server uses [CloudFlare Ethereum Gateway](https://developers.cloudflare.com/distributed-web/ethereum-gateway) as a source of truth by default.
An exponential backoff retry policy is used on cache miss.

Proxy has the partial cache invalidation routine for chain fork (reorder) event.
On fork event the predefined number of blocks (20 by default) are invalidated from cache.
As the proxy might not be in sync with the "real" latest block, blocks are invalidated on the left and on the right side of the latest one (by default, 41 in total might be cleaned up).

The cache layer is the LRU-based cache with fixed length (100 by default).

## REST API

By default, server starts at 2412 HTTP port with the following endpoints:

- `/v1/block/latest` - the latest block request;
- `/v1/block/{blockNumber}` - block request by block number [uint64 integer];
- `/v1/block/latest/{txHash}` - Tx request for the latest block by transaction hash [HEX string];
    - `{txHash}` might be prefixed with `0x`.
- `/v1/block/{blockNumber}/{txHash}` - same as above, but for a particular block number;
- `/ping` - request server to response `pong`;
- `/panic` - request server to panic (panic handling debug);

The `v1` (version 1) API response format is identical to [Eth JSON RPC](https://eth.wiki/json-rpc/API).

## Code

### Packages used

App architecture and configuration:

- [viper](https://github.com/spf13/viper) - app configuration with default values and ENVs override;
- [cobra](https://github.com/spf13/cobra) - CLI;
- [golang-lru](github.com/hashicorp/golang-lru) - LRU cache implementation;
- [go-ethereum](github.com/ethereum/go-ethereum) - Ethereum client and Go-structs source;
- [zerolog](github.com/rs/zerolog) - structural logs provider;

Networking:

- [backoff](github.com/cenkalti/backoff/v4) - exponential or no-retry policy implementation for request retry mechanism;
- [go-chi](github.com/go-chi/chi) - HTTP router with convenient middlewares;
- [renderer](github.com/go-chi/render) - HTTP response builder;

Testing:

- [testify](github.com/stretchr/testify) - tests build helper tool with TestSuite support;
- [httpexpect](github.com/gavv/httpexpect/v2) - REST API testing tool;

### Code structure

- `/cmd` - CLI entry point;
- `/app` - the main app dependencies setup and API server start / stop;
- `/api/rest` - the API REST layer with router and middlewares setup:
    - `/api/rest/handlers` - versioned REST API handlers;
- `/provider` - external data providers:
    - `/provider/block` - Ethereum block info providers:
        - `/provider/block/cloudflare` - CloudFlare Ethereum Gateway block provider implementation;
- `/service` - application logic services:
    - `/service/block` - Ethereum blocks services:
        - `/service/block/reader` - Ethereum blocks reader with backoff retry mechanism;
        - `/service/block/notifier` - a new Ethereum block event and chain forked event provider;
    - `/service/cache` - Ethereum blocks caching layer with `/service/block/notifier` events handling;
- `pkg` - common objects and functions;

## CLI

    # Print CLI help (arguments and flags description)
    eth-block-proxy -h

    # Help for a particular command
    eth-block-proxy server -h

    # Print the app version
    eth-block-proxy version

    # Save default values into a config file (to overview or modify values)
    eth-block-proxy default-config --config {config_file_path}

    # Start the proxy server
    eth-block-proxy server --log-level "debug"

## Configuration

### Default config

```toml
[app]
  host = "0.0.0.0"
  serviceport = "2412"

[service]

  [service.block]

    [service.block.readerv1]
      maxretrydur = "500ms"
      minretrydur = "50ms"
      requesttimeoutdur = "5s"

  [service.cache_v1]
    cachesize = 100
    forklength = 20
```

### EVNs override

Each config (or default) value can be overridden using ENV variable.
This might be useful running the app with Docker container and defining particular values rather than providing a full config file.

ENVs prefix - `ETH_BLOCK_PROXY_`.

Example:

    # Redefine server port
    export ETH_BLOCK_PROXY_APP_SERVICEPORT=1234

## How to run

### Binary

    make install
    ${GOPATH}/bin/eth-block-proxy server

The `eth-block-proxy` is installed to the `${GOPATH}/bin` folder.

### Docker

    make build-docker
    docker run -p 2412:2412 eth-block-proxy

## Example queries

    # Ping server
    curl http://127.0.0.1:2412/ping

    # Test panic handling
    curl http://127.0.0.1:2412/panic

    # Request the latest block
    curl http://127.0.0.1:2412/v1/block/latest

    # Request block by number [integer]
    curl http://127.0.0.1:2412/v1/block/12662132

    # Request block's Tx by hash [HEX string with or without "0x" prefix]
    curl http://127.0.0.1:2412/v1/block/12662132/txs/0da33e7d7536844f6d261bc11c349dcb2a715e39f1060fd99bba444e7c577669

## Points of improvement

Monitoring:
- Prometheus metrics collector:
    - Simple errors counter makes it is easy to implement Alerting;
    - Stats on cache hits and misses would help to adjust the settings;
- Sentry.io:
    - Crash and error notifications and reports would help to react faster;
    
Chain fork event criteria:
    - ATM the simplest algo is used and I'm not sure whether it is robust enough;

New block event:
    - ATM the simplest polling mechanism with polling duration readjust is used to get the new blocks;
    - The Cloudflare solution doesn't have (ot I didn't find it) WebSockets, so it is hard to build and event-driven system;

CI/CD:
    - Auto tests on Git push events;
    - Auto Docker build on tag events;
