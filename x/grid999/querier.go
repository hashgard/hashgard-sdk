package grid999

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, keeper)
		case QueryDappList:
			return queryDappList(ctx, path[1], path[2], keeper)
		case QueryDapp:
			return queryDapp(ctx, path[1], keeper)
		case QueryGrid999:
			return queryGrid999(ctx, path[1], path[2], keeper)
		case QueryGrid999List:
			return queryGrid999List(ctx, path[1], path[2], path[3], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown scenes query endpoint")
		}
	}
}
func queryParams(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	grid := keeper.GetParam(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), grid)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryGrid999(ctx sdk.Context, dappIDStr, gridId string, keeper Keeper) ([]byte, sdk.Error) {
	id, err := strconv.ParseUint(gridId, 0, 64)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	dappID, err := strconv.Atoi(dappIDStr)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	grid, err := keeper.GetGrid(ctx, uint(dappID), id)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), grid)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
func queryDapp(ctx sdk.Context, generateId string, keeper Keeper) ([]byte, sdk.Error) {
	id, err := strconv.Atoi(generateId)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	grid := keeper.GetDapp(ctx, uint(id))

	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), grid)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryGrid999List(ctx sdk.Context, dappIDStr string, startID, limitStr string, keeper Keeper) ([]byte, sdk.Error) {
	id, err := strconv.ParseUint(startID, 0, 64)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	dappID, err := strconv.Atoi(dappIDStr)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	grid := keeper.GetGrids(ctx, uint(dappID), id, limit)

	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), grid)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryDappList(ctx sdk.Context, startID, limitStr string, keeper Keeper) ([]byte, sdk.Error) {
	id, err := strconv.ParseUint(startID, 0, 64)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	list := keeper.GetDapps(ctx, uint(id), limit)

	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), list)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
