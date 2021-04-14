package main

import (
	"encoding/json"
	"fmt"

	types "github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenAmount = sdk.Int

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

	r, err := transact(c.cosmosClient.ctx, &msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r.raw))
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
