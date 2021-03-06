// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/kava-labs/kava/x/auction/types/
// ALIASGEN: github.com/kava-labs/kava/x/auction/keeper/
package auction

import (
	"github.com/kava-labs/kava/x/auction/keeper"
	"github.com/kava-labs/kava/x/auction/types"
)

const (
	DefaultCodespace                      = types.DefaultCodespace
	CodeInvalidInitialAuctionID           = types.CodeInvalidInitialAuctionID
	CodeUnrecognizedAuctionType           = types.CodeUnrecognizedAuctionType
	CodeAuctionNotFound                   = types.CodeAuctionNotFound
	CodeAuctionHasNotExpired              = types.CodeAuctionHasNotExpired
	CodeAuctionHasExpired                 = types.CodeAuctionHasExpired
	CodeInvalidBidDenom                   = types.CodeInvalidBidDenom
	CodeInvalidLotDenom                   = types.CodeInvalidLotDenom
	CodeBidTooSmall                       = types.CodeBidTooSmall
	CodeBidTooLarge                       = types.CodeBidTooLarge
	CodeLotTooLarge                       = types.CodeLotTooLarge
	CodeCollateralAuctionIsInReversePhase = types.CodeCollateralAuctionIsInReversePhase
	CodeCollateralAuctionIsInForwardPhase = types.CodeCollateralAuctionIsInForwardPhase
	ModuleName                            = types.ModuleName
	StoreKey                              = types.StoreKey
	RouterKey                             = types.RouterKey
	DefaultParamspace                     = types.DefaultParamspace
	DefaultMaxAuctionDuration             = types.DefaultMaxAuctionDuration
	DefaultBidDuration                    = types.DefaultBidDuration
	QueryGetAuction                       = types.QueryGetAuction
	DefaultNextAuctionID                  = types.DefaultNextAuctionID
)

var (
	// functions aliases
	NewSurplusAuction    = types.NewSurplusAuction
	NewDebtAuction       = types.NewDebtAuction
	NewCollateralAuction = types.NewCollateralAuction
	NewWeightedAddresses = types.NewWeightedAddresses
	RegisterCodec        = types.RegisterCodec
	NewGenesisState      = types.NewGenesisState
	DefaultGenesisState  = types.DefaultGenesisState
	GetAuctionKey        = types.GetAuctionKey
	GetAuctionByTimeKey  = types.GetAuctionByTimeKey
	Uint64FromBytes      = types.Uint64FromBytes
	Uint64ToBytes        = types.Uint64ToBytes
	NewMsgPlaceBid       = types.NewMsgPlaceBid
	NewParams            = types.NewParams
	DefaultParams        = types.DefaultParams
	ParamKeyTable        = types.ParamKeyTable
	NewKeeper            = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier

	// variable aliases
	ModuleCdc              = types.ModuleCdc
	AuctionKeyPrefix       = types.AuctionKeyPrefix
	AuctionByTimeKeyPrefix = types.AuctionByTimeKeyPrefix
	NextAuctionIDKey       = types.NextAuctionIDKey
	KeyAuctionBidDuration  = types.KeyAuctionBidDuration
	KeyAuctionDuration     = types.KeyAuctionDuration
)

type (
	Auction           = types.Auction
	BaseAuction       = types.BaseAuction
	SurplusAuction    = types.SurplusAuction
	DebtAuction       = types.DebtAuction
	CollateralAuction = types.CollateralAuction
	WeightedAddresses = types.WeightedAddresses
	SupplyKeeper      = types.SupplyKeeper
	GenesisAuctions   = types.GenesisAuctions
	GenesisAuction    = types.GenesisAuction
	GenesisState      = types.GenesisState
	MsgPlaceBid       = types.MsgPlaceBid
	Params            = types.Params
	Keeper            = keeper.Keeper
)
