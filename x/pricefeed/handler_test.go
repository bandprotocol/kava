package pricefeed_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/pricefeed"
	"github.com/kava-labs/kava/x/pricefeed/types"
)

func parsePK(pkStr string) secp256k1.PubKeySecp256k1 {
	bytes, _ := hex.DecodeString(pkStr)
	var pk secp256k1.PubKeySecp256k1
	copy(pk[:], bytes[:])
	return pk
}

func hexToBytes(str string) []byte {
	bz, _ := hex.DecodeString(str)
	return bz
}

func TestHandleReportPrice(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, abci.Header{})
	keeper := tApp.GetPriceFeedKeeper()

	mp := types.Params{
		Markets: types.Markets{
			types.Market{MarketID: "btc:usd", BaseAsset: "btc", QuoteAsset: "usd", Oracles: []sdk.AccAddress{}, Active: true},
		},
	}
	keeper.SetParams(ctx, mp)

	pubkeys := []secp256k1.PubKeySecp256k1{
		parsePK("03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909"),
		parsePK("02724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885"),
		parsePK("03a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33"),
		parsePK("03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f"),
	}

	keeper.SetValidatorPubkeys(ctx, pubkeys)

	proof := types.Proof{
		BlockHeight: 4065,
		OracleDataProof: types.OracleDataProof{
			Version:   3612,
			RequestID: 5,
			CodeHash:  hexToBytes("6B7BE61B150AEC5EB853AFB3B53E41438959554580D31259A1095E51645BCD28"),
			Params:    hexToBytes("00000004"),
			Data:      hexToBytes("00000000000EEDB0000000005E53CCB1"),
			MerklePaths: []types.IAVLMerklePath{
				types.IAVLMerklePath{
					IsDataOnRight:  false, // isDataOnRight
					SubtreeHeight:  1,     // subtreeHeight
					SubtreeSize:    2,     // subtreeSize
					SubtreeVersion: 4034,  // subtreeVersion
					SiblingHash:    hexToBytes("785071F43D33029DAFA723F065C0D92B00D31CB36CBC2E00D2A379D4D2BECFE2"),
				},
				types.IAVLMerklePath{
					IsDataOnRight:  true, // isDataOnRight
					SubtreeHeight:  2,    // subtreeHeight
					SubtreeSize:    3,    // subtreeSize
					SubtreeVersion: 4034, // subtreeVersion
					SiblingHash:    hexToBytes("C5A20E25BA2D930132447F1D6C15311288D5789DBDF90A8841400CB04E4E9DCF"),
				},
				types.IAVLMerklePath{
					IsDataOnRight:  true, // isDataOnRight
					SubtreeHeight:  3,    // subtreeHeight
					SubtreeSize:    5,    // subtreeSize
					SubtreeVersion: 4034, // subtreeVersion
					SiblingHash:    hexToBytes("E8FAAA55A02072B13B2FA780EAC32F17AF4FC7C7091F9AC295905EB8EA33BB46"),
				},
				types.IAVLMerklePath{
					IsDataOnRight:  true, // isDataOnRight
					SubtreeHeight:  4,    // subtreeHeight
					SubtreeSize:    11,   // subtreeSize
					SubtreeVersion: 4034, // subtreeVersion
					SiblingHash:    hexToBytes("48FAF5FF10D2FD97BD16A8472F7FE63EAD3EEBCDB644CE54AEE200041FA31FE3"),
				},
				types.IAVLMerklePath{
					IsDataOnRight:  true,
					SubtreeHeight:  5,
					SubtreeSize:    19,
					SubtreeVersion: 4034,
					SiblingHash:    hexToBytes("E120EC735159E95B7929B5EC855508C665FF68446CCAE2C78249BA4325260500"),
				},
				types.IAVLMerklePath{
					IsDataOnRight:  true,
					SubtreeHeight:  6,
					SubtreeSize:    30,
					SubtreeVersion: 4064,
					SiblingHash:    hexToBytes("1E14724C694844290433D09125B44B7C0C72024A8CBD061CA314B1733D571E98"),
				},
			},
		},
		BlockRelayProof: types.BlockRelayProof{
			hexToBytes("08A8B029E86AEFC31F86349B5694F7EBF242971D8549DFEAAC68B083502092DB"),
			hexToBytes("10A94B07B203521A9C0E14AAF1C720EC0B3F53DD86C30BA2C8AC7CD19D1804C1"),
			hexToBytes("70252288AF697C727A6A3670549B4B4EF99748F5D28D8FAFBBC28414B3D0A41D"),
			types.BlockHeaderMerkleParts{
				hexToBytes("32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E"),
				hexToBytes("E19CAD181721A668C0D37758DF915ED7A0872D1FDD53827436D0E511EFB20085"),
				hexToBytes("67771F6C96DCF89D7FCEE28304CB37F92FBA6DC27747C5594E81C01F7138984B"),
				hexToBytes("B6F005C12B317FC853AEAE6851352D96343A31F79D46E63204A6470C5E934060"),
				hexToBytes("6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D"),
				hexToBytes("0CF1E6ECE60E49D19BB57C1A432E805F39BB4F65C366741E4F03FA54FBD90714"),
			},
			hexToBytes(""),
			[]types.TMSignature{
				types.TMSignature{
					hexToBytes("447DB67B461EF5AD8245F0B63A4EF5F2E0AD87758F0B0FA14EA681713B92BD3D"), // r
					hexToBytes("777BD0A0BBB4A8A7BB34B5608E99E0D52D20F0294AEEE641319DE019C1444DB0"), // s
					27, // v
					hexToBytes("12240A204207677E4D9D3CE97C3904D4705BC4D31ADF5E7F2B68FAAAB21DE25A6ED80A9910012A0C088A9DCFF20510DEC7D1BA01320962616E64636861696E"), // _signedDataSuffix
				},
			},
		},
	}
	msg := types.NewMsgReportPrice(sdk.AccAddress("relayer"), proof)
	got := pricefeed.HandleMsgReportPrice(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set report to be ok, got %v", got)

	price, err := keeper.GetCurrentPrice(ctx, "btc:usd")
	require.Nil(t, err)
	expect, _ := sdk.NewDecFromStr("9783.52")
	require.Equal(t, expect, price.Price)
}
