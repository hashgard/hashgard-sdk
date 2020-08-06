package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Dapp struct {
	ID                              uint           `json:"id" yaml:"id"`
	Owner                           sdk.AccAddress `json:"owner" yaml:"owner"`
	Name                            string         `json:"name" yaml:"name"`
	Icon                            string         `json:"icon" yaml:"icon"`
	Desc                            string         `json:"desc" yaml:"desc"`
	Voted                           bool           `json:"voted" yaml:"voted"`
	EndHeight                       int64          `json:"end_height" yaml:"end_height"`
	MaxBlocksGridCreate             int64          `json:"max_blocks_grid_create" yaml:"max_blocks_grid_create"`
	GridCreateCanDeposit            bool           `json:"grid_create_can_deposit" yaml:"grid_create_can_deposit"`
	OnlyOwnerCanCreateGrid          bool           `json:"only_owner_can_create_grid" yaml:"only_owner_can_create_grid"`
	MaxBlocksGridDeposit            int64          `json:"max_blocks_grid_deposit" yaml:"max_blocks_grid_deposit"`
	MaxBlocksGridRewardsWithdraw    int64          `json:"max_blocks_grid_rewards_withdraw" yaml:"max_blocks_grid_rewards_withdraw"`
	MaxBlocksGridDepositWithdraw    int64          `json:"max_blocks_grid_deposit_withdraw" yaml:"max_blocks_grid_deposit_withdraw"`
	RandNumberNegativeCriticalValue uint           `json:"rand_number_negative_critical_value" yaml:"rand_number_negative_critical_value"`
	MinDepositAmount                sdk.Coins      `json:"min_deposit_amount" yaml:"min_deposit_amount"`
	MaxDepositAmount                sdk.Coins      `json:"max_deposit_amount" yaml:"max_deposit_amount"`
	PerGridMaxDeposits              int            `json:"per_grid_max_deposits" yaml:"per_grid_max_deposits"`
	OwnerMinDeposit                 sdk.Coin       `json:"owner_min_deposit" yaml:"owner_min_deposit"`
	MemberMinDeposit                sdk.Coin       `json:"member_min_deposit" yaml:"member_min_deposit"`
	MaxPerDeposit                   sdk.Coin       `json:"max_per_deposit" yaml:"max_per_deposit"`
	OwnerRewardsRatio               sdk.Dec        `json:"owner_rewards_ratio" yaml:"owner_rewards_ratio"`
	FeeRatio                        sdk.Dec        `json:"ree_ratio" yaml:"ree_ratio"`
	LuckyPoolRatio                  sdk.Dec        `json:"lucky_pool_ratio" yaml:"lucky_pool_ratio"`
	DepositToLuckyPoolDigit         int            `json:"deposit_to_lucky_pool_digit" yaml:"deposit_to_lucky_pool_digit"`
	LuckyPoolRewardsDigit           int            `json:"lucky_pool_rewards_digit" yaml:"lucky_pool_rewards_digit"`
	RankType                        uint           `json:"rank_type" yaml:"rank_type"`
	Ranks                           Ranks          `json:"ranks" yaml:"ranks"`
	WinnerRewards                   WinnerRewards  `json:"winner_rewards" yaml:"winner_rewards"`
}
type Dapps []Dapp

func DefaultDapp() Dapp {
	return Dapp{
		MaxBlocksGridCreate:             99,
		GridCreateCanDeposit:            true,
		OnlyOwnerCanCreateGrid:          false,
		MaxBlocksGridDeposit:            99,
		PerGridMaxDeposits:              9,
		MaxBlocksGridRewardsWithdraw:    0,
		MaxBlocksGridDepositWithdraw:    0,
		RandNumberNegativeCriticalValue: 50,
		MinDepositAmount:                sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(0))),
		MaxDepositAmount:                sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(0))),
		OwnerMinDeposit:                 sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10000)),
		MemberMinDeposit:                sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000)),
		MaxPerDeposit:                   sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(1000000)),
		OwnerRewardsRatio:               sdk.NewDecWithPrec(20, 2),
		FeeRatio:                        sdk.NewDecWithPrec(20, 2),
		LuckyPoolRatio:                  sdk.NewDecWithPrec(20, 2),
		LuckyPoolRewardsDigit:           2,
		DepositToLuckyPoolDigit:         2,
		RankType:                        1,
		Ranks:                           DefaultRanks(),
		WinnerRewards:                   DefaultWinnerRewards(),
	}
}
func (c Dapp) String() string {
	return fmt.Sprintf(`Grid:
		MaxBlocksGridCreate:    %d
		MaxBlocksGridDeposit:   %d
		MemberMinDeposit:     %s`, c.MaxBlocksGridCreate, c.MaxBlocksGridDeposit, c.OwnerMinDeposit.String())
}

func (c Dapps) String() string {
	out := ""
	for _, v := range c {
		out += v.String()
	}
	return strings.TrimSpace(out)
}
