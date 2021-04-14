package main

import (
	"encoding/json"

	types "github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func instantiateContract(ctx client.Context, codeID CodeID) ContractAddress {
	msg, err := genInstantiateMsg(ctx.GetFromAddress(), codeID)
	if err != nil {
		panic(err)
	}

	r, err := transact(ctx, &msg)
	if err != nil {
		panic(err)
	}

	return NewInstantiateResponse(r).ContractAddress()
}

type InitMsg struct {
	Denom string `json:"denom"`
}

func genInstantiateMsg(sender sdk.AccAddress, codeID CodeID) (msg types.MsgInstantiateContract, err error) {
	initMsg, err := json.Marshal(InitMsg{
		Denom: DENOMINATION,
	})
	if err != nil {
		return
	}

	msg = types.MsgInstantiateContract{
		Sender:  sender.String(),
		CodeID:  codeID,
		Label:   "Perun",
		Funds:   sdk.NewCoins(),
		InitMsg: []byte(initMsg),
		Admin:   "",
	}
	return
}

type InstantiateResponse struct {
	*TxResponse
}

func NewInstantiateResponse(c *TxResponse) *InstantiateResponse {
	return &InstantiateResponse{c}
}

type ContractAddress = string

func (r *InstantiateResponse) ContractAddress() ContractAddress {
	a, ok := r.EventAttributeValue("contract_address")
	if !ok {
		panic("not found")
	}

	return a
}
