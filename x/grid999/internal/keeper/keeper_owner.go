package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

func (k Keeper) CheckOwner(ctx sdk.Context, sender sdk.AccAddress) (err sdk.Error) {
	params := k.GetParam(ctx)
	if !strings.Contains(params.Owners, sender.String()) {
		return sdk.ErrInternal("Permission denied")
	}
	return nil
}

func (k Keeper) WithdrawFees(ctx sdk.Context, dappId uint, sender sdk.AccAddress, to sdk.AccAddress) (withdraw, fees sdk.Coins, err sdk.Error) {
	dapp := k.GetDapp(ctx, dappId)
	if dapp == nil {
		return sdk.Coins{}, sdk.Coins{}, sdk.ErrInternal("dapp is not found")
	}
	if !dapp.Owner.Equals(sender) {
		return sdk.Coins{}, sdk.Coins{}, sdk.ErrInternal("Permission denied")
	}
	withdraw = k.GetBankKeeper().GetCoins(ctx, k.GetFeeAddr(dappId))
	if withdraw.IsZero() {
		return sdk.Coins{}, sdk.Coins{}, sdk.ErrInsufficientCoins("zero")
	}
	fees = k.getDefaultCoin(dapp)

	params := k.GetParam(ctx)
	if !params.FeeWithdrawFee.IsZero() {
		fees = k.MulRate(dapp, withdraw, params.FeeWithdrawFee)
		if err := k.SendGrid999FeeCoins(ctx, k.GetFeeAddr(dappId), fees); err != nil {
			return sdk.Coins{}, sdk.Coins{}, err
		}
		withdraw = withdraw.Sub(fees)
	}

	if err = k.GetBankKeeper().SendCoins(ctx, k.GetFeeAddr(dappId), to, withdraw); err != nil {
		return sdk.Coins{}, sdk.Coins{}, err
	}
	return
}
func (k Keeper) WithdrawLucky(ctx sdk.Context, sender sdk.AccAddress, dappId uint, luckyCoin sdk.Coins) (fees sdk.Coins, err sdk.Error) {
	dapp := k.GetDapp(ctx, dappId)
	if dapp == nil {
		return sdk.Coins{}, sdk.ErrInternal("dapp not found")
	}
	if !dapp.Owner.Equals(sender) {
		return sdk.Coins{}, sdk.ErrInternal("Permission denied")
	}

	if dapp.EndHeight == 0 || dapp.EndHeight > ctx.BlockHeight() {
		return sdk.Coins{}, sdk.ErrInternal(fmt.Sprintf("can't withdraw lucky deposit from height:%d,you must disabled the dapp first", ctx.BlockHeight()))
	}

	fees = k.getDefaultCoin(dapp)
	params := k.GetParam(ctx)
	if !params.FeeWithdrawFee.IsZero() {
		fees = k.MulRate(dapp, luckyCoin, params.FeeWithdrawFee)
		if err = k.SendGrid999FeeCoins(ctx, k.GetFeeAddr(dappId), fees); err != nil {
			return
		}
		luckyCoin = luckyCoin.Sub(fees)
	}
	if err = k.GetBankKeeper().SendCoins(ctx, k.GetLuckyPoolAddr(dappId), k.GetCommunityPoolAddr(), luckyCoin); err != nil {
		return sdk.Coins{}, err
	}
	return
}
