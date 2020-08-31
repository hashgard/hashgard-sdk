package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// HashGard
type StakeIssueTokenConfig struct {
	LockCoins         sdk.Coins `json:"lock_coins" yaml:"lock_coins"`
	LockPeriodHeight  int64     `json:"lock_period_height" yaml:"lock_period_height"`
	MinSelfDelegation sdk.Int   `json:"min_self_delegation"  yaml:"min_self_delegation"`
}

type PerBlockMint struct {
	StartHeight        int64   `json:"start_height" yaml:"start_height"`
	ProposerNodeAmount sdk.Int `json:"proposer_node_amount"  yaml:"proposer_node_amount"`
	VoterNodeAmount    sdk.Int `json:"voter_node_amount"  yaml:"voter_node_amount"`
}
type StakeIssueTokenDescription struct {
	WholeName string `json:"whole_name"  yaml:"whole_name"`
	Website   string `json:"website" yaml:"website"`
	Icon      string `json:"icon" yaml:"icon"`
	Details   string `json:"details" yaml:"details"`
}

// StakeIssueToken - struct for stake issue token transactions
type StakeIssueToken struct {
	Owner          sdk.AccAddress             `json:"owner"  yaml:"owner"`
	Val            sdk.ValAddress             `json:"val"  yaml:"val"`
	Denom          string                     `json:"denom"  yaml:"denom"`
	TotalSupply    sdk.Int                    `json:"total_supply"  yaml:"total_supply"`
	PreMintAddress sdk.AccAddress             `json:"pre_mint_address" yaml:"pre_mint_address"`
	PreMintAmount  sdk.Int                    `json:"pre_mint_amount" yaml:"pre_mint_amount"`
	GenesisHeight  int64                      `json:"genesis_height"  yaml:"genesis_height"`
	PerBlockMint   []PerBlockMint             `json:"per_block_mint" yaml:"per_block_mint"`
	LockedAddr     sdk.AccAddress             `json:"locked_addr"  yaml:"locked_addr"`
	Description    StakeIssueTokenDescription `json:"description"  yaml:"description"`
	Height         int64                      `json:"height"  yaml:"height"`
	Rule           StakeIssueTokenConfig      `json:"rule"  yaml:"rule"`
}

type StakeIssueTokens []StakeIssueToken

func (p StakeIssueToken) String() string {
	return fmt.Sprintf(`StakeIssueToken:
  Denom:    %s
  Owner:    %s
  Val:    %s
  TotalSupply:       %s
  Description: %s`, p.Denom, p.Owner.String(), p.Val.String(), p.TotalSupply.String(), p.Description)
}
func (p StakeIssueTokenConfig) String() string {
	return fmt.Sprintf(`StakeIssueTokenConfig:
  LockCoins:    %s
  MinSelfDelegation:       %s`, p.LockCoins, p.MinSelfDelegation)
}
func (c StakeIssueTokens) String() string {
	out := ""
	for _, v := range c {
		out += v.String()
	}
	return strings.TrimSpace(out)
}
