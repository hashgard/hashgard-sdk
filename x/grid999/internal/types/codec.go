package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDappCreateGrid{}, "cosmos-sdk/MsgDappCreateGrid", nil)
	cdc.RegisterConcrete(MsgDappDeposit{}, "cosmos-sdk/MsgDappDeposit", nil)
	cdc.RegisterConcrete(MsgDappWithdraw{}, "cosmos-sdk/MsgDappWithdraw", nil)
	cdc.RegisterConcrete(MsgDappGenerate{}, "cosmos-sdk/MsgDappGenerate", nil)
	cdc.RegisterConcrete(MsgDisableDapp{}, "cosmos-sdk/MsgDisableDapp", nil)
	cdc.RegisterConcrete(MsgDappWithdrawFees{}, "cosmos-sdk/MsgDappWithdrawFees", nil)
	cdc.RegisterConcrete(MsgDappWithdrawLucky{}, "cosmos-sdk/MsgDappWithdrawLucky", nil)
	cdc.RegisterConcrete(MsgGridParams{}, "cosmos-sdk/MsgGridParams", nil)

}

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
