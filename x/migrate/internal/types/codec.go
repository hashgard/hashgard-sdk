package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgERC20MigrateExchange{}, "cosmos-sdk/MsgERC20MigrateExchange", nil)
	cdc.RegisterConcrete(MsgERC20Migrate{}, "cosmos-sdk/MsgERC20Migrate", nil)
}

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
