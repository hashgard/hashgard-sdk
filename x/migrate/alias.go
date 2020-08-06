package migrate

import (
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QueryParameters   = types.QueryParameters
	QueryExchange     = types.QueryExchange
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
)

type (
	Keeper = keeper.Keeper
	Params = types.Params

	MsgERC20Migrate         = types.MsgERC20Migrate
	MsgERC20MigrateExchange = types.MsgERC20MigrateExchange
	ERC20MigrateExchange    = types.ERC20MigrateExchange
)

var (
	ModuleCdc      = types.ModuleCdc
	NewKeeper      = keeper.NewKeeper
	RegisterCodec  = types.RegisterCodec
	ParamKeyTable  = types.ParamKeyTable
	NewParams      = types.NewParams
	DefaultParams  = types.DefaultParams
	ValidateParams = types.ValidateParams
)
