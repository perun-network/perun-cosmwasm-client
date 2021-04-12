package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/CosmWasm/wasmd/app"
	types "github.com/CosmWasm/wasmd/x/wasm"
	wasmUtils "github.com/CosmWasm/wasmd/x/wasm/client/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	LOCAL_NODE_DIR = ".wasmd"
	NODE_URL       = "http://localhost:26657" // "https://rpc.musselnet.cosmwasm.com:443"
	CHAIN_ID       = "localnet"               // "musselnet-4"
	WASMD_VERSION  = "v0.16.0-alpha"          // "v0.15.1"

	GAS_PRICES     = "0.01ucosm"
	GAS_SETTING    = "auto"
	GAS_ADJUSTMENT = 1.3

	OUTPUT_FORMAT = "text" // | json

	ContractFilePath = "perun_cosmwasm.wasm"
	OperatorName     = "operator"
	OperatorPassword = ""

	KeyringBackend   = "test"
	KeyringDirectory = ".keyring"
	START_NODE_FILE  = "start_node"

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

	acc, err := newAccount(kr, OperatorName, OperatorPassword)
	if err != nil {
		panic(err)
	}
	log.Println("Account created")

	startNode(acc)
	log.Println("Node started")

	ctx, err := createClientContext(kr, acc)
	if err != nil {
		panic(err)
	}
	log.Println("Client context created")

	err = deploy(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Contract deployed")
}

func startNode(acc keyring.Info) {
	script := fmt.Sprintf(`set -xe

APP_HOME="%s"
CHAIN_ID="%s"
OWNER="%s"
OWNER_NAME="%s"
KEYRING_BACKEND="%s"
KEYRING_DIR="%s"

# initialize wasmd configuration files
wasmd init localnet --chain-id ${CHAIN_ID} --home ${APP_HOME}

# add minimum gas prices config to app configuration file
sed -i -r 's/minimum-gas-prices = ""/minimum-gas-prices = "0.01ucosm"/' ${APP_HOME}/config/app.toml

# add your wallet addresses to genesis
wasmd add-genesis-account $OWNER 10000000000ucosm,10000000000stake --home ${APP_HOME}

# add fred's address as validator's address
wasmd gentx $OWNER_NAME 1000000000stake --home ${APP_HOME} --chain-id ${CHAIN_ID} --keyring-backend $KEYRING_BACKEND --keyring-dir $KEYRING_DIR

# collect gentxs to genesis
wasmd collect-gentxs --home ${APP_HOME}

# validate the genesis file
wasmd validate-genesis --home ${APP_HOME}

# run the node
wasmd start --home ${APP_HOME}`, LOCAL_NODE_DIR, CHAIN_ID, acc.GetAddress(), acc.GetName(), KeyringBackend, KeyringDirectory)

	func() {
		f, err := os.Create(START_NODE_FILE)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(script)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Run 'sh %s' to start a local node.\nOnce the node is running, press enter to continue.", START_NODE_FILE)
	os.Stdin.Read([]byte{0})
}

func newAccount(kr keyring.Keyring, uid, pwd string) (info keyring.Info, err error) {
	algos, _ := kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return
	}

	m, err := newMnemonic()
	if err != nil {
		return
	}

	hdPath := ""
	return kr.NewAccount(uid, m, pwd, hdPath, algo)
}

func newMnemonic() (m string, err error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return
	}
	return bip39.NewMnemonic(entropySeed)
}

func deploy(ctx client.Context) (err error) {
	msg, err := genStoreMsg(ContractFilePath, ctx.FromAddress)
	if err != nil {
		return errors.WithMessage(err, "generating store message")
	}

	tf, err := createTransactionFactory(ctx)
	if err != nil {
		return errors.WithMessage(err, "creating transaction factory")
	}

	return tx.BroadcastTx(ctx, tf, &msg)
}

func createTransactionFactory(ctx client.Context) (tf tx.Factory, err error) {
	gasSetting, err := flags.ParseGasSetting(GAS_SETTING)
	if err != nil {
		return
	}

	return tx.Factory{}.
			WithTxConfig(ctx.TxConfig).
			WithAccountRetriever(ctx.AccountRetriever).
			WithKeybase(ctx.Keyring).
			WithChainID(CHAIN_ID).
			WithGas(gasSetting.Gas).
			WithSimulateAndExecute(gasSetting.Simulate).
			WithAccountNumber(0).
			WithSequence(0).
			WithTimeoutHeight(0).
			WithGasAdjustment(GAS_ADJUSTMENT).
			WithMemo("").
			WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
			WithFees("").
			WithGasPrices(GAS_PRICES),
		nil
}

func createClientContext(kr keyring.Keyring, acc keyring.Info) (ctx client.Context, err error) {
	httpClient, err := rpchttp.New(NODE_URL, "/websocket")
	if err != nil {
		return
	}

	encodingConfig := app.MakeEncodingConfig()

	return client.Context{
		FromAddress:       acc.GetAddress(),
		Client:            httpClient,
		ChainID:           CHAIN_ID,
		JSONMarshaler:     encodingConfig.Marshaler,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Input:             os.Stdin,
		Keyring:           kr,
		Output:            os.Stdout,
		OutputFormat:      OUTPUT_FORMAT,
		Height:            0,
		HomeDir:           app.DefaultNodeHome,
		KeyringDir:        "",
		From:              acc.GetName(),
		BroadcastMode:     "sync",
		FromName:          acc.GetName(),
		SignModeStr:       "",
		UseLedger:         false,
		Simulate:          false,
		GenerateOnly:      false,
		Offline:           false,
		SkipConfirm:       true,
		TxConfig:          encodingConfig.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		NodeURI:           NODE_URL,
	}, nil
}

func genStoreMsg(file string, sender sdk.AccAddress) (types.MsgStoreCode, error) {
	wasm, err := ioutil.ReadFile(file)
	if err != nil {
		return types.MsgStoreCode{}, err
	}

	// gzip the wasm file
	if wasmUtils.IsWasm(wasm) {
		wasm, err = wasmUtils.GzipIt(wasm)

		if err != nil {
			return types.MsgStoreCode{}, err
		}
	} else if !wasmUtils.IsGzip(wasm) {
		return types.MsgStoreCode{}, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	// Optional: allow contract instantiation only by specified address
	// var perm *types.AccessConfig
	// x := types.AccessTypeOnlyAddress.With(allowedAddr)
	// perm = &x

	// Optional: reference source code by URI and docker build system by tag
	// var source = ""
	// var builder = ""

	msg := types.MsgStoreCode{
		Sender:       sender.String(),
		WASMByteCode: wasm,
		// Optional
		// Source:                source,
		// Builder:               builder,
		// InstantiatePermission: perm,
	}
	return msg, nil
}
