package grid999

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - scene state
type GenesisState struct {
	Dapps       []Dapp   `json:"dapps" yaml:"dapps"`
	Grids       []Grid   `json:"grids" yaml:"grids"`
	Params      Params   `json:"params" yaml:"params"`
	LastDappsId uint     `json:"last_dapps_id" yaml:"last_dapps_id"`
	LastGridId  []string `json:"last_grid_id" yaml:"last_grid_id"`
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Dapps:      make([]Dapp, 0),
		Grids:      make([]Grid, 0),
		Params:     DefaultParams(),
		LastGridId: make([]string, 0),
	}
}

// InitGenesis new scene genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if len(data.Params.Owners) > 0 {
		keeper.SetParam(ctx, data.Params)
	}

	for _, v := range data.Grids {
		keeper.SetGrid(ctx, &v)
	}
	for _, v := range data.Dapps {
		keeper.SetDapp(ctx, &v)
	}
	keeper.ImportGridIDs(ctx, data.LastGridId)
	if data.LastDappsId > 0 {
		keeper.SetNewDappID(ctx, data.LastDappsId)
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
	genesisState := GenesisState{
		LastDappsId: keeper.GetLastDappID(ctx),
		Dapps:       keeper.ExportDapps(ctx),
		Params:      keeper.GetParam(ctx),
		Grids:       keeper.ExportGrids(ctx),
		LastGridId:  keeper.ExportGridIDs(ctx),
	}

	return genesisState
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return nil
}
