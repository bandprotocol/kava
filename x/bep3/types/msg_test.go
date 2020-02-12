package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var (
	coinsSingle  = sdk.Coins{sdk.Coin{Denom: "bnb", Amount: sdk.NewInt(50000)}}
	coinsZero    = sdk.Coins{sdk.Coin{}}
	binanceAddrs = []sdk.AccAddress{
		sdk.AccAddress(crypto.AddressHash([]byte("BinanceTest1"))),
		sdk.AccAddress(crypto.AddressHash([]byte("BinanceTest2"))),
	}
	kavaAddrs = []sdk.AccAddress{
		sdk.AccAddress(crypto.AddressHash([]byte("KavaTest1"))),
		sdk.AccAddress(crypto.AddressHash([]byte("KavaTest2"))),
	}
	randomNumberBytes = []byte{15}
	timestampInt64    = int64(9988776655)
	randomNumberHash  = CalculateRandomHash(randomNumberBytes, timestampInt64)
	ethAddrs          = []common.Address{
		common.HexToAddress("0x6f456B7F0b1658Be2683375159E7f09a8831CBe5"),
		common.HexToAddress("0x3a6CEef76Fd677332Dc0bA09604bD6acB1BeF613"),
	}
)

func TestHTLTMsg(t *testing.T) {
	tests := []struct {
		description         string
		from                sdk.AccAddress
		to                  sdk.AccAddress
		recipientOtherChain string
		senderOtherChain    string
		randomNumberHash    SwapBytes
		timestamp           int64
		amount              sdk.Coins
		expectedIncome      string
		heightSpan          int64
		crossChain          bool
		expectPass          bool
	}{
		{"create htlt", binanceAddrs[0], kavaAddrs[0], "", "", randomNumberHash, timestampInt64, coinsSingle, "bnb50000", 80000, false, true},
		{"create htlt cross-chain", binanceAddrs[0], kavaAddrs[0], kavaAddrs[0].String(), binanceAddrs[0].String(), randomNumberHash, timestampInt64, coinsSingle, "bnb50000", 80000, true, true},
		{"create htlt with other chain fields", binanceAddrs[0], kavaAddrs[0], kavaAddrs[0].String(), binanceAddrs[0].String(), randomNumberHash, timestampInt64, coinsSingle, "bnb50000", 80000, false, false},
		{"create htlt cross-cross no other chain fields", binanceAddrs[0], kavaAddrs[0], "", "", randomNumberHash, timestampInt64, coinsSingle, "bnb50000", 80000, true, false},
		{"create htlt zero coins", binanceAddrs[0], kavaAddrs[0], "", "", randomNumberHash, timestampInt64, coinsZero, "bnb50000", 80000, true, false},
	}

	for i, tc := range tests {
		msg := NewHTLTMsg(
			tc.from,
			tc.to,
			tc.recipientOtherChain,
			tc.senderOtherChain,
			tc.randomNumberHash,
			tc.timestamp,
			tc.amount,
			tc.expectedIncome,
			tc.heightSpan,
			tc.crossChain,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgDepositHTLT(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		swapID      SwapBytes
		amount      sdk.Coins
		expectPass  bool
	}{
		{"deposit htlt", binanceAddrs[0], CalculateSwapID(randomNumberHash, binanceAddrs[0], ""), coinsSingle, true},
	}

	for i, tc := range tests {
		msg := NewMsgDepositHTLT(
			tc.from,
			tc.swapID,
			tc.amount,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgClaimHTLT(t *testing.T) {
	tests := []struct {
		description  string
		from         sdk.AccAddress
		swapID       SwapBytes
		randomNumber SwapBytes
		expectPass   bool
	}{
		{"claim htlt", binanceAddrs[0], CalculateSwapID(randomNumberHash, binanceAddrs[0], ""), randomNumberHash, true},
	}

	for i, tc := range tests {
		msg := NewMsgClaimHTLT(
			tc.from,
			tc.swapID,
			tc.randomNumber,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgRefundHTLT(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		swapID      SwapBytes
		expectPass  bool
	}{
		{"claim htlt", binanceAddrs[0], CalculateSwapID(randomNumberHash, binanceAddrs[0], ""), true},
	}

	for i, tc := range tests {
		msg := NewMsgRefundHTLT(
			tc.from,
			tc.swapID,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgCalculateSwapID(t *testing.T) {
	tests := []struct {
		description      string
		from             sdk.AccAddress
		randomNumberHash []byte
		sender           string
		senderOtherChain string
		expectPass       bool
	}{
		{"normal", kavaAddrs[0], CalculateRandomHash(randomNumberBytes, timestampInt64), ethAddrs[0].String(), ethAddrs[1].String(), true},
		{"no sender", kavaAddrs[0], CalculateRandomHash(randomNumberBytes, timestampInt64), "", ethAddrs[1].String(), false},
		{"short hash", kavaAddrs[0], []byte("hash_is_too_short"), "", ethAddrs[1].String(), false},
	}

	for i, tc := range tests {
		msg := NewMsgCalculateSwapID(
			tc.from,
			tc.randomNumberHash,
			tc.sender,
			tc.senderOtherChain,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}