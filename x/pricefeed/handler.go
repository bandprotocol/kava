package pricefeed

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler handles all pricefeed type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPostPrice:
			return HandleMsgPostPrice(ctx, k, msg)
		case MsgReportPrice:
			return HandleMsgReportPrice(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized pricefeed message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// price feed questions:
// do proposers need to post the round in the message? If not, how do we determine the round?

// HandleMsgPostPrice handles prices posted by oracles
func HandleMsgPostPrice(
	ctx sdk.Context,
	k Keeper,
	msg MsgPostPrice) sdk.Result {

	_, err := k.GetOracle(ctx, msg.MarketID, msg.From)
	if err != nil {
		return err.Result()
	}
	_, err = k.SetPrice(ctx, msg.From, msg.MarketID, msg.Price, msg.Expiry)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleMsgReportPrice(
	ctx sdk.Context,
	k Keeper,
	msg MsgReportPrice,
) sdk.Result {

	err := k.RelayOracleState(ctx, msg.Proof.BlockHeight, msg.Proof.BlockRelayProof)
	if err != nil {
		return err.Result()
	}
	err = k.VerifyOracleData(ctx, msg.Proof.BlockHeight, msg.Proof.OracleDataProof)
	if err != nil {
		return err.Result()
	}

	// TODO: Check codehash of this id with verified codehash set in genesis state

	// TODO: Find good way to parse result data
	if len(msg.Proof.OracleDataProof.Params) != 4 || len(msg.Proof.OracleDataProof.Data) != 16 {
		return sdk.ErrInvalidCoins("Invalid params").Result()
	}
	price := binary.BigEndian.Uint64(msg.Proof.OracleDataProof.Data[:8])
	decPrice, err := sdk.NewDecFromStr(fmt.Sprintf("%.2f", float64(price)/100.0))
	if err != nil {
		return err.Result()
	}
	// TODO: find market id that match our parameter
	switch uint8(msg.Proof.OracleDataProof.Params[3]) {
	case 4:
		{
			err := k.SetCurrentPriceByBand(ctx, "btc:usd", decPrice)
			if err != nil {
				return err.Result()
			}
		}
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
