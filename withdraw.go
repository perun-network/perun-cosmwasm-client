package main

import (
	"crypto/sha256"
	"encoding/json"

	types "github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WithdrawMsg struct {
	Params       ChannelParameters `json:"params"`
	AccountIndex uint16            `json:"account_index"`
	Signature    Signature         `json:"sig"`
}

func (c *ChannelClient) withdraw(ch *Channel, idx uint16, sig Signature) {
	msg, err := genWithdrawMsg(
		c.ctx.FromAddress,
		c.contractAddress,
		ch,
		idx,
		sig,
	)
	if err != nil {
		panic(err)
	}

	validateMessageJSON(msg.Msg)

	_, err = transact(c.ctx, &msg)
	if err != nil {
		panic(err)
	}
}

func genWithdrawMsg(
	sender sdk.AccAddress,
	contract ContractAddress,
	ch *Channel,
	idx uint16,
	sig Signature,
) (msg types.MsgExecuteContract, err error) {
	_msg, err := json.Marshal(
		map[string]interface{}{
			"withdraw": WithdrawMsg{
				Params:       ch.params,
				AccountIndex: idx,
				Signature:    sig,
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
		Funds:    nil,
	}
	return
}

func (c *ChannelClient) signWithdrawal(ch *Channel) Signature {
	hasher := sha256.New()
	hasher.Write(ch.params.hash())
	return c.signHash(hasher.Sum(nil))
}
