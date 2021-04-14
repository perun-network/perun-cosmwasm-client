package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	types "github.com/CosmWasm/wasmd/x/wasm"
	wasmUtils "github.com/CosmWasm/wasmd/x/wasm/client/utils"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeID = uint64

func storeCode(ctx client.Context) CodeID {
	msg, err := genStoreMsg(ContractFilePath, ctx.FromAddress)
	if err != nil {
		panic(err)
	}

	r, err := transact(ctx, &msg)
	if err != nil {
		panic(err)
	}

	return NewStoreResponse(r).CodeID()
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

type StoreResponse struct {
	*TxResponse
}

func NewStoreResponse(c *TxResponse) *StoreResponse {
	return &StoreResponse{c}
}

func (r *StoreResponse) CodeID() uint64 {
	v, ok := r.EventAttributeValue("code_id")
	if !ok {
		panic("not found")
	}

	codeID, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		panic(err)
	}

	return codeID
}
