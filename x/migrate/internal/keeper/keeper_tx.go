package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func (k Keeper) AddTxMigrate(ctx sdk.Context, migrateKey string) (migrateId uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	migrateId = k.getNewTxMigrateID(store, migrateKey)
	store.Set(KeyTxMigrate(migrateKey, migrateId), k.cdc.MustMarshalBinaryLengthPrefixed(fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes()))))
	return
}
func (k Keeper) getTxMigrate(store sdk.KVStore, migrateKey string, migrateId uint64) (txHash string) {
	bz := store.Get(KeyTxMigrate(migrateKey, migrateId))
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &txHash)
	}
	return
}

func (k Keeper) GetTxMigrateTxs(ctx sdk.Context, migrateKey string, numLatest uint64, limit int, sort string) []*types.TxId {
	store := ctx.KVStore(k.storeKey)
	list := make([]*types.TxId, 0, limit)
	maxMigrateId := k.getLastTxMigrateID(store, migrateKey) - 1
	if "asc" == sort {
		if numLatest == 0 {
			numLatest = 1
		}
		for txID := numLatest; txID <= maxMigrateId; txID++ {
			list = append(list, &types.TxId{ID: txID, TxHash: k.getTxMigrate(store, migrateKey, txID)})
			if len(list) >= limit {
				break
			}
		}
		return list
	}
	if numLatest == 0 {
		numLatest = maxMigrateId
	}
	for txID := numLatest; txID > 0; txID-- {
		list = append(list, &types.TxId{ID: txID, TxHash: k.getTxMigrate(store, migrateKey, txID)})
		if len(list) >= limit {
			break
		}
	}
	return list
}

func (k Keeper) getNewTxMigrateID(store sdk.KVStore, migrateKey string) (migrateId uint64) {
	migrateId = k.getLastTxMigrateID(store, migrateKey)
	store.Set(KeyNextTxMigrateID(migrateKey), k.cdc.MustMarshalBinaryLengthPrefixed(migrateId+1))
	return
}

func (k Keeper) getLastTxMigrateID(store sdk.KVStore, migrateKey string) (migrateId uint64) {
	bz := store.Get(KeyNextTxMigrateID(migrateKey))
	migrateId = 1
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &migrateId)
	}
	return
}
