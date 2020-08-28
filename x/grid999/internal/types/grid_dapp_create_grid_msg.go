package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"
)

// MsgDappCreateGrid
type MsgDappCreateGrid struct {
	Sender     sdk.AccAddress `json:"sender" yaml:"sender"`
	DappID     uint           `json:"dapp_id" yaml:"dapp_id"`
	Deposit    sdk.Coin       `json:"deposit" yaml:"deposit"`
	GridType   string         `json:"grid_type" yaml:"grid_type"`
	ZeroValued bool           `json:"zero_valued" yaml:"zero_valued"`
	Prepaid    string         `json:"prepaid" yaml:"prepaid"`
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
	if len(msg.Prepaid) > 0 {
		values := strings.Split(msg.Prepaid, ",")
		if len(values) > 100 {
			return sdk.ErrInvalidAddress("prepaid max length is 100")
		}
		for _, v := range values {
			prepaid := strings.Split(v, ":")
			if len(prepaid) != 3 {
				return sdk.ErrInvalidAddress("Prepaid is invalid,The format is like this 'address:height:coins,address:height:coins'")
			}
			_, err := sdk.AccAddressFromBech32(prepaid[0])
			if err != nil {
				return sdk.ErrInvalidAddress("Prepaid is invalid,The format is like this 'address:height:coins,address:height:coins'")
			}
			_, err = strconv.ParseInt(prepaid[1], 10, 64)
			if err != nil {
				return sdk.ErrInvalidAddress("Prepaid is invalid,The format is like this 'address:height:coins,address:height:coins'")
			}
			_, err = sdk.ParseCoins(prepaid[2])
			if err != nil {
				return sdk.ErrInvalidAddress("Prepaid is invalid,The format is like this 'address:height:coins,address:height:coins'")
			}
		}
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
