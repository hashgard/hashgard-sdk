package migrate

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, keeper)
		case QueryExchange:
			return queryExchange(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown migrates query endpoint")
		}
	}
}
func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(k.GetCodec(), params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}
	return res, nil
}
func queryExchange(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	err, migrate := keeper.GetErc20MigrateExchange(ctx)
	if err != nil {
		return nil, err
	}
	bz, err1 := codec.MarshalJSONIndent(keeper.GetCodec(), migrate)
	if err1 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err1.Error()))
	}
	return bz, nil
}
