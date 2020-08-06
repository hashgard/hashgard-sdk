package params

import sdk "github.com/cosmos/cosmos-sdk/types"

type ERC20MigrateParams struct {
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
	Erc20TxHash  string         `json:"erc20_tx_hash" yaml:"erc20_tx_hash"`
	EthTxHash    string         `json:"eth_tx_hash" yaml:"eth_tx_hash"`
	Erc20Address string         `json:"sender" yaml:"sender"`
	To           sdk.AccAddress `json:"to" yaml:"to"`
	Amount       sdk.Coins      `json:"amount" yaml:"amount"`
	Channel      sdk.AccAddress `json:"channel" yaml:"channel"`
}
