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

	operator := createCosmosClient(OperatorName, kr)
	log.Printf("Operator created with address=%s\n", operator.ctx.FromAddress)

	startNode(operator.acc)
	log.Println("Node started")

	codeID := storeCode(operator.ctx)
	log.Printf("Contract code stored with code_id=%d\n", codeID)

	contractAddress := instantiateContract(operator.ctx, codeID)
	log.Printf("Contract instatiated with address=%s\n", contractAddress)

	client1 := createChannelClient(Client1Name, kr, contractAddress)
	log.Printf("Client1 created with address=%s\n", client1.cosmosClient.ctx.FromAddress)
	client2 := createChannelClient(Client2Name, kr, contractAddress)
	log.Printf("Client2 created with address=%s\n", client2.cosmosClient.ctx.FromAddress)

	amount := sdk.NewInt64Coin(DENOMINATION, 100000)
	receiver := client1.cosmosClient.acc.GetAddress()
	operator.transfer(receiver, amount)
	log.Printf("Transfered %v from %s to %s\n", amount, operator.ctx.GetFromAddress(), receiver)

	amount = sdk.NewInt64Coin(DENOMINATION, 100000)
	receiver = client2.cosmosClient.acc.GetAddress()
	operator.transfer(receiver, amount)
	log.Printf("Transfered %v from %s to %s\n", amount, operator.ctx.GetFromAddress(), receiver)

	ch := createChannel(client1, client2)
	log.Printf("Channel created with id=%x\n", ch.ID())

	client1.deposit(ch, sdk.NewInt(10))
	log.Printf("Channel funds deposited by %s\n", client1.cosmosClient.ctx.GetFromName())

	client2.deposit(ch, sdk.NewInt(10))
	log.Printf("Channel funds deposited by %s\n", client2.cosmosClient.ctx.GetFromName())
}
