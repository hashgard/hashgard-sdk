package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDisableDapp
type MsgDisableDapp struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	DappID uint           `json:"dapp_id" yaml:"dapp_id"`
	Height int64          `json:"height" yaml:"height"`
}

var _ sdk.Msg = MsgDisableDapp{}

// Route Implements Msg.
func (msg MsgDisableDapp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDisableDapp) Type() string { return "grid999_dapp_disable" }

// ValidateBasic Implements Msg.
func (msg MsgDisableDapp) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDisableDapp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDisableDapp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
