package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ERC20MigrateExchangeAllowType = "erc20-migrate-exchange"
)

// MsgERC20MigrateExchange
type MsgERC20MigrateExchange struct {
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
	ExchangeFrom sdk.AccAddress `json:"exchange_from" yaml:"exchange_from"`
	Allows       string         `json:"allows" yaml:"allows"`
}

var _ sdk.Msg = MsgERC20MigrateExchange{}

func NewERC20MigrateExchange(sender, exchangeFrom sdk.AccAddress, allows string) MsgERC20MigrateExchange {
	return MsgERC20MigrateExchange{
		Sender:       sender,
		ExchangeFrom: exchangeFrom,
		Allows:       allows}
}

// Route Implements Msg.
func (msg MsgERC20MigrateExchange) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgERC20MigrateExchange) Type() string { return ERC20MigrateExchangeAllowType }

// ValidateBasic Implements Msg.
func (msg MsgERC20MigrateExchange) ValidateBasic() sdk.Error {
	if len(msg.Allows) == 0 {
		return sdk.ErrInvalidAddress("missing allows")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgERC20MigrateExchange) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgERC20MigrateExchange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
