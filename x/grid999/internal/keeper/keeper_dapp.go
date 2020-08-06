package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	Grid999FeeAddr      = sdk.AccAddress(crypto.AddressHash([]byte("Grid999FeeAddr")))
	MinDappID      uint = 1000000
	MaxDappID      uint = 9999999
)

func (k Keeper) GenerateDapp(ctx sdk.Context, dapp *types.Dapp) (err sdk.Error) {
	params := k.GetParam(ctx)

	if params.PerGridMaxDeposits < dapp.PerGridMaxDeposits {
		return sdk.ErrInternal(fmt.Sprintf("per grid max deposits is %d", params.PerGridMaxDeposits))
	}
	if err := k.SendGrid999FeeCoins(ctx, dapp.Owner, sdk.NewCoins(params.GenerateDappFee)); err != nil {
		return sdk.ErrInsufficientFee("generate a dapp requires payment:" + params.GenerateDappFee.String())
	}

	store := ctx.KVStore(k.storeKey)
	dapp.ID = k.getLastDappID(store)
	if dapp.ID >= MaxDappID {
		return sdk.ErrInternal("Exceed the maximum number of dapps")
	}
	k.SetDapp(ctx, dapp)
	k.setNewDappID(store, dapp.ID+1)
	return
}
func (k Keeper) SetDapp(ctx sdk.Context, dapp *types.Dapp) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyDapp(dapp.ID), k.cdc.MustMarshalBinaryLengthPrefixed(dapp))
	return
}
func (k Keeper) DisableDapp(ctx sdk.Context, sender sdk.AccAddress, dappID uint, height int64) (err sdk.Error) {

	dapp := k.GetDapp(ctx, dappID)
	if dapp == nil {
		return sdk.ErrInternal("not found thd dapp")
	}
	if !dapp.Owner.Equals(sender) {
		return sdk.ErrInternal("permission denied")
	}
	if dapp.EndHeight > 0 {
		return sdk.ErrInternal("already disabled")
	}
	dapp.EndHeight = height
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyDapp(dapp.ID), k.cdc.MustMarshalBinaryLengthPrefixed(dapp))
	return nil
}

func (k Keeper) GetDapp(ctx sdk.Context, id uint) (dapp *types.Dapp) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyDapp(id))
	if len(bz) == 0 {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dapp)
	return
}
func (k Keeper) GetDapps(ctx sdk.Context, numLatest uint, limit int) []*types.Dapp {
	store := ctx.KVStore(k.storeKey)
	list := make([]*types.Dapp, 0, limit)
	maxId := k.getLastDappID(store)
	if numLatest == 0 {
		numLatest = maxId
	}
	for txID := numLatest; txID > 0; txID-- {
		grid := k.GetDapp(ctx, txID)
		if grid == nil {
			continue
		}
		list = append(list, grid)
		if len(list) >= limit {
			break
		}
	}
	return list
}
func (k Keeper) ExportDapps(ctx sdk.Context) []types.Dapp {
	store := ctx.KVStore(k.storeKey)
	list := make([]types.Dapp, 0)
	maxId := k.getLastDappID(store)
	for txID := MinDappID; txID <= maxId; txID++ {
		dapp := k.GetDapp(ctx, txID)
		if dapp == nil {
			continue
		}
		list = append(list, *dapp)
	}
	return list
}
func (k Keeper) setNewDappID(store sdk.KVStore, id uint) {
	store.Set(KeyNextDappID, k.cdc.MustMarshalBinaryLengthPrefixed(id))
	return
}
func (k Keeper) SetNewDappID(ctx sdk.Context, id uint) {
	store := ctx.KVStore(k.storeKey)
	k.setNewDappID(store, id)
}
func (k Keeper) GetLastDappID(ctx sdk.Context) (id uint) {
	store := ctx.KVStore(k.storeKey)
	return k.getLastDappID(store)
}

func (k Keeper) getLastDappID(store sdk.KVStore) (id uint) {
	bz := store.Get(KeyNextDappID)
	id = MinDappID
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	}
	return
}
