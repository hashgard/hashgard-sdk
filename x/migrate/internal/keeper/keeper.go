package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/migrate/internal/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	bankKeeper   BankKeeper
	supplyKeeper SupplyKeeper
	sceneKeeper  SceneKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, bankKeeper BankKeeper, supplyKeeper SupplyKeeper, sceneKeeper SceneKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper:   bankKeeper,
		supplyKeeper: supplyKeeper,
		sceneKeeper:  sceneKeeper,
	}
}
func (k Keeper) GetCodec() *codec.Codec {
	return k.cdc
}

func (k Keeper) GetBankKeeper() BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetSupplyKeeper() SupplyKeeper {
	return k.supplyKeeper
}
func (k Keeper) GetSceneKeeper() SceneKeeper {
	return k.sceneKeeper
}

//______________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
