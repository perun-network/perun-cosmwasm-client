package main

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	ChannelID = [32]byte

	ChannelMemberAddress     = [20]byte
	ChannelNonce             = [32]byte
	ChannelChallengeDuration = uint64

	ChannelParameters struct {
		members           [2]ChannelMemberAddress
		nonce             ChannelNonce
		challengeDuration ChannelChallengeDuration
	}

	ChannelVersion = uint64
	ChannelBalance = sdk.Uint

	ChannelState struct {
		version   ChannelVersion
		balance   [2]ChannelBalance
		finalized bool
	}

	Channel struct {
		params ChannelParameters
		state  ChannelState
	}
)

type ChannelClient struct {
	cosmosClient    *CosmosClient
	contractAddress ContractAddress
}

func createChannelClient(name string, kr keyring.Keyring, contractAddress ContractAddress) *ChannelClient {
	cosmosClient := createCosmosClient(name, kr)
	return &ChannelClient{
		cosmosClient:    cosmosClient,
		contractAddress: contractAddress,
	}
}

func (c *ChannelClient) ChannelMemberAddress() (addr ChannelMemberAddress) {
	h := sha256.Sum256(c.cosmosClient.acc.GetPubKey().Bytes())
	copy(addr[:], h[:])
	return
}

func createChannel(c1, c2 *ChannelClient) *Channel {
	return &Channel{
		params: ChannelParameters{
			members:           [2][20]byte{c1.ChannelMemberAddress(), c2.ChannelMemberAddress()},
			nonce:             ChannelNonce{},
			challengeDuration: 60,
		},
		state: ChannelState{},
	}
}

func (c *Channel) ID() (id ChannelID) {
	return c.params.ID()
}

func (p *ChannelParameters) ID() (id ChannelID) {
	hasher := sha256.New()
	hasher.Write(p.members[0][:])
	hasher.Write(p.members[1][:])
	hasher.Write(p.nonce[:])

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, p.challengeDuration)
	hasher.Write(buf)

	copy(id[:], hasher.Sum(nil))
	return
}
