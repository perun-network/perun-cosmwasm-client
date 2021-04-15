package main

import (
	"crypto/sha256"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type (
	ChannelID = [32]byte

	ChannelMemberAddress     = [20]byte
	ChannelNonce             = [32]byte
	ChannelChallengeDuration = uint64

	ChannelParameters struct {
		Members           [2]ChannelMemberAddress  `json:"participants"`
		Nonce             ChannelNonce             `json:"nonce"`
		ChallengeDuration ChannelChallengeDuration `json:"challenge_duration"`
	}

	ChannelVersion = uint64
	ChannelBalance = sdk.Uint

	ChannelState struct {
		Version   ChannelVersion    `json:"version"`
		Balance   [2]ChannelBalance `json:"balance"`
		Finalized bool              `json:"finalized"`
	}

	Channel struct {
		params ChannelParameters
		state  ChannelState
	}
)

type ChannelClient struct {
	*CosmosClient
	contractAddress ContractAddress
}

func createChannelClient(name string, kr keyring.Keyring, contractAddress ContractAddress) *ChannelClient {
	cosmosClient := createCosmosClient(name, kr)
	return &ChannelClient{
		CosmosClient:    cosmosClient,
		contractAddress: contractAddress,
	}
}

func (c *ChannelClient) ChannelMemberAddress() (addr ChannelMemberAddress) {
	pub := c.acc.GetPubKey().Bytes()
	pubDecompressed, err := crypto.DecompressPubkey(pub)
	if err != nil {
		panic(err)
	}

	h := sha256.Sum256(crypto.FromECDSAPub(pubDecompressed))
	copy(addr[:], h[:])
	return
}

func createChannel(c1, c2 *ChannelClient, challengeDuration uint64) *Channel {
	return &Channel{
		params: ChannelParameters{
			Members:           [2][20]byte{c1.ChannelMemberAddress(), c2.ChannelMemberAddress()},
			Nonce:             ChannelNonce{},
			ChallengeDuration: challengeDuration,
		},
		state: ChannelState{
			Version:   0,
			Balance:   [2]sdk.Uint{sdk.NewUint(0), sdk.NewUint(0)},
			Finalized: false,
		},
	}
}

func (c *Channel) ID() (id ChannelID) {
	copy(id[:], c.params.hash())
	return
}

func (p *ChannelParameters) hash() []byte {
	hasher := sha256.New()
	hasher.Write(p.Members[0][:])
	hasher.Write(p.Members[1][:])
	hasher.Write(p.Nonce[:])
	hasher.Write(sdk.Uint64ToBigEndian(p.ChallengeDuration))
	return hasher.Sum(nil)
}

func (s *ChannelState) hash() []byte {
	hasher := sha256.New()
	hasher.Write(sdk.Uint64ToBigEndian(s.Version))
	buf := make([]byte, 16)
	s.Balance[0].BigInt().FillBytes(buf)
	hasher.Write(buf)
	buf = make([]byte, 16)
	s.Balance[1].BigInt().FillBytes(buf)
	hasher.Write(buf)

	b := func() byte {
		if s.Finalized {
			return 1
		} else {
			return 0
		}
	}()
	hasher.Write([]byte{b})

	return hasher.Sum(nil)
}

func (c *Channel) hash() []byte {
	hasher := sha256.New()
	hasher.Write(c.params.hash())
	hasher.Write(c.state.hash())
	return hasher.Sum(nil)
}
