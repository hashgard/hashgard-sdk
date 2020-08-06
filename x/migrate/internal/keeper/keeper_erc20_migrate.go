package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
	"strings"
)

func (k Keeper) SetErc20MigrateExchange(ctx sdk.Context, migrate types.ERC20MigrateExchange) (err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(KeyErc20MigrateExchange) {
		return sdk.ErrInternal("erc20 migrate exchange already set")
	}
	store.Set(KeyErc20MigrateExchange, k.cdc.MustMarshalBinaryLengthPrefixed(migrate))
	return
}
func (k Keeper) GetErc20MigrateExchange(ctx sdk.Context) (err sdk.Error, migrate types.ERC20MigrateExchange) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyErc20MigrateExchange)
	if len(bz) == 0 {
		err = sdk.ErrInternal("erc20 migrate exchange not set")
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &migrate)
	return
}
func (k Keeper) AddErc20Migrate(ctx sdk.Context, erc20TxHash string) (err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(KeyERC20TxMigrateKey(erc20TxHash)) {
		return sdk.ErrInternal("erc20 txhash already migrate:" + erc20TxHash)
	}
	store.Set(KeyERC20TxMigrateKey(erc20TxHash), []byte{1})
	return
}
func (k Keeper) ExportErc20Migrate(ctx sdk.Context) (list []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyERC20TxMigratePrefix())
	defer iterator.Close()
	list = make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		list = append(list, strings.ReplaceAll(string(iterator.Key()), Erc20TxMigrateKeyPrefix, ""))
	}
	return
}
