package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
)

var (
	KeyStakeIssueTokenConfig = []byte(fmt.Sprintf("sicfg"))
	StakeIssueLockedAdr      = sdk.AccAddress(crypto.AddressHash([]byte("StakeIssueLockedAdr")))
)

// HashGard
func KeyStakeIssueToken(sender sdk.ValAddress) []byte {
	return []byte(fmt.Sprintf("sit:%s", sender.String()))
}
func KeyStakeIssueTokenPrefix() []byte {
	return []byte(fmt.Sprintf("sit:"))
}
func KeyStakeIssueTokenSymbol(denom string) []byte {
	return []byte(fmt.Sprintf("sits:%s", denom))
}
func KeyStakeIssueTokenUnbondHeight(denom string) []byte {
	return []byte(fmt.Sprintf("situh:%s", denom))
}
func GetStakeIssueLockedAddr(denom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte("StakeIssueLockedAdr" + denom)))
}
func (k Keeper) SetStakeIssueTokenConfig(ctx sdk.Context, sender sdk.AccAddress, config types.StakeIssueTokenConfig) sdk.Error {
	if !strings.Contains(types.StakeIssueTokenConfigAllows, sender.String()) {
		return sdk.ErrInternal("Permission denied")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyStakeIssueTokenConfig, k.cdc.MustMarshalBinaryLengthPrefixed(config))
	return nil
}
func (k Keeper) GetStakeIssueTokenConfig(ctx sdk.Context) (config types.StakeIssueTokenConfig) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyStakeIssueTokenConfig)
	if len(bz) == 0 {
		return types.StakeIssueTokenConfig{
			MinSelfDelegation: sdk.TokensFromConsensusPower(10000000),
			LockPeriodHeight:  10512000,
			LockCoins:         sdk.NewCoins(sdk.NewCoin("uggt", sdk.TokensFromConsensusPower(200000))),
		}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &config)
	return
}
func (k Keeper) UnlockLockedCoins(ctx sdk.Context, stakeIssueToken types.StakeIssueToken) sdk.Error {
	if stakeIssueToken.Rule.LockCoins.IsZero() {
		return nil
	}
	k.sceneKeeper.AddSenderTxScene(ctx, StakeIssueLockedAdr)

	if stakeIssueToken.Height+stakeIssueToken.Rule.LockPeriodHeight > ctx.BlockHeight() {
		return k.bankKeeper.SendCoins(ctx, StakeIssueLockedAdr, stakeIssueToken.LockedAddr, stakeIssueToken.Rule.LockCoins)
	}

	return k.bankKeeper.SendCoins(ctx, StakeIssueLockedAdr, stakeIssueToken.Owner, stakeIssueToken.Rule.LockCoins)
}
func (k Keeper) UndelegateCoinsToLockedAddr(ctx sdk.Context, lockedAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.NotBondedPoolName, lockedAddr, amt); err != nil {
		return err
	}
	k.sceneKeeper.AddSenderTxScene(ctx, lockedAddr)
	return nil
}
func (k Keeper) addLockCoins(ctx sdk.Context, sender sdk.AccAddress, lockCoins sdk.Coins) sdk.Error {
	if lockCoins.IsZero() {
		return nil
	}
	return k.bankKeeper.SendCoins(ctx, sender, StakeIssueLockedAdr, lockCoins)
}

func (k Keeper) AddStakeIssueToken(ctx sdk.Context, sender sdk.AccAddress, valAddress sdk.ValAddress, minSelfDelegation sdk.Int, stakeIssueToken types.StakeIssueToken) sdk.Error {
	if ctx.BlockHeight() < types.StakeIssueTokenStartBlock {
		return sdk.ErrInternal("issuance of tokens is not allowed now")
	}
	if stakeIssueToken.GenesisHeight-ctx.BlockHeight() < 20 {
		return sdk.ErrInternal(fmt.Sprintf("genesis height must be greater than the current 20 blocks,current block is %d", ctx.BlockHeight()))
	}
	store := ctx.KVStore(k.storeKey)
	if store.Has(KeyStakeIssueToken(valAddress)) {
		return sdk.ErrInternal(valAddress.String() + " have issued a token")
	}
	if store.Has(KeyStakeIssueTokenSymbol(stakeIssueToken.Denom)) {
		return sdk.ErrInternal(stakeIssueToken.Denom + " already issued")
	}
	config := k.GetStakeIssueTokenConfig(ctx)
	if minSelfDelegation.LT(config.MinSelfDelegation) {
		return sdk.ErrInternal("nodes must self-staking at least " + config.MinSelfDelegation.String() + " to issue a token")
	}
	stakeIssueToken.LockedAddr = GetStakeIssueLockedAddr(stakeIssueToken.Denom)

	if err := k.addLockCoins(ctx, sender, config.LockCoins); err != nil {
		return err
	}

	k.bankKeeper.AddCoins(ctx, stakeIssueToken.LockedAddr, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())))

	stakeIssueToken.Val = valAddress
	stakeIssueToken.Owner = sender
	stakeIssueToken.Rule = config
	stakeIssueToken.Height = ctx.BlockHeight()

	if err := k.SetStakeIssueToken(ctx, stakeIssueToken); err != nil {
		return err
	}

	return nil
}
func (k Keeper) SetStakeIssueToken(ctx sdk.Context, stakeIssueToken types.StakeIssueToken) sdk.Error {

	k.EditStakeIssueToken(ctx, stakeIssueToken)

	if stakeIssueToken.PreMintAmount.GT(sdk.ZeroInt()) && stakeIssueToken.PreMintAddress != nil {
		if k.bankKeeper != nil {
			coins := sdk.NewCoins(sdk.NewCoin(stakeIssueToken.Denom, stakeIssueToken.PreMintAmount))
			if err := k.supplyKeeper.MintCoins(ctx, "mint", coins); err != nil {
				return err
			}
			_, err := k.bankKeeper.AddCoins(ctx, stakeIssueToken.PreMintAddress, coins)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (k Keeper) EditStakeIssueToken(ctx sdk.Context, stakeIssueToken types.StakeIssueToken) {
	store := ctx.KVStore(k.storeKey)

	store.Set(KeyStakeIssueToken(stakeIssueToken.Val), k.cdc.MustMarshalBinaryLengthPrefixed(stakeIssueToken))
	store.Set(KeyStakeIssueTokenSymbol(stakeIssueToken.Denom), k.cdc.MustMarshalBinaryLengthPrefixed(stakeIssueToken.Owner))
}
func (k Keeper) GetStakeIssueTokenByAddress(ctx sdk.Context, address sdk.AccAddress, valAddr sdk.ValAddress) (stakeIssueToken types.StakeIssueToken, issued bool) {
	stakeIssueToken, issued = k.GetStakeIssueToken(ctx, sdk.ValAddress(address))
	if issued && !stakeIssueToken.Val.Equals(valAddr) {
		issued = false
	}
	return
}
func (k Keeper) GetStakeIssueToken(ctx sdk.Context, valAddress sdk.ValAddress) (stakeIssueToken types.StakeIssueToken, issued bool) {
	if ctx.BlockHeight() < types.StakeIssueTokenStartBlock {
		return
	}
	if k.bankKeeper == nil {
		return
	}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyStakeIssueToken(valAddress))
	if len(bz) == 0 {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &stakeIssueToken)
	issued = true
	return
}
func (k Keeper) GetStakeIssueTokenPerBlockMint(ctx sdk.Context, consAddr sdk.ConsAddress) (proposerNodeAmount, voterNodeAmount sdk.Coin, minted bool) {
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return
	}
	stakeIssueToken, issued := k.GetStakeIssueToken(ctx, validator.GetOperator())
	if !issued {
		return
	}
	if stakeIssueToken.GenesisHeight > ctx.BlockHeight() {
		return
	}
	for _, v := range stakeIssueToken.PerBlockMint {
		if ctx.BlockHeight() < v.StartHeight {
			break
		}
		proposerNodeAmount = sdk.NewCoin(stakeIssueToken.Denom, v.ProposerNodeAmount)
		voterNodeAmount = sdk.NewCoin(stakeIssueToken.Denom, v.VoterNodeAmount)
	}

	if proposerNodeAmount.IsZero() && voterNodeAmount.IsZero() {
		return
	}
	supply := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(stakeIssueToken.Denom)
	if supply.Add(proposerNodeAmount.Amount).Add(voterNodeAmount.Amount).Sub(stakeIssueToken.TotalSupply).GT(sdk.ZeroInt()) {
		remain := stakeIssueToken.TotalSupply.Sub(supply)
		if remain.LTE(sdk.ZeroInt()) {
			return
		}
		proposerNodeAmount = sdk.NewCoin(stakeIssueToken.Denom, remain)
		voterNodeAmount = sdk.NewCoin(stakeIssueToken.Denom, sdk.ZeroInt())
	}

	recipientAcc := k.supplyKeeper.GetModuleAccount(ctx, "distribution")
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s isn't able to be created", "distribution"))
	}
	coins := sdk.NewCoins(proposerNodeAmount).Add(sdk.NewCoins(voterNodeAmount))

	k.supplyKeeper.MintCoins(ctx, "mint", coins)
	k.bankKeeper.AddCoins(ctx, recipientAcc.GetAddress(), coins)

	minted = true
	return
}
func (k Keeper) StakeIssueLockedSpend(ctx sdk.Context, denom string, Recipient sdk.AccAddress, amount sdk.Coins) (err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyStakeIssueTokenSymbol(denom))
	if len(bz) == 0 {
		return sdk.ErrInternal(denom + "not issued")
	}
	return k.bankKeeper.SendCoins(ctx, GetStakeIssueLockedAddr(denom), Recipient, amount)
}

func (k Keeper) GetStakeIssueTokens(ctx sdk.Context) (list []types.StakeIssueToken) {
	if k.bankKeeper == nil {
		return
	}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyStakeIssueTokenPrefix())
	defer iterator.Close()
	list = make([]types.StakeIssueToken, 0)
	for ; iterator.Valid(); iterator.Next() {
		var item types.StakeIssueToken
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &item)
		list = append(list, item)
	}
	return
}
func (k Keeper) AddStakeIssueTokenUnbondHeight(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyStakeIssueTokenUnbondHeight(denom), k.cdc.MustMarshalBinaryLengthPrefixed(ctx.BlockHeight()))
}

func (k Keeper) GetStakeIssueTokenUnbondHeight(ctx sdk.Context, denom string) (height int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyStakeIssueTokenUnbondHeight(denom))
	if len(bz) == 0 {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)
	return
}
