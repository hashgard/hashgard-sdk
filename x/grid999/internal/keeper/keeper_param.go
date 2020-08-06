package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
)

func (k Keeper) GetParam(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyParams)
	if len(bz) == 0 {
		return types.DefaultParams()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &params)
	return
}
func (k Keeper) SetParam(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyParams, k.cdc.MustMarshalBinaryLengthPrefixed(params))
	return
}
