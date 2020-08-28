package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDappGenerate
type MsgDappGenerate struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Dapp   Dapp           `json:"dapp" yaml:"dapp"`
}

var _ sdk.Msg = MsgDappGenerate{}

// Route Implements Msg.
func (msg MsgDappGenerate) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDappGenerate) Type() string { return "grid999_dapp_generate" }

// ValidateBasic Implements Msg.
func (msg MsgDappGenerate) ValidateBasic() sdk.Error {
	if !msg.Dapp.OwnerMinDeposit.IsValid() {
		return sdk.ErrInsufficientCoins("owner min deposit invalid")
	}
	if !msg.Dapp.MemberMinDeposit.IsValid() {
		return sdk.ErrInsufficientCoins("member min deposit invalid")
	}
	if msg.Dapp.MemberMinDeposit.IsZero() {
		return sdk.ErrInsufficientCoins("member min deposit cannot be 0")
	}
	if msg.Dapp.FeeRatio.IsNil() || msg.Dapp.FeeRatio.LT(sdk.ZeroDec()) {
		return sdk.ErrInsufficientCoins("fee ratio is invalid")
	}
	if msg.Dapp.OwnerRewardsRatio.IsNil() || msg.Dapp.OwnerRewardsRatio.LT(sdk.ZeroDec()) {
		return sdk.ErrInsufficientCoins("owner rewards ratio is invalid")
	}
	if msg.Dapp.LuckyPoolRatio.IsNil() || msg.Dapp.LuckyPoolRatio.LT(sdk.ZeroDec()) {
		return sdk.ErrInsufficientCoins("lucky pool ratio is invalid")
	}
	if msg.Dapp.LuckyPoolRewardsDigit > 10 || msg.Dapp.DepositToLuckyPoolDigit > 10 {
		return sdk.ErrInternal("digit is invalid")
	}
	if checkDec(msg.Dapp.FeeRatio) {
		return sdk.ErrInternal("fee ratio is invalid")
	}
	if checkDec(msg.Dapp.OwnerRewardsRatio) {
		return sdk.ErrInternal("owner rewards ratio is invalid")
	}
	if checkDec(msg.Dapp.LuckyPoolRatio) {
		return sdk.ErrInternal("lucky pool ratio is invalid")
	}
	if msg.Dapp.RandNumberNegativeCriticalValue < 0 || msg.Dapp.RandNumberNegativeCriticalValue > 99 {
		return sdk.ErrInternal("rand number negative critical value is [0-99]")
	}
	if len(msg.Dapp.Name) > 20 {
		return sdk.ErrInternal("name max length is 20")
	}
	if len(msg.Dapp.DappType) > 20 {
		return sdk.ErrInternal("dapp type max length is 20")
	}
	if len(msg.Dapp.Icon) > 100 {
		return sdk.ErrInternal("icon max length is 100")
	}
	if len(msg.Dapp.Desc) > 300 {
		return sdk.ErrInternal("desc max length is 300")
	}
	if err := msg.Dapp.Ranks.ValidateBasic(); err != nil {
		return err
	}
	if err := msg.Dapp.WinnerRewards.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDappGenerate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDappGenerate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
