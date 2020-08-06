package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Params struct {
	Owners             string   `json:"owners" yaml:"owners"`
	PerGridMaxDeposits int      `json:"per_grid_max_deposits" yaml:"per_grid_max_deposits"`
	GenerateDappFee    sdk.Coin `json:"generate_dapp_fee" yaml:"generate_dapp_fee"`
	CreateGridFee      sdk.Coin `json:"create_grid_fee" yaml:"create_grid_fee"`
	DepositFee         sdk.Coin `json:"deposit_fee" yaml:"deposit_fee"`
	WithdrawRewardsFee sdk.Coin `json:"withdraw_rewards_fee" yaml:"withdraw_rewards_fee"`
	FeeWithdrawFee     sdk.Dec  `json:"fee_withdraw_fee" yaml:"fee_withdraw_fee"`
}

func DefaultParams() Params {
	return Params{
		Owners:             "gard1rns8r0rzs629avtajcttkydjhcfy3n7na0cjge",
		PerGridMaxDeposits: 100,
		GenerateDappFee:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000)),
		FeeWithdrawFee:     sdk.NewDecWithPrec(2, 2),
		CreateGridFee:      sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10)),
		DepositFee:         sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(0)),
		WithdrawRewardsFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10)),
	}
}
func (c Params) String() string {
	return fmt.Sprintf(`Params:
		Owners:    %s
		GenerateDappFee:    %s`, c.Owners, c.GenerateDappFee.String())
}
