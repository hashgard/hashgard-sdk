package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgGridParams
type MsgGridParams struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Params Params         `json:"params" yaml:"params"`
}

var _ sdk.Msg = MsgGridParams{}

// Route Implements Msg.
func (msg MsgGridParams) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgGridParams) Type() string { return "grid999_params" }

// ValidateBasic Implements Msg.
func (msg MsgGridParams) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgGridParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgGridParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
