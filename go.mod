module perun.network/perun-cosmwasm-golang

go 1.16

require github.com/cosmos/cosmos-sdk v0.42.0

require (
	github.com/CosmWasm/wasmd v0.16.0-alpha1
	github.com/cosmos/go-bip39 v1.0.0
	github.com/ethereum/go-ethereum v1.10.2
	github.com/tendermint/tendermint v0.34.8
	github.com/xeipuuv/gojsonschema v1.2.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
