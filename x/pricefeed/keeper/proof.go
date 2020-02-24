package keeper

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/kava-labs/kava/x/pricefeed/types"
)

func (k Keeper) SetValidatorPubkeys(ctx sdk.Context, pubkeys []secp256k1.PubKeySecp256k1) {
	store := ctx.KVStore(k.key)
	store.Set(
		[]byte(types.BandValidatorsPrefix),
		k.cdc.MustMarshalBinaryBare(pubkeys),
	)
}

func (k Keeper) GetValidatorPubkeys(ctx sdk.Context) []secp256k1.PubKeySecp256k1 {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(types.BandValidatorsPrefix))
	var pubkeys []secp256k1.PubKeySecp256k1
	k.cdc.MustUnmarshalBinaryBare(bz, &pubkeys)
	return pubkeys
}

func (k Keeper) RelayOracleState(
	ctx sdk.Context, blockHeight uint64, blockProof types.BlockRelayProof,
) sdk.Error {
	appHash := blockProof.GetAppHash()
	blockHeader := blockProof.BlockHeaderMerkleParts.GetBlockHeader(
		blockHeight,
		appHash,
	)

	_ = blockHeader
	// TODO: Check validator signatures
	// pubkeys := k.GetValidatorPubkeys(ctx)
	// for _, pk := range pubkeys {
	// 	found := false
	// 	for i := range blockProof.Signatures {
	// 		if blockProof.Signatures[i].VerifyWithPub(&pk, blockHeader, blockProof.SignedDataPrefix) {
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		return types.ErrInvalidOracle(k.Codespace(), sdk.AccAddress(pk.Address()))
	// 	}
	// }
	k.SetOracleStateHash(ctx, blockHeight, blockProof.OracleIAVLStateHash)
	return nil
}

func (k Keeper) GetOracleStateHash(ctx sdk.Context, blockHeight uint64) ([]byte, sdk.Error) {
	store := ctx.KVStore(k.key)
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, blockHeight)
	key := append([]byte(types.OracleStateHashPrefix), buf[:n]...)
	if !store.Has(key) {
		return []byte{}, types.ErrEmptyInput(k.Codespace())
	}
	return store.Get(key), nil
}

func (k Keeper) SetOracleStateHash(ctx sdk.Context, blockHeight uint64, oracleStateHash []byte) {
	store := ctx.KVStore(k.key)
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, blockHeight)
	store.Set(
		append([]byte(types.OracleStateHashPrefix), buf[:n]...),
		oracleStateHash,
	)
}

func (k Keeper) VerifyOracleData(
	ctx sdk.Context, blockHeight uint64, dataProof types.OracleDataProof,
) sdk.Error {
	oracleStateRoot, err := k.GetOracleStateHash(ctx, blockHeight)
	if err != nil {
		return err
	}
	if !bytes.Equal(dataProof.GetOracleStateRoot(blockHeight), oracleStateRoot) {
		return types.ErrEmptyInput(k.Codespace())
	}
	return nil
}
