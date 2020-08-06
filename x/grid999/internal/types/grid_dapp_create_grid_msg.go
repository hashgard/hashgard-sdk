package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappCreateGrid
type MsgDappCreateGrid struct {
	Sender     sdk.AccAddress `json:"sender" yaml:"sender"`
	DappID     uint           `json:"dapp_id" yaml:"dapp_id"`
	Deposit    sdk.Coin       `json:"deposit" yaml:"deposit"`
	Locked     bool           `json:"locked" yaml:"locked"`
	ZeroValued bool           `json:"zero_valued" yaml:"zero_valued"`
}

var _ sdk.Msg = MsgDappCreateGrid{}

// Route Implements Msg.
func (msg MsgDappCreateGrid) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappCreateGrid) Type() string { return "grid999_dapp_create_grid" }

// ValidateBasic Implements Msg.
func (msg MsgDappCreateGrid) ValidateBasic() sdk.Error {
	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidAddress("deposit is invalid")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappCreateGrid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappCreateGrid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
