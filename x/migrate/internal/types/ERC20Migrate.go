package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ERC20MigrateExchange struct {
	ExchangeFrom sdk.AccAddress `json:"exchange_from" yaml:"exchange_from"`
	RewardsFrom  sdk.AccAddress `json:"rewards_from" yaml:"rewards_from"`
	Allows       string         `json:"allows" yaml:"allows"`
}

func (item ERC20MigrateExchange) String() string {
	return fmt.Sprintf(`ERC20MigrateExchange:
		ExchangeFrom:    %s
		RewardsFrom:    %s
		Allows:       %s`, item.ExchangeFrom.String(), item.RewardsFrom, item.Allows)
}
