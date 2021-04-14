package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

type TxResponse struct {
	raw []byte
}

func transact(ctx client.Context, msg sdk.Msg) (r *TxResponse, err error) {
	tf, err := createTransactionFactory(ctx)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	ctx.OutputFormat = "json"
	ctx.Output = &buf
	err = tx.BroadcastTx(ctx, tf, msg)
	if err != nil {
		return
	}
	return &TxResponse{raw: buf.Bytes()}, nil
}

func createTransactionFactory(ctx client.Context) (tf tx.Factory, err error) {
	return tx.Factory{}.
			WithTxConfig(ctx.TxConfig).
			WithAccountRetriever(ctx.AccountRetriever).
			WithKeybase(ctx.Keyring).
			WithChainID(CHAIN_ID).
			WithGas(0).
			WithSimulateAndExecute(true).
			WithAccountNumber(0).
			WithSequence(0).
			WithTimeoutHeight(0).
			WithGasAdjustment(1.3).
			WithMemo("").
			WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
			WithFees("").
			WithGasPrices("0.01" + DENOMINATION),
		nil
}

func (r *TxResponse) EventAttributeValue(key string) (string, bool) {
	fmt.Println(string(r.raw))
	var response map[string]interface{}
	json.Unmarshal(r.raw, &response)
	logs := response["logs"].([]interface{})[0].(map[string]interface{})
	events := logs["events"].([]interface{})[0].(map[string]interface{})
	attributes := events["attributes"].([]interface{})
	for _, a := range attributes {
		a := a.(map[string]interface{})
		if a["key"] == key {
			return a["value"].(string), true
		}
	}
	return "", false
}
