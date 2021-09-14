# Perun CosmWasm Study: Client

This repository contains a proof-of-concept implementation of a client for the Perun CosmWasm contracts.

The repository uses the contracts from [perun-cosmwasm-study_contracts](https://github.com/perun-network/perun-cosmwasm-study_contracts).

## Organization

The source code is organized as follows.

* `perun_cosmwasm.wasm`: The precompiled contract.

* `main.go`: This file contains a test that demonstrates the functionality of the client.

* `client.go`: This file contains basic client functionality, like creating a new client, transferring coins, and generating signatures.

* `node.go`: This file contains functionality for starting a local cosmwasm blockchain node.

* `store.go`: This file contains functions for storing contract code on the blockchain.

* `instantiate.go`: This file contains functions for instantiating a contract on the blockchain. Storing and instantiating a contract is a two-step process. A stored contract can be instantiated multiple times.

* `channel.go`: This file contains the basic types and functions for channels.

* `deposit.go`: This file contains functions for depositing funds into a channel.

* `register.go`: This file contains functions for registering a channel state on the blockchain.

* `withdraw.go`: This file contains functions withdrawing funds from a channel.

## Execution

To run the test locally, you need to have [wasmd v0.16.0](https://github.com/CosmWasm/wasmd/tree/v0.16.0) installed. You can also run the test via a remote node. The configuration can be found in file `main.go`.

To execute the test, run
```bash
go run .
```

The test will generate a script `start_node`, which when executed starts a local cosmwasm blockchain node with prefunded accounts. Once the node is running, press a key to signal that the test can continue running.

## License

The source code is published under Apache License Version 2.0.
