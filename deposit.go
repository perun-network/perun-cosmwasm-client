package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	types "github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xeipuuv/gojsonschema"
)

type TokenAmount = sdk.Uint

type DepositMsg struct {
	Channel ChannelID            `json:"channel"`
	Account ChannelMemberAddress `json:"account"`
}

func (c *ChannelClient) deposit(ch *Channel, a TokenAmount) {
	msg, err := genDepositMsg(
		c.cosmosClient.ctx.FromAddress,
		c.contractAddress,
		ch.ID(),
		c.ChannelMemberAddress(),
		a,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Msg))
	validateMessageJSON(msg.Msg)

	r, err := transact(c.cosmosClient.ctx, &msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r.raw))
}

func validateMessageJSON(msg []byte) {
	b, err := ioutil.ReadFile("schema/handle_msg.json")
	if err != nil {
		panic(err)
	}
	schemaLoader := gojsonschema.NewBytesLoader(b)
	documentLoader := gojsonschema.NewBytesLoader(msg)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}

func genDepositMsg(
	sender sdk.AccAddress,
	contract ContractAddress,
	ch ChannelID,
	acc ChannelMemberAddress,
	amount TokenAmount,
) (msg types.MsgExecuteContract, err error) {
	_msg, err := json.Marshal(
		map[string]interface{}{
			"deposit": DepositMsg{
				Channel: ch,
				Account: acc,
			},
		},
	)
	if err != nil {
		return
	}

	msg = types.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: contract,
		Msg:      _msg,
		Funds:    sdk.NewCoins(sdk.NewCoin(DENOMINATION, sdk.Int(amount))),
	}
	return
}
