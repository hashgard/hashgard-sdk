package migrate

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - migrate state
type GenesisState struct {
	Params              Params               `json:"params" yaml:"params"`
	Exchange            ERC20MigrateExchange `json:"exchange" yaml:"exchange"`
	ERC20MigratedTxHash []string             `json:"erc20_migrated_tx_hash" yaml:"erc20_migrated_tx_hash"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, exchange ERC20MigrateExchange, erc20MigratedTxHash []string) GenesisState {
	return GenesisState{
		Params:              params,
		Exchange:            exchange,
		ERC20MigratedTxHash: erc20MigratedTxHash,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:              DefaultParams(),
		Exchange:            ERC20MigrateExchange{},
		ERC20MigratedTxHash: make([]string, 0),
	}
}

// InitGenesis new migrate genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	if len(data.Exchange.Allows) > 0 {
		keeper.SetErc20MigrateExchange(ctx, data.Exchange)
	}
	for _, tx := range data.ERC20MigratedTxHash {
		keeper.AddErc20Migrate(ctx, tx)
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	_, exchange := keeper.GetErc20MigrateExchange(ctx)
	return NewGenesisState(keeper.GetParams(ctx), exchange, keeper.ExportErc20Migrate(ctx))
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
