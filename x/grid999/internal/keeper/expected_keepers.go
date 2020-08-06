package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected supply keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
type SceneKeeper interface {
	AddSenderTxScene(ctx sdk.Context, sender sdk.AccAddress) (err sdk.Error)
}
