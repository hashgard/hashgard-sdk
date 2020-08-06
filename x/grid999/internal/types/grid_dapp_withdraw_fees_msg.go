package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappWithdrawFees
type MsgDappWithdrawFees struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	DappID uint           `json:"dapp_id" yaml:"dapp_id"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
}

var _ sdk.Msg = MsgDappWithdrawFees{}

// Route Implements Msg.
func (msg MsgDappWithdrawFees) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappWithdrawFees) Type() string { return "grid999_dapp_withdraw_fees" }

// ValidateBasic Implements Msg.
func (msg MsgDappWithdrawFees) ValidateBasic() sdk.Error {
	if msg.To == nil {
		return sdk.ErrInvalidAddress("missing to")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappWithdrawFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappWithdrawFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
