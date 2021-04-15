package main

import (
	"log"
	"os"
	"time"

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
	log.Printf("Client1 created with address=%s\n", client1.ctx.FromAddress)
	client2 := createChannelClient(Client2Name, kr, contractAddress)
	log.Printf("Client2 created with address=%s\n", client2.ctx.FromAddress)

	client1.PrintBalance()
	client2.PrintBalance()

	asset := sdk.NewInt64Coin(DENOMINATION, 100000)
	receiver := client1.acc.GetAddress()
	operator.transfer(receiver, asset)
	log.Printf("Transfered %v from %s to %s\n", asset, operator.ctx.GetFromAddress(), receiver)

	asset = sdk.NewInt64Coin(DENOMINATION, 10000)
	receiver = client2.acc.GetAddress()
	operator.transfer(receiver, asset)
	log.Printf("Transfered %v from %s to %s\n", asset, operator.ctx.GetFromAddress(), receiver)

	client1.PrintBalance()
	client2.PrintBalance()

	challengeDuration := uint64(5)
	ch := createChannel(client1, client2, challengeDuration)
	log.Printf("Channel created with id=%x\n", ch.ID())

	amount := sdk.NewUint(50000)
	client1.deposit(ch, amount)
	ch.state.Balance[0] = amount
	log.Printf("Funds deposited at channel %x by %s\n", ch.ID(), client1.ctx.GetFromName())

	// amount = sdk.NewUint(10000)
	// client2.deposit(ch, amount)
	// ch.state.Balance[1] = amount
	// log.Printf("Funds deposited at channel %x by %s\n", ch.ID(), client2.ctx.GetFromName())

	client1.PrintBalance()
	client2.PrintBalance()

	transferAmount := sdk.NewUint(30000)
	ch.state.Balance[0] = ch.state.Balance[0].Sub(transferAmount)
	ch.state.Balance[1] = ch.state.Balance[1].Add(transferAmount)

	sig1 := client1.signChannel(ch)
	sig2 := client2.signChannel(ch)
	client1.register(ch, sig1, sig2)
	log.Printf("Channel %x registered by %s\n", ch.ID(), client2.ctx.GetFromName())

	time.Sleep(time.Duration(challengeDuration+10) * time.Second)

	idx := uint16(0)
	sig := client1.signWithdrawal(ch)
	client1.withdraw(ch, idx, sig)
	log.Printf("Withdrawn funds from channel %x by %s\n", ch.ID(), client1.ctx.GetFromName())

	idx = uint16(1)
	sig = client2.signWithdrawal(ch)
	client2.withdraw(ch, idx, sig)
	log.Printf("Withdrawn funds from channel %x by %s\n", ch.ID(), client2.ctx.GetFromName())

	client1.PrintBalance()
	client2.PrintBalance()
}
