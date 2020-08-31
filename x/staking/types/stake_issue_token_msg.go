package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// HashGard
// MsgStakeIssueToken - struct for stake issue token transactions
type MsgStakeIssueToken struct {
	Sender          sdk.AccAddress  `json:"sender" yaml:"sender"`
	StakeIssueToken StakeIssueToken `json:"stake_issue_token"  yaml:"stake_issue_token"`
}

//nolint
func (msg MsgStakeIssueToken) Route() string { return RouterKey }
func (msg MsgStakeIssueToken) Type() string  { return "stake_issue_token" }
func (msg MsgStakeIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// get the bytes for the message signer to sign on
func (msg MsgStakeIssueToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgStakeIssueToken) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if len(msg.StakeIssueToken.Denom) == 0 || !strings.HasPrefix(msg.StakeIssueToken.Denom, "u") {
		return sdk.ErrInternal("denom must not empty and start with u")
	}
	if strings.Contains(msg.StakeIssueToken.Denom, "gard") || strings.Contains(msg.StakeIssueToken.Denom, "ggt") {
		return sdk.ErrInternal("denom cannot contains gard or ggt")
	}
	coin := sdk.Coin{Denom: msg.StakeIssueToken.Denom, Amount: msg.StakeIssueToken.TotalSupply}
	if !coin.IsValid() {
		return sdk.ErrInternal("invalid coin: " + coin.String())
	}
	if len(msg.StakeIssueToken.Description.WholeName) < 2 || len(msg.StakeIssueToken.Description.WholeName) > 10 {
		return sdk.ErrInternal("whole name length is 2-10")
	}
	if msg.StakeIssueToken.PreMintAmount.LT(sdk.ZeroInt()) {
		return sdk.ErrUnknownRequest("invalid pre mint amount")
	}
	if len(msg.StakeIssueToken.Description.Details) > 255 {
		return sdk.ErrInternal("description max length is 255")
	}
	if len(msg.StakeIssueToken.Description.Website) > 100 {
		return sdk.ErrInternal("website max length is 100")
	}
	if len(msg.StakeIssueToken.Description.Icon) > 100 {
		return sdk.ErrInternal("icon max length is 100")
	}
	if len(msg.StakeIssueToken.PerBlockMint) == 0 {
		return sdk.ErrInternal("invalid per block mint")
	}
	if msg.StakeIssueToken.GenesisHeight < StakeIssueTokenStartBlock {
		return sdk.ErrInternal(fmt.Sprintf("genesis height must greater than %d", StakeIssueTokenStartBlock))
	}
	if sdk.ZeroInt().GTE(msg.StakeIssueToken.TotalSupply) {
		return sdk.ErrInternal("invalid total supply")
	}
	if msg.StakeIssueToken.PreMintAmount.IsNegative() {
		return sdk.ErrInternal("invalid pre mint amount")
	}
	if !msg.StakeIssueToken.PreMintAmount.IsZero() && msg.StakeIssueToken.PreMintAmount.GT(msg.StakeIssueToken.TotalSupply) {
		return sdk.ErrInternal("invalid pre mint amount")
	}
	preBlockHeight := int64(-1)
	for k, v := range msg.StakeIssueToken.PerBlockMint {
		if v.ProposerNodeAmount.IsNegative() {
			return sdk.ErrInternal("invalid proposer node amount")
		}
		if v.VoterNodeAmount.IsNegative() {
			return sdk.ErrInternal("invalid voter node amount")
		}
		if v.VoterNodeAmount.Add(v.ProposerNodeAmount).GT(msg.StakeIssueToken.TotalSupply.Sub(msg.StakeIssueToken.PreMintAmount)) {
			return sdk.ErrInternal("per block mint amount cannot be greater than totalSupply")
		}
		if k == 0 {
			if v.StartHeight != msg.StakeIssueToken.GenesisHeight {
				return sdk.ErrInternal("first mint block must eq genesis height")
			}
			preBlockHeight = v.StartHeight
			continue
		}
		if preBlockHeight >= v.StartHeight {
			return sdk.ErrInternal("per block mint height should be incremental sequence")
		}
	}
	return nil
}

// MsgStakeIssueTokenConfig - struct for stake issue token transactions
type MsgStakeIssueTokenConfig struct {
	Sender sdk.AccAddress        `json:"sender" yaml:"sender"`
	Config StakeIssueTokenConfig `json:"config"  yaml:"config"`
}

//nolint
func (msg MsgStakeIssueTokenConfig) Route() string { return RouterKey }
func (msg MsgStakeIssueTokenConfig) Type() string  { return "stake_issue_token_config" }
func (msg MsgStakeIssueTokenConfig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// get the bytes for the message signer to sign on
func (msg MsgStakeIssueTokenConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgStakeIssueTokenConfig) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if !strings.Contains(StakeIssueTokenConfigAllows, msg.Sender.String()) {
		return sdk.ErrInternal("Permission denied")
	}
	return nil
}
