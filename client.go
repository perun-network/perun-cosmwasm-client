package main

import (
	"os"

	"github.com/CosmWasm/wasmd/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Client struct {
	acc keyring.Info
	ctx client.Context
}

func createClient(name string, kr keyring.Keyring) Client {
	acc, err := kr.Key(name)
	if err != nil {
		acc, err = newAccount(kr, name)
		if err != nil {
			panic(err)
		}
	}
	ctx, err := createClientContext(kr, acc)
	if err != nil {
		panic(err)
	}
	return Client{acc, ctx}
}

func newAccount(kr keyring.Keyring, uid string) (info keyring.Info, err error) {
	algos, _ := kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return
	}

	m, err := newMnemonic()
	if err != nil {
		return
	}

	accountPassword := ""
	hdPath := ""
	return kr.NewAccount(uid, m, accountPassword, hdPath, algo)
}

func newMnemonic() (m string, err error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return
	}
	return bip39.NewMnemonic(entropySeed)
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
		OutputFormat:      "json",
		Height:            0,
		HomeDir:           app.DefaultNodeHome,
		KeyringDir:        "",
		From:              acc.GetName(),
		BroadcastMode:     "block",
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
