package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyERC20MigrateDisabled = []byte("ERC20MigrateDisabled")
)

// mint parameters
type Params struct {
	ERC20MigrateDisabled bool `json:"erc20_migrate_disabled" yaml:"erc20_migrate_disabled"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(erc20MigrateDisabled bool) Params {
	return Params{
		ERC20MigrateDisabled: erc20MigrateDisabled,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		ERC20MigrateDisabled: false,
	}
}

// validate params
func ValidateParams(params Params) error {
	return nil
}
func (p Params) String() string {
	return fmt.Sprintf(`Migrate Params:
  ERC20 Migrate Disabled:             %t`,
		p.ERC20MigrateDisabled,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyERC20MigrateDisabled, &p.ERC20MigrateDisabled},
	}
}
