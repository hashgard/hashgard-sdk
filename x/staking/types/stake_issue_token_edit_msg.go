package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// HashGard
// MsgStakeIssueTokenEdit - struct for stake issue token transactions
type MsgStakeIssueTokenEdit struct {
	Sender      sdk.AccAddress             `json:"sender" yaml:"sender"`
	Denom       string                     `json:"denom"  yaml:"denom"`
	Description StakeIssueTokenDescription `json:"description"  yaml:"description"`
}

//nolint
func (msg MsgStakeIssueTokenEdit) Route() string { return RouterKey }
func (msg MsgStakeIssueTokenEdit) Type() string  { return "stake_issue_token_edit" }
func (msg MsgStakeIssueTokenEdit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// get the bytes for the message signer to sign on
func (msg MsgStakeIssueTokenEdit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgStakeIssueTokenEdit) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.Denom) == 0 || !strings.HasPrefix(msg.Denom, "u") {
		return sdk.ErrInternal("denom must not empty and start with u")
	}
	if len(msg.Description.Details) > 255 {
		return sdk.ErrInternal("description max length is 255")
	}
	if len(msg.Description.Website) > 100 {
		return sdk.ErrInternal("website max length is 100")
	}
	if len(msg.Description.Icon) > 100 {
		return sdk.ErrInternal("icon max length is 100")
	}
	if len(msg.Description.WholeName) < 2 || len(msg.Description.WholeName) > 10 {
		return sdk.ErrInternal("whole name length is 2-10")
	}
	return nil
}
