package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/migrate/params"
)

const (
	ERC20MigrateType = "erc20-migrate"
)

// MsgERC20Migrate
type MsgERC20Migrate struct {
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
	Erc20TxHash  string         `json:"erc20_tx_hash" yaml:"erc20_tx_hash"`
	EthTxHash    string         `json:"eth_tx_hash" yaml:"eth_tx_hash"`
	Erc20Address string         `json:"erc20_address" yaml:"erc20_address"`
	To           sdk.AccAddress `json:"to" yaml:"to"`
	Amount       sdk.Coins      `json:"amount" yaml:"amount"`
	Channel      sdk.AccAddress `json:"channel" yaml:"channel"`
}

var _ sdk.Msg = MsgERC20Migrate{}

func NewMsgERC20Migrate(sender sdk.AccAddress, para params.ERC20MigrateParams) MsgERC20Migrate {
	return MsgERC20Migrate{
		Sender:       sender,
		Erc20TxHash:  para.Erc20TxHash,
		EthTxHash:    para.EthTxHash,
		Erc20Address: para.Erc20Address,
		To:           para.To,
		Amount:       para.Amount,
		Channel:      para.Channel}
}

// Route Implements Msg.
func (msg MsgERC20Migrate) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgERC20Migrate) Type() string { return ERC20MigrateType }

// ValidateBasic Implements Msg.
func (msg MsgERC20Migrate) ValidateBasic() sdk.Error {
	if len(msg.Erc20TxHash) == 0 {
		return sdk.ErrInvalidAddress("missing erc20_tx_hash")
	}
	if len(msg.EthTxHash) == 0 {
		return sdk.ErrInvalidAddress("missing eth_tx_hash")
	}
	if len(msg.Erc20Address) == 0 {
		return sdk.ErrInvalidAddress("missing erc20_address")
	}
	if msg.To == nil {
		return sdk.ErrInvalidAddress("missing to")
	}
	if msg.Amount == nil || msg.Amount.IsZero() {
		return sdk.ErrInvalidAddress("missing amount")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgERC20Migrate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgERC20Migrate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
