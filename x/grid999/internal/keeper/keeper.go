package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the grid999 store
type Keeper struct {
	cdc         *codec.Codec
	storeKey    sdk.StoreKey
	bankKeeper  BankKeeper
	sceneKeeper SceneKeeper
}

// NewKeeper creates a new grid999 Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper BankKeeper, sceneKeeper SceneKeeper) Keeper {
	return Keeper{
		cdc:         cdc,
		storeKey:    key,
		bankKeeper:  bankKeeper,
		sceneKeeper: sceneKeeper,
	}
}
func (k Keeper) GetCodec() *codec.Codec {
	return k.cdc
}

func (k Keeper) GetBankKeeper() BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSceneKeeper() SceneKeeper {
	return k.sceneKeeper
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SendCoins(ctx sdk.Context, internal bool, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if amt.IsZero() {
		return nil
	}
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, amt); err != nil {
		return err
	}
	if internal {
		k.sceneKeeper.AddSenderTxScene(ctx, fromAddr)
	}
	return nil
}
func (k Keeper) SendGrid999FeeCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.SendCoins(ctx, true, fromAddr, Grid999FeeAddr, amt)
}
func (k Keeper) WithdrawGrid999Fee(ctx sdk.Context, to sdk.AccAddress) (withdraw sdk.Coins, err sdk.Error) {
	withdraw = k.GetBankKeeper().GetCoins(ctx, Grid999FeeAddr)

	if withdraw.IsZero() {
		return withdraw, sdk.ErrInsufficientCoins("zero")
	}

	if err = k.GetBankKeeper().SendCoins(ctx, Grid999FeeAddr, to, withdraw); err != nil {
		return
	}
	return
}
