package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappWithdraw
type MsgDappWithdraw struct {
	DappID uint           `json:"dapp_id" yaml:"dapp_id"`
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	GridId uint64         `json:"grid_id" yaml:"grid_id"`
}

var _ sdk.Msg = MsgDappWithdraw{}

// Route Implements Msg.
func (msg MsgDappWithdraw) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappWithdraw) Type() string { return "grid999_dapp_withdraw" }

// ValidateBasic Implements Msg.
func (msg MsgDappWithdraw) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
