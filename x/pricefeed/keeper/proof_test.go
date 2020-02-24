package keeper_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/kava-labs/kava/app"
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

func TestRelayOracleState(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, abci.Header{})
	keeper := tApp.GetPriceFeedKeeper()

	pubkeys := []secp256k1.PubKeySecp256k1{
		parsePK("03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909"),
		parsePK("02724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885"),
		parsePK("03a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33"),
		parsePK("03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f"),
	}

	keeper.SetValidatorPubkeys(ctx, pubkeys)

	blockRelayProof := types.BlockRelayProof{
		hexToBytes("ACAB016EC9FB3AA28A6A4BE8A364AEDAA9A42866E2957C5C267E340CE67C55EE"),
		hexToBytes("6ABB9CA0E0AC77A3B7C7F94D56E181DE954B92A19389829CE0E5A95B74BE0B7D"),
		hexToBytes("2726178BFFB0D462C15AB546DE7B4CA86588A98FF0F629DB7CA7E318AA61A846"),
		types.BlockHeaderMerkleParts{
			hexToBytes("32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E"),
			hexToBytes("2157913D927F4249C52FAB326E9E0E83FACFAF167FA038A88173FA42ADF2452C"),
			hexToBytes("7EECC6A0EE0136DE143C92370E4BE8FA6F545C02C23DAFA62CC4AA0A14701787"),
			hexToBytes("DEF482CDA986470C27374601EC716E9853DE47D72828AE0131CF8EF98E2972C5"),
			hexToBytes("6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D"),
			hexToBytes("0EFE3E12F46363C7779140D4CE659925DB52F19053E114D7CC4EFD666B37F79F"),
		},
		hexToBytes("6E080211880400000000000022480A20"),
		[]types.TMSignature{
			types.TMSignature{
				hexToBytes("B88E0A2054A96A6775A9F5D1FA23B6FFA41274DD35C6431DAB0977F8CE4FB480"), // r
				hexToBytes("3D759EFF85E17601624D560A8ACD70E782EA23B58C2E718FAC98EBF488750A86"), // s
				28, // v
				hexToBytes("12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E"), // _signedDataSuffix
			},
			types.TMSignature{
				hexToBytes("17A66FF70C81C6A9C3040C1037CCC4EE9319E184D40956DC0DC30C1318901D36"), // r
				hexToBytes("4A4C0C9BF150967CE25C724E613DF6BE0C401B84AA29DE8599963F52A7DFA940"), // s
				27, // v
				hexToBytes("12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E"), // _signedDataSuffix
			},
			types.TMSignature{
				hexToBytes("453498042685AB34C627B5652E2F1FAD839C21DB3CEE4E01822F00885F1E0321"), // r
				hexToBytes("679781C8F2E3597ED3DF15E8E44B9CAF17D74894F0D4E22CD0F8C7CC1CB43963"), // s
				27, // v
				hexToBytes("12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E"), // _signedDataSuffix
			},
			types.TMSignature{
				hexToBytes("174505557E61260C06C7FD8962FF485BEBAD68E91B00C225452962B1FCBF1114"), // r
				hexToBytes("39B37ACD1759D47B09D18C6C3144EAB5B2D2CA34347DC60A4D58B369730C0DB9"), // s
				27, // v
				hexToBytes("12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E"), // _signedDataSuffix
			},
		},
	}
	err := keeper.RelayOracleState(ctx, 1160, blockRelayProof)
	require.Nil(t, err)
}

func TestVerifyOracleData(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, abci.Header{})
	keeper := tApp.GetPriceFeedKeeper()

	keeper.SetOracleStateHash(
		ctx,
		1456,
		hexToBytes("93C06B3C6357C6CD299552B4CD683B95F99695F1285C877D2CC481CE92FF0ACA"),
	)

	dataProof := types.OracleDataProof{
		Version:   51,
		RequestID: 1,
		CodeHash:  hexToBytes("0874EE3E5ABA7AE0EB8CB43BFEBED358826D111CECE0EF0F804E99EEA9264060"),
		Params:    hexToBytes("706172616D73"),
		Data:      hexToBytes("00000000000AEAD8"),
		MerklePaths: []types.IAVLMerklePath{
			types.IAVLMerklePath{
				IsDataOnRight:  true, // isDataOnRight
				SubtreeHeight:  1,    // subtreeHeight
				SubtreeSize:    2,    // subtreeSize
				SubtreeVersion: 51,   // subtreeVersion
				SiblingHash:    hexToBytes("81C8ADD6C13611F5416CA722C98815846504D7D2A77777B15D4D4212902F7785"),
			},
			types.IAVLMerklePath{
				IsDataOnRight:  true, // isDataOnRight
				SubtreeHeight:  2,    // subtreeHeight
				SubtreeSize:    3,    // subtreeSize
				SubtreeVersion: 51,   // subtreeVersion
				SiblingHash:    hexToBytes("29BB2C7201C2ECC2B0D1524D7E51827BC6C79E1E6E53249A630485A19876A173"),
			},
			types.IAVLMerklePath{
				IsDataOnRight:  true, // isDataOnRight
				SubtreeHeight:  3,    // subtreeHeight
				SubtreeSize:    5,    // subtreeSize
				SubtreeVersion: 51,   // subtreeVersion
				SiblingHash:    hexToBytes("8ED79B88CF56E86D5FA467274468ADE8FCA0BCCD5C797DE745AEB3D04F6B94DD"),
			},
			types.IAVLMerklePath{
				IsDataOnRight:  true, // isDataOnRight
				SubtreeHeight:  4,    // subtreeHeight
				SubtreeSize:    9,    // subtreeSize
				SubtreeVersion: 1455, // subtreeVersion
				SiblingHash:    hexToBytes("3E3734FD33F6C1DC10C2A0979FF3E36E31D1039331B19041BCE34BD4155E6757"),
			},
		},
	}
	err := keeper.VerifyOracleData(ctx, 1456, dataProof)
	require.Nil(t, err)
}
