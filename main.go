package main

import (
	"log"
	"os"

	"github.com/CosmWasm/wasmd/app"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LOCAL_NODE_DIR = ".wasmd"
	NODE_URL       = "http://localhost:26657" // "https://rpc.musselnet.cosmwasm.com:443"
	CHAIN_ID       = "localnet"               // "musselnet-4"
	DENOMINATION   = "ucosm"

	ContractFilePath = "perun_cosmwasm.wasm"
	OperatorName     = "operator"
	Client1Name      = "alice"
	Client2Name      = "bob"

	KeyringBackend   = "test"
	KeyringDirectory = ".keyring"
	StartNodeFile    = "start_node"

	mnemonicEntropySize = 256
)

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.Seal()

	kr, err := keyring.New("perun", KeyringBackend, KeyringDirectory, os.Stdin)
	if err != nil {
		panic(err)
	}
	log.Println("Keyring created")

	operator := createClient(OperatorName, kr)
	log.Println("Operator created")
	client1 := createClient(Client1Name, kr)
	log.Println("Client1 created")
	client2 := createClient(Client2Name, kr)
	log.Println("Client2 created")

	startNode(operator.acc)
	log.Println("Node started")

	codeID := storeCode(operator.ctx)
	log.Printf("Contract code stored with code_id=%d\n", codeID)

	contractAddress := instantiateContract(operator.ctx, codeID)
	log.Printf("Contract instatiated with address=%s\n", contractAddress)

	// err = deposit(client1, "20ucosm")
	_ = contractAddress
	_, _ = client1, client2
}
