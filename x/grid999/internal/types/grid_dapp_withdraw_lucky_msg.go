package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappWithdrawLucky
type MsgDappWithdrawLucky struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	DappID uint           `json:"dapp_id" yaml:"dapp_id"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgDappWithdrawLucky{}

// Route Implements Msg.
func (msg MsgDappWithdrawLucky) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappWithdrawLucky) Type() string { return "grid999_dapp_withdraw_lucky" }

// ValidateBasic Implements Msg.
func (msg MsgDappWithdrawLucky) ValidateBasic() sdk.Error {
	if msg.Amount.IsZero() {
		return sdk.ErrInsufficientCoins("amount")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappWithdrawLucky) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappWithdrawLucky) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
