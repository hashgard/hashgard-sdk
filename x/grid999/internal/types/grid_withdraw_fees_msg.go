package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgGridWithdrawFees
type MsgGridWithdrawFees struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
}

var _ sdk.Msg = MsgGridWithdrawFees{}

// Route Implements Msg.
func (msg MsgGridWithdrawFees) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgGridWithdrawFees) Type() string { return "grid999_withdraw_fees" }

// ValidateBasic Implements Msg.
func (msg MsgGridWithdrawFees) ValidateBasic() sdk.Error {
	if msg.To == nil {
		return sdk.ErrInvalidAddress("missing to")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgGridWithdrawFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgGridWithdrawFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
