package grid999

import (
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryDapp         = types.QueryDapp
	QueryParams       = types.QueryParams
	QueryGrid999      = types.QueryGrid999
	QueryGrid999List  = types.QueryGrid999List
	QueryDappList     = types.QueryDappList
	RouterKey         = types.RouterKey
)

type (
	Keeper = keeper.Keeper
	Params = types.Params
	Dapp   = types.Dapp
	Grid   = types.Grid
)

var (
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	DefaultParams = types.DefaultParams
)
