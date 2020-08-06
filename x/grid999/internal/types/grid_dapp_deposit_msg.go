package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappDeposit
type MsgDappDeposit struct {
	DappID  uint           `json:"dapp_id" yaml:"dapp_id"`
	Sender  sdk.AccAddress `json:"sender" yaml:"sender"`
	Deposit sdk.Coin       `json:"deposit" yaml:"deposit"`
	GridId  uint64         `json:"grid_id" yaml:"grid_id"`
	Index   uint           `json:"index" yaml:"index"`
}

var _ sdk.Msg = MsgDappDeposit{}

// Route Implements Msg.
func (msg MsgDappDeposit) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappDeposit) Type() string { return "grid999_dapp_deposit" }

// ValidateBasic Implements Msg.
func (msg MsgDappDeposit) ValidateBasic() sdk.Error {
	if msg.Deposit.IsZero() {
		return sdk.ErrInternal("deposit cannot be 0")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
