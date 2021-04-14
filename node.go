package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

func startNode(acc keyring.Info) {
	script := fmt.Sprintf(`set -xe

APP_HOME="%s"
CHAIN_ID="%s"
OWNER="%s"
OWNER_NAME="%s"
KEYRING_BACKEND="%s"
KEYRING_DIR="%s"
DENOMINATION="%s"

# initialize wasmd configuration files
wasmd init localnet --chain-id ${CHAIN_ID} --home ${APP_HOME}

# add minimum gas prices config to app configuration file
sed -i -r "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.01$DENOMINATION\"/" ${APP_HOME}/config/app.toml

# add your wallet addresses to genesis
wasmd add-genesis-account $OWNER 10000000000$DENOMINATION,10000000000stake --home ${APP_HOME}

# add fred's address as validator's address
wasmd gentx $OWNER_NAME 1000000000stake --home ${APP_HOME} --chain-id ${CHAIN_ID} --keyring-backend $KEYRING_BACKEND --keyring-dir $KEYRING_DIR

# collect gentxs to genesis
wasmd collect-gentxs --home ${APP_HOME}

# validate the genesis file
wasmd validate-genesis --home ${APP_HOME}

# run the node
wasmd start --home ${APP_HOME}`, LOCAL_NODE_DIR, CHAIN_ID, acc.GetAddress(), acc.GetName(), KeyringBackend, KeyringDirectory, DENOMINATION)

	func() {
		f, err := os.Create(StartNodeFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(script)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Run 'sh %s' to start a local node.\nOnce the node is running, press enter to continue.", StartNodeFile)
	_, err := os.Stdin.Read([]byte{0})
	if err != nil {
		panic(err)
	}
}
